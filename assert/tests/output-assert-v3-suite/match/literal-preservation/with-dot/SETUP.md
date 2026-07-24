# Scenario

**Feature**: V3S-M10 — escaped dot is literal under v3 raw RE

```
# version 1\.0 — author escapes metachar so '.' is not "any char"
Matcher <- exact literal version 1.0
```

## Steps
1. Set escaped pattern line `version 1\.0` and identical actual `version 1.0`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Template = v3Template("", "version 1\\.0")
	req.Actual = "version 1.0"
	return nil
}
```
