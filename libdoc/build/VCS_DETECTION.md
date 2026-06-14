# VCS Detection in Doctest Builds

## Background

### What is `-buildvcs`?

Go's `go build` and `go test -c` commands support a `-buildvcs` flag that controls whether version control information (commit hash, branch, etc.) is embedded (stamped) into the resulting binary. The flag accepts three values:

- `true` — always stamp; error out if VCS info is unavailable
- `false` — never stamp
- `auto` (default) — stamp only if the main package, main module, and current directory are all in the same repository

Since Go 1.18, this flag has existed. The default changed from `true` to `auto` in Go 1.21.

### What is Git Safe Directory?

Git's `safe.directory` configuration controls which repositories Git trusts. When enabled (the default in modern Git), operations like `git status` or `git rev-parse` fail with "detected dubious ownership" unless the directory is explicitly added to `safe.directory`.

This becomes a problem in containerized CI environments (Docker, GitHub Actions) where the filesystem user differs from the repository owner, causing Git to reject all operations.

## The Problem

Doctest creates temporary build directories (via `os.MkdirTemp`) to compile generated test code. These directories are **never** inside a Git working tree. When Go invokes `git` internally during `go build`/`go test -c` for VCS stamping, it can fail with:

```
error obtaining VCS status: exit status 128
Use -buildvcs=false to disable VCS stamping.
```

This occurs in CI environments where Git's `safe.directory` check fails, but the fundamental issue is that the build directory has no VCS information to stamp.

## The Solution

### Detection Logic

`needsBuildVCSFlag(dir string) bool` checks two conditions:

1. **Is `git` available on the system?** (`exec.LookPath("git")`)
2. **Is the build directory inside a Git working tree?** (`git -C <dir> rev-parse --is-inside-work-tree`)

If either check fails, the function returns `true` — meaning `-buildvcs=false` is needed.

### Where It's Used

- `libdoc/build/build.go` — `go build` command
- `libdoc/build/test.go` — `go test -c` command

Both insert `-buildvcs=false` into their argument lists when `needsBuildVCSFlag` returns true.

### Why Not Always Add `-buildvcs=false`?

Always adding the flag would work, but it's unnecessarily defensive. When the build directory *is* inside a Git tree (e.g., a user intentionally places a `GenDir` inside their repo), we want to preserve normal VCS stamping behavior to avoid surprising behavior.

### Why Not Try-First-Then-Retry?

Try-and-retry is error-prone and slow. It can mask real build errors, pollute caches, and adds latency. A static check based on Git availability is deterministic and fast.
