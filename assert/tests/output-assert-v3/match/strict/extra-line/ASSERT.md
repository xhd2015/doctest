## Expected
- Template parses as v3 (content line, not YAML header as literal).
- Match fails for extra actual line.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	// Parse must succeed on the v3 path — summary must not treat --- as body literals only.
	requireParseOK(t, resp)
	// v3 content lines are RegexLine (or at least not LiteralLine×4 from YAML dump).
	if resp.ParsedSummary == "LiteralLine×4" || resp.ParsedSummary == "LiteralLine×5" {
		t.Fatalf("expected v3 parse, got v1-like summary %q", resp.ParsedSummary)
	}
	requireMatchError(t, resp)
}
```
