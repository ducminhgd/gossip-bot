# Project overview

## Dependencies

1. Go version 1.24 or upper
2. Github action

## Project layout

- `.ai` directory contains intructions files for AI Agent.
    - `main.md` is the main instruction file for AI Agent.
    - `features.md` contains the list of features for this projects.
- `.github` directory contains github action files.
- `go.mod` and `go.sum` are the Go module files.
- `.gitignore` is the git ignore file.
- `cmd` (required): Contains commands to run the application. For example, `cmd/bot/main.go` is the entry point for the bot
- `internal` (required): Contains internal packages used exclusively within this project.
    - `models` (optional): Contains data object definitions, data models from data sources, or ORMs. Model files MUST be named like `{singular_name}.go`.
    - `repositories` (optional): Contains repository structs and functions to interact with databases. Repository files MUST be named like `{singular_name}.go`.
    - `require` (required): Contains essential packages that the services cannot run without.
    - `services` (required): Defines the services of this API, such as user services and order services. Service files MUST be named like `{singular_name}.go`. 
    - `utils` (optional): Contains utility packages.
- `pkg` (optional): Contains packages that are exposed for use by other services or clients.