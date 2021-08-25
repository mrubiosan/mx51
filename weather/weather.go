package weather

import (
	"log"
)

// Location is a place to generate reports for. Australian cities only thus far.
type Location string

// A Report has the relevant weather information.
type Report struct {
	WindSpeed   int  `json:"wind_speed"`          // in meters per second.
	Temperature int  `json:"temperature_degrees"` // in celsius.
	Stale       bool `json:"stale,omitempty"`     // whether the report is stale.
}

// Source is a provider of weather data.
type Source interface {
	// Report generates a weather report for the given Location.
	Report(loc Location) (Report, error)
}

// ChainedSource allows chaining sources. If a request to Main returns an error,
// then the result of a call to Next is returned.
type ChainedSource struct {
	Main Source // main source
	Next Source // fallback source
}

func (c *ChainedSource) Report(loc Location) (Report, error) {
	r, err := c.Main.Report(loc)
	if err != nil && c.Next != nil {
		log.Printf("Weather Source %T failed. Falling back to %T", c.Main, c.Next)
		return c.Next.Report(loc)
	}
	return r, err
}
