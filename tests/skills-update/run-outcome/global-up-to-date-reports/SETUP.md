# Scenario

**Feature**: `--global` batch update reports installed global skills

```
doctest skill tdd --install --global -> ~/.agents/skills/doctest-tdd/SKILL.md
doctest skills update --global -> Skill is up to date for global target
```

## Preconditions

- `HOME` is an isolated temp directory.

## Steps

1. Set `HOME` to a fresh temp dir.
2. Pre-install `tdd` with `--global`.
3. Run `skills update --global`.

```go
func Setup(t *testing.T, req *Request) error {
	home := t.TempDir()
	// Child-only HOME — never t.Setenv (Parallel-incompatible).
	req.Home = home
	req.Env = append(req.Env, "HOME="+home)
	req.PreInstalls = []PreInstallCLI{{
		Args: []string{"skill", "tdd", "--install", "--global"},
	}}
	req.Args = []string{"skills", "update", "--global"}
	return nil
}
```