# How to run
```
docker build -t mrubiosan/mx51 .
docker run --rm -p 8080:8080 -ti mrubiosan/mx51
```

## Environment variables
These can be modified if needed
* OPENWEATHERMAP_KEY: OpenWeatherMap API Key
* WEATHERSTACK_KEY: WeatherStack API Key
* CACHE_TTL: cache TTL in seconds
