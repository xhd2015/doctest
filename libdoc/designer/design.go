package designer

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/doctest/doc/snippets"

	"github.com/xhd2015/agent-pro/agent/subagent"
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
	if opts.SessionBase == "" {
		homeDir, _ := os.UserHomeDir()
		opts.SessionBase = filepath.Join(homeDir, ".doctest")
	}
	return subagent.Run(context.Background(), subagent.Config{
		RoleName:         "designer",
		Cmd:              "design",
		PromptContent:    PromptContent(),
		SessionEnvVar:    "DOCTEST_AGENT_DESIGNER_SESSION_ID",
		SessionMetaField: "doctest_agent_designer_session_id",
		DebugSessionEnv:  "DOCTEST_DEBUG_SESSION_HOME",
		AgentRunnerEnv:   "DOCTEST_SUBAGENT_AGENT_RUNNER",
		ModelEnv:         "DOCTEST_SUBAGENT_MODEL",
	}, opts)
}
