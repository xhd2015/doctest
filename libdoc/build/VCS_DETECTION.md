# Go `-buildvcs` and doctest

## What is `-buildvcs`?

Go’s `go build` / `go test` support `-buildvcs` (`true` | `false` | `auto`).
Default since Go 1.21 is **`auto`**: stamp commit / dirty state into the binary
when the main package, main module, and current directory share one local repo.

When Go tries to stamp and `git status` / `git log` fail (exit 128), you see:

```text
error obtaining VCS status: exit status 128
Use -buildvcs=false to disable VCS stamping.
```

Common causes (see also `tests/build/buildvcs-auto/`):

| Situation | Result under `auto` |
|-----------|---------------------|
| Healthy full clone | OK, stamps `vcs.*` |
| Healthy shallow (`--depth 1`) | OK, stamps tip (shallow alone is **not** the fail axis) |
| Broken / untrusted git (`safe.directory`, corrupt HEAD, …) | **Fails** with the message above |

## What doctest does

Doctest **does not** inject `-buildvcs=false` into `go build` / `go test`.

Stamping follows the Go toolchain and the user’s environment (`GOFLAGS`, etc.).
If go reports `error obtaining VCS status`, doctest appends guidance:

```text
Error: go could not obtain VCS status (git failed while stamping; not caused by shallow clone alone)
hint: set GOFLAGS=-buildvcs=false and re-run
hint: CI (GitHub Actions) example:
hint:   env:
hint:     GOFLAGS: -buildvcs=false
hint: or fix git trust: git config --global --add safe.directory '*'
```

## Fix options

### Local shell

```bash
export GOFLAGS=-buildvcs=false
doctest test ./...
```

### GitHub Actions

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    env:
      GOFLAGS: -buildvcs=false
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
      - run: go test ./...
      - run: doctest test ./...
```

### Fix git trust (containers / mismatched UID)

```bash
git config --global --add safe.directory '*'
# or the specific workspace path
```

## Related tests

- Toolchain matrix (full / shallow / broken HEAD): `tests/build/buildvcs-auto/`
- Hint unit tests: `libdoc/build/vcs_test.go`
