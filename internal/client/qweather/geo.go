package qweather

import (
	"context"
	"fmt"
	"net/url"
)

// CitySearchResponse represents the city search API response
type CitySearchResponse struct {
	Code     string     `json:"code"`
	Location []Location `json:"location"`
	Refer    Refer      `json:"refer"`
}

// Location contains city/location information
type Location struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Lat       string `json:"lat"`
	Lon       string `json:"lon"`
	Adm2      string `json:"adm2"`
	Adm1      string `json:"adm1"`
	Country   string `json:"country"`
	Tz        string `json:"tz"`
	UtcOffset string `json:"utcOffset"`
	IsDst     string `json:"isDst"`
	Type      string `json:"type"`
	Rank      string `json:"rank"`
	FxLink    string `json:"fxLink"`
}

// SearchCity searches for cities by name or location
func (c *Client) SearchCity(ctx context.Context, query string) (*CitySearchResponse, error) {
	params := url.Values{}
	params.Set("location", query)

	var result CitySearchResponse
	if err := c.doRequest(ctx, "/geo/v2/city/lookup", params, &result); err != nil {
		return nil, fmt.Errorf("search city: %w", err)
	}

	// Check API response code
	if result.Code != "200" {
		return nil, fmt.Errorf("API returned error code: %s", result.Code)
	}

	return &result, nil
}
