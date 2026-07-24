## Expected
- `resp.RootOkResult` is `false` (no DOCTEST.md or SETUP.md found).

```go
func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	if resp.RootOkResult {
		t.Fatalf("expected ok == false, got root %q", resp.RootResult)
	}
}
```
