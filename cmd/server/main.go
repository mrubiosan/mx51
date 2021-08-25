package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.cwx.io/mrubiosan/mx51/weather"
	"github.cwx.io/mrubiosan/mx51/weather/cache"
	"github.cwx.io/mrubiosan/mx51/weather/openweathermap"
	"github.cwx.io/mrubiosan/mx51/weather/weatherstack"
)

func main() {
	source := weatherSource()
	http.HandleFunc("/v1/weather", func(writer http.ResponseWriter, request *http.Request) {
		log.Printf("Incoming request: %s", request.RequestURI)
		queryParams := request.URL.Query()
		report, err := source.Report(weather.Location(queryParams.Get("city")))

		if err != nil {
			writer.WriteHeader(http.StatusServiceUnavailable)
			_, _ = writer.Write([]byte(err.Error()))
		} else {
			payload, _ := json.Marshal(report)
			_, _ = writer.Write(payload)
		}

	})

	log.Print("Starting server")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Print(err)
	}
}

func weatherSource() weather.Source {
	weatherStackKey := os.Getenv("WEATHERSTACK_KEY")
	if weatherStackKey == "" {
		log.Fatal("Missing WEATHERSTACK_KEY")
	}
	openWeatherMapKey := os.Getenv("OPENWEATHERMAP_KEY")
	if openWeatherMapKey == "" {
		log.Fatal("Missing OPENWEATHERMAP_KEY")
	}

	cacheTtl := 3
	if envCacheTtl := os.Getenv("CACHE_TTL"); envCacheTtl != "" {
		res, err := strconv.Atoi(envCacheTtl)
		if err != nil {
			log.Fatal(err)
		}
		cacheTtl = res
	}

	ws := &weatherstack.WeatherStack{ApiKey: weatherStackKey}
	ow := &openweathermap.OpenWeatherMap{ApiKey: openWeatherMapKey}
	cs := &weather.ChainedSource{
		Main: ws,
		Next: ow,
	}

	return cache.New(cs, time.Duration(cacheTtl)*time.Second)
}
