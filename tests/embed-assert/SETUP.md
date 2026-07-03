# Scenario

**Feature**: doctest materializes embedded assert as a cached local module at test/build time

```
# assert import detected in SETUP/ASSERT Go blocks
doctest test/build <tree> -> CasesImportAssertPackage -> MaterializeAssertModule -> cache

# legacy nested module (no internal/)
assert import -> replace github.com/xhd2015/doctest/assert => <cache> in testcase go.mod

# internal compile path
internal + assert -> temp -modfile (parent go.mod + assert replace) -> go test -modfile=...
```

## Preconditions

- The doctest module root is two levels above this tree (`DOCTEST_ROOT/../..`).
- Each leaf builds a fresh doctest binary and runs it in a temp Go module.
- `GOWORK=off` is set for subprocess invocations.
- Assert cache lives at `$CACHE/doctest/assert-mod/<content-md5>/`.

## Steps

1. Build the doctest binary from the module root.
2. Create a temp module and doctest tree per leaf scenario.
3. Execute the doctest binary with leaf-specific args and inspect outputs/cache.

## Context

- Package-level vars `moduleRoot`, `testDir`, `genDir`, and `outsideGenDir` are
  set by shared helpers for nested-module scenarios.
- `internalModuleRoot` is set for internal-compile fixture copies.
- Helpers use `assertmod.RawSourceCacheKeyMD5()` for the cache path that
  `MaterializeAssertModule` writes (same contract whether cache is warm or cold).

```go
import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
	"testing"
	"time"
	libdocbuild "github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/assertmod"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

const (
	modPath		= "example.com/app"
	assertModPath	= "github.com/xhd2015/doctest/assert"
)
var (
	moduleRoot	string
	testDir		string
	bt		= string([]byte{96, 96, 96})
)
func lockCacheTests(t *testing.T) {
	t.Helper()
	lockPath := filepath.Join(os.TempDir(), "doctest-embed-assert-cache-tests.lock")
	f, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("open cache test lock: %v", err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
		f.Close()
		t.Fatalf("acquire cache test lock: %v", err)
	}
	t.Cleanup(func() {
		_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
		_ = f.Close()
	})
}
func Setup(t *testing.T, req *Request) error {
	req.Timeout = 120 * time.Second

	tmp := t.TempDir()
	doctestBin := filepath.Join(tmp, "doctest")
	buildDir := filepath.Join(DOCTEST_ROOT, "..", "..")
	buildArgs := []string{"build", "-o", doctestBin}
	if libdocbuild.NeedsBuildVCSFlag(buildDir) {
		buildArgs = append(buildArgs, "-buildvcs=false")
	}
	buildArgs = append(buildArgs, "./cmd/doctest")
	build := exec.Command("go", buildArgs...)
	build.Dir = buildDir
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build doctest: %v\n%s", err, string(out))
	}
	req.Bin = doctestBin
	return nil
}
func doctestGoBlock(code string) string {
	return "## Test\n\n" + bt + "go\n" + code + bt + "\n"
}
func createDoctestRoot(dir string, extraImports string, runCode string) error {
	goBody := "import (\n" +
		"\t\"testing\"\n" +
		extraImports +
		")\n\n" +
		"type Request struct{}\n" +
		"type Response struct{ Message string }\n\n" +
		runCode
	return os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(testtree.MinimalDOCTEST(goBody)), 0644)
}
func createDoctestLeaf(dir string, setupGo string, assertGo string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if setupGo == "" {
		setupGo = "import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { _ = req; return nil }"
	}
	if assertGo == "" {
		assertGo = "import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n" +
			"\tif err != nil { t.Fatal(err) }\n" +
			"\tif resp.Message != \"hi\" { t.Fatalf(\"expected hi, got %q\", resp.Message) }\n" +
			"}"
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(doctestGoBlock(setupGo)), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "ASSERT.md"), []byte(doctestGoBlock(assertGo)), 0644)
}
func copyDir(dst string, src string) error {
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		return os.WriteFile(target, data, 0644)
	})
}
func createPublicModuleProject(t *testing.T, leafSetupGo string, leafAssertGo string) {
	t.Helper()

	moduleRoot = t.TempDir()
	if err := os.WriteFile(
		filepath.Join(moduleRoot, "go.mod"),
		[]byte("module "+modPath+"\n\ngo 1.21\n"),
		0644,
	); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	greetDir := filepath.Join(moduleRoot, "pkg", "greet")
	if err := os.MkdirAll(greetDir, 0755); err != nil {
		t.Fatalf("mkdir pkg/greet: %v", err)
	}
	if err := os.WriteFile(
		filepath.Join(greetDir, "greet.go"),
		[]byte("package greet\n\nfunc Hello() string { return \"hi\" }\n"),
		0644,
	); err != nil {
		t.Fatalf("write greet.go: %v", err)
	}

	testDir = filepath.Join(moduleRoot, "tests")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("mkdir tests: %v", err)
	}
	extraImports := "\t\"" + modPath + "/pkg/greet\"\n"
	runCode := "func Run(t *testing.T, req *Request) (*Response, error) {\n" +
		"\treturn &Response{Message: greet.Hello()}, nil\n" +
		"}"
	if err := createDoctestRoot(testDir, extraImports, runCode); err != nil {
		t.Fatalf("create doctest root: %v", err)
	}
	if err := createDoctestLeaf(filepath.Join(testDir, "leaf"), leafSetupGo, leafAssertGo); err != nil {
		t.Fatalf("create doctest leaf: %v", err)
	}
}
func copyFixtureModule(t *testing.T, fixtureRel string) {
	t.Helper()

	fixtureSrc := filepath.Join(DOCTEST_ROOT, fixtureRel)
	moduleRoot = t.TempDir()
	if err := copyDir(moduleRoot, fixtureSrc); err != nil {
		t.Fatalf("copy fixture %s: %v", fixtureRel, err)
	}
	testDir = filepath.Join(moduleRoot, "tests")
}
func createInternalOnlyProject(t *testing.T) {
	copyFixtureModule(t, "testdata/internal-only-module")
}
func createInternalAssertProject(t *testing.T) {
	copyFixtureModule(t, "testdata/internal-assert-module")
}
func setupModuleEnv(t *testing.T, req *Request) {
	t.Helper()
	req.WorkDir = moduleRoot
	req.Env = append(req.Env, "GOWORK=off")
}
func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file %s to exist: %v", path, err)
	}
}
func assertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("expected file %s to not exist", path)
	}
}
func generatedLeafTestPath(genRoot string) string {
	return filepath.Join(genRoot, "tests", "leaf", "leaf_test.go")
}
func findDoctestRunDirs(root string) ([]string, error) {
	var found []string
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), ".doctest_run_") {
			found = append(found, filepath.Join(root, entry.Name()))
		}
	}
	return found, nil
}
func assertNoDoctestRunDirs(t *testing.T, root string) {
	t.Helper()
	dirs, err := findDoctestRunDirs(root)
	if err != nil {
		t.Fatalf("scan for .doctest_run_* dirs: %v", err)
	}
	if len(dirs) > 0 {
		t.Fatalf("expected no .doctest_run_* dirs under %s, found: %v", root, dirs)
	}
}
func assertStderrUsesTempCompile(t *testing.T, resp *Response) {
	t.Helper()
	combined := resp.Stdout + resp.Stderr
	if !strings.Contains(combined, ".doctest_run_") {
		t.Fatalf("expected stderr/stdout to reference .doctest_run_ temp compile dir, got:\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
}
func assertStderrUsesModfile(t *testing.T, resp *Response) {
	t.Helper()
	combined := resp.Stdout + resp.Stderr
	if !strings.Contains(combined, "-modfile=") {
		t.Fatalf("expected stderr/stdout to include -modfile= for internal-compile + assert, got:\nstdout:\n%s\nstderr:\n%s", resp.Stdout, resp.Stderr)
	}
}
func assertDumpNoNestedGoMod(t *testing.T, dumpRoot string) {
	t.Helper()
	assertFileNotExists(t, filepath.Join(dumpRoot, "go.mod"))
}
func assertNestedGoMod(t *testing.T, dir string) {
	t.Helper()
	nestedGoMod := filepath.Join(dir, "go.mod")
	assertFileExists(t, nestedGoMod)
	goModData, readErr := os.ReadFile(nestedGoMod)
	if readErr != nil {
		t.Fatalf("read nested go.mod: %v", readErr)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "module testcase") {
		t.Fatalf("expected nested module testcase, got:\n%s", goMod)
	}
	if !strings.Contains(goMod, "replace "+modPath+" =>") {
		t.Fatalf("expected replace directive for parent module, got:\n%s", goMod)
	}
}
func assertModCacheRoot() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "doctest", "assert-mod"), nil
}
func assertAssertReplaceInGoMod(t *testing.T, goModPath string) {
	t.Helper()
	data, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	goMod := string(data)
	needle := "replace " + assertModPath + " =>"
	if !strings.Contains(goMod, needle) {
		t.Fatalf("expected assert replace in go.mod, got:\n%s", goMod)
	}
	cacheRoot, cacheErr := assertModCacheRoot()
	if cacheErr != nil {
		t.Fatalf("assert mod cache root: %v", cacheErr)
	}
	if !strings.Contains(goMod, cacheRoot) {
		t.Fatalf("expected assert replace to point under %s, got:\n%s", cacheRoot, goMod)
	}
}
func assertNoAssertReplaceInGoMod(t *testing.T, goModPath string) {
	t.Helper()
	data, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if strings.Contains(string(data), "replace "+assertModPath) {
		t.Fatalf("expected no assert replace in go.mod, got:\n%s", string(data))
	}
}
func expectedAssertCacheDir(t *testing.T) string {
	t.Helper()
	root, err := assertModCacheRoot()
	if err != nil {
		t.Fatalf("cache root: %v", err)
	}
	return filepath.Join(root, assertmod.RawSourceCacheKeyMD5())
}
func listAssertModCacheEntries(t *testing.T) []string {
	t.Helper()
	root, err := assertModCacheRoot()
	if err != nil {
		t.Fatalf("cache root: %v", err)
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		t.Fatalf("read assert-mod cache: %v", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names
}
func assertCacheLayout(t *testing.T, cacheDir string) {
	t.Helper()
	assertFileExists(t, filepath.Join(cacheDir, "assert.go"))
	assertFileExists(t, filepath.Join(cacheDir, "go.mod"))
	goModData, err := os.ReadFile(filepath.Join(cacheDir, "go.mod"))
	if err != nil {
		t.Fatalf("read cached go.mod: %v", err)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "module "+assertModPath) {
		t.Fatalf("expected cached module path %s, got:\n%s", assertModPath, goMod)
	}
	if !strings.Contains(goMod, "go 1.18") {
		t.Fatalf("expected cached go.mod to require go 1.18, got:\n%s", goMod)
	}
	cachedAssert, err := os.ReadFile(filepath.Join(cacheDir, "assert.go"))
	if err != nil {
		t.Fatalf("read cached assert.go: %v", err)
	}
	if !strings.Contains(string(cachedAssert), "package assert") {
		t.Fatalf("cached assert.go missing package assert declaration")
	}
	if strings.Contains(string(cachedAssert), "_test.go") {
		t.Fatalf("cached assert.go must not include test file markers")
	}
	legacyNames, err := assertmod.LegacyV1Filenames()
	if err != nil {
		t.Fatalf("legacy_v1 filenames: %v", err)
	}
	if len(legacyNames) == 0 {
		t.Fatal("expected embedded legacy_v1 sources")
	}
	for _, name := range legacyNames {
		legacyPath := filepath.Join(cacheDir, "legacy_v1", name)
		assertFileExists(t, legacyPath)
		data, readErr := os.ReadFile(legacyPath)
		if readErr != nil {
			t.Fatalf("read cached legacy_v1/%s: %v", name, readErr)
		}
		if !strings.Contains(string(data), "package legacy_v1") {
			t.Fatalf("cached legacy_v1/%s missing package legacy_v1", name)
		}
	}
}
func snapshotFileState(t *testing.T, path string) (int64, [16]byte) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	sum := md5.New()
	if _, err := io.Copy(sum, f); err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var digest [16]byte
	copy(digest[:], sum.Sum(nil))
	return info.ModTime().UnixNano(), digest
}
func defaultAssertAssertGo() string {
	return "import (\n" +
		"\t\"testing\"\n" +
		"\t\"github.com/xhd2015/doctest/assert\"\n" +
		")\n\n" +
		"func Assert(t *testing.T, req *Request, resp *Response, err error) {\n" +
		"\tif err != nil { t.Fatal(err) }\n" +
		"\tassert.Output(t, \"hello\\n\", \"hello\\n\")\n" +
		"}"
}
func aliasedAssertAssertGo() string {
	return "import (\n" +
		"\t\"testing\"\n" +
		"\toutputassert \"github.com/xhd2015/doctest/assert\"\n" +
		")\n\n" +
		"func Assert(t *testing.T, req *Request, resp *Response, err error) {\n" +
		"\tif err != nil { t.Fatal(err) }\n" +
		"\toutputassert.Output(t, \"hello\\n\", \"hello\\n\")\n" +
		"}"
}
func defaultPublicAssertGo() string {
	return "import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n" +
		"\tif err != nil { t.Fatal(err) }\n" +
		"\tif resp.Message != \"hi\" { t.Fatalf(\"expected hi, got %q\", resp.Message) }\n" +
		"}"
}
```