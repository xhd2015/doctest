# test.config.json go.min / go.max

## Version

0.0.1

# DSN (Domain Specific Notion)

## Participants

- **test.config.json** — project-root config (alongside go.mod) may declare
  `"go":{"min":"…","max":"…"}` or `"go":"1.18"` (string → min only).
- **Shared parser** — `github.com/xhd2015/xgo/support/testconfig` parses the
  shared schema and validates host Go against min/max (major.minor, patch ignored).
- **Doctest loader** — `core.LoadXgoTestConfig` / `ValidateXgoTestConfigGoVersion`
  wraps the shared package; production `applyProjectTestConfig` fails before
  pre_test and `go`/`xgo test` when out of range.

## Behaviors

1. Missing / empty go constraints → validation is a no-op.
2. Host Go &lt; min → error containing `&lt; <min>`.
3. Host Go &gt; max (minor) → error containing `&gt; <max>`.
4. Wide range (min 1.0, max 99.0) → OK on any reasonable host.
5. Object and string forms of `"go"` are parsed into min/max fields.

## Decision Tree

```text
go-version-config/
├── parse/
│   ├── object-min-max/     go:{min,max} preserved on load
│   └── string-min-only/    "go":"1.18" → Min only
├── validate/
│   ├── no-constraint/      empty/nil → OK
│   ├── in-range/           min 1.0 max 99.0 → OK
│   ├── below-min/          min 99.0 → error
│   └── above-max/          max 1.0 → error
```

## Test Index

| Leaf | Contract |
|---|---|
| `parse/object-min-max` | Load preserves min and max |
| `parse/string-min-only` | String form becomes Min |
| `validate/no-constraint` | No go block → no error |
| `validate/in-range` | Wide range accepts host go |
| `validate/below-min` | Impossible min fails with `<` |
| `validate/above-max` | Impossible max fails with `>` |

## How to Run

```sh
doctest vet ./tests/test/go-version-config
doctest test ./tests/test/go-version-config --label-all
```

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/session"
)

// Request selects which load/validate path Run exercises.
type Request struct {
	// ConfigJSON is written to ModRoot/test.config.json when non-empty.
	ConfigJSON string
	// SkipWrite: do not write a config file (nil load).
	SkipWrite bool
	// OnlyLoad: load only; do not validate (parse leaves).
	OnlyLoad bool
}

type Response struct {
	// Loaded is true when Load returned a non-nil config.
	Loaded bool
	Min    string
	Max    string
	// ErrMsg is the validation (or load) error text.
	ErrMsg string
}

func writeConfig(t *testing.T, dir, body string) {
	t.Helper()
	path := filepath.Join(dir, core.DefaultXgoTestConfigName)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}
	dir := t.TempDir()
	var cfg *core.XgoTestConfig
	var err error
	if !req.SkipWrite {
		if req.ConfigJSON == "" {
			req.ConfigJSON = "{}"
		}
		writeConfig(t, dir, req.ConfigJSON)
		cfg, err = core.LoadXgoTestConfig(filepath.Join(dir, core.DefaultXgoTestConfigName))
		if err != nil {
			resp.ErrMsg = err.Error()
			return resp, nil
		}
	} else {
		cfg, err = core.LoadXgoTestConfig(filepath.Join(dir, "missing.json"))
		if err != nil {
			resp.ErrMsg = err.Error()
			return resp, nil
		}
	}
	if cfg != nil {
		resp.Loaded = true
		if cfg.Go != nil {
			resp.Min = cfg.Go.Min
			resp.Max = cfg.Go.Max
		}
	}
	if req.OnlyLoad {
		return resp, nil
	}
	if err := core.ValidateXgoTestConfigGoVersion(cfg, "go"); err != nil {
		resp.ErrMsg = err.Error()
	}
	return resp, nil
}
```
