package implementer

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/xhd2015/agent-pro/agent/subagent"
)

//go:embed PROMPT.md
var promptContent string

func PromptContent() string {
	return promptContent
}

type Options = subagent.Options

func Run(opts Options) error {
	if opts.SessionBase == "" {
		homeDir, _ := os.UserHomeDir()
		opts.SessionBase = filepath.Join(homeDir, ".doctest")
	}
	return subagent.Run(subagent.Config{
		RoleName:         "implementer",
		Cmd:              "implement",
		PromptContent:    promptContent,
		SessionEnvVar:    "DOCTEST_AGENT_IMPLEMENTER_SESSION_ID",
		SessionMetaField: "doctest_agent_implementer_session_id",
		DebugSessionEnv:  "DOCTEST_DEBUG_SESSION_HOME",
	}, opts)
}
