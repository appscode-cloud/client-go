# AGENTS.md

This file provides guidance to coding agents (e.g. Claude Code, claude.ai/code) when working with code in this repository.

## Repository purpose

Go module `go.bytebuilders.dev/client` — the Go SDK for the [ByteBuilders](https://byte.builders) (b3) cloud platform API. Library only; consumed by ACE-side controllers, internal admin tools, and integration tests.

Note: the GitHub repo is `bytebuilders/client-go`, but the Go module path is `go.bytebuilders.dev/client` (no `-go`). Use the module path in imports.

## Architecture

- Top-level Go files: `client.go` (`Client` constructor + HTTP plumbing), `auth.go` (auth helpers), `cluster.go` (cluster API), `ace_license.go` (license endpoints), plus matching `_test.go` files for each.
- `api/` — generated request/response types for the b3 API.
- `examples/` — usage samples.
- `hack/` — codegen / linting helpers.
- `doc.go` — package documentation.

This is a **library**: no binaries, no Docker images.

## Common commands

- `make ci` — full CI pipeline.
- `make fmt`, `make lint`, `make unit-tests` / `make test` — standard.
- `make verify` — codegen + module-tidy verification.
- `make add-license` / `make check-license` — manage license headers.

Run a single Go test:

```
go test . -run TestName -v
```

## Conventions

- Module path is `go.bytebuilders.dev/client` (vanity URL, **not** `client-go`); imports must use that.
- License: see `LICENSE`. Sign off commits (`git commit -s`); contributions follow the DCO (`DCO`).
- Every exported symbol is API — don't break downstream consumers (`b3`, ACE controllers) without coordinating.
- `Client` methods follow the resource-noun + verb pattern; new resource families go in their own top-level file (`cluster.go`, `auth.go`, …) plus matching types under `api/`.
- The README badges point at `byte.builders` — the consumer-facing product brand stays even if internal names drift.
