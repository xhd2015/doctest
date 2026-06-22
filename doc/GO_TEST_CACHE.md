# Go Test Cache Invalidation: Directory mtime Side Effect

## Summary

When generating Go test files into a package directory, creating and deleting a
temporary file inside that directory (even one that is immediately removed)
changes the directory's modification time (mtime). If the generated test
binaries use `os.Chdir` (or otherwise access the directory), Go's test caching
infrastructure tracks the directory's stat information — including mtime — as
part of the cache key. This causes all cached test results for that package to
be invalidated on every run.

## Root Cause

### The sequence

1. `Doctest` generates a `_test.go` file into a leaf directory
2. For unchanged leaves, `WriteGeneratedCase` determines the content has not
   changed and skips the write
3. But the implementation used `os.CreateTemp(leafDir, ".doctest-gen-*")` to
   stage formatted source — **inside the package directory**
4. Even though the temp file is deleted after comparison, the directory's mtime
   is updated (create + delete = two mtime changes)
5. Go's `testcache` records the directory stat (including mtime) for packages
   whose tests access the filesystem
6. On the next `go test` invocation, the changed mtime produces a different
   cache key → cache miss → full recompilation and re-execution

### Why mtime matters

From `cmd/go/internal/test/test.go` in the Go source, the test result cache
uses a two-level lookup:

- First level: hash of the test binary + command-line flags
- Second level: list of environment variables and **files consulted** during the
  test run (including their `stat` information)

When a test calls `os.Chdir(leafDir)`, Go records `stat(leafDir)`. The `stat`
includes mtime. A temp file created and deleted inside that directory changes
its mtime, producing a different cache key.

### Cacheable flags constraint

The cache also requires that test flags come from a restricted set:

```
-benchtime, -coverprofile, -cpu, -failfast, -fullpath, -list,
-outputdir, -parallel, -run, -short, -skip, -timeout, -v
```

Flags outside this set (including non-test build flags like `-buildvcs`) are
passed through as test binary arguments. If they do not start with `-test.`
and contain `=`, they disable caching. However, `-buildvcs=false` is a build
flag (parsed before the package list), so it does **not** disable caching.

### Other cache-disabling conditions

Caching is **disabled** when:

- `go test` is invoked with no package arguments (local directory mode)
- The package is outside any module, GOPATH, or GOROOT
- A `testexpire.txt` file exists in the cache directory with a future timestamp
  (created by `go clean -testcache`)
- `GOCACHE=off` is set

Caching is **not** affected by:

- Shell vs direct invocation (`/bin/sh -c` vs `exec.Command`)
- `-buildvcs=false` flag (it is a build flag, not a test flag)
- `-mod=mod` flag
- File content comparison (same content = same hash, regardless of how written)
- `replace` directives in `go.mod`

## The Fix

In `libdoc/core/discover.go`, `WriteGeneratedCase`:

**Before (broken):**
```go
tmpFile, err := os.CreateTemp(leafDir, ".doctest-gen-*")
```

**After (fixed):**
```go
tmpFile, err := os.CreateTemp("", ".doctest-gen-*")
```

Creating the temp file in the system temporary directory (`""` = `os.TempDir()`)
instead of inside the package directory prevents the directory mtime from
changing. The temp file is then renamed into the leaf directory only when the
content actually differs.

## Session ID and environment variables

`doctest test` sets `DOCTEST_SESSION_ID` to a new UUID for each invocation and
injects it into generated tests as:

```go
DOCTEST_SESSION_ID, ok := syscall.Getenv("DOCTEST_SESSION_ID")
if !ok || DOCTEST_SESSION_ID == "" {
    t.Fatalf("DOCTEST_SESSION_ID not set")
}
```

Using `syscall.Getenv` (not `os.Getenv`) keeps the session value out of Go's
test-result cache key while still giving every package in the run the same
session id. See `tests/test/go-test-cache/env-getenv/` for proof tests.

## Key Lesson

Any file creation or deletion **inside a Go package directory** between
consecutive `go test` invocations can invalidate the test cache for that
package, if the test binary accesses the directory. This applies not just to
source files (`_test.go`) but to **any** file — including temporary files,
lock files, or log files.

The robust pattern for atomic file updates without cache invalidation is:

1. Create the temp file in the system temp directory (`os.TempDir()`)
2. Prepare the content
3. Only `os.Rename` into the target directory if the content differs from
   the existing file
4. Never create-and-delete files inside the package directory

## Proof from Go 1.25.10 Source

The full chain from `os.Chdir` to cache invalidation:

### 1. Test binary records filesystem operations

When `os.Chdir` is called inside a test binary, the `testlog` hook records it:

`/Users/xhd2015/installed/go1.25.10/src/os/file.go:376-380`
```go
if log := testlog.Logger(); log != nil {
    wd, err := Getwd()
    if err == nil {
        log.Chdir(wd)
    }
}
```

The `testlog.Logger` interface is defined at:

`/Users/xhd2015/installed/go1.25.10/src/internal/testlog/log.go:19-20`
```go
type Interface interface {
    ...
    Chdir(dir string)
}
```

The implementation that records `chdir` operations is at:

`/Users/xhd2015/installed/go1.25.10/src/testing/internal/testdeps/deps.go:88-90`
```go
func (l *testLog) Chdir(name string) {
    l.add("chdir", name)
}
```

This writes `"chdir /path/to/dir\n"` into the test log file.

### 2. Cache key computation includes directory mtime

After the test completes, `go test` reads the test log and computes a "test inputs
ID" that becomes part of the cache key. For each `chdir` entry, it calls
`hashStat`:

`/Users/xhd2015/installed/go1.25.10/src/cmd/go/internal/test/test.go:1999-2001`
```go
case "chdir":
    pwd = name // always absolute
    fmt.Fprintf(h, "chdir %s %x\n", name, hashStat(name))
```

`hashStat` (line 2088) calls `os.Stat` and `hashWriteStat`:

`/Users/xhd2015/installed/go1.25.10/src/cmd/go/internal/test/test.go:2088-2101`
```go
func hashStat(name string) cache.ActionID {
    h := cache.NewHash("stat")
    if info, err := os.Stat(name); err != nil {
        fmt.Fprintf(h, "err %v\n", err)
    } else {
        hashWriteStat(h, info)
    }
    ...
    return h.Sum()
}
```

`hashWriteStat` (line 2103) writes the FULL directory stat information — including
`Size()`, `Mode()`, `ModTime()`, and `IsDir()`:

`/Users/xhd2015/installed/go1.25.10/src/cmd/go/internal/test/test.go:2103-2105`
```go
func hashWriteStat(h io.Writer, info fs.FileInfo) {
    fmt.Fprintf(h, "stat %d %x %v %v\n",
        info.Size(), uint64(info.Mode()), info.ModTime(), info.IsDir())
}
```

**This is the smoking gun: `info.ModTime()` is written into the cache key hash.**
Any file created or deleted in the directory changes its `ModTime()`, producing a
different hash and therefore a cache miss.

### 3. Two-level cache key structure

`/Users/xhd2015/installed/go1.25.10/src/cmd/go/internal/test/test.go:2108-2110`
```go
func testAndInputKey(testID, testInputsID cache.ActionID) cache.ActionID {
    return cache.Subkey(testID, fmt.Sprintf("inputs:%x", testInputsID))
}
```

The cache key = `Subkey(testBinaryHash, "inputs:<testInputsID>")`. Since the
test binary hash depends on source files and the `testInputsID` depends on
`chdir` target directory stats (including mtime), any mtime change invalidates
both levels.

### 4. `modTimeCutoff` — files newer than 2s are rejected

`/Users/xhd2015/installed/go1.25.10/src/cmd/go/internal/test/test.go:2045-2084`
```go
const modTimeCutoff = 2 * time.Second
...
if time.Since(info.ModTime()) < modTimeCutoff {
    return cache.ActionID{}, errFileTooNew
}
```

This exists specifically because Go relies on mtime for caching of `open`'d
files. Files modified less than 2 seconds ago are rejected from caching to
account for filesystem timestamp granularity. This further confirms that mtime
is central to Go's cache key computation.

### 5. Scope: only files inside the package/module root

`/Users/xhd2015/installed/go1.25.10/src/cmd/go/internal/test/test.go:2006-2008`
```go
if a.Package.Root == "" || search.InDir(name, a.Package.Root) == "" {
    // Do not recheck files outside the module, GOPATH, or GOROOT root.
    break
}
```

`stat` and `open` entries are only tracked for paths inside the module root.
Files outside (e.g., in `/tmp/`) do not affect the cache key. This is why
creating temp files in the system temp directory avoids cache invalidation.

## References

- Go 1.25.10 source: `cmd/go/internal/test/test.go:1970-2110` — cache key computation
- Go 1.25.10 source: `cmd/go/internal/test/test.go:2103-2105` — `hashWriteStat` includes `ModTime()`
- Go 1.25.10 source: `os/file.go:376-380` — `os.Chdir` records to testlog
- Go 1.25.10 source: `testing/internal/testdeps/deps.go:88-90` — `Chdir` testlog implementation
- Go 1.25.10 source: `internal/testlog/log.go:19-20` — Logger interface with `Chdir`
- Go doc: `go help testflag` — cacheable test flags documentation
- Go doc: `go help cache` — build and test caching overview
