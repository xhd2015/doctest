package implementer

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/doctest/version"

	"github.com/xhd2015/agent-pro/agent/subagent"
)

//go:embed PROMPT.md
var promptContent string

func PromptContent() string {
	return strings.ReplaceAll(promptContent, "__DOCTEST_VERSION__", version.Version())
}

type Options = subagent.Options

func Run(opts Options) error {
	if opts.SessionBase == "" {
		homeDir, _ := os.UserHomeDir()
		opts.SessionBase = filepath.Join(homeDir, ".doctest")
	}
	return subagent.Run(context.Background(), subagent.Config{
		RoleName:         "implementer",
		Cmd:              "implement",
		PromptContent:    PromptContent(),
		SessionEnvVar:    "DOCTEST_AGENT_IMPLEMENTER_SESSION_ID",
		SessionMetaField: "doctest_agent_implementer_session_id",
		DebugSessionEnv:  "DOCTEST_DEBUG_SESSION_HOME",
		AgentRunnerEnv:   "DOCTEST_SUBAGENT_AGENT_RUNNER",
		ModelEnv:         "DOCTEST_SUBAGENT_MODEL",
	}, opts)
}
