## Expected

- Exactly three decoded events: run_start, leaf_start, leaf_end.
- No event has `type=run_end`.
- Each line is valid JSON with `schema_version` 1.

```go
import "testing"

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v (resp.Err=%s)", err, resp.Err)
	}
	if len(resp.Decoded) != 3 {
		t.Fatalf("decoded count = %d, want 3\n%s", len(resp.Decoded), resp.FileData)
	}
	want := []string{"run_start", "leaf_start", "leaf_end"}
	for i, w := range want {
		if resp.Decoded[i]["type"] != w {
			t.Fatalf("event[%d].type = %v, want %s", i, resp.Decoded[i]["type"], w)
		}
		if _, bad := resp.Decoded[i]["_parse_error"]; bad {
			t.Fatalf("event[%d] failed to parse: %v", i, resp.Decoded[i])
		}
	}
	for i, m := range resp.Decoded {
		if m["type"] == "run_end" {
			t.Fatalf("unexpected run_end at %d", i)
		}
	}
}
```
