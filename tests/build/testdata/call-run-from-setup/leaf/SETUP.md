# Scenario

**Feature**: tests for leaf

```
# parse test tree, generate Go code, compile binary
doctest build <test-dir> -> .md files -> Go code -> go build -> binary

# gen-dir controls output layout
gen-dir -> per-leaf packages -> file system
```

## Steps

- Call Run(t, req) from inside Setup.
- The generator emits Run as lowercase `run`, so uppercase `Run`
  in the Setup body will be undefined.

```go
func Setup(t *testing.T, req *Request) error {
    req.Name = "leaf"
    resp, runErr := Run(t, req)
    _ = resp
    _ = runErr
    return nil
}
```
