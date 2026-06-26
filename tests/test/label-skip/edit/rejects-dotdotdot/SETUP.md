# Scenario

**Feature**: edit rejects `...` patterns

```
# invalid target
doctest edit ./mod/... -> error exit
```

## Steps

1. Run `doctest edit ./any/... --add-label ui`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"edit", "./any/...", "--add-label", "ui-automation"}
	return nil
}
```