# Scenario

**Feature**: package helper that needs paths takes `d *session.Doctest`

```
# fixture DOCTEST.md helper
func joinCase(d *session.Doctest, name string) string
  -> filepath.Join(d.DOCTEST_CASE, name)

Setup(t, d, req) -> req.Joined = joinCase(d, "marker.txt")
Assert -> Joined == filepath.Join(d.DOCTEST_CASE, "marker.txt")

# outer
doctest test <fixture> -> exit 0
```

## Preconditions

- Helpers must not close over free `DOCTEST_ROOT` / `DOCTEST_CASE` package vars.
- Helper signature includes `d *session.Doctest` when paths are required.

## Steps

1. Write fixture with root helper `joinCase`.
2. Setup calls `joinCase(d, "marker.txt")`.
3. Assert compares against `filepath.Join(d.DOCTEST_CASE, "marker.txt")`.
4. Outer `doctest test` expects PASS.

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
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

type Request struct {
	Joined string
}
type Response struct {
	Joined string
}

// joinCase joins a relative name under the leaf case dir via d.
func joinCase(d *session.Doctest, name string) string {
	return filepath.Join(d.DOCTEST_CASE, name)
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = d
	return &Response{Joined: req.Joined}, nil
}`)

	writeScenarioSetup(t, filepath.Join(leaf, "SETUP.md"),
		"helper joinCase takes d",
		"Setup -> joinCase(d, \"marker.txt\") -> req.Joined",
		`import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Joined = joinCase(d, "marker.txt")
	return nil
}`)

	writeAssert(t, filepath.Join(leaf, "ASSERT.md"), `import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = req
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(d.DOCTEST_CASE, "marker.txt")
	if resp.Joined != want {
		t.Fatalf("Joined = %q, want %q", resp.Joined, want)
	}
}`)

	setFixtureTest(req, root)
	return nil
}
```
