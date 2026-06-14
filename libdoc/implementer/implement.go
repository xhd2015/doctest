package implementer

import (
	_ "embed"

	"github.com/xhd2015/doctest/libdoc/subagent"
)

//go:embed PROMPT.md
var promptContent string

func PromptContent() string {
	return promptContent
}

type Options = subagent.Options

func Run(opts Options) error {
	return subagent.Run(subagent.Config{
		RoleName:      "implementer",
		Cmd:           "implement",
		PromptContent: promptContent,
	}, opts)
}
