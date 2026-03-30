package cmd

import (
	"context"
	"os"
	"time"

	"github.com/pangu-studio/awesome-skills/internal/client/qweather"
	"github.com/pangu-studio/awesome-skills/internal/config"
	"github.com/pangu-studio/awesome-skills/internal/output"
	"github.com/spf13/cobra"
)

var (
	forecastLocation string
	forecastCity     string
	forecastDays     int
)

// forecastCmd represents the forecast command
var forecastCmd = &cobra.Command{
	Use:   "forecast",
	Short: "Get weather forecast",
	Long: `Get daily weather forecast for a specified location.

Location can be specified using either --location or --city:
  - Location ID (e.g., "101010100")
  - Coordinates (e.g., "116.41,39.92")
  - City name (use --city, e.g., "北京", "Shanghai")

Forecast days can be: 3, 7, 10, 15, or 30

Examples:
  qweather forecast --location "101010100" --days 7
  qweather forecast --location "116.41,39.92" --days 15 --format json
  qweather forecast --city "北京" --days 3
  qweather forecast --city "shanghai" --days 7 --format table`,
	RunE: runForecast,
}

func init() {
	rootCmd.AddCommand(forecastCmd)
	forecastCmd.Flags().StringVarP(&forecastLocation, "location", "l", "", "Location ID or coordinates (required, mutually exclusive with --city)")
	forecastCmd.Flags().StringVarP(&forecastCity, "city", "c", "", "City name (auto-resolve to location ID, required, mutually exclusive with --location)")
	forecastCmd.Flags().IntVarP(&forecastDays, "days", "d", 3, "Forecast days: 3, 7, 10, 15, or 30")
	forecastCmd.MarkFlagsMutuallyExclusive("location", "city")
}

func runForecast(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		printError(err)
		return err
	}

	// Create API client
	client := qweather.NewClient(cfg.QWeather.APIKey, cfg.QWeather.APIHost, qweather.WithLogger(logger))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Resolve location ID (from --location or --city)
	location, err := resolveLocation(ctx, client, forecastCity, forecastLocation)
	if err != nil {
		printError(err)
		return err
	}

	if verboseFlag && forecastCity != "" {
		logger.Debug("resolved city to location ID", "city", forecastCity, "locationID", location)
	}

	// Get weather forecast
	forecastData, err := client.GetDailyForecast(ctx, location, forecastDays)
	if err != nil {
		printError(err)
		return err
	}

	// Format and print output
	formatter, err := output.NewFormatter(formatFlag)
	if err != nil {
		printError(err)
		return err
	}

	if err := formatter.FormatWeatherDaily(forecastData, os.Stdout); err != nil {
		printError(err)
		return err
	}

	return nil
}
