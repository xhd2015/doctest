# Session Doctest — public context type (fields only)

## Version
0.0.2

# DSN (Domain Specific Notion)

**Participants**

- **Caller** — integration leaf, harness, or product code that needs an
  immutable inject contract for the current doctest run.
- **Session Doctest** (`github.com/xhd2015/doctest/session`) — public type
  `session.Doctest` with three string fields only (no methods):
  - `DOCTEST_ROOT` — absolute doctest tree root
  - `DOCTEST_CASE` — absolute leaf case directory
  - `DOCTEST_SESSION_ID` — session id for this doctest run
- **Inject contract** — field names are ALL_CAPS on purpose, matching the
  former free variables of the same names. These are **struct fields**, not
  process environment variables.

**Behaviors**

- Caller constructs `session.Doctest` (or takes a pointer) and reads fields
  by name; each field returns exactly the string that was set.
- The zero value of `session.Doctest` has empty strings for all three fields.
- Fields are independent: assigning one field never changes the other two.
- No methods, no env reads, no side effects on construction or field access.

## Decision Tree

```
doctest-context/
├── construct-and-read/     S1: composite literal; all three fields readable
├── zero-value/             S2: zero value → empty strings for all fields
└── field-independence/     S3: setting one field leaves the others empty
```

## Test Index

| Leaf | Scenario |
|------|----------|
| `construct-and-read` | S1 — construct with all three fields set; each reads back the set value |
| `zero-value` | S2 — zero `session.Doctest` has `""` for ROOT, CASE, and SESSION_ID |
| `field-independence` | S3 — set-only-root / set-only-case / set-only-session leave others empty |

## How to Run

```sh
doctest vet ./session/tests/doctest-context/
doctest test ./session/tests/doctest-context/   # expect RED until session.Doctest lands
go test ./session/
```

```go
import (
	"fmt"
	"testing"

	"github.com/xhd2015/doctest/session"
)

// Request configures how Run builds a session.Doctest observation.
// Leaves set Mode and optional Want* values via Setup.
type Request struct {
	// Mode selects construction path:
	//   "construct"    — composite literal with all Want* fields
	//   "zero"         — zero value only
	//   "independence" — three instances, each with only one field set
	Mode string
	// WantRoot is the value for DOCTEST_ROOT when constructing.
	WantRoot string
	// WantCase is the value for DOCTEST_CASE when constructing.
	WantCase string
	// WantSessionID is the value for DOCTEST_SESSION_ID when constructing.
	WantSessionID string
}

// FieldView is a snapshot of the three session.Doctest fields.
type FieldView struct {
	Root      string
	Case      string
	SessionID string
}

// Response carries observed field values for Assert.
type Response struct {
	// View is used by construct and zero modes (single instance).
	View FieldView
	// OnlyRoot / OnlyCase / OnlySession are used by independence mode:
	// each snapshot is from a separate instance with only that field set.
	OnlyRoot    FieldView
	OnlyCase    FieldView
	OnlySession FieldView
}

func viewOf(d session.Doctest) FieldView {
	return FieldView{
		Root:      d.DOCTEST_ROOT,
		Case:      d.DOCTEST_CASE,
		SessionID: d.DOCTEST_SESSION_ID,
	}
}

func Run(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	resp := &Response{}
	switch req.Mode {
	case "construct":
		d := session.Doctest{
			DOCTEST_ROOT:       req.WantRoot,
			DOCTEST_CASE:       req.WantCase,
			DOCTEST_SESSION_ID: req.WantSessionID,
		}
		resp.View = viewOf(d)
	case "zero":
		var d session.Doctest
		resp.View = viewOf(d)
	case "independence":
		// Three separate values: only one field set on each.
		resp.OnlyRoot = viewOf(session.Doctest{
			DOCTEST_ROOT: req.WantRoot,
		})
		resp.OnlyCase = viewOf(session.Doctest{
			DOCTEST_CASE: req.WantCase,
		})
		resp.OnlySession = viewOf(session.Doctest{
			DOCTEST_SESSION_ID: req.WantSessionID,
		})
	default:
		return nil, fmt.Errorf("unknown Request.Mode %q", req.Mode)
	}
	// Never fail Run on product observations — Assert owns expectations.
	return resp, nil
}
```
