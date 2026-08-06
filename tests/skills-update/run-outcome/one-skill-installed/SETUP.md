# Scenario

**Feature**: batch update touches only installed registry skills

```
doctest skill tdd --install -> doctest skills update -> "tdd  up to date" + not-installed for others
```

## Preconditions

- Only the `tdd` skill is installed to `.agents/skills/doctest-tdd`.

## Steps

1. Pre-install with `skill tdd --install`.
2. Run `skills update`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.PreInstalls = []PreInstallCLI{{
		Args: []string{"skill", "tdd", "--install"},
	}}
	req.Args = []string{"skills", "update"}
	return nil
}
```