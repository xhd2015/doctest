# Scenario

**Feature**: N1 — MatchContains contiguous subregion

```
# match options alter comparison policy
Matcher <- actual (+ Contains option or CRLF normalization)
```

## Steps
1. Set template/actual fields for N1 — MatchContains contiguous subregion.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = "..\n  (2 Run, 2 Pass, 1 Cached, 0 Fail)"
	req.Actual = "[info] start\n[info] compiling\n..\n  (2 Run, 2 Pass, 1 Cached, 0 Fail)\n[info] done"
	return nil
}
```
