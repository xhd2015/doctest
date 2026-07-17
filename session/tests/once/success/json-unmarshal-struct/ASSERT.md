## Expected

- Once succeeds with valid JSON.
- `json.Unmarshal` into `struct { Path string \`json:"path"\` }` succeeds.
- `Path` equals `/tmp/doctest-session-once-bin`.

```go
import (
	"encoding/json"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run harness error: %v", err)
	}
	if resp.Err != nil {
		t.Fatalf("Once error: %v", resp.Err)
	}
	if !json.Valid(resp.Value) {
		t.Fatalf("not valid JSON: %s", resp.Value)
	}
	var got struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(resp.Value, &got); err != nil {
		t.Fatalf("unmarshal: %v (raw=%s)", err, resp.Value)
	}
	if got.Path != "/tmp/doctest-session-once-bin" {
		t.Fatalf("path=%q want /tmp/doctest-session-once-bin", got.Path)
	}
}
```
