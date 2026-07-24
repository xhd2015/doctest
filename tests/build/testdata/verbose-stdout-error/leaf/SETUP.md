# Scenario

**Feature**: tests for leaf

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Steps

- Set the request name.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Name = "test"
    return nil
}
```
