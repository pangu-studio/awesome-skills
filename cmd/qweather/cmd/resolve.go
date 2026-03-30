package cmd

import (
	"context"
	"fmt"

	"github.com/pangu-studio/awesome-skills/internal/client/qweather"
)

// resolveLocation returns the QWeather location ID to use for API calls.
// If locationID is provided directly, it is returned as-is.
// If city is provided, it performs a city search and returns the first match's ID.
// If neither is provided, an error is returned.
func resolveLocation(ctx context.Context, svc qweather.GeoService, city, locationID string) (string, error) {
	if locationID != "" {
		return locationID, nil
	}

	if city == "" {
		return "", fmt.Errorf("either --location or --city is required")
	}

	searchData, err := svc.SearchCity(ctx, city)
	if err != nil {
		return "", fmt.Errorf("failed to search city %q: %w", city, err)
	}

	if len(searchData.Location) == 0 {
		return "", fmt.Errorf("no location found for city %q, try a different spelling or use --location with a location ID directly", city)
	}

	return searchData.Location[0].ID, nil
}
