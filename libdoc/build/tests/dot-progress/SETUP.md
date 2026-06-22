## Preconditions
- The `build` package is importable (`github.com/xhd2015/doctest/libdoc/build`).
- This test creates a temporary sub-tree with 2 leaves (fast + slow) and
  runs `build.Test` on it, capturing stdout to verify dot progress timing.
- Backtick characters in embedded Go strings use `\x60` to avoid
  conflicting with the outer markdown code fence.

## Steps
1. Create sub-tree under a temp dir with `a_fast` (no sleep) and `z_slow`
   (`time.Sleep(5s)` in Setup).
2. Redirect `os.Stdout` to a pipe.
3. In a background goroutine, read the pipe byte-by-byte, recording when the
   first `"."` appears.
4. Call `build.Test(subRoot, core.Options{RemoveTemp: true})`.
5. After `build.Test` returns, close the pipe and collect the result.
6. Return `Incremental` (true if first dot appeared within 4s) and
   `DotCount` (number of dots before the summary line) in the Response.

## Context
- The fast leaf finishes quickly; the slow leaf takes ~5s. If dots are
  incremental, the first dot appears within ~1s (when a_fast completes).
  If batched, all dots appear after ~5s (when z_slow completes).

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
)

type Request struct{}

type Response struct {
	Incremental bool
	DotCount    int
	Output      string
}

func Run(t *testing.T, req *Request) (*Response, error) {
	subRoot := t.TempDir()

	// Build the sub-tree by writing markdown files directly.
	// Use \x60 for backtick to avoid breaking the outer code fence.
	writeFile(t, subRoot, "SETUP.md", ""+
		"## Preconditions\n"+
		"- Minimal test tree for dot progress.\n\n"+
		"## Steps\n"+
		"1. Run returns immediately.\n\n"+
		"\x60\x60\x60go\n"+
		"import \"testing\"\n"+
		"type Request struct{}\n"+
		"type Response struct{}\n"+
		"func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"+
		"\x60\x60\x60\n")

	// Fast leaf — no delay.
	writeFile(t, subRoot, "a_fast/SETUP.md", ""+
		"## Steps\n"+
		"1. No setup needed.\n\n"+
		"\x60\x60\x60go\n"+
		"import \"testing\"\n"+
		"func Setup(t *testing.T, req *Request) error { _ = req; return nil }\n"+
		"\x60\x60\x60\n")
	writeFile(t, subRoot, "a_fast/ASSERT.md", ""+
		"## Expected\n"+
		"- Always passes.\n\n"+
		"\x60\x60\x60go\n"+
		"func Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"+
		"\x60\x60\x60\n")

	// Slow leaf — sleeps 5s in Setup.
	writeFile(t, subRoot, "z_slow/SETUP.md", ""+
		"## Steps\n"+
		"1. Sleep 5 seconds to simulate a long-running test.\n\n"+
		"\x60\x60\x60go\n"+
		"import (\"testing\"; \"time\")\n"+
		"func Setup(t *testing.T, req *Request) error { time.Sleep(5 * time.Second); return nil }\n"+
		"\x60\x60\x60\n")
	writeFile(t, subRoot, "z_slow/ASSERT.md", ""+
		"## Expected\n"+
		"- Always passes.\n\n"+
		"\x60\x60\x60go\n"+
		"func Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"+
		"\x60\x60\x60\n")

	// Redirect stdout.
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	// Read dots in background goroutine.
	type dotInfo struct {
		firstDot time.Duration // -1 means never
		output   string
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

	// Run build.Test on the sub-tree.
	testErr := build.Test(subRoot, core.Options{RemoveTemp: true})
	w.Close()
	info := <-ch
	os.Stdout = oldStdout

	// Build Response.
	incremental := info.firstDot >= 0 && info.firstDot < 4*time.Second

	inlineIdx := strings.Index(info.output, "  (")
	dotCount := 0
	if inlineIdx > 0 {
		dotCount = strings.Count(info.output[:inlineIdx], ".")
	}

	if testErr != nil {
		return &Response{
			Incremental: incremental,
			DotCount:    dotCount,
			Output:      info.output,
		}, testErr
	}

	return &Response{
		Incremental: incremental,
		DotCount:    dotCount,
		Output:      info.output,
	}, nil
}

func writeFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
```
