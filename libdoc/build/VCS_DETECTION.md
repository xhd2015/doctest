# VCS Detection in Doctest Builds

## Background

### What is `-buildvcs`?

Go's `go build` and `go test -c` commands support a `-buildvcs` flag that controls whether version control information (commit hash, branch, etc.) is stamped into the resulting binary:

- `true` — always stamp; error out if VCS info is unavailable
- `false` — never stamp
- `auto` (default since Go 1.21) — stamp only if the main package, main module, and current directory are all in the same repository

### What is Git Safe Directory?

Git's `safe.directory` config controls which repositories Git trusts. When enabled (the default in modern Git), operations like `git status` or `git rev-parse` fail with "detected dubious ownership" unless the directory is added to `safe.directory`.

This is a problem in containerized CI (Docker, GitHub Actions) where the filesystem user differs from the repo owner. Notably, `actions/checkout@v4` sets `safe.directory` in a **temporary HOME override** that is discarded after checkout.

## The Problem

Doctest invokes `go build` / `go test -c` in two contexts, both vulnerable:

1. **Internal temp build dirs**: Created by `os.MkdirTemp` to compile generated test code — never inside a Git working tree.

2. **Project root builds**: Integration tests resolve a shared doctest binary via
   `libdoc/testbin.Ensure` (`go build -o $CACHE/doctest/selftest-bin/<key>/doctest ./cmd/doctest`).
   The module root IS a Git repo, but in CI containers Git may reject all operations due to "dubious ownership".

Both cases produce the same error:

```
error obtaining VCS status: exit status 128
Use -buildvcs=false to disable VCS stamping.
```

## The Solution

### Detection Logic

`NeedsBuildVCSFlag(dir string) bool` (in `libdoc/build/vcs.go`) checks two conditions:

1. **Is `git` available?** (`exec.LookPath("git")`)
2. **Does `git rev-parse --is-inside-work-tree` succeed *and* return true?**

If either check fails — git not found, directory not in a work tree, or git commands rejected (dubious ownership) — the function returns `true`.

### Where It's Used

**Engine code** (adds `-buildvcs=false` when the temp build dir has no git):

- `libdoc/build/build.go` — `go build` in generated temp dir
- `libdoc/build/test.go` — `go test -c` in generated temp dir

**Test specs** (adds `-buildvcs=false` when the project root cannot be accessed by git):

- `tests/SETUP.md` — builds doctest for integration tests
- `tests/implementer/SETUP.md` — builds doctest for implementer tests
- `tests/main-orchestrator/SETUP.md` — builds doctest for orchestrator tests
- `tests/test/nested-workdir/SETUP.md` — builds doctest for nested workdir tests

Each uses `libdocbuild.NeedsBuildVCSFlag(buildDir)` to conditionally add `-buildvcs=false`.

**CI workflow** (hardcoded in `.github/workflows/test.yml` since this is a shell command, not Go code):

```
go build -buildvcs=false -o doctest ./cmd/doctest
```

### Why Not Always Add `-buildvcs=false`?

When the build directory IS in a healthy Git tree (e.g., a user places `GenDir` inside their repo, or runs on a properly configured local machine), we preserve normal VCS stamping to avoid surprising behavior.

### Why Not Try-First-Then-Retry?

Try-and-retry is error-prone and slow. It can mask real build errors, pollute caches, and adds latency. A static check based on Git availability is deterministic and fast.
