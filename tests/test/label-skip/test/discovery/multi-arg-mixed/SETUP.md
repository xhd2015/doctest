# Scenario

**Feature**: multi-arg invocation runs explicit labeled leaf and skips others

```
# discovery pattern plus explicit leaf in one command
doctest test ./mod/... <explicit-leaf> -> PASS(2/2) + skip other labeled
```

## Steps

1. Create `mod/` tree with fast, skip_labeled, and explicit_labeled leaves.
2. Run `doctest test ./mod/... <explicit_labeled>` from work dir parent.

```go
func Setup(t *testing.T, req *Request) error {
	workDir, explicitLeaf := writeMultiArgModTree(t)
	req.WorkDir = workDir
	req.Args = []string{"test", "./mod/...", explicitLeaf}
	return nil
}
```