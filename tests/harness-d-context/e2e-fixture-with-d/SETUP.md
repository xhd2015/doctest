# Scenario

**Feature**: e2e fixture Setup/Run/Assert take `d` and use `d.DOCTEST_ROOT`

```
# temp fixture tree
Setup(t, d, req) -> req.RootSeen = d.DOCTEST_ROOT
Run(t, d, req)   -> Response{OK: RootSeen != ""}
Assert(t, d, ...) -> require OK; d.DOCTEST_ROOT and d.DOCTEST_CASE non-empty

# outer leaf
write fixture -> doctest test <fixture> -> exit 0
```

## Preconditions

- Fixture author signatures include `d *session.Doctest` (not free path vars).
- P2 inject constructs and passes `d` so the fixture compiles without package free vars.

## Steps

1. Write a one-leaf temp tree under `t.TempDir()`.
2. Root Setup already set `req.Bin` via `d.DOCTEST_ROOT`.
3. Run `doctest test <fixture> -v`.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	root := t.TempDir()
	leaf := filepath.Join(root, "leaf")

	writeDOCTEST(t, filepath.Join(root, "DOCTEST.md"), `import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

type Request struct {
	RootSeen string
}
type Response struct {
	OK bool
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	return &Response{OK: req.RootSeen != ""}, nil
}`)

	writeScenarioSetup(t, filepath.Join(leaf, "SETUP.md"),
		"setup records d.DOCTEST_ROOT",
		"Setup -> d.DOCTEST_ROOT -> req.RootSeen",
		`import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.RootSeen = d.DOCTEST_ROOT
	return nil
}`)

	writeAssert(t, filepath.Join(leaf, "ASSERT.md"), `import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if !resp.OK {
		t.Fatalf("expected RootSeen non-empty, got %q", req.RootSeen)
	}
	if d.DOCTEST_ROOT == "" {
		t.Fatal("d.DOCTEST_ROOT is empty")
	}
	if d.DOCTEST_CASE == "" {
		t.Fatal("d.DOCTEST_CASE is empty")
	}
	if d.DOCTEST_SESSION_ID == "" {
		t.Fatal("d.DOCTEST_SESSION_ID is empty")
	}
}`)

	setFixtureTest(req, root)
	return nil
}
```
