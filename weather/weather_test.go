package weather_test

import (
	"errors"
	"testing"

	"github.cwx.io/mrubiosan/mx51/weather"
)

type WeatherSourceMock struct {
	reportFn func(loc weather.Location) (weather.Report, error)
	calls    int
}

func (w *WeatherSourceMock) Report(loc weather.Location) (weather.Report, error) {
	defer func() { w.calls++ }()
	return w.reportFn(loc)
}

func TestItFallsBackOnFailure(t *testing.T) {
	main := &WeatherSourceMock{
		reportFn: func(loc weather.Location) (weather.Report, error) {
			return weather.Report{}, errors.New("main mock failure")
		},
	}

	expectedReport := weather.Report{
		WindSpeed:   1,
		Temperature: 3,
	}
	next := &WeatherSourceMock{
		reportFn: func(loc weather.Location) (weather.Report, error) {
			return expectedReport, nil
		},
	}

	cs := weather.ChainedSource{
		Main: main,
		Next: next,
	}

	actualReport, err := cs.Report("foo")
	if err != nil {
		t.Fatalf("Unexpected failure: %s", err)
	}

	if actualReport != expectedReport {
		t.Errorf("Expected report %v but got %v", expectedReport, actualReport)
	}
}
