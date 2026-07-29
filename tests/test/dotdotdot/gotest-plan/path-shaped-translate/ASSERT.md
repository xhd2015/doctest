---
label: heavy
---

## Expected

- CLI still succeeds (workspace suite for single-mod).
- `gotestmap.TranslatePath("./tree/mid/...", layout)` yields mid pattern + nested mod, not `./tree/...`.

```go
import (
	"reflect"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/gotestmap"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	out := gotestPlanOut(resp)
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d\n%s", resp.ExitCode, out)
	}
	// Current single-mod CLI path: workspace suite (not path-shaped multi-cmd).
	assertContainsGoTestLine(t, out, "__workspace/suite")

	layout := gotestmap.Layout{ModuleRoots: []string{".", "tree/mid/nestedmod"}}
	got, err := gotestmap.TranslatePath("./tree/mid/...", layout)
	if err != nil {
		t.Fatal(err)
	}
	want := []gotestmap.Cmd{
		{Dir: ".", Pattern: "./tree/mid/..."},
		{Dir: "tree/mid/nestedmod", Pattern: "./..."},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("TranslatePath:\n got %v\nwant %v", got, want)
	}
	// Must not widen to whole tree.
	for _, c := range got {
		if c.Pattern == "./tree/..." {
			t.Fatalf("must not translate mid to ./tree/...: %v", got)
		}
	}
	if !strings.Contains(out, "PASS") {
		t.Fatalf("missing PASS\n%s", out)
	}
}
```
