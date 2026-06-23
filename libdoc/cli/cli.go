package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	lessflags "github.com/xhd2015/less-flags"

	"github.com/xhd2015/doctest/libdoc/agent"
	"github.com/xhd2015/doctest/libdoc/designer"
	"github.com/xhd2015/doctest/libdoc/implementer"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/libdoc/spec"
	"github.com/xhd2015/agent-pro/agent/subagent"
)

const usage = `Usage: doctest <command> [options]

Commands:
  vet [-v|--verbose] <dir...>
  build <dir>
  test <dir>

Agents:
  agent generate <idea> [-d|--dir <target-dir>] [--agent-runner RUNNER]
  agent fill-code <target-dir>
  agent implement <prompt> [--session-id ID] [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH] [--timeout DURATION] [--trace] [--status] [--list-sessions]
  agent design <prompt> [--session-id ID] [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH] [--timeout DURATION] [--trace] [--status] [--list-sessions]
  agent with --agent-runner=RUNNER [--model=MODEL] <prog> [args...]

Skills:
  skill --list
  skill doc-spec show|install
  skill code-spec show|install
  skill tdd show|install
  skill implementer show|install
  skill designer show|install

Run doctest <command> --help for command-specific options.

Examples:
  doc test -v ./...
`

const agentGenerateUsage = `Usage: doctest agent generate <idea> [-d|--dir <target-dir>] [--agent-runner RUNNER]

Generate a doc-style test tree from an idea.

Options:
  -d, --dir <target-dir>       Directory to write
  --agent-runner RUNNER        opencode, codex, or fake-codex
  -h, --help                   Show help
`

const skillUsage = `Usage: doctest skill --list
       doctest skill doc-spec show|install
       doctest skill code-spec show|install
       doctest skill tdd show|install
       doctest skill tdd-lite show|install
       doctest skill implementer show|install
       doctest skill designer show|install
`

const vetUsage = `Usage: doctest vet [-v|--verbose] [--changed] <dir...>

Validate doc-test tree structure and anti-patterns, allow ./... patterns like go build.

Options:
  -v, --verbose     Show directories and files being validated
  --changed         Only validate doctest files affected by git working-tree changes
  -h, --help        Show help

Examples:
  doctest vet -v ./
  doctest vet -v ./...
  doctest vet -v ./sub-module/...
`

const buildUsage = `Usage: doctest build [-v|--verbose] [--rm] [--gen-dir DIR] [-count=N] [--changed] <dir>

Validate generated snippets compile without executing behavior,,allow ./... patterns like go build.

Options:
  -v, --verbose     Show generated files and build command output
  --rm              Remove the temporary generated test directory
  --gen-dir DIR     Write generated Go test files to DIR
  -count=N          Forward Go test count option to generated build
  --changed         Only build doctest leaves affected by git working-tree changes
  -h, --help        Show help

Examples:
  doctest build -v ./
  doctest build -v ./...
  doctest build -v ./sub-module/...
`

const testUsage = `Usage: doctest test [-v|--verbose] [--rm] [--gen-dir DIR] [-count=N] [--timeout DURATION] [--color] [--no-color] [--changed] <dir>

Run executable Go snippets from a doc-style test directory, allow ./... patterns like go test.

Options:
  -v, --verbose     Show generated test names and runner output
  --rm              Remove the temporary generated test directory
  --gen-dir DIR     Write generated Go test files to DIR
  -count=N          Forward Go test count option to generated test binary
  --timeout DURATION
                    Forward Go test timeout to generated test binary
                    (e.g. 30s, 5m, 1h); omitted uses go test default (10m)
  --color           Force ANSI color in non-verbose progress output
  --no-color        Disable ANSI color in non-verbose progress output
  --changed         Only run doctest leaves affected by git working-tree changes
  -h, --help        Show help

Examples:
  doctest test -v ./
  doctest test -v ./...
  doctest test -v ./sub-module/...
`

const agentImplementUsage = `Usage: doctest agent implement [--session-id ID] [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH] [--timeout DURATION] [--trace] [--status] [--list-sessions] <prompt>

Spawn a sub-agent to implement code that makes doctests pass.
Blocks until the sub-agent completes or yields questions via
yield-pending-questions.

With --trace, follow and print the events of an existing session
instead of spawning a sub-agent.

With --status, display the current status of an existing session
(requires --session-id).

With --list-sessions, list all sessions from the last 7 days.

Options:
  --session-id ID         session to use or resume; for --trace, the session to follow
  --agent-runner RUNNER   opencode, codex, or fake-codex (default: opencode)
  --mock-config PATH      mock config JSON for fake-codex
  --requirement PATH      read requirement from file (useful for long prompts
                          or prompts with shell special characters)
  --timeout DURATION      max duration for the sub-agent (default: 1h);
                          accepts Go-style durations like 30s, 5m, 1h, 1h30m
  --trace                 follow and print events from an existing session
  --status                display session status (requires --session-id)
  --list-sessions         list all sessions from the last 7 days
  -h, --help              Show help
`

const agentDesignerUsage = `Usage: doctest agent design [--session-id ID] [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH] [--timeout DURATION] [--trace] [--status] [--list-sessions] <prompt>

Spawn a sub-agent to design doctest trees for new features or update
existing test suites. The sub-agent analyzes requirements, breaks down
parameters into mutually exclusive cases, and writes comprehensive
doctest files.

With --trace, follow and print the events of an existing session
instead of spawning a sub-agent.

With --status, display the current status of an existing session
(requires --session-id).

With --list-sessions, list all sessions from the last 7 days.

Options:
  --session-id ID         session to use or resume; for --trace, the session to follow
  --agent-runner RUNNER   opencode, codex, or fake-codex (default: opencode)
  --mock-config PATH      mock config JSON for fake-codex
  --requirement PATH      read requirement from file (useful for long prompts
                          or prompts with shell special characters)
  --timeout DURATION      max duration for the sub-agent (default: 1h);
                          accepts Go-style durations like 30s, 5m, 1h, 1h30m
  --trace                 follow and print events from an existing session
  --status                display session status (requires --session-id)
  --list-sessions         list all sessions from the last 7 days
  -h, --help              Show help
`

func Run(args []string) error {
	if filepath.Base(os.Args[0]) == "yield-pending-questions" {
		return subagent.HandleYieldPendingQuestions(args)
	}
	if filepath.Base(os.Args[0]) == "report-progress" {
		return subagent.HandleReportProgress(args)
	}
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(usage)
		return nil
	}
	switch args[0] {
	case "agent":
		return runAgent(args[1:])
	case "vet":
		return runRunner(args[1:], vetUsage, runner.VetArgs)
	case "build":
		return runRunner(args[1:], buildUsage, runner.BuildArgs)
	case "test":
		return runRunner(args[1:], testUsage, runner.Test)
	case "skill":
		return runSkill(args[1:])
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runAgent(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(`Usage: doctest agent <command> [options]

Commands:
  generate <idea> [-d|--dir <target-dir>]
  fill-code <target-dir>
  implement <prompt> [--session-id ID] [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH] [--trace] [--status] [--list-sessions]
  design <prompt> [--session-id ID] [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH] [--trace] [--status] [--list-sessions]
  with --agent-runner=RUNNER [--model=MODEL] <prog> [args...]
`)
		return nil
	}
	switch args[0] {
	case "generate":
		return runAgentGenerate(args[1:])
	case "fill-code":
		if len(args) > 1 && (args[1] == "-h" || args[1] == "--help") {
			fmt.Print("Usage: doctest agent fill-code <target-dir>\n")
			return nil
		}
		if len(args) != 2 {
			return fmt.Errorf("agent fill-code requires <target-dir>")
		}
		return agent.FillCode(args[1])
	case "implement":
		return runAgentImplement(args[1:])
	case "design":
		return runAgentDesign(args[1:])
	case "with":
		return runAgentWith(args[1:])
	default:
		return fmt.Errorf("unknown agent command: %s", args[0])
	}
}

func runAgentGenerate(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(agentGenerateUsage)
		return nil
	}
	opts := agent.GenerateOptions{AgentRunner: "opencode"}
	remainArgs, err := lessflags.String("-d,--dir", &opts.Dir).
		String("--agent-runner", &opts.AgentRunner).
		Help("-h,--help", agentGenerateUsage).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return err
	}
	opts.Idea = strings.Join(remainArgs, " ")
	if strings.TrimSpace(opts.Idea) == "" {
		return fmt.Errorf("agent generate requires <idea>")
	}
	return agent.Generate(opts)
}

func runAgentImplement(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(agentImplementUsage)
		return nil
	}
	opts := implementer.Options{}
	var timeoutStr string
	remainArgs, err := lessflags.String("--session-id", &opts.SessionID).
		String("--agent-runner", &opts.AgentRunner).
		String("--mock-config", &opts.MockConfig).
		String("--requirement", &opts.Requirement).
		String("--timeout", &timeoutStr).
		Bool("--trace", &opts.CatchUp).
		Bool("--status", &opts.Status).
		Bool("--list-sessions", &opts.ListSessions).
		Help("-h,--help", agentImplementUsage).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return err
	}
	if timeoutStr != "" {
		d, err := subagent.ParseTimeoutDuration(timeoutStr)
		if err != nil {
			return fmt.Errorf("--timeout: %w", err)
		}
		opts.Timeout = d
	}
	opts.Prompt = strings.Join(remainArgs, " ")
	if opts.Prompt == "" {
		var err error
		opts.Prompt, err = readStdinIfPresent()
		if err != nil {
			return err
		}
	}
	return implementer.Run(opts)
}

func runAgentDesign(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(agentDesignerUsage)
		return nil
	}
	opts := designer.Options{}
	var timeoutStr string
	remainArgs, err := lessflags.String("--session-id", &opts.SessionID).
		String("--agent-runner", &opts.AgentRunner).
		String("--mock-config", &opts.MockConfig).
		String("--requirement", &opts.Requirement).
		String("--timeout", &timeoutStr).
		Bool("--trace", &opts.CatchUp).
		Bool("--status", &opts.Status).
		Bool("--list-sessions", &opts.ListSessions).
		Help("-h,--help", agentDesignerUsage).
		HelpNoExit().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return err
	}
	if timeoutStr != "" {
		d, err := subagent.ParseTimeoutDuration(timeoutStr)
		if err != nil {
			return fmt.Errorf("--timeout: %w", err)
		}
		opts.Timeout = d
	}
	opts.Prompt = strings.Join(remainArgs, " ")
	if opts.Prompt == "" {
		var err error
		opts.Prompt, err = readStdinIfPresent()
		if err != nil {
			return err
		}
	}
	return designer.Run(opts)
}

func runAgentWith(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(`Usage: doctest agent with --agent-runner=RUNNER [--model=MODEL] <prog> [args...]

Execute a program with DOCTEST_SUBAGENT_AGENT_RUNNER and optionally DOCTEST_SUBAGENT_MODEL set in its environment.

Options:
  --agent-runner=RUNNER   opencode, pi, crush, or codex (required)
  --model=MODEL           Model override for the agent runner
  -h, --help              Show help
`)
		return nil
	}

	var agentRunner string
	var model string
	remainArgs, err := lessflags.String("--agent-runner", &agentRunner).
		String("--model", &model).
		StopOnFirstArg().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return err
	}

	if agentRunner == "" {
		return fmt.Errorf("--agent-runner requires a value")
	}
	if len(remainArgs) == 0 {
		return fmt.Errorf("agent with requires <prog>")
	}

	prog := remainArgs[0]
	progArgs := remainArgs[1:]

	cmd := exec.Command(prog, progArgs...)
	cmd.Env = append(os.Environ(),
		"DOCTEST_SUBAGENT_AGENT_RUNNER="+agentRunner,
	)
	if model != "" {
		cmd.Env = append(cmd.Env, "DOCTEST_SUBAGENT_MODEL="+model)
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr
		}
		return err
	}
	return nil
}

func runRunner(args []string, usage string, fn func([]string) error) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(usage)
		return nil
	}
	err := fn(args)
	if errors.Is(err, runner.ErrNoTestsFound) {
		fmt.Fprintln(os.Stderr, "no tests")
		return nil
	}
	return err
}

func runSkill(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(skillUsage)
		return nil
	}
	var listFlag bool
	remainArgs, err := lessflags.Bool("--list", &listFlag).
		Help("-h,--help", skillUsage).
		HelpNoExit().
		StopOnFirstArg().
		Parse(args)
	if err != nil {
		if errors.Is(err, lessflags.ErrHelp) {
			return nil
		}
		return err
	}
	if listFlag {
		fmt.Println("doc-spec")
		fmt.Println("code-spec")
		fmt.Println("tdd")
		fmt.Println("tdd-lite")
		fmt.Println("implementer")
		fmt.Println("designer")
		return nil
	}
	if len(remainArgs) < 2 {
		return fmt.Errorf("skill requires doc-spec, code-spec, tdd, tdd-lite, implementer, or designer plus show or install")
	}
	switch remainArgs[1] {
	case "show":
		content, err := spec.Content(remainArgs[0])
		if err != nil {
			return err
		}
		fmt.Print(content)
		return nil
	case "install":
		return spec.Install(remainArgs[0], remainArgs[2:])
	default:
		return fmt.Errorf("unknown skill action: %s", remainArgs[1])
	}
}

func readStdinIfPresent() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Main() {
	if err := Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
