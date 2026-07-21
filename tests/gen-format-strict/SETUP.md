# Scenario

**Feature**: Phase A strict gen-format — no stdlib auto-add, prune unused, harness-only injects

```
# fixture tree
author SETUP (explicit imports only)
  -> Assemble unified packages
  -> WriteFormattedGo (no stdlibByPkgName; unused prune OK)
  -> go test suite

# library surface
runner.RunTest(dir, Options{GenDir})
  -> Response.RunErr + generated .go under GenDir
core.WriteFormattedGo
  -> bytes on disk without session/time auto-inject from usage maps
```

## Preconditions

- Package `github.com/xhd2015/doctest/libdoc/runner` exposes `RunTest`.
- Package `github.com/xhd2015/doctest/libdoc/core` exposes `WriteFormattedGo`.
- Each leaf uses an isolated fixture tree and gen root under `t.TempDir()`.
- Classic TDD: **RED** expected for A1 / A3 / A5 until implementer removes
  `format.Source` requirement and `stdlibByPkgName` auto-add.
- Out of scope: gen-manifest, tree-stamps, frontier skip.

## Steps

1. Leaf Setup sets `req.Op` and `FixtureKind` / `Source`.
2. Root `Run` builds the fixture (if needed), runs generate+test or WriteFormattedGo.
3. Leaf Assert checks `RunErr` / generated source / build outcome.

## Context

- Shared helpers live in root `DOCTEST.md` (`writeFixture`, `fillGenSources`, …).
- Prefer inspecting generated `.go` under `GenDir` over CLI stdout scraping.
- Do not assert gofmt equality for success leaves (A6).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Defaults; leaves override Op / FixtureKind / Source.
	if req.Op == "" {
		req.Op = "run_fixture"
	}
	return nil
}
```
