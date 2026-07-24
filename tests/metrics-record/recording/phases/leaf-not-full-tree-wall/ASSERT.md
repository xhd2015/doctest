## Expected

- leaf_end with result pass has elapsed_ns > 0 for single-leaf (package-attributed).
- go_test phase elapsed_ns > 0.
- leaf elapsed_ns <= go_test elapsed_ns * 3 (slack for measurement noise; not a multi-leaf clone of full suite).

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if len(resp.Decoded) == 0 {
		t.Fatalf("no events; RunErr=%s", resp.RunErr)
	}
	var goTestNs int64
	var leafNs int64
	for _, m := range resp.Decoded {
		switch m["type"] {
		case "phase":
			if m["phase"] == "go_test" {
				if f, ok := m["elapsed_ns"].(float64); ok {
					goTestNs = int64(f)
				}
			}
		case "leaf_end":
			if m["result"] == "pass" || m["result"] == "PASS" {
				if f, ok := m["elapsed_ns"].(float64); ok {
					leafNs = int64(f)
				}
			}
		}
	}
	if goTestNs <= 0 {
		t.Fatalf("go_test phase elapsed_ns=%d want >0", goTestNs)
	}
	if leafNs <= 0 {
		t.Fatalf("leaf_end elapsed_ns=%d want >0 (package-attributed)", leafNs)
	}
	if leafNs > goTestNs*3 {
		t.Fatalf("leaf elapsed %d far exceeds go_test %d (looks like full-tree clone)", leafNs, goTestNs)
	}
}
```
