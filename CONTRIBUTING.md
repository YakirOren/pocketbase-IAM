# Contributing

Thanks for your interest in contributing to pocketbase-IAM!

## Prerequisites

- Go 1.26+
- Node 22+

## Setup

```bash
git clone https://github.com/YakirOren/pocketbase-IAM.git
cd pocketbase-IAM
go run . serve
```

For frontend development:

```bash
cd ui
npm install
npm run dev
```

The Vite dev server runs at `http://localhost:5173/_/iam/` and proxies API requests to PocketBase at `:8090`.

## Running Tests

```bash
# All Go tests
go test ./iam/...

# A specific test
go test ./iam/... -run TestEvaluateStatements

# Frontend lint
cd ui && npm run lint
```

## Building the Dashboard

The dashboard UI is embedded in the Go binary. After making frontend changes:

```bash
cd ui
npm run build
```

This outputs to `iam/dashboard/`. Commit the built files along with your source changes.

## Project Structure

- `iam/` — Go library (the core package)
- `ui/` — React admin dashboard (Refine + Tailwind + shadcn/ui)
- `main.go` — Thin consumer that calls `iam.Setup()`
- `docs/` — VitePress documentation site

## Conventions

- Short commit messages
- Table-driven tests with `{name, input, expected}` structs
- Unit tests cover pure logic (evaluation, parsing, matching)

## Submitting Changes

1. Fork the repo and create a branch from `main`
2. Make your changes
3. Run `go test ./iam/...` and `cd ui && npm run lint`
4. If you changed the dashboard, run `cd ui && npm run build` and commit the output
5. Open a pull request

## Releases

Releases are automated via GitHub Actions. When a PR is merged to `main` with a `release:patch`, `release:minor`, or `release:major` label, the workflow tags and publishes a new release.
