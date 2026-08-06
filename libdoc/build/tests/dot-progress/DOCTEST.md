# Dot Progress Tests

## Version
0.0.2

## DSN (Domain Specific Notion)

### Participants
- **system under test** — behavior covered by this tree.

### Behaviors
- **run** — executes the scenarios in this suite.


Tests that `build.Test` prints dot progress **incrementally** — one dot per
test package as it completes — rather than batching them all after `go test`
finishes.

Timing-sensitive incremental behavior is covered by `TestDotProgressIncremental`
in `libdoc/build/build_engine_test.go` (skipped with `-short`). This doctest
verifies dot structure on a fast, cache-friendly fixture.

## How to Run

```sh
doctest test ./libdoc/build/tests/dot-progress/...
go test ./libdoc/build -run TestDotProgressIncremental
```

```go
import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

type Request struct{}
type Response struct {
	DotCount	int
	Output		string
}
func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	fixtureBase := filepath.Join(os.TempDir(), "doctest-libdoc-build-dot-progress-fast")
	subRoot := filepath.Join(fixtureBase, "tree")
	genDir := filepath.Join(fixtureBase, "gendir")
	initMarker := filepath.Join(fixtureBase, ".initialized")
	if _, err := os.Stat(initMarker); err != nil {
		testtree.WriteMinimalRunnableTree(t, subRoot, []testtree.LeafSpec{
			{Name: "a_fast", Steps: "No setup needed.", Expected: "Always passes."},
			{Name: "b_fast", Steps: "No setup needed.", Expected: "Always passes."},
		})
		if err := os.WriteFile(initMarker, []byte("ok\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	var buf bytes.Buffer
	testErr := build.Test(subRoot, core.Options{GenDir: genDir, RemoveTemp: false, Stdout: &buf})
	output := buf.String()
	inlineIdx := strings.Index(output, "  (")
	dotCount := 0
	if inlineIdx > 0 {
		dotCount = strings.Count(output[:inlineIdx], ".")
	}

	if testErr != nil {
		return &Response{DotCount: dotCount, Output: output}, testErr
	}
	return &Response{DotCount: dotCount, Output: output}, nil
}
```