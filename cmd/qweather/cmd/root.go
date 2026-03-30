package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	formatFlag  string
	verboseFlag bool

	// logger is the global structured logger; level is set during PersistentPreRun.
	logger *slog.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qweather",
	Short: "Weather CLI powered by QWeather API",
	Long: `A command-line tool to query weather information using QWeather API.

Supports current weather, weather forecast, and city search.
Requires QWeather API key to be configured.`,
	Version:           "0.1.0",
	PersistentPreRunE: initLogger,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&formatFlag, "format", "f", "text", "Output format: text, json, table")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")
}

// initLogger sets up the global slog logger based on the --verbose flag.
// All diagnostic output goes to stderr so it never pollutes stdout.
func initLogger(cmd *cobra.Command, args []string) error {
	level := slog.LevelWarn
	if verboseFlag {
		level = slog.LevelDebug
	}
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	}))
	return nil
}

// printError prints error message to stderr
func printError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}
