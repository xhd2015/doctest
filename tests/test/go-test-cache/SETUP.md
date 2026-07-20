# Scenario

**Feature**: the doctest binary is built by the root Setup

```
# build and run test binary, report results
doctest test <dir> -> build -> run binary -> pass/fail per leaf -> exit code

# path patterns
.../ -> walk tree | subdir -> run subtree | multi-dir -> aggregate results

# output
progress dots -> . F | verbose -> go test -v | count -> N tests
```

## Preconditions
- The doctest binary is built by the root Setup.
- Tests run with an extended timeout to allow for Go compilation (first run is slow).

## Steps
1. Set a generous timeout (120s) for tests that compile Go code.
2. Provide shared helpers for temp doctest trees.
3. Distinguish **observed** Setup/Run effects (go testcache miss) vs **unread** writes (may stay cached).

```go
import (
    "os"
    "path/filepath"
    "strings"
    "testing"
    "time"

    "github.com/xhd2015/doctest/libdoc/testtree"
)

var bt = "`" + "`" + "`"

func doctestGoBlock(code string) string {
    return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}

func doctestBody(extraRunCode string) string {
    return "import \"testing\"\n\ntype Request struct{ Args []string; WorkDir string }\ntype Response struct{ ExitCode int; Stdout string; Stderr string }\n\n" + extraRunCode
}

func defaultRunCode() string {
    return "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }"
}

func rootSetupContent() string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
}

func rootSetupWorkDir(tag string) string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { req.WorkDir = \"" + tag + "\"; return nil }")
}

// intermediateSetupWorkDir sets WorkDir so edits change a real field write.
func intermediateSetupWorkDir(tag string) string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { req.WorkDir = \"" + tag + "\"; return nil }")
}

// intermediateSetupDiscardString: only discards a string — typically DCE'd.
func intermediateSetupDiscardString(tag string) string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error {\n\t_ = req\n\t_ = \"" + tag + "\"\n\treturn nil\n}")
}

// intermediateSetupTLog: t.Log keeps the tag live in the test binary.
func intermediateSetupTLog(tag string) string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error {\n\t_ = req\n\tt.Log(\"" + tag + "\")\n\treturn nil\n}")
}

func leafSetupContent() string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }")
}

func leafSetupWorkDir(tag string) string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { req.WorkDir = \"" + tag + "\"; return nil }")
}

func leafAssertContent() string {
    return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}")
}

// leafAssertObserveWorkDir reads WorkDir so Setup writes stay live in the test binary
// (prevents DCE). Accepts any non-empty value so V1→V2 edits still pass.
func leafAssertObserveWorkDir() string {
    return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n\tif req.WorkDir == \"\" {\n\t\tt.Fatal(\"expected non-empty WorkDir from Setup\")\n\t}\n}")
}

// leafAssertTLogWorkDir always logs WorkDir (stronger live use than non-empty check alone).
func leafAssertTLogWorkDir() string {
    return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n\tt.Log(req.WorkDir)\n\tif req.WorkDir == \"\" {\n\t\tt.Fatal(\"expected non-empty WorkDir from Setup\")\n\t}\n}")
}

// leafAssertStdoutNonEmpty only checks non-empty — Run string swaps can still DCE.
func leafAssertStdoutNonEmpty() string {
    return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n\tif resp == nil || resp.Stdout == \"\" {\n\t\tt.Fatal(\"expected non-empty Stdout\")\n\t}\n}")
}

func runCodeWithStdout(tag string) string {
    return "func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Stdout: \"" + tag + "\"}, nil }"
}

func modifiedDiscardStringSetup(tag string) string {
    return intermediateSetupDiscardString(tag)
}

func modifiedTLogSetup(tag string) string {
    return intermediateSetupTLog(tag)
}



func writeFile(path, content string) error {
    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        return err
    }
    return os.WriteFile(path, []byte(content), 0644)
}

type treeOpts struct {
    // ObserveWorkDir: leaf ASSERT requires non-empty WorkDir; intermediates/root
    // that set WorkDir use workDir tags.
    ObserveWorkDir bool
    // TLogWorkDir: leaf ASSERT always t.Log(req.WorkDir) (strong live use).
    TLogWorkDir bool
    // RootWorkDirTag: if non-empty, root SETUP sets WorkDir to this tag.
    RootWorkDirTag string
    // LeafWorkDirTag: if non-empty, leaf SETUP sets WorkDir (for dead-write cases).
    LeafWorkDirTag string
    // MidWorkDirPrefix: intermediates set WorkDir to prefix+"-"+pathTag+"-v1".
    MidWorkDirPrefix string
    // MidDiscardString: intermediates use `_ = "tag"` only (dead discard).
    MidDiscardString bool
    // MidTLog: intermediates use t.Log("tag") (live side effect).
    MidTLog bool
    RunCode string
}

func createTestTree(dir string, extraRunCode string) error {
    return createTestTreeOpts(dir, treeOpts{RunCode: extraRunCode})
}

func createTestTreeOpts(dir string, opts treeOpts) error {
    runCode := opts.RunCode
    if runCode == "" {
        runCode = defaultRunCode()
    }
    if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(doctestBody(runCode))), 0644); err != nil {
        return err
    }
    root := rootSetupContent()
    if opts.RootWorkDirTag != "" {
        root = rootSetupWorkDir(opts.RootWorkDirTag)
    }
    if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(root), 0644); err != nil {
        return err
    }
    leafDir := filepath.Join(dir, "simple")
    if err := os.MkdirAll(leafDir, 0755); err != nil {
        return err
    }
    leafSetup := leafSetupContent()
    switch {
    case opts.LeafWorkDirTag != "":
        leafSetup = leafSetupWorkDir(opts.LeafWorkDirTag)
    case (opts.ObserveWorkDir || opts.TLogWorkDir) && opts.RootWorkDirTag == "":
        // Flat tree: leaf Setup owns WorkDir when root does not.
        leafSetup = leafSetupWorkDir("leaf-v1")
    }
    if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(leafSetup), 0644); err != nil {
        return err
    }
    assert := leafAssertContent()
    switch {
    case opts.TLogWorkDir:
        assert = leafAssertTLogWorkDir()
    case opts.ObserveWorkDir:
        assert = leafAssertObserveWorkDir()
    }
    if err := os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(assert), 0644); err != nil {
        return err
    }
    return nil
}

// createCustomTestTree builds a tree with optional intermediate SETUP.md dirs and leaf paths.
func createCustomTestTree(dir string, runCode string, intermediates []string, leaves []string) error {
    return createCustomTestTreeOpts(dir, treeOpts{RunCode: runCode}, intermediates, leaves)
}

func createCustomTestTreeOpts(dir string, opts treeOpts, intermediates []string, leaves []string) error {
    runCode := opts.RunCode
    if runCode == "" {
        runCode = defaultRunCode()
    }
    if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(doctestBody(runCode))), 0644); err != nil {
        return err
    }
    root := rootSetupContent()
    if opts.RootWorkDirTag != "" {
        root = rootSetupWorkDir(opts.RootWorkDirTag)
    }
    if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(root), 0644); err != nil {
        return err
    }
    prefix := opts.MidWorkDirPrefix
    if prefix == "" {
        prefix = "mid"
    }
    for _, mid := range intermediates {
        mid = filepath.FromSlash(mid)
        tag := strings.ReplaceAll(filepath.ToSlash(mid), "/", "_")
        full := prefix + "-" + tag + "-v1"
        var content string
        switch {
        case opts.MidTLog:
            content = intermediateSetupTLog(full)
        case opts.MidDiscardString:
            content = intermediateSetupDiscardString(full)
        default:
            content = intermediateSetupWorkDir(full)
        }
        if err := writeFile(filepath.Join(dir, mid, "SETUP.md"), content); err != nil {
            return err
        }
    }
    assert := leafAssertContent()
    switch {
    case opts.TLogWorkDir:
        assert = leafAssertTLogWorkDir()
    case opts.ObserveWorkDir:
        assert = leafAssertObserveWorkDir()
    }
    for _, leaf := range leaves {
        leaf = filepath.FromSlash(leaf)
        if err := writeFile(filepath.Join(dir, leaf, "SETUP.md"), leafSetupContent()); err != nil {
            return err
        }
        if err := writeFile(filepath.Join(dir, leaf, "ASSERT.md"), assert); err != nil {
            return err
        }
    }
    return nil
}

func createTempTestProject(t *testing.T, dirName string) string {
    t.Helper()
    return createTempTestProjectOpts(t, dirName, treeOpts{})
}

func createTempTestProjectOpts(t *testing.T, dirName string, opts treeOpts) string {
    t.Helper()
    tmp := t.TempDir()
    if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }
    testDir := filepath.Join(tmp, dirName)
    if err := os.MkdirAll(testDir, 0755); err != nil {
        t.Fatalf("mkdir test dir: %v", err)
    }
    if err := createTestTreeOpts(testDir, opts); err != nil {
        t.Fatalf("create test tree: %v", err)
    }
    return testDir
}

// createTempTestProjectObserveWorkDir: leaf Setup sets WorkDir; ASSERT reads it.
func createTempTestProjectObserveWorkDir(t *testing.T, dirName string) string {
    t.Helper()
    return createTempTestProjectOpts(t, dirName, treeOpts{ObserveWorkDir: true})
}

func createTempCustomProject(t *testing.T, dirName string, intermediates, leaves []string) string {
    t.Helper()
    return createTempCustomProjectOpts(t, dirName, treeOpts{}, intermediates, leaves)
}

func createTempCustomProjectOpts(t *testing.T, dirName string, opts treeOpts, intermediates, leaves []string) string {
    t.Helper()
    tmp := t.TempDir()
    if err := os.WriteFile(filepath.Join(tmp, "go.mod"), []byte("module testproj\ngo 1.21\n"), 0644); err != nil {
        t.Fatalf("write go.mod: %v", err)
    }
    testDir := filepath.Join(tmp, dirName)
    if err := os.MkdirAll(testDir, 0755); err != nil {
        t.Fatalf("mkdir test dir: %v", err)
    }
    if err := createCustomTestTreeOpts(testDir, opts, intermediates, leaves); err != nil {
        t.Fatalf("create custom tree: %v", err)
    }
    return testDir
}

// Observing variants: intermediate Setup sets WorkDir; leaf ASSERT reads it.
func createDeepChainProject(t *testing.T, dirName string) string {
    t.Helper()
    return createTempCustomProjectOpts(t, dirName, treeOpts{ObserveWorkDir: true},
        []string{"mid-a", "mid-a/mid-b", "mid-a/mid-b/mid-c"},
        []string{"mid-a/mid-b/mid-c/leaf-x"},
    )
}

func createSharedMidTwoLeavesProject(t *testing.T, dirName string) string {
    t.Helper()
    return createTempCustomProjectOpts(t, dirName, treeOpts{ObserveWorkDir: true},
        []string{"mid-a"},
        []string{"mid-a/leaf-x", "mid-a/leaf-y"},
    )
}

func createL1Project(t *testing.T, dirName string) string {
    t.Helper()
    return createTempCustomProjectOpts(t, dirName, treeOpts{ObserveWorkDir: true},
        []string{"mid-a"},
        []string{"mid-a/leaf-x"},
    )
}

func createL2Project(t *testing.T, dirName string) string {
    t.Helper()
    return createTempCustomProjectOpts(t, dirName, treeOpts{ObserveWorkDir: true},
        []string{"mid-a", "mid-a/mid-b"},
        []string{"mid-a/mid-b/leaf-x"},
    )
}

// Dead (unread WorkDir) trees — Setup writes; ASSERT empty → go may keep (cached).
func createL1ProjectDead(t *testing.T, dirName string) string {
    t.Helper()
    return createTempCustomProjectOpts(t, dirName, treeOpts{ObserveWorkDir: false},
        []string{"mid-a"},
        []string{"mid-a/leaf-x"},
    )
}

// createL1ProjectDeadDiscard: mid Setup only `_ = "tag"` (DCE-friendly).
func createL1ProjectDeadDiscard(t *testing.T, dirName string) string {
    t.Helper()
    return createTempCustomProjectOpts(t, dirName, treeOpts{MidDiscardString: true},
        []string{"mid-a"},
        []string{"mid-a/leaf-x"},
    )
}

// createL1ProjectMidTLog: mid Setup uses t.Log (live → cache miss on tag change).
func createL1ProjectMidTLog(t *testing.T, dirName string) string {
    t.Helper()
    return createTempCustomProjectOpts(t, dirName, treeOpts{MidTLog: true},
        []string{"mid-a"},
        []string{"mid-a/leaf-x"},
    )
}

// createL1ProjectWorkDirTLogAssert: mid sets WorkDir; ASSERT always t.Log(WorkDir).
func createL1ProjectWorkDirTLogAssert(t *testing.T, dirName string) string {
    t.Helper()
    return createTempCustomProjectOpts(t, dirName, treeOpts{TLogWorkDir: true},
        []string{"mid-a"},
        []string{"mid-a/leaf-x"},
    )
}

// createTempTestProjectLeafWorkDirDead: leaf Setup sets WorkDir; empty ASSERT.
func createTempTestProjectLeafWorkDirDead(t *testing.T, dirName string) string {
    t.Helper()
    return createTempTestProjectOpts(t, dirName, treeOpts{LeafWorkDirTag: "leaf-dead-v1"})
}

func createTempTestProjectRootWorkDir(t *testing.T, dirName, tag string) string {
    t.Helper()
    return createTempTestProjectOpts(t, dirName, treeOpts{
        ObserveWorkDir: true,
        RootWorkDirTag: tag,
    })
}

func createTempTestProjectRootWorkDirDead(t *testing.T, dirName, tag string) string {
    t.Helper()
    return createTempTestProjectOpts(t, dirName, treeOpts{
        ObserveWorkDir: false,
        RootWorkDirTag: tag,
    })
}

func modifiedSetupContent(tag string) string {
    return doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { req.WorkDir = \"" + tag + "\"; return nil }")
}

func modifiedAssertContent() string {
    return doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n    if resp.Stdout != \"modified\" {\n        t.Log(\"stdout was not modified\")\n    }\n}")
}

// runCodeWithLog uses t.Log so the tag string stays live in the test binary
// (plain Stdout constants can be DCE'd when only tested for non-empty).
func runCodeWithLog(tag string) string {
    return "func Run(t *testing.T, req *Request) (*Response, error) {\n\tt.Log(\"" + tag + "\")\n\treturn &Response{}, nil\n}"
}

func modifiedRunCode() string {
    return runCodeWithLog("run-edited")
}

func Setup(t *testing.T, req *Request) error {
    req.Timeout = 120 * time.Second
    return nil
}
```
