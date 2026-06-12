package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xhd2015/doctest/libdoc/agent"
	"github.com/xhd2015/doctest/libdoc/implementer"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/libdoc/spec"
	"github.com/xhd2015/doctest/libdoc/validate"
)

const usage = `Usage: doctest <command> [options]

Commands:
  agent generate <idea> [-d|--dir <target-dir>] [--agent-runner RUNNER]
  agent fill-code <target-dir>
   agent implement <prompt> [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH]
  validate <dir>
  build <dir>
  test <dir>
   skill --list
   skill doc-spec show|install
   skill code-spec show|install
   skill tdd show|install
   skill implementer show|install

Run doctest <command> --help for command-specific options.
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
       doctest skill implementer show|install
`

const buildUsage = `Usage: doctest build [-v|--verbose] [--rm] [--gen-dir DIR] [-count=N] <dir>

Validate generated snippets compile without executing behavior,,allow ./... patterns like go build.

Options:
  -v, --verbose     Show generated files and build command output
  --rm              Remove the temporary generated test directory
  --gen-dir DIR     Write generated Go test files to DIR
  -count=N          Forward Go test count option to generated build
  -h, --help        Show help

Examples:
  doctest build -v ./
  doctest build -v ./...
  doctest build -v ./sub-module/...
`

const testUsage = `Usage: doctest test [-v|--verbose] [--rm] [-count=N] <dir>

Run executable Go snippets from a doc-style test directory, allow ./... patterns like go test.

Options:
  -v, --verbose     Show generated test names and runner output
  --rm              Remove the temporary generated test directory
  -count=N          Forward Go test count option to generated test binary
  -h, --help        Show help

Examples:
  doctest test -v ./
  doctest test -v ./...
  doctest test -v ./sub-module/...
`

const agentImplementUsage = `Usage: doctest agent implement [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH] <prompt>

Spawn a sub-agent to implement code that makes doctests pass.
Blocks until the sub-agent completes or yields questions via
yield-pending-questions.

Options:
  --agent-runner RUNNER   opencode, codex, or fake-codex (default: opencode)
  --mock-config PATH      mock config JSON for fake-codex
  --requirement PATH      read requirement from file (useful for long prompts
                          or prompts with shell special characters)
  -h, --help              Show help
`

func Run(args []string) error {
	if filepath.Base(os.Args[0]) == "yield-pending-questions" {
		return implementer.HandleYieldPendingQuestions(args)
	}
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(usage)
		return nil
	}
	switch args[0] {
	case "agent":
		return runAgent(args[1:])
	case "validate":
		return runOneDir("validate", args[1:], validate.Run)
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
  implement <prompt> [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH]
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
	var idea []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			fmt.Print(agentGenerateUsage)
			return nil
		case "-d", "--dir":
			if i+1 >= len(args) {
				return fmt.Errorf("%s requires value", arg)
			}
			opts.Dir = args[i+1]
			i++
		case "--agent-runner":
			if i+1 >= len(args) {
				return fmt.Errorf("--agent-runner requires value")
			}
			opts.AgentRunner = args[i+1]
			i++
		default:
			idea = append(idea, arg)
		}
	}
	opts.Idea = strings.Join(idea, " ")
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
	opts := implementer.Options{AgentRunner: "opencode"}
	var promptParts []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			fmt.Print(agentImplementUsage)
			return nil
		case "--agent-runner":
			if i+1 >= len(args) {
				return fmt.Errorf("--agent-runner requires value")
			}
			opts.AgentRunner = args[i+1]
			i++
		case "--session-id":
			if i+1 >= len(args) {
				return fmt.Errorf("--session-id requires value")
			}
			opts.SessionID = args[i+1]
			i++
		case "--mock-config":
			if i+1 >= len(args) {
				return fmt.Errorf("--mock-config requires value")
			}
			opts.MockConfig = args[i+1]
			i++
		case "--requirement":
			if i+1 >= len(args) {
				return fmt.Errorf("--requirement requires value")
			}
			opts.Requirement = args[i+1]
			i++
		default:
			promptParts = append(promptParts, arg)
		}
	}
	opts.Prompt = strings.Join(promptParts, " ")
	return implementer.Run(opts)
}

func runOneDir(name string, args []string, fn func(string) error) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Printf("Usage: doctest %s <dir>\n", name)
		return nil
	}
	if len(args) != 1 {
		return fmt.Errorf("%s requires <dir>", name)
	}
	return fn(args[0])
}

func runRunner(args []string, usage string, fn func([]string) error) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Print(usage)
		return nil
	}
	err := fn(args)
	if errors.Is(err, runner.ErrNoTestsFound) {
		fmt.Fprintln(os.Stderr, "no tests found")
		return nil
	}
	return err
}

func runSkill(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(skillUsage)
		return nil
	}
	if args[0] == "--list" {
		fmt.Println("doc-spec")
		fmt.Println("code-spec")
		fmt.Println("tdd")
		fmt.Println("implementer")
		return nil
	}
	if len(args) < 2 {
		return fmt.Errorf("skill requires doc-spec, code-spec, tdd, or implementer plus show or install")
	}
	switch args[1] {
	case "show":
		content, err := spec.Content(args[0])
		if err != nil {
			return err
		}
		fmt.Print(content)
		return nil
	case "install":
		return spec.Install(args[0], args[2:])
	default:
		return fmt.Errorf("unknown skill action: %s", args[1])
	}
}

func Main() {
	if err := Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
