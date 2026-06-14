package designer

import (
	_ "embed"
	"strings"

	"github.com/xhd2015/doctest/doc/snippets"

	"github.com/xhd2015/doctest/libdoc/subagent"
)

//go:embed PROMPT.md
var promptContent string

func PromptContent() string {
	s := strings.Replace(promptContent, "__DOCTEST_SPEC__", snippets.DocTestSpec, 1)
	s = strings.Replace(s, "__DOCTEST_DESIGN_SPEC__", snippets.DocTestDesignSpec, 1)
	return s
}

type Options = subagent.Options

func Run(opts Options) error {
	return subagent.Run(subagent.Config{
		RoleName:      "designer",
		Cmd:           "design",
		PromptContent: promptContent,
	}, opts)
}
