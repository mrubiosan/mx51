package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.cwx.io/mrubiosan/mx51/weather"
	"github.cwx.io/mrubiosan/mx51/weather/cache"
)

type WeatherSourceMock struct {
	reportFn func(loc weather.Location) (weather.Report, error)
	calls    int
}

func (w *WeatherSourceMock) Report(loc weather.Location) (weather.Report, error) {
	defer func() { w.calls++ }()
	return w.reportFn(loc)
}

func TestItDelegatesOnMiss(t *testing.T) {
	expectedReport := weather.Report{
		WindSpeed:   3,
		Temperature: 5,
	}
	sourceMock := &WeatherSourceMock{reportFn: func(loc weather.Location) (weather.Report, error) {
		return expectedReport, nil
	}}
	cs := cache.New(sourceMock, time.Second)
	r, err := cs.Report("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if r != expectedReport {
		t.Errorf("expected %v but got %v", expectedReport, r)
	}

	if sourceMock.calls != 1 {
		t.Error("Expected one call to source")
	}
}

func TestItReusesResultOnHit(t *testing.T) {
	expectedReport := weather.Report{
		WindSpeed:   3,
		Temperature: 5,
	}
	sourceMock := &WeatherSourceMock{reportFn: func(loc weather.Location) (weather.Report, error) {
		return expectedReport, nil
	}}
	// Get a report to warm up the cache
	cs := cache.New(sourceMock, time.Second)
	r, err := cs.Report("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if r != expectedReport {
		t.Errorf("expected %v but got %v", expectedReport, r)
	}

	// Get the report a second time
	r, err = cs.Report("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if r != expectedReport {
		t.Errorf("expected %v but got %v", expectedReport, r)
	}

	if sourceMock.calls != 1 {
		t.Errorf("Expected one call to source but got %d", sourceMock.calls)
	}
}

func TestItFlagsStaleReports(t *testing.T) {
	expectedReport := weather.Report{
		WindSpeed:   3,
		Temperature: 5,
	}

	var sourceMock *WeatherSourceMock

	sourceMock = &WeatherSourceMock{reportFn: func(loc weather.Location) (weather.Report, error) {
		if sourceMock.calls > 0 { // Make further calls fail
			return weather.Report{}, errors.New("test simulated failure")
		}
		return expectedReport, nil
	}}
	// Get a report to warm up the cache
	cs := cache.New(sourceMock, time.Nanosecond)
	r, err := cs.Report("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if r != expectedReport {
		t.Errorf("expected %v but got %v", expectedReport, r)
	}

	time.Sleep(time.Nanosecond) // Let the cache expire
	// Reset the mock and get the report a second time
	r, err = cs.Report("foo")
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	expectedReport.Stale = true
	if r != expectedReport {
		t.Errorf("expected %v but got %v", expectedReport, r)
	}

	if sourceMock.calls != 2 {
		t.Errorf("Expected 2 calls to source but got %d", sourceMock.calls)
	}
}
