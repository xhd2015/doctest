---
explanation: nested prepare/go test via runner.RunTest (multi-second)
---

## Expected

- At least one new run JSONL under MetricsRoot.
- Decoded events include `run_start` and `run_end`.
- Those events include `schema_version` == 1.
- First event type is `run_start`; last is `run_end`.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if len(resp.RunFiles) == 0 {
		t.Fatalf("expected run JSONL under MetricsRoot=%s; RunErr=%s stderr=%s",
			resp.MetricsRoot, resp.RunErr, resp.Stderr)
	}
	if !hasEventType(resp.Decoded, "run_start") {
		t.Fatalf("missing run_start in %v types=%v", resp.RunFiles[0], eventTypes(resp.Decoded))
	}
	if !hasEventType(resp.Decoded, "run_end") {
		t.Fatalf("missing run_end in %v types=%v", resp.RunFiles[0], eventTypes(resp.Decoded))
	}
	types := eventTypes(resp.Decoded)
	if types[0] != "run_start" {
		t.Fatalf("first event type=%q, want run_start; types=%v", types[0], types)
	}
	if types[len(types)-1] != "run_end" {
		t.Fatalf("last event type=%q, want run_end; types=%v", types[len(types)-1], types)
	}
	for i, m := range resp.Decoded {
		ty, _ := m["type"].(string)
		if ty != "run_start" && ty != "run_end" {
			continue
		}
		sv := m["schema_version"]
		switch v := sv.(type) {
		case float64:
			if int(v) != 1 {
				t.Fatalf("event[%d] (%s) schema_version=%v, want 1", i, ty, sv)
			}
		case int:
			if v != 1 {
				t.Fatalf("event[%d] (%s) schema_version=%v, want 1", i, ty, sv)
			}
		default:
			t.Fatalf("event[%d] (%s) schema_version missing/type %T=%v", i, ty, sv, sv)
		}
	}
}
```
