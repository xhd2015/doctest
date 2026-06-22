# Dot Progress Tests

## Version
0.0.2


Tests that `build.Test` prints dot progress **incrementally** — one dot per
test package as it completes — rather than batching them all after `go test`
finishes.

## How to Run

```sh
doctest test ./libdoc/build/tests/dot-progress/...
```

```go
import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
	"github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

type Request struct{}
type Response struct {
	Incremental	bool
	DotCount	int
	Output		string
}
func Run(t *testing.T, req *Request) (*Response, error) {
	subRoot := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, subRoot, []testtree.LeafSpec{
		{Name: "a_fast", Steps: "No setup needed.", Expected: "Always passes."},
		{Name: "z_slow", Steps: "Sleep 5 seconds to simulate a long-running test.", Expected: "Always passes.",
			SetupGo: "import (\"testing\"; \"time\")\n\nfunc Setup(t *testing.T, req *Request) error { time.Sleep(5 * time.Second); return nil }"},
	})

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	// Read dots in background goroutine.
	type dotInfo struct {
		firstDot	time.Duration	// -1 means never
		output		string
	}
	ch := make(chan dotInfo, 1)
	start := time.Now()
	go func() {
		var buf bytes.Buffer
		firstDot := time.Duration(-1)
		tmp := make([]byte, 1)
		for {
			n, readErr := r.Read(tmp)
			if n > 0 {
				buf.WriteByte(tmp[0])
				if tmp[0] == '.' && firstDot < 0 {
					firstDot = time.Since(start)
				}
			}
			if readErr != nil {
				break
			}
		}
		ch <- dotInfo{firstDot, buf.String()}
	}()

	testErr := build.Test(subRoot, core.Options{RemoveTemp: true})
	w.Close()
	info := <-ch
	os.Stdout = oldStdout

	incremental := info.firstDot >= 0 && info.firstDot < 4*time.Second

	inlineIdx := strings.Index(info.output, "  (")
	dotCount := 0
	if inlineIdx > 0 {
		dotCount = strings.Count(info.output[:inlineIdx], ".")
	}

	if testErr != nil {
		return &Response{
			Incremental:	incremental,
			DotCount:	dotCount,
			Output:		info.output,
		}, testErr
	}

	return &Response{
		Incremental:	incremental,
		DotCount:	dotCount,
		Output:		info.output,
	}, nil
}
```
