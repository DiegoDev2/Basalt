# Contributing to Basalt

We love contributions! Here's how you can help:

## Adding a New Framework Adapter

1. Add a new template in `internal/generator/templates.go`.
2. Update `internal/generator/generator.go` to handle the new framework option.
3. Update the `init` wizard in `internal/tui/init.go` to include the new framework.

## Development

1. Clone the repo.
2. Run `make build` to compile the CLI.
3. Run `make test` to ensure everything is working.
