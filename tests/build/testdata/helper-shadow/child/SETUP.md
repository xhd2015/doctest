# Scenario

**Feature**: tests for child

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Setup

- Redefines `myHelper` with the same name as the root.
- The generator emits both as closures using `:=`,
  causing "no new variables on left side of :=".

```go
func myHelper(s string) string {
    return "child: " + s
}

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Name = myHelper("name")
    return nil
}
```
