## Expected

- Run JSONL exists.
- At least one `type=phase` event.
- Core phases present: `discover`, `generate`, `go_test`.
- Each core phase has `elapsed_ns >= 0`.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
	if len(resp.RunFiles) == 0 {
		t.Fatalf("expected run JSONL; RunErr=%s stderr=%s", resp.RunErr, resp.Stderr)
	}
	phases := map[string]int64{}
	for _, m := range resp.Decoded {
		if m["type"] != "phase" {
			continue
		}
		name, _ := m["phase"].(string)
		if name == "" {
			t.Fatalf("phase event missing phase name: %v", m)
		}
		var ns int64
		if f, ok := m["elapsed_ns"].(float64); ok {
			ns = int64(f)
		}
		if ns < 0 {
			t.Fatalf("phase %s negative elapsed: %v", name, m)
		}
		phases[name] += ns
	}
	if len(phases) == 0 {
		t.Fatalf("no type=phase events; types=%v", eventTypes(resp.Decoded))
	}
	for _, want := range []string{"discover", "generate", "go_test"} {
		if _, ok := phases[want]; !ok {
			var have []string
			for k := range phases {
				have = append(have, k)
			}
			t.Fatalf("missing phase %q; have %v", want, have)
		}
	}
}
```
