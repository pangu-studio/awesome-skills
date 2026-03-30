package qweather

import "context"

// WeatherService defines the contract for weather data operations.
type WeatherService interface {
	GetNowWeather(ctx context.Context, location string) (*WeatherNowResponse, error)
	GetDailyForecast(ctx context.Context, location string, days int) (*WeatherDailyResponse, error)
}

// GeoService defines the contract for location/city search operations.
type GeoService interface {
	SearchCity(ctx context.Context, query string) (*CitySearchResponse, error)
}

// Service combines all weather and geo capabilities.
type Service interface {
	WeatherService
	GeoService
}

// Compile-time check: Client must satisfy Service.
var _ Service = (*Client)(nil)
