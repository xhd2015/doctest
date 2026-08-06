## Expected

- No error; exactly one `-overlay=` flag.
- Final Replace at project-source path equals **hook** value `from-hook`, not `from-user`.
- Seed-only key `active-vendor` remains `seed-only` (proves user seed was applied).

## Errors

- Classic TDD **RED** until user seed is applied before hooks and hooks overwrite.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run err: %v", err)
	}
	if resp.ExitCode != 0 || resp.ErrMsg != "" {
		t.Fatalf("unexpected error: %s", errText(resp))
	}
	if resp.OverlayFlagN != 1 || len(resp.GoFlags) != 1 || !strings.HasPrefix(resp.GoFlags[0], "-overlay=") {
		t.Fatalf("want single -overlay=, got n=%d flags=%v", resp.OverlayFlagN, resp.GoFlags)
	}
	got := resp.OverlayReplace[resp.ProjectSource]
	if got != "from-hook" {
		t.Fatalf("same key must be hook-wins: got %q want from-hook; map=%v project=%q",
			got, resp.OverlayReplace, resp.ProjectSource)
	}
	if got == "from-user" {
		t.Fatal("user seed must not win over later hook on the same key")
	}
	// Seed-only key must survive: proves seed layer ran (not hook-only false green).
	if got := resp.OverlayReplace[resp.ActiveVendorSource]; got != "seed-only" {
		t.Fatalf("seed-only key missing (seed not applied?): got %q map=%v active=%q",
			got, resp.OverlayReplace, resp.ActiveVendorSource)
	}
}
```
