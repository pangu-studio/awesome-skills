package cmd

import (
	"context"
	"fmt"
	"testing"

	"github.com/pangu-studio/awesome-skills/internal/client/qweather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockGeoService is a test double for qweather.GeoService.
type mockGeoService struct {
	result *qweather.CitySearchResponse
	err    error
}

func (m *mockGeoService) SearchCity(_ context.Context, _ string) (*qweather.CitySearchResponse, error) {
	return m.result, m.err
}

func TestResolveLocation_DirectLocationID(t *testing.T) {
	// Arrange
	svc := &mockGeoService{}

	// Act
	id, err := resolveLocation(context.Background(), svc, "", "101010100")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "101010100", id)
}

func TestResolveLocation_DirectCoordinates(t *testing.T) {
	// Arrange
	svc := &mockGeoService{}

	// Act
	id, err := resolveLocation(context.Background(), svc, "", "116.41,39.92")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "116.41,39.92", id)
}

func TestResolveLocation_CityFound(t *testing.T) {
	// Arrange
	svc := &mockGeoService{
		result: &qweather.CitySearchResponse{
			Code: "200",
			Location: []qweather.Location{
				{ID: "101010100", Name: "北京"},
				{ID: "101010200", Name: "北京朝阳"},
			},
		},
	}

	// Act
	id, err := resolveLocation(context.Background(), svc, "北京", "")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "101010100", id, "should return the first match")
}

func TestResolveLocation_CityNotFound(t *testing.T) {
	// Arrange
	svc := &mockGeoService{
		result: &qweather.CitySearchResponse{
			Code:     "200",
			Location: []qweather.Location{},
		},
	}

	// Act
	id, err := resolveLocation(context.Background(), svc, "nonexistent", "")

	// Assert
	require.Error(t, err)
	assert.Empty(t, id)
	assert.Contains(t, err.Error(), "no location found")
}

func TestResolveLocation_SearchError(t *testing.T) {
	// Arrange
	svc := &mockGeoService{
		err: fmt.Errorf("network error"),
	}

	// Act
	id, err := resolveLocation(context.Background(), svc, "北京", "")

	// Assert
	require.Error(t, err)
	assert.Empty(t, id)
	assert.Contains(t, err.Error(), "failed to search city")
}

func TestResolveLocation_NeitherProvided(t *testing.T) {
	// Arrange
	svc := &mockGeoService{}

	// Act
	id, err := resolveLocation(context.Background(), svc, "", "")

	// Assert
	require.Error(t, err)
	assert.Empty(t, id)
	assert.Contains(t, err.Error(), "--location or --city is required")
}

func TestMaskAPIKey_NormalKey(t *testing.T) {
	testCases := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "standard 32-char key",
			key:      "abcd1234efgh5678ijkl9012mnop3456",
			expected: "abcd" + "************************" + "3456",
		},
		{
			name:     "minimum visible key (9 chars)",
			key:      "abc123xyz",
			expected: "abc1" + "*" + "3xyz",
		},
		{
			name:     "exactly 8 chars — fully masked",
			key:      "12345678",
			expected: "********",
		},
		{
			name:     "less than 8 chars",
			key:      "abc",
			expected: "***",
		},
		{
			name:     "empty key",
			key:      "",
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result := maskAPIKey(tc.key)

			// Assert
			assert.Equal(t, tc.expected, result)
		})
	}
}
