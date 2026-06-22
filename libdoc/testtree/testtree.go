package testtree

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/version"
)

const MinimalDSN = `## DSN (Domain Specific Notion)

### Participants
- **system** — under test.

### Behaviors
- **run** — executes the scenario.
`

func WriteFile(t *testing.T, root, rel, content string) {
	t.Helper()
	path := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func MinimalDOCTEST(goBody string) string {
	return fmt.Sprintf("# Tests\n\n## Version\n%s\n\n%s\n\n```go\n%s\n```\n", version.Version(), MinimalDSN, goBody)
}

func MinimalRunGo() string {
	return `import "testing"

type Request struct{}
type Response struct{}

func Run(t *testing.T, req *Request) (*Response, error) {
	return &Response{}, nil
}`
}

type LeafSpec struct {
	Name     string
	Steps    string
	Expected string
	SetupGo  string
	AssertGo string
}

func WriteMinimalRunnableTree(t *testing.T, root string, leaves []LeafSpec) {
	t.Helper()
	WriteFile(t, root, "DOCTEST.md", MinimalDOCTEST(MinimalRunGo()))
	for _, leaf := range leaves {
		setup := leaf.SetupGo
		if setup == "" {
			setup = `import "testing"

func Setup(t *testing.T, req *Request) error { _ = req; return nil }`
		}
		steps := leaf.Steps
		if steps == "" {
			steps = "leaf setup"
		}
		expected := leaf.Expected
		if expected == "" {
			expected = "passes"
		}
		WriteFile(t, root, leaf.Name+"/SETUP.md", fmt.Sprintf("## Steps\n1. %s\n\n```go\n%s\n```\n", steps, setup))
		assert := leaf.AssertGo
		if assert == "" {
			assert = `func Assert(t *testing.T, req *Request, resp *Response, err error) {}`
		}
		WriteFile(t, root, leaf.Name+"/ASSERT.md", fmt.Sprintf("## Expected\n- %s\n\n```go\n%s\n```\n", expected, assert))
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b [20]byte
	pos := len(b)
	for i > 0 {
		pos--
		b[pos] = byte('0' + i%10)
		i /= 10
	}
	return string(b[pos:])
}

func WritePassFailTree(t *testing.T, root string, passCount, failCount int) {
	t.Helper()
	var leaves []LeafSpec
	for i := 0; i < passCount; i++ {
		leaves = append(leaves, LeafSpec{Name: "a_pass_" + itoa(i)})
	}
	for i := 0; i < failCount; i++ {
		name := "z_fail_" + itoa(i)
		leaves = append(leaves, LeafSpec{
			Name: name,
			AssertGo: `import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) { t.Fatal("forced failure") }`,
		})
	}
	if len(leaves) == 0 {
		leaves = append(leaves, LeafSpec{Name: "a_pass_0"})
	}
	WriteMinimalRunnableTree(t, root, leaves)
}

func VetDOCTEST() string {
	return MinimalDOCTEST(MinimalRunGo())
}

func VetDOCTESTWithoutRun() string {
	return MinimalDOCTEST(`import "testing"

type Request struct{}
type Response struct{}`)
}