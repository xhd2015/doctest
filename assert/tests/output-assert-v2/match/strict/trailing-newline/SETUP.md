# Scenario

**Feature**: V2-M9 — strict trailing newline mismatch

```
# template without trailing \n vs actual with trailing \n
Matcher rejects trailing newline drift (same as v1 M2)
```

## Steps
1. Set template without trailing newline; actual has trailing newline.

```go
func Setup(t *testing.T, req *Request) error {
	req.Template = v2Template("", "OK")
	req.Actual = "OK\n"
	return nil
}
```