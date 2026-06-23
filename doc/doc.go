package doc

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/xhd2015/doctest/doc/snippets"
	"github.com/xhd2015/doctest/version"
)

//go:embed DOC_STYLE_TEST_SPECIFICATION.md
var docStyleTestSpecification string

//go:embed DOC_STYLE_TEST_CODE_SPECIFICATION.md
var docStyleTestCodeSpecification string

//go:embed DOC_STYLE_TEST_BASED_TDD.md
var docStyleTestBasedTDD string

//go:embed DOC_STYLE_TEST_BASED_TDD_LITE.md
var docStyleTestBasedTDDLITE string

func Content(fileName string) (string, error) {
	var content string
	switch fileName {
	case "DOC_STYLE_TEST_SPECIFICATION.md":
		content = docStyleTestSpecification
	case "DOC_STYLE_TEST_CODE_SPECIFICATION.md":
		content = docStyleTestCodeSpecification
	case "DOC_STYLE_TEST_BASED_TDD.md":
		content = docStyleTestBasedTDD
	case "DOC_STYLE_TEST_BASED_TDD_LITE.md":
		content = docStyleTestBasedTDDLITE
	default:
		return "", fmt.Errorf("unknown file: %s", fileName)
	}
	content = strings.ReplaceAll(content, "__DOCTEST_SPEC__", snippets.ResolvedDocTestSpec())
	content = strings.ReplaceAll(content, "__DOCTEST_DESIGN_SPEC__", snippets.DocTestDesignSpec)
	content = strings.ReplaceAll(content, "__DOCTEST_VERSION__", version.Version())
	return content, nil
}
