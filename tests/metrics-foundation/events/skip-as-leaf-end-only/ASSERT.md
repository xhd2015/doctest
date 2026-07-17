## Expected

- Three lines: run_start, leaf_end, run_end (no leaf_start).
- leaf_end has `result=skip` and no `ts_start` key (or null/empty — prefer absent).
- Labels include `heavy`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v (resp.Err=%s)", err, resp.Err)
	}
	if len(resp.Decoded) != 3 {
		t.Fatalf("decoded count = %d, want 3\n%s", len(resp.Decoded), resp.FileData)
	}
	if resp.Decoded[0]["type"] != "run_start" {
		t.Fatalf("first type = %v", resp.Decoded[0]["type"])
	}
	if resp.Decoded[1]["type"] != "leaf_end" {
		t.Fatalf("second type = %v", resp.Decoded[1]["type"])
	}
	if resp.Decoded[2]["type"] != "run_end" {
		t.Fatalf("third type = %v", resp.Decoded[2]["type"])
	}
	le := resp.Decoded[1]
	if le["result"] != "skip" {
		t.Fatalf("result = %v, want skip", le["result"])
	}
	if _, hasStart := le["ts_start"]; hasStart {
		// allow empty string but prefer omission; non-empty is a hard fail for this leaf
		if s, ok := le["ts_start"].(string); ok && s != "" {
			t.Fatalf("skip leaf_end must omit ts_start or leave empty, got %q", s)
		}
	}
	// ensure no leaf_start anywhere
	for i, m := range resp.Decoded {
		if m["type"] == "leaf_start" {
			t.Fatalf("unexpected leaf_start at index %d", i)
		}
	}
}
```
