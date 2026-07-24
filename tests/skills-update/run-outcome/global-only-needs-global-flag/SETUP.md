# Scenario

**Feature**: global-only install is not visible to default-scope batch update

```
# skill installed under $HOME/.agents/skills only
doctest skill tdd --install --global -> ~/.agents/skills/doctest-tdd/SKILL.md

# default-scope update sees no local installs for any registry skill
doctest skills update -> skill not installed line per registry name
```

## Preconditions

- `HOME` is an isolated temp directory (no pre-existing skill installs).

## Steps

1. Set `HOME` to a fresh temp dir.
2. Pre-install `tdd` with `--global`.
3. Run `skills update` without `--global`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	home := t.TempDir()
	// Child-only HOME — never t.Setenv (Parallel-incompatible).
	req.Home = home
	req.Env = append(req.Env, "HOME="+home)
	req.PreInstalls = []PreInstallCLI{{
		Args: []string{"skill", "tdd", "--install", "--global"},
	}}
	req.Args = []string{"skills", "update"}
	return nil
}
```