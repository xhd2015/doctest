## Expected

- `SizeBeforeClose` > 0 (threshold flush occurred mid-run).
- `SizeAfterClose` >= `SizeBeforeClose`.
- File contains at least the run_start line (and pad line after full flush).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/libdoc/metrics"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v (resp.Err=%s)", err, resp.Err)
	}
	if metrics.FlushThreshold != 128*1024 {
		t.Fatalf("FlushThreshold = %d, want %d", metrics.FlushThreshold, 128*1024)
	}
	if resp.SizeBeforeClose <= 0 {
		t.Fatalf("size before close = %d, want > 0 after ≥128KiB write", resp.SizeBeforeClose)
	}
	if resp.SizeAfterClose < resp.SizeBeforeClose {
		t.Fatalf("size after close %d < before %d", resp.SizeAfterClose, resp.SizeBeforeClose)
	}
	if len(resp.Decoded) < 1 {
		t.Fatal("expected at least one decoded JSON line")
	}
	foundStart := false
	for _, m := range resp.Decoded {
		if m["type"] == "run_start" {
			foundStart = true
			break
		}
	}
	if !foundStart {
		t.Fatalf("run_start missing in %#v", resp.Decoded)
	}
}
```
