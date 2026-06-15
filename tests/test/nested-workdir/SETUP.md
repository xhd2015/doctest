## Preconditions
- A valid doc-style test tree exists in the module.
- The command is invoked from inside a nested module directory.

## Steps
1. Set the process working directory to the module root (`DOCTEST_ROOT/..`).
2. Run `doctest test <absolute-dir>`.

```go
import (
    "os/exec"
    "path/filepath"
    "testing"

    libdocbuild "github.com/xhd2015/doctest/libdoc/build"
)

func Setup(t *testing.T, req *Request) error {
    exampleDir := filepath.Join(DOCTEST_ROOT, "testdata/basic-request-runner")
    req.WorkDir = filepath.Join(DOCTEST_ROOT, "..")
    req.Args = []string{"test", exampleDir}

    tmp := t.TempDir()
    doctestBin := filepath.Join(tmp, "doctest")
    buildDTDir := filepath.Join(DOCTEST_ROOT, "..")
    buildDTArgs := []string{"build", "-o", doctestBin}
    if libdocbuild.NeedsBuildVCSFlag(buildDTDir) {
        buildDTArgs = append(buildDTArgs, "-buildvcs=false")
    }
    buildDTArgs = append(buildDTArgs, "./cmd/doctest")
    buildDT := exec.Command("go", buildDTArgs...)
    buildDT.Dir = buildDTDir
    if out, err := buildDT.CombinedOutput(); err != nil {
        t.Fatalf("build doctest: %v\n%s", err, string(out))
    }
    req.Bin = doctestBin
    return nil
}
```
