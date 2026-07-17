## Expected

- Run file exists with decoded events.
- At least one `leaf_start` and one `leaf_end`.
- Some `leaf_end` has `result` of `pass` (or equivalent success string).

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if len(resp.RunFiles) == 0 {
		t.Fatalf("expected run JSONL; RunErr=%s stderr=%s", resp.RunErr, resp.Stderr)
	}
	if !hasEventType(resp.Decoded, "leaf_start") {
		t.Fatalf("missing leaf_start; types=%v", eventTypes(resp.Decoded))
	}
	if !hasEventType(resp.Decoded, "leaf_end") {
		t.Fatalf("missing leaf_end; types=%v", eventTypes(resp.Decoded))
	}
	foundPass := false
	for _, m := range resp.Decoded {
		if m["type"] != "leaf_end" {
			continue
		}
		r, _ := m["result"].(string)
		if r == "pass" || r == "PASS" || r == "ok" {
			foundPass = true
			break
		}
	}
	if !foundPass {
		t.Fatalf("no leaf_end with pass result; events=%v", resp.Decoded)
	}
}
```
