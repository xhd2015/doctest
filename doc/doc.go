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

//go:embed DOCTEST_TDD.md
var doctestTDD string

//go:embed DOCTEST_TDD_CLI_AGENT.md
var doctestTDDCLIAgent string

//go:embed DOCTEST_TDD_LITE.md
var doctestTDDLITE string

//go:embed DOCTEST_REPRODUCE.md
var doctestReproduce string

//go:embed DOCTEST_REVIEW.md
var doctestReview string

//go:embed DOCTEST_REVIEW_PERF.md
var doctestReviewPerf string

//go:embed DOCTEST_ANALYSE_PERF.md
var doctestAnalysePerf string

//go:embed DOCTEST_OUTPUT_ASSERT.md
var doctestOutputAssert string

//go:embed DOCTEST_DESIGN_PRINCIPLE.md
var doctestDesignPrinciple string

//go:embed DOCTEST_LINT.md
var doctestLint string

//go:embed DOCTEST_MIGRATE.md
var doctestMigrate string

func Content(fileName string) (string, error) {
	var content string
	switch fileName {
	case "DOC_STYLE_TEST_SPECIFICATION.md":
		content = docStyleTestSpecification
	case "DOC_STYLE_TEST_CODE_SPECIFICATION.md":
		content = docStyleTestCodeSpecification
	case "DOCTEST_TDD.md":
		content = doctestTDD
	case "DOCTEST_TDD_CLI_AGENT.md":
		content = doctestTDDCLIAgent
	case "DOCTEST_TDD_LITE.md":
		content = doctestTDDLITE
	case "DOCTEST_REPRODUCE.md":
		content = doctestReproduce
	case "DOCTEST_REVIEW.md":
		content = doctestReview
	case "DOCTEST_REVIEW_PERF.md":
		content = doctestReviewPerf
	case "DOCTEST_ANALYSE_PERF.md":
		content = doctestAnalysePerf
	case "DOCTEST_OUTPUT_ASSERT.md":
		content = doctestOutputAssert
	case "DOCTEST_DESIGN_PRINCIPLE.md":
		content = doctestDesignPrinciple
	case "DOCTEST_LINT.md":
		content = doctestLint
	case "DOCTEST_MIGRATE.md":
		content = doctestMigrate
	default:
		return "", fmt.Errorf("unknown file: %s", fileName)
	}
	content = strings.ReplaceAll(content, "__DOCTEST_SPEC__", snippets.ResolvedDocTestSpec())
	content = strings.ReplaceAll(content, "__DOCTEST_DESIGN_SPEC__", snippets.DocTestDesignSpec)
	content = strings.ReplaceAll(content, "__DOCTEST_VERSION__", version.Version())
	return content, nil
}
