# Scenario

**Feature**: tests for helper shadow

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Setup

- Defines a helper function `myHelper` at the root level.
- A child SETUP.md also defines `myHelper`, causing a name collision
  when both are emitted as closures with `:=`.

```go
func myHelper(s string) string {
	return "root: " + s
}

func Setup(t *testing.T, req *Request) error {
	req.Name = "root"
	return nil
}
```