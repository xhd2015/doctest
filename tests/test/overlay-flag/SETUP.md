# Scenario

**Feature**: `doctest test` accepts user `-overlay`/`--overlay` and merges it as a seed with internal overlay layers into one go-test flag

```
# parse
doctest test -overlay FILE | --overlay FILE <dir>
  -> Options.Overlay abs path (relative resolved against process cwd)

# materialize (single -overlay=)
user Replace seed
  -> pre_test hooks overwrite same keys
  -> vendor-gomod bridge merge overwrites same keys
  -> GoFlags: exactly one -overlay=<driver file> when Replace non-empty

# scope
doctest build -overlay ... -> still rejected (test-only)
```

## Preconditions

- Nested self-contained root: does **not** inherit parent `tests/DOCTEST.md` Run.
- Classic TDD: `Options.Overlay`, `-overlay/--overlay` parse, and
  `ApplyPreTestHooksWithUserOverlay` / `MaterializeUserVendorOverlay` may be
  missing — leaves expected **RED** until implementer lands them.
- Prefer L2 in-process over product-binary e2e. No process `Setenv`/`Chdir`.
- Use `d.DOCTEST_ROOT` / `d.DOCTEST_CASE` / `d.DOCTEST_SESSION_ID` only (never
  `os.Getenv` for session id).
- Help documentation is asserted in `tests/help/test-options` (sibling tree).

## Steps

1. Root Setup sets default Mode empty (leaves set Mode).
2. `parse/*` leaves set `Mode=parse` and `ParseArgs`.
3. `materialize/*` leaves set `Mode=materialize` and layer fixtures.
4. `build-reject` sets `Mode=build_reject` and build argv.
5. `Run` dispatches; Assert checks parse Overlay path, Replace map, GoFlags, or build error.

## Context

- Fixture aliases for hook overlays match `tests/test/pre-test-hooks`:
  `project-source`, `active-vendor`, `inactive-vendor`.
- Vendor bridge uses xgo-style placeholder go.mod under
  `generated/vendor-gomod-overlay/<mod>/go.mod`.
- Helpers: `errText`, `wantAbsOverlay` (optional leaf-local).

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Leaves set Mode; default empty so a misconfigured leaf fails loudly in Run.
	return nil
}

func errText(resp *Response) string {
	if resp == nil {
		return ""
	}
	return strings.TrimSpace(resp.ErrMsg + "\n" + resp.ParseErr + "\n" + resp.Stderr)
}

// absLike reports whether p is absolute (and non-empty).
func absLike(p string) bool {
	return p != "" && filepath.IsAbs(p)
}
```
