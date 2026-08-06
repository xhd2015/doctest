# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Steps
- Input is `"/foo/bar/..."` (absolute path with `/...` suffix).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Input = "/foo/bar/..."
	return nil
}
```