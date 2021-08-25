package weatherstack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.cwx.io/mrubiosan/mx51/weather"
)

const baseUrl = "http://api.weatherstack.com"

// WeatherStack is an API consumer of https://weatherstack.com/.
type WeatherStack struct {
	ApiKey string
	Client *http.Client // an optional http client.
}

type response struct {
	Current responseCurrent `json:"current"`
	Error   responseError   `json:"error"`
}

type responseCurrent struct {
	WindSpeed   int `json:"wind_speed"`
	Temperature int `json:"temperature"`
}

type responseError struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

func (w *WeatherStack) Report(loc weather.Location) (weather.Report, error) {
	httpClient := w.Client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	report := weather.Report{}

	httpResp, err := httpClient.Get(w.buildReportUrl(loc))
	if err != nil {
		return report, err
	}
	defer func() { _ = httpResp.Body.Close() }()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return report, fmt.Errorf("invalid HTTP status code %d received", httpResp.StatusCode)
	}

	res, err := parseResponse(httpResp.Body)
	if err != nil {
		return report, fmt.Errorf("could not parse response: %w", err)
	}
	if res.Error.Code > 0 {
		return report, fmt.Errorf("API error: %s", res.Error.Info)
	}

	report.Temperature = res.Current.Temperature
	// Convert km/h to m/s
	report.WindSpeed = int((float32(res.Current.WindSpeed) * 1000 / 3600) + 0.5)

	return report, nil
}

func (w *WeatherStack) buildReportUrl(loc weather.Location) string {
	query := url.Values{
		"access_key": []string{w.ApiKey},
		"query":      []string{string(loc)},
	}

	return baseUrl + "/current?" + query.Encode()
}

func parseResponse(r io.Reader) (response, error) {
	var parsedResponse response
	rawResponse, err := io.ReadAll(r)
	if err == nil {
		err = json.Unmarshal(rawResponse, &parsedResponse)
	}

	return parsedResponse, err
}
