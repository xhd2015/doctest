# Scenario

**Feature**: batch update with no local installs reports each registry skill

```
fresh WorkDir -> doctest skills update -> exit 0, not-installed line per skill
```

## Preconditions

- No `PreInstalls`.

## Steps

1. Run `skills update` only.

```go
func Setup(t *testing.T, req *Request) error {
	req.PreInstalls = nil
	req.Args = []string{"skills", "update"}
	return nil
}
```