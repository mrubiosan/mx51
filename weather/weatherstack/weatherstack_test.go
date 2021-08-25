package weatherstack

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

type TransportMock struct {
	roundTripFn func(request *http.Request) (*http.Response, error)
}

func (t *TransportMock) RoundTrip(request *http.Request) (*http.Response, error) {
	return t.roundTripFn(request)
}

func TestItGetsReport(t *testing.T) {
	expectedWindSpeed := 123
	expectedTemperature := 25
	client := &http.Client{Transport: &TransportMock{
		roundTripFn: func(request *http.Request) (*http.Response, error) {
			b := strings.NewReader(fmt.Sprintf(`{"current":{"wind_speed": %d, "temperature": %d}}`, expectedWindSpeed*3600/1000, expectedTemperature))
			res := &http.Response{
				Status:     "OK",
				StatusCode: 200,
				Body:       io.NopCloser(b),
			}

			return res, nil
		},
	}}

	ws := &WeatherStack{ApiKey: "fooKey", Client: client}
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
