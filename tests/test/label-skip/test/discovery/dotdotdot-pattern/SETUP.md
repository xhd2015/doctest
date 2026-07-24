# Scenario

**Feature**: `./mod/...` discovery pattern skips labeled leaves

```
# relative ... pattern from work dir parent
doctest test ./mod/... -> skip labeled leaf under mod/
```

## Steps

1. Create `mod/` tree under a work dir with fast + labeled leaves.
2. Set `req.WorkDir` to parent and run `doctest test ./mod/...`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkDir = treeInWorkDir(t, "mod", true, "ui-automation", "dotdotdot skip")
	req.Args = []string{"test", "./mod/..."}
	return nil
}
```