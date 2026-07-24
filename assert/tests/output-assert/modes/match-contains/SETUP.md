# Scenario

**Feature**: MatchContains contiguous subregion

```
# match options alter comparison policy
Matcher <- actual (+ Contains option or CRLF normalization)
```

## Steps
1. Set `req.Options = []string{"contains"}`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Options = []string{"contains"}
	return nil
}
```
