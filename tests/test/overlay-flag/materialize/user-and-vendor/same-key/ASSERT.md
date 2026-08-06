## Expected

- No error; exactly one `-overlay=` flag.
- Final Replace at active vendor go.mod path equals **placeholder** path
  (`resp.ActiveBridgeSource`), not `from-user-gomod`.
- Seed-only `project-source` remains `seed-only` (proves user seed applied).

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
	got := resp.OverlayReplace[resp.ActiveVendorGoMod]
	if got != resp.ActiveBridgeSource {
		t.Fatalf("vendor must win same go.mod key: got %q want placeholder %q; map=%v",
			got, resp.ActiveBridgeSource, resp.OverlayReplace)
	}
	if got == "from-user-gomod" {
		t.Fatal("user seed must not win over vendor on the same key")
	}
	if got := resp.OverlayReplace[resp.ProjectSource]; got != "seed-only" {
		t.Fatalf("seed-only key missing (user seed not applied?): got %q map=%v", got, resp.OverlayReplace)
	}
}
```
