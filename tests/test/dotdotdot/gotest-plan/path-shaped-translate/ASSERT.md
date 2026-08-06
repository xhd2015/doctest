## Expected

- CLI still succeeds with **single** workspace suite plan (not multi-cmd path-shaped execution).
- Pure `TranslatePath` / `Plan(ModePathShaped)` mid+nestedmod: two fixture cmds, not `./tree/...`.
- Phase 2 only: wiring those multi cmds into production go-test finish.

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
	// Phase 1 CLI: single workspace suite — not path-shaped multi-cmd finish.
	assertExactlyOneGoTestPlanFamily(t, out, "__workspace/suite")

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
	// Same contract via Plan(ModePathShaped) — multi-cmd is pure plan only (Phase 2 exec).
	viaPlan, err := gotestmap.Plan(gotestmap.PlanInput{
		Mode:    gotestmap.ModePathShaped,
		UserArg: "./tree/mid/...",
		Layout:  layout,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(viaPlan, want) {
		t.Fatalf("Plan(ModePathShaped):\n got %v\nwant %v", viaPlan, want)
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
