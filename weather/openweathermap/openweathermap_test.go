package openweathermap_test

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.cwx.io/mrubiosan/mx51/weather/openweathermap"
)

type TransportMock struct {
	roundTripFn func(request *http.Request) (*http.Response, error)
}

func (t *TransportMock) RoundTrip(request *http.Request) (*http.Response, error) {
	return t.roundTripFn(request)
}

func TestItGetsReport(t *testing.T) {
	expectedWindSpeed := 30
	expectedTemperature := 25
	client := &http.Client{Transport: &TransportMock{
		roundTripFn: func(request *http.Request) (*http.Response, error) {
			b := strings.NewReader(fmt.Sprintf(`{"main": {"temp": %f},"wind":{"speed": %f}}`, 24.6, 30.3))
			res := &http.Response{
				Status:     "OK",
				StatusCode: 200,
				Body:       io.NopCloser(b),
			}

			return res, nil
		},
	}}

	ws := &openweathermap.OpenWeatherMap{Client: client}
	r, err := ws.Report("Sydney")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if r.WindSpeed != expectedWindSpeed {
		t.Errorf("expected wind speed %d but got %d", expectedWindSpeed, r.WindSpeed)
	}

	if r.Temperature != expectedTemperature {
		t.Errorf("expected Temperature %d but got %d", expectedTemperature, r.Temperature)
	}

}
