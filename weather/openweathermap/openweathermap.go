package openweathermap

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.cwx.io/mrubiosan/mx51/weather"
)

const baseUrl = "http://api.openweathermap.org/data/2.5/weather"

// OpenWeatherMap is an API consumer of https://openweathermap.org.
type OpenWeatherMap struct {
	ApiKey string
	Client *http.Client // an optional http client.
}

type response struct {
	Main responseMain `json:"main"`
	Wind responseWind `json:"wind"`
}

type responseMain struct {
	Temperature float32 `json:"temp"`
}

type responseWind struct {
	Speed float32 `json:"speed"`
}

func (o *OpenWeatherMap) Report(loc weather.Location) (weather.Report, error) {
	httpClient := o.Client
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	report := weather.Report{}

	httpResp, err := httpClient.Get(o.buildReportUrl(loc))
	if err != nil {
		return report, err
	}
	defer func() { _ = httpResp.Body.Close() }()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return report, fmt.Errorf("invalid http status code %d received", httpResp.StatusCode)
	}

	res, err := parseResponse(httpResp.Body)
	if err != nil {
		return report, fmt.Errorf("could not parse response: %w", err)
	}

	report.Temperature = int(res.Main.Temperature + 0.5) // Easy rounding
	report.WindSpeed = int(res.Wind.Speed + 0.5)

	return report, nil
}

func (o *OpenWeatherMap) buildReportUrl(loc weather.Location) string {
	query := url.Values{
		"appid": []string{o.ApiKey},
		"q":     []string{string(loc) + ",AU"},
		"units": []string{"metric"},
	}

	return baseUrl + "?" + query.Encode()
}

func parseResponse(r io.Reader) (response, error) {
	var parsedResponse response
	rawResponse, err := io.ReadAll(r)
	if err == nil {
		err = json.Unmarshal(rawResponse, &parsedResponse)
	}

	return parsedResponse, err
}
