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
	searchQuery string
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search for cities",
	Long: `Search for cities by name or location.

The search supports fuzzy matching and returns multiple results.

Examples:
  qweather search --query "北京"
  qweather search --query "beijing" --format table
  qweather search --query "london" --format json`,
	RunE: runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().StringVarP(&searchQuery, "query", "q", "", "Search query (required)")
	searchCmd.MarkFlagRequired("query")
}

func runSearch(cmd *cobra.Command, args []string) error {
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

	// Search for cities
	searchData, err := client.SearchCity(ctx, searchQuery)
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

	if err := formatter.FormatCitySearch(searchData, os.Stdout); err != nil {
		printError(err)
		return err
	}

	return nil
}
