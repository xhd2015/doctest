package runner

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeEmptyFile(t *testing.T, dir, name string) {
	t.Helper()
	writeFile(t, dir, name, "")
}

func TestResolveRoot(t *testing.T) {
	t.Run("dir has DOCTEST.md", func(t *testing.T) {
		dir := t.TempDir()
		writeEmptyFile(t, dir, "DOCTEST.md")
		writeEmptyFile(t, dir, "SETUP.md")

		root, ok := ResolveRoot(dir)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != dir {
			t.Fatalf("root = %q, want %q", root, dir)
		}
	})

	t.Run("ancestor has DOCTEST.md", func(t *testing.T) {
		dir := t.TempDir()
		writeEmptyFile(t, dir, "DOCTEST.md")
		writeEmptyFile(t, dir, "SETUP.md")
		leaf := filepath.Join(dir, "group-a", "leaf-1")
		writeEmptyFile(t, leaf, "ASSERT.md")

		root, ok := ResolveRoot(leaf)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != dir {
			t.Fatalf("root = %q, want %q", root, dir)
		}
	})

	t.Run("no DOCTEST, dir has SETUP.md", func(t *testing.T) {
		dir := t.TempDir()
		writeEmptyFile(t, dir, "SETUP.md")
		leaf := filepath.Join(dir, "leaf")
		if err := os.MkdirAll(leaf, 0755); err != nil {
			t.Fatal(err)
		}

		root, ok := ResolveRoot(leaf)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != dir {
			t.Fatalf("root = %q, want %q", root, dir)
		}
	})

	t.Run("no DOCTEST, ancestor has SETUP.md with go.mod above", func(t *testing.T) {
		dir := t.TempDir()
		writeEmptyFile(t, dir, "go.mod")
		writeEmptyFile(t, dir, "SETUP.md")
		sub := filepath.Join(dir, "sub")
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatal(err)
		}

		root, ok := ResolveRoot(sub)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != dir {
			t.Fatalf("root = %q, want %q", root, dir)
		}
	})

	t.Run("farthest ancestor with SETUP.md wins", func(t *testing.T) {
		dir := t.TempDir()
		writeEmptyFile(t, dir, "SETUP.md")
		mid := filepath.Join(dir, "mid")
		writeEmptyFile(t, mid, "SETUP.md")
		leaf := filepath.Join(mid, "leaf")
		if err := os.MkdirAll(leaf, 0755); err != nil {
			t.Fatal(err)
		}

		root, ok := ResolveRoot(leaf)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != dir {
			t.Fatalf("root = %q (farthest ancestor), want %q", root, dir)
		}
	})

	t.Run("go.mod between dir and DOCTEST.md ancestor", func(t *testing.T) {
		dir := t.TempDir()
		writeEmptyFile(t, dir, "DOCTEST.md")
		writeEmptyFile(t, dir, "SETUP.md")

		inner := filepath.Join(dir, "inner")
		writeEmptyFile(t, inner, "go.mod")
		writeEmptyFile(t, inner, "SETUP.md")
		sub := filepath.Join(inner, "sub")
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatal(err)
		}

		root, ok := ResolveRoot(sub)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != inner {
			t.Fatalf("root = %q (should be inner because go.mod boundary), want %q", root, inner)
		}
	})

	t.Run("abs dir is module root with go.mod and SETUP.md", func(t *testing.T) {
		dir := t.TempDir()
		writeEmptyFile(t, dir, "go.mod")
		writeEmptyFile(t, dir, "SETUP.md")
		sub := filepath.Join(dir, "sub")
		if err := os.MkdirAll(sub, 0755); err != nil {
			t.Fatal(err)
		}

		root, ok := ResolveRoot(sub)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != dir {
			t.Fatalf("root = %q, want %q", root, dir)
		}
	})

	t.Run("no test markers anywhere", func(t *testing.T) {
		dir := t.TempDir()

		_, ok := ResolveRoot(dir)
		if ok {
			t.Fatal("expected not ok for empty dir")
		}
	})

	t.Run("SETUP.md in leaf but no ancestor SETUP.md", func(t *testing.T) {
		dir := t.TempDir()
		leaf := filepath.Join(dir, "leaf")
		writeEmptyFile(t, leaf, "SETUP.md")

		root, ok := ResolveRoot(leaf)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != leaf {
			t.Fatalf("root = %q, want %q (dir itself has SETUP.md, no ancestor)", root, leaf)
		}
	})

	t.Run("SETUP.md at dir itself with ancestor also having SETUP.md", func(t *testing.T) {
		dir := t.TempDir()
		writeEmptyFile(t, dir, "SETUP.md")
		parent := filepath.Dir(dir)
		writeEmptyFile(t, parent, "SETUP.md")

		root, ok := ResolveRoot(dir)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != parent {
			t.Fatalf("root = %q, want %q (farthest ancestor)", root, parent)
		}
	})

	t.Run("non-existent dir", func(t *testing.T) {
		_, ok := ResolveRoot("/nonexistent/path/12345")
		if ok {
			t.Fatal("expected not ok for missing dir")
		}
	})

	t.Run("DOCTEST.md wins over SETUP.md", func(t *testing.T) {
		dir := t.TempDir()
		writeEmptyFile(t, dir, "DOCTEST.md")
		writeEmptyFile(t, filepath.Dir(dir), "SETUP.md")

		root, ok := ResolveRoot(dir)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != dir {
			t.Fatalf("root = %q, want %q (DOCTEST.md at dir takes priority over ancestor SETUP.md)", root, dir)
		}
	})

	t.Run("target in testdata does not use grandparent DOCTEST.md", func(t *testing.T) {
		dir := t.TempDir()

		rootDir := filepath.Join(dir, "tests")
		writeEmptyFile(t, rootDir, "DOCTEST.md")
		writeEmptyFile(t, rootDir, "SETUP.md")

		testdataDir := filepath.Join(rootDir, "build", "testdata")
		fixtureDir := filepath.Join(testdataDir, "helper-shadow")
		writeEmptyFile(t, fixtureDir, "SETUP.md")
		childDir := filepath.Join(fixtureDir, "child")
		if err := os.MkdirAll(childDir, 0755); err != nil {
			t.Fatal(err)
		}

		root, ok := ResolveRoot(fixtureDir)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != fixtureDir {
			t.Fatalf("root = %q, want %q (testdata boundary — should use fixture dir with its own SETUP.md, not grandparent with DOCTEST.md)", root, fixtureDir)
		}
	})

	t.Run("target in testdata child resolves to fixture root", func(t *testing.T) {
		dir := t.TempDir()

		rootDir := filepath.Join(dir, "tests")
		writeEmptyFile(t, rootDir, "DOCTEST.md")
		writeEmptyFile(t, rootDir, "SETUP.md")

		testdataDir := filepath.Join(rootDir, "build", "testdata")
		fixtureDir := filepath.Join(testdataDir, "helper-shadow")
		writeEmptyFile(t, fixtureDir, "SETUP.md")
		childDir := filepath.Join(fixtureDir, "child")
		writeEmptyFile(t, childDir, "SETUP.md")

		root, ok := ResolveRoot(childDir)
		if !ok {
			t.Fatal("expected ok")
		}
		if root != fixtureDir {
			t.Fatalf("root = %q, want %q (farthest ancestor with SETUP.md within testdata boundary)", root, fixtureDir)
		}
	})
}

func TestFindDOCTestDirsFromSubdirs_Basic(t *testing.T) {
	dir := t.TempDir()
	exec.Command("git", "-C", dir, "init").Run()

	createModule := func(name string) {
		mdir := filepath.Join(dir, name)
		writeFile(t, mdir, "go.mod", "module "+name+"\ngo 1.21\n")
		writeFile(t, mdir, "DOCTEST.md", "# "+name+"\n")
		writeFile(t, mdir, "SETUP.md", doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { return nil }\nfunc Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"))
		leafDir := filepath.Join(mdir, "simple")
		writeFile(t, leafDir, "SETUP.md", doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { return nil }\n"))
		writeFile(t, leafDir, "ASSERT.md", doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"))
	}
	createModule("mod_a")
	createModule("mod_b")
	createModule("ign_a")
	createModule("ign_b")
	writeFile(t, dir, ".gitignore", "ign_a/\nign_b/\n")

	result, err := FindDOCTestDirsFromSubdirs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 dirs, got %d: %v", len(result), result)
	}
	for _, name := range []string{"mod_a", "mod_b"} {
		found := false
		for _, r := range result {
			if filepath.Base(r) == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("%s not found in results", name)
		}
	}
	for _, r := range result {
		if filepath.Base(r) == "ign_a" || filepath.Base(r) == "ign_b" {
			t.Errorf("gitignored %s found in results", filepath.Base(r))
		}
	}
}

func TestFindDOCTestDirsFromSubdirs_DeepModule(t *testing.T) {
	dir := t.TempDir()
	mdir := filepath.Join(dir, "a", "b", "c")
	if err := os.MkdirAll(mdir, 0755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, mdir, "go.mod", "module c\ngo 1.21\n")
	writeFile(t, mdir, "DOCTEST.md", "# c\n")
	writeFile(t, mdir, "SETUP.md", doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { return nil }\nfunc Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"))
	leafDir := filepath.Join(mdir, "simple")
	writeFile(t, leafDir, "SETUP.md", doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { return nil }\n"))
	writeFile(t, leafDir, "ASSERT.md", doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"))

	result, err := FindDOCTestDirsFromSubdirs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 dir, got %d: %v", len(result), result)
	}
}

func TestFindDOCTestDirsFromSubdirs_SiblingModules(t *testing.T) {
	dir := t.TempDir()

	createModule := func(name string) {
		mdir := filepath.Join(dir, "deep", name)
		writeFile(t, mdir, "go.mod", "module "+name+"\ngo 1.21\n")
		writeFile(t, mdir, "DOCTEST.md", "# "+name+"\n")
		writeFile(t, mdir, "SETUP.md", doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { return nil }\nfunc Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"))
		leafDir := filepath.Join(mdir, "simple")
		writeFile(t, leafDir, "SETUP.md", doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { return nil }\n"))
		writeFile(t, leafDir, "ASSERT.md", doctestGoBlock("import \"testing\"\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {}\n"))
	}
	createModule("pkg1")
	createModule("pkg2")

	result, err := FindDOCTestDirsFromSubdirs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 dirs, got %d: %v", len(result), result)
	}
}

func TestFindDOCTestDirsFromSubdirs_Empty(t *testing.T) {
	dir := t.TempDir()
	result, err := FindDOCTestDirsFromSubdirs(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Fatalf("expected 0 dirs, got %d: %v", len(result), result)
	}
}

func TestFindDOCTestDirsWithBase_Filter(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module testproj\ngo 1.21\n")

	createDOCTest := func(name string) {
		mdir := filepath.Join(dir, name)
		writeFile(t, mdir, "DOCTEST.md", "# "+name+"\n")
		writeFile(t, mdir, "SETUP.md", doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { return nil }\nfunc Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"))
	}
	createDOCTest("alpha")
	createDOCTest("beta")

	result, err := FindDOCTestDirsWithBase(dir, filepath.Join(dir, "alpha"))
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 dir, got %d: %v", len(result), result)
	}
	if filepath.Base(result[0]) != "alpha" {
		t.Fatalf("expected alpha, got %s", filepath.Base(result[0]))
	}
}

func TestFindDOCTestDirsWithBase_All(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module testproj\ngo 1.21\n")

	createDOCTest := func(name string) {
		mdir := filepath.Join(dir, name)
		writeFile(t, mdir, "DOCTEST.md", "# "+name+"\n")
		writeFile(t, mdir, "SETUP.md", doctestGoBlock("import \"testing\"\nfunc Setup(t *testing.T, req *Request) error { return nil }\nfunc Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }\n"))
	}
	createDOCTest("alpha")
	createDOCTest("beta")

	result, err := FindDOCTestDirsWithBase(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 dirs, got %d: %v", len(result), result)
	}
}

func TestIsDotDotDotPattern(t *testing.T) {
	cases := []struct {
		arg      string
		expected bool
	}{
		{"./...", true},
		{"./alpha/...", true},
		{"./group/subgroup/...", true},
		{"alpha/...", true},
		{"beta/...", true},
		{"...", false},
		{"/...", false},
		{".../", false},
		{"./alpha", false},
	}
	for _, tc := range cases {
		if got := isDotDotDotPattern(tc.arg); got != tc.expected {
			t.Errorf("isDotDotDotPattern(%q) = %v, want %v", tc.arg, got, tc.expected)
		}
	}
}

func TestExtractBasePath(t *testing.T) {
	cases := []struct {
		arg      string
		expected string
	}{
		{"./...", "."},
		{"./alpha/...", "alpha"},
		{"./group/subgroup/...", "group/subgroup"},
		{"alpha/...", "alpha"},
		{"beta/...", "beta"},
	}
	for _, tc := range cases {
		if got := extractBasePath(tc.arg); got != tc.expected {
			t.Errorf("extractBasePath(%q) = %q, want %q", tc.arg, got, tc.expected)
		}
	}
}

func doctestGoBlock(code string) string {
	return "```go\n" + code + "\n```\n"
}
