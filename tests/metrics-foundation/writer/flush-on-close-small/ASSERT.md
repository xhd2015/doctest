## Expected

- `SizeBeforeClose` is 0 (buffer holds the only event).
- `SizeAfterClose` is > 0 and file decodes to one `run_start` line.
- Raw file ends with `\n`.

```go
import (
	"bytes"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v (resp.Err=%s)", err, resp.Err)
	}
	if resp.SizeBeforeClose != 0 {
		t.Fatalf("size before close = %d, want 0 (still buffered)", resp.SizeBeforeClose)
	}
	if resp.SizeAfterClose <= 0 {
		t.Fatalf("size after close = %d, want > 0", resp.SizeAfterClose)
	}
	if !bytes.HasSuffix(resp.FileData, []byte("\n")) {
		t.Fatal("expected trailing newline after flush")
	}
	if len(resp.Decoded) != 1 || resp.Decoded[0]["type"] != "run_start" {
		t.Fatalf("decoded = %#v", resp.Decoded)
	}
}
```
