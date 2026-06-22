package build

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/core"
)

func TestChildCannotRedefineRun(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{ Source string }
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Source: "root"}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Source: "leaf"}, nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	err := Test(root, core.Options{RemoveTemp: true})
	if err == nil {
		t.Fatal("expected child Run redefine error")
	}
	if !strings.Contains(err.Error(), "cannot redefine Run") {
		t.Fatalf("expected cannot redefine Run, got %v", err)
	}
}

func TestExecutionOrderSetupRunAssert(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{ Order []string }
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { req.Order = append(req.Order, "run"); return &Response{}, nil }
`, `
func Setup(t *testing.T, req *Request) error { req.Order = append(req.Order, "root setup"); return nil }
`)
	writeTreeFile(t, root, "parent/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Order = append(req.Order, "parent setup"); return nil }
`))
	writeTreeFile(t, root, "parent/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Order = append(req.Order, "leaf setup"); return nil }
`))
	writeTreeFile(t, root, "parent/leaf/ASSERT.md", assertDoc(`
import "reflect"
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	req.Order = append(req.Order, "assert")
	want := []string{"root setup", "parent setup", "leaf setup", "run", "assert"}
	if !reflect.DeepEqual(req.Order, want) { t.Fatalf("order = %#v, want %#v", req.Order, want) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestSetupErrorFailsBeforeRunAndAssert(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
import "fmt"
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { t.Fatal("run should not execute"); return nil, nil }
`, `
func Setup(t *testing.T, req *Request) error { return fmt.Errorf("setup failed") }
`)
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	t.Fatal("assert should not execute")
}
`))

	err := Test(root, core.Options{RemoveTemp: true})
	if err == nil {
		t.Fatal("expected setup failure")
	}
	if !strings.Contains(err.Error(), "setup failed") && !strings.Contains(err.Error(), "exit status 1") {
		t.Fatalf("expected setup failure in error, got %v", err)
	}
}

func TestRunErrorPassedToAssert(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
import "fmt"
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return nil, fmt.Errorf("run failed") }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil || err.Error() != "run failed" { t.Fatalf("expected run failed error, got %v", err) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestRequestMutatedThroughSetupChain(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{ Value int }
type Response struct{ Value int }
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Value: req.Value}, nil }
`, `
func Setup(t *testing.T, req *Request) error { req.Value += 1; return nil }
`)
	writeTreeFile(t, root, "parent/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Value += 2; return nil }
`))
	writeTreeFile(t, root, "parent/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Value += 3; return nil }
`))
	writeTreeFile(t, root, "parent/leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if req.Value != 6 || resp.Value != 6 { t.Fatalf("req=%d resp=%d, want 6", req.Value, resp.Value) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestResponsePassedFromRunToAssert(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{ Message string }
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Message: "ok"}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if resp == nil || resp.Message != "ok" { t.Fatalf("unexpected response: %#v", resp) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestDuplicateSetupHooksAcrossLevelsAllowed(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{ Count int }
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, `
func Setup(t *testing.T, req *Request) error { req.Count++; return nil }
`)
	writeTreeFile(t, root, "parent/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Count++; return nil }
`))
	writeTreeFile(t, root, "parent/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Count++; return nil }
`))
	writeTreeFile(t, root, "parent/leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if req.Count != 3 { t.Fatalf("count = %d, want 3", req.Count) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestDuplicateNamesAcrossDifferentLeavesAllowed(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{ Name string }
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Name: "root"}, nil }
`, "")
	writeTreeFile(t, root, "a/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "a/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	var duplicate = "a"
	if duplicate != "a" { t.Fatal(duplicate) }
}
`))
	writeTreeFile(t, root, "b/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "b/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	var duplicate = "b"
	if duplicate != "b" { t.Fatal(duplicate) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestChildCannotRedefineRequestOrResponse(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
type Request struct{ Bad bool }
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	err := Test(root, core.Options{RemoveTemp: true})
	if err == nil {
		t.Fatal("expected child Request redefinition error")
	}
	if !strings.Contains(err.Error(), "Request") {
		t.Fatalf("expected Request conflict error, got %v", err)
	}
}

func TestGeneratedCodeHasDoctestRootConst(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	leafTestData, err := os.ReadFile(filepath.Join(genDir, "leaf", "leaf_test.go"))
	if err != nil {
		t.Fatalf("read leaf_test.go: %v", err)
	}
	code := string(leafTestData)

	absRoot, _ := filepath.Abs(root)
	if !strings.Contains(code, "const DOCTEST_ROOT = `"+absRoot+"`") {
		t.Fatalf("expected DOCTEST_ROOT const with path %q, got:\n%s", absRoot, code)
	}
	if !strings.Contains(code, "DOCTEST_SESSION_ID, __sessionOk := syscall.Getenv(\"DOCTEST_SESSION_ID\")") {
		t.Fatalf("expected DOCTEST_SESSION_ID from syscall.Getenv, got:\n%s", code)
	}
	if !strings.Contains(code, "t.Fatalf(\"DOCTEST_SESSION_ID not set\")") {
		t.Fatalf("expected DOCTEST_SESSION_ID missing fatal, got:\n%s", code)
	}
	if !strings.Contains(code, "os.Chdir(filepath.Join(DOCTEST_ROOT, \"leaf\"))") {
		t.Fatalf("expected os.Chdir(filepath.Join(DOCTEST_ROOT, \"leaf\")), got:\n%s", code)
	}
	if !strings.Contains(code, "__origWd, __wdErr := os.Getwd()") {
		t.Fatalf("expected os.Getwd() before chdir, got:\n%s", code)
	}
	if !strings.Contains(code, "defer os.Chdir(__origWd)") {
		t.Fatalf("expected defer os.Chdir(__origWd), got:\n%s", code)
	}
	if strings.Contains(code, "func init()") {
		t.Fatal("expected no func init() in generated code")
	}
}

func TestGeneratedRootCaseChdirsToDoctestRoot(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	rootTestData, err := os.ReadFile(filepath.Join(genDir, "root_test.go"))
	if err != nil {
		t.Fatalf("read root_test.go: %v", err)
	}
	code := string(rootTestData)

	absRoot, _ := filepath.Abs(root)
	if !strings.Contains(code, "const DOCTEST_ROOT = `"+absRoot+"`") {
		t.Fatalf("expected DOCTEST_ROOT const with path %q, got:\n%s", absRoot, code)
	}
	if !strings.Contains(code, "DOCTEST_SESSION_ID, __sessionOk := syscall.Getenv(\"DOCTEST_SESSION_ID\")") {
		t.Fatalf("expected DOCTEST_SESSION_ID from syscall.Getenv, got:\n%s", code)
	}
	if !strings.Contains(code, "t.Fatalf(\"DOCTEST_SESSION_ID not set\")") {
		t.Fatalf("expected DOCTEST_SESSION_ID missing fatal, got:\n%s", code)
	}
	if !strings.Contains(code, "os.Chdir(DOCTEST_ROOT)") {
		t.Fatalf("expected os.Chdir(DOCTEST_ROOT) for root-level case, got:\n%s", code)
	}
	if strings.Contains(code, "func init()") {
		t.Fatal("expected no func init() in generated code")
	}
}

func TestCompileGoModGeneratedWithModule(t *testing.T) {
	srcRoot := t.TempDir()
	writeTreeFile(t, srcRoot, "go.mod", "module example.com/test\n\ngo 1.21\n")
	writeTreeFile(t, srcRoot, "tests/README.md", "# tree")
	writeRootHarness(t, filepath.Join(srcRoot, "tests"), `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, srcRoot, "tests/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, srcRoot, "tests/leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(filepath.Join(srcRoot, "tests"), core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	goModData, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("expected go.mod in gen dir: %v", err)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "module testcase") {
		t.Fatalf("expected module testcase, got %q", goMod)
	}
	if !strings.Contains(goMod, "example.com/test") {
		t.Fatalf("expected require example.com/test, got %q", goMod)
	}
	if !strings.Contains(goMod, "replace example.com/test => "+srcRoot) {
		t.Fatalf("expected replace with abs path %q, got %q", srcRoot, goMod)
	}
}

func TestCompileNoGoModWhenNoSourceModule(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	goModData, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("expected go.mod in gen dir: %v", err)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "module testcase") {
		t.Fatalf("expected module testcase, got %q", goMod)
	}
	if strings.Contains(goMod, "replace ") {
		t.Fatalf("expected no replace directive when no source module, got %q", goMod)
	}
}

func TestCompileWithExternalPackageImport(t *testing.T) {
	srcRoot := t.TempDir()
	writeTreeFile(t, srcRoot, "go.mod", "module example.com/mylib\n\ngo 1.21\n")
	writeTreeFile(t, srcRoot, "pkg/helper.go", "package pkg\n\nfunc Name() string { return \"hello\" }\n")

	testsDir := filepath.Join(srcRoot, "tests")
	writeTreeFile(t, testsDir, "README.md", "# tree")
	writeTreeFile(t, testsDir, "DOCTEST.md", doctestDoc(`
import "example.com/mylib/pkg"

type Request struct{}
type Response struct{ Val string }
func Run(t *testing.T, req *Request) (*Response, error) {
	return &Response{Val: pkg.Name()}, nil
}
`))
	writeTreeFile(t, testsDir, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, testsDir, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if resp.Val != "hello" { t.Fatalf("expected hello, got %q", resp.Val) }
}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(testsDir, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	goModData, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("expected go.mod in gen dir: %v", err)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "replace example.com/mylib => "+srcRoot) {
		t.Fatalf("expected replace with path %q, got %q", srcRoot, goMod)
	}

	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Dir = genDir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test in gen dir failed: %v\n%s", err, out)
	}
}

func TestSourceFilesCopiedWithPackageRename(t *testing.T) {
	srcRoot := t.TempDir()
	writeTreeFile(t, srcRoot, "go.mod", "module example.com/mylib\n\ngo 1.21\n")
	writeTreeFile(t, srcRoot, "pkg/helper.go", "package pkg\n\nfunc Exported() string { return \"exported\" }\nfunc private() string { return \"private\" }\n")
	writeTreeFile(t, srcRoot, "pkg/types.go", "package pkg\n\ntype MyType struct{ Val string }\n")

	testsDir := filepath.Join(srcRoot, "tests")
	writeTreeFile(t, testsDir, "README.md", "# tree")
	writeTreeFile(t, testsDir, "DOCTEST.md", doctestDoc(`import "example.com/mylib/pkg"

type Request struct{}
type Response struct{ Msg string }
func Run(t *testing.T, req *Request) (*Response, error) {
	return &Response{Msg: Exported() + " " + private()}, nil
}`))
	writeTreeFile(t, testsDir, "SETUP.md", `# Setup
- Go module: example.com/mylib
- Package under test: pkg

`+setupDoc(`func Setup(t *testing.T, req *Request) error { _ = req; return nil }`))

	writeTreeFile(t, testsDir, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, testsDir, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if resp.Msg != "exported private" { t.Fatalf("got %q", resp.Msg) }
}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(testsDir, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	helperPath := filepath.Join(genDir, "tests", "leaf", "helper.go")
	helperData, err := os.ReadFile(helperPath)
	if err != nil {
		t.Fatalf("expected helper.go in gen dir: %v", err)
	}
	if !strings.Contains(string(helperData), "package pkg_tc") {
		t.Fatalf("expected package pkg_tc in helper.go, got:\n%s", helperData)
	}
	if strings.Contains(string(helperData), "package pkg\n") {
		t.Fatal("expected original package name to be replaced")
	}
	if !strings.Contains(string(helperData), "func Exported()") {
		t.Fatal("expected Exported func in copied helper.go")
	}
	if !strings.Contains(string(helperData), "func private()") {
		t.Fatal("expected private func in copied helper.go")
	}

	leafTestPath := filepath.Join(genDir, "tests", "leaf", "leaf_test.go")
	testData, err := os.ReadFile(leafTestPath)
	if err != nil {
		t.Fatalf("expected leaf_test.go: %v", err)
	}
	if !strings.Contains(string(testData), "package pkg_tc") {
		t.Fatalf("expected package pkg_tc in test file, got:\n%s", testData)
	}
}

func TestNoTestGoFilesCopied(t *testing.T) {
	srcRoot := t.TempDir()
	writeTreeFile(t, srcRoot, "go.mod", "module example.com/mylib\n\ngo 1.21\n")
	writeTreeFile(t, srcRoot, "pkg/helper.go", "package pkg\n\nfunc Ok() string { return \"ok\" }\n")
	writeTreeFile(t, srcRoot, "pkg/helper_test.go", "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n")

	testsDir := filepath.Join(srcRoot, "tests")
	writeTreeFile(t, testsDir, "README.md", "# tree")
	writeTreeFile(t, testsDir, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, testsDir, "SETUP.md", `# Setup
- Go module: example.com/mylib
- Package under test: pkg

`+setupDoc(`func Setup(t *testing.T, req *Request) error { _ = req; return nil }`))

	writeTreeFile(t, testsDir, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, testsDir, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(testsDir, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	if _, err := os.Stat(filepath.Join(genDir, "tests", "leaf", "helper.go")); os.IsNotExist(err) {
		t.Fatal("expected helper.go to be copied")
	}
	if _, err := os.Stat(filepath.Join(genDir, "tests", "leaf", "helper_test.go")); !os.IsNotExist(err) {
		t.Fatal("expected helper_test.go NOT to be copied")
	}
}

func TestBackwardCompatNoPackageUnderTest(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	leafTestData, err := os.ReadFile(filepath.Join(genDir, "leaf", "leaf_test.go"))
	if err != nil {
		t.Fatalf("read leaf_test.go: %v", err)
	}
	if !strings.Contains(string(leafTestData), "package testcase") {
		t.Fatalf("expected package testcase when no metadata, got:\n%s", leafTestData)
	}
}

func TestCompileDefaultKeepsTempDir(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	var stderr bytes.Buffer
	if err := Build(root, core.Options{Stderr: &stderr}); err != nil {
		t.Fatalf("build: %v", err)
	}

	out := stderr.String()
	i := strings.Index(out, "\u2192 ")
	if i < 0 {
		t.Fatalf("expected gen dir in output, got %q", out)
	}
	rest := out[i+len("\u2192 "):]
	genDir := strings.TrimSpace(strings.SplitN(rest, "\n", 2)[0])

	if _, err := os.Stat(genDir); os.IsNotExist(err) {
		t.Fatalf("expected temp dir %s to be kept, but it was removed", genDir)
	}
	defer os.RemoveAll(genDir)
}

func TestCompileRmRemovesTempDir(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	var stderr bytes.Buffer
	if err := Build(root, core.Options{RemoveTemp: true, Stderr: &stderr}); err != nil {
		t.Fatalf("build: %v", err)
	}

	out := stderr.String()
	i := strings.Index(out, "\u2192 ")
	if i < 0 {
		t.Fatalf("expected gen dir in output, got %q", out)
	}
	rest := out[i+len("\u2192 "):]
	genDir := strings.TrimSpace(strings.SplitN(rest, "\n", 2)[0])

	if _, err := os.Stat(genDir); !os.IsNotExist(err) {
		t.Fatalf("expected temp dir %s to be removed, but it still exists", genDir)
	}
}

func TestCompileRmDoesNotRemoveGenDir(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir, RemoveTemp: true}); err != nil {
		t.Fatalf("build: %v", err)
	}

	if _, err := os.Stat(genDir); os.IsNotExist(err) {
		t.Fatalf("expected gen dir %s to exist with --rm, but it was removed", genDir)
	}
}

func TestAtLeastOneRunRequiredInSetupChain(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
`))
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	err := Test(root, core.Options{RemoveTemp: true})
	if err == nil {
		t.Fatal("expected missing Run error")
	}
	if !strings.Contains(err.Error(), "Run") {
		t.Fatalf("expected missing Run error, got %v", err)
	}
}
