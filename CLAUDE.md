# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

```bash
make build              # Build binary to bin/qweather
make build-linux        # Cross-compile for Linux amd64
make test               # Run all tests (go test -v ./...)
make test-cover         # Tests with coverage report (coverage.html)
make fmt                # Format code
make vet                # Run go vet
make lint               # fmt + vet + golangci-lint
make clean              # Remove bin/ and coverage files

# Run a single test
go test -run TestFunctionName ./internal/client/qweather/

# Run the CLI after building
./bin/qweather now --location "北京"
./bin/qweather forecast --city "上海" --days 7
./bin/qweather search --query "beijing"
```

## Architecture

Go CLI tool (`qweather`) providing weather data via the QWeather API. Built as a skill for AI agent platforms (OpenClaw) using the Cobra framework. Module path: `github.com/pangu-studio/awesome-skills`.

- **`cmd/qweather/`** — CLI entry point and Cobra commands (`now`, `forecast`, `search`, `config`). `root.go` defines global flags (`--format`, `--verbose`).
- **`internal/client/qweather/`** — QWeather API client. `client.go` handles HTTP with gzip and `X-QW-Api-Key` header auth. `weather.go` and `geo.go` define API methods and response types.
- **`internal/config/`** — Config loading with priority: env vars > `.env` file (godotenv) > config files at `~/.config/awesome-skill/qweather/`. Supports `XDG_CONFIG_HOME`.
- **`internal/output/`** — `Formatter` interface (strategy pattern) with `TextFormatter`, `JSONFormatter`, `TableFormatter`. Factory function `NewFormatter(format)` in `formatter.go`.
- **`skills/qweather/`** — Skill metadata (`skill.toml`) and usage docs (`SKILL.md`) for agent platform integration.

### Key Design Details

- `--location` and `--city` flags on `now`/`forecast` are **mutually exclusive** (enforced via `MarkFlagsMutuallyExclusive`). `--city` auto-resolves to a location ID via the search API before querying weather.
- All API calls use `context.WithTimeout(context.Background(), 30*time.Second)`.
- Default API base URL: `https://devapi.qweather.com` (free tier). Use `api.qweather.com` for paid tier.
- Config env vars: `QWEATHER_API_KEY` (required), `QWEATHER_API_HOST` (optional).
- Tests use `httptest.NewServer` for API mocking and `t.TempDir()` for filesystem isolation.

## Code Conventions

See AGENTS.md for full Go style guidelines. Key points:

- Import order: stdlib, external deps, internal packages (blank-line separated)
- Wrap errors: `fmt.Errorf("context: %w", err)`
- Tests: `testify/assert` + `testify/require`, table-driven tests with `t.Run()`
- Commits: conventional commits (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`)
