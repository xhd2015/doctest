# Scenario

**Feature**: `doctest test` label skip behavior

```
# discovery vs explicit leaf invocation
doctest test <target> -> skip or run labeled leaves
```

## Steps

1. Configure `req.Args` with `doctest test` and a temp tree target.

```go
func Setup(t *testing.T, req *Request) error {
	if req.Timeout == 0 {
		req.Timeout = 120 * time.Second
	}
	return nil
}
```