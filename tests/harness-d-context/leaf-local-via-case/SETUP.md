# Scenario

**Feature**: leaf-local fixture file is read via `d.DOCTEST_CASE` (not cwd / free vars)

```
# temp fixture
leaf/fixture.txt ("hello-from-case")
Run -> os.ReadFile(filepath.Join(d.DOCTEST_CASE, "fixture.txt"))
Assert -> content matches

# outer
doctest test <fixture> -> exit 0
```

## Preconditions

- Process cwd is undetermined for generated tests; leaf-local I/O must use `d.DOCTEST_CASE`.
- Fixture file lives next to the leaf `ASSERT.md` / `SETUP.md`.

## Steps

1. Write temp tree with `leaf/fixture.txt`.
2. Fixture `Run` reads only through `filepath.Join(d.DOCTEST_CASE, "fixture.txt")`.
3. Outer runs `doctest test` and expects PASS.

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
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/session"
)

type Request struct{}
type Response struct {
	Content string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	_ = req
	p := filepath.Join(d.DOCTEST_CASE, "fixture.txt")
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return &Response{Content: string(b)}, nil
}`)

	writeScenarioSetup(t, filepath.Join(leaf, "SETUP.md"),
		"leaf-local path via d.DOCTEST_CASE",
		"Run -> read filepath.Join(d.DOCTEST_CASE, \"fixture.txt\")",
		`import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	_ = req
	return nil
}`)

	writeAssert(t, filepath.Join(leaf, "ASSERT.md"), `import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	if err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(resp.Content)
	if got != "hello-from-case" {
		t.Fatalf("fixture content = %q, want hello-from-case", resp.Content)
	}
}`)

	writeFile(t, filepath.Join(leaf, "fixture.txt"), "hello-from-case\n")
	setFixtureTest(req, root)
	return nil
}
```
