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

	"github.com/xhd2015/agent-pro/agent/subagent"
	"github.com/xhd2015/doctest/libdoc/agent"
	"github.com/xhd2015/doctest/libdoc/designer"
	"github.com/xhd2015/doctest/libdoc/edit"
	"github.com/xhd2015/doctest/libdoc/implementer"
	"github.com/xhd2015/doctest/libdoc/runner"
	"github.com/xhd2015/doctest/libdoc/spec"
	"github.com/xhd2015/skills/skill_file"
)

const usage = `Usage: doctest <command> [options]

Commands:
  vet [-v|--verbose] <dir...>
  build <dir>
  test <dir>
  edit <leaf-path> [--add-label NAME] [--add-explanation TEXT]

Agents:
  agent generate <idea> [-d|--dir <target-dir>] [--agent-runner RUNNER]
  agent fill-code <target-dir>
  agent implement <prompt> [--session-id ID] [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH] [--timeout DURATION] [--trace] [--status] [--list-sessions]
  agent design <prompt> [--session-id ID] [--agent-runner RUNNER] [--mock-config PATH] [--requirement PATH] [--timeout DURATION] [--trace] [--status] [--list-sessions]
  agent with --agent-runner=RUNNER [--model=MODEL] <prog> [args...]

Skills:
  skills update [--global] [--cursor] [--codex] [--dry-run]
  skill --list
  skill --show <name>
  skill <name> --show
  skill --install <name> [OPTIONS]
  skill <name> --install [OPTIONS]

Metrics:
  metrics path|last|top|phases|summary|show|prune

Run doctest <command> --help for command-specific options.
Run doctest skill --help and doctest skill --install <name> --help for skill flags.

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

const skillsUsage = `Usage: doctest skills update [OPTIONS] [<dir>]
       doctest skills --help

Update already-installed doctest skills (SKILL.md must exist at target paths).

Commands:
  update    Refresh installed skills from the doctest skill registry

Options for update match doctest skill --install locations:
  --global, --cursor, --codex, --opencode, --general-agents, --dry-run

Run doctest skills update --help for update flags.
`

const skillUsage = `Usage: doctest skill --list
       doctest skill --show <name> [--header]
       doctest skill <name> --show [--header]
       doctest skill --install <name> [OPTIONS] [<dir>]
       doctest skill <name> --install [OPTIONS] [<dir>]

List, print, or install a registered doctest skill.

Registered skills:
  doc-spec, code-spec, tdd, tdd-cli-agent, tdd-lite,
  reproduce, review, review-perf, output-assert, implementer, designer

Both flag orders are valid (--show/--install before or after <name>).

Options:
  --list, -l     List registered skill names
  --show         Print skill content to stdout
  --header       With --show: print YAML frontmatter only
  --install      Install skill SKILL.md (see --install --help for targets)
  -h, --help     Show this help

Install example:
  doctest skill --install tdd --codex --opencode
  doctest skill tdd --install --global

Run doctest skill --install <name> --help for install target flags.
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

const editUsage = `Usage: doctest edit <leaf-path> [--add-label NAME] [--add-explanation TEXT]

Update YAML frontmatter on a single concrete doctest leaf ASSERT.md.

Options:
  --add-label NAME          Append a label (idempotent; warns if already set)
  --add-explanation TEXT    Set or append explanation (appends with "; " when present)
  -h, --help                Show help

Examples:
  doctest edit ./tests/feature/ui-leaf --add-label ui-automation --add-explanation "AX window test"
  doctest edit ./tests/feature/ui-leaf/ASSERT.md --add-label manual
`

const testUsage = `Usage: doctest test [-v|--verbose] [--rm] [--gen-dir DIR] [-count=N] [--timeout DURATION] [--color] [--no-color] [--changed] [--label EXPR]... [--label-all] [--metrics-on] [--cold-cache] [--experiment-ref-instead-of-inline] [--experiment-unified-package-per-doctest-tree] [-cpuprofile FILE] [-memprofile FILE] [-memprofilerate N] [-blockprofile FILE] [-blockprofilerate N] [-mutexprofile FILE] [-mutexprofilefraction N] [-trace FILE] [-outputdir DIR] [-coverprofile FILE] [-cover] <dir>

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
  --label EXPR      Run only leaves whose ASSERT.md labels match EXPR (repeatable;
                    multiple flags are OR'd). Supports &&, ||, and parentheses.
                    Unlabeled leaves are skipped when this flag is set.
  --label-all       Discovery mode: run all leaves including labeled ones (full
                    suite). Mutually exclusive with --label.
  --metrics-on      Opt in to suite metrics JSONL recording (off by default)
  --cold-cache      Reproducible cold run: wipe mapping-gen root on startup
                    (auto: $CacheHome/doctest/mapping-gen-cold; refuses --gen-dir
                    equal/under warm mapping-gen), force -count=1 when unset,
                    isolate empty GOCACHE for go test; leftover gen kept after finish
  -cpuprofile FILE  Forward CPU profile path to go test (relative paths abs-resolved)
  -memprofile FILE  Forward memory profile path to go test (relative paths abs-resolved)
  -memprofilerate N Forward memprofilerate to go test (including 0)
  -blockprofile FILE
                    Forward block profile path to go test (relative paths abs-resolved)
  -blockprofilerate N
                    Forward blockprofilerate to go test (including 0)
  -mutexprofile FILE
                    Forward mutex profile path to go test (relative paths abs-resolved)
  -mutexprofilefraction N
                    Forward mutexprofilefraction to go test (including 0)
  -trace FILE       Forward execution trace path to go test (relative paths abs-resolved)
  -outputdir DIR    Forward profile/output directory to go test (relative paths abs-resolved)
  -coverprofile FILE
                    Forward cover profile path to go test (relative paths abs-resolved)
  -cover            Enable coverage analysis (forward -cover to go test)
  --experiment-ref-instead-of-inline
                    Experimental: generate a shared root package + thin leaf tests
                    that import it (ref-instead-of-inline) instead of classic inline
  --experiment-unified-package-per-doctest-tree
                    Experimental: one go test package/binary per DOCTEST tree
                    (implies --experiment-ref-instead-of-inline; suite iterator
                    over registered leaf RunTestLeaf funcs)
  -h, --help        Show help

Examples:
  doctest test -v ./
  doctest test -v ./...
  doctest test -v ./... --label-all
  doctest test -v ./sub-module/...
  doctest test ./mod --label slow
  doctest test ./mod --label 'slow && ui-automation'
  doctest test ./mod --label slow --label heavy
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
	case "edit":
		return runEdit(args[1:])
	case "skill":
		return runSkill(args[1:])
	case "skills":
		return runSkills(args[1:])
	case "metrics":
		return runMetrics(args[1:])
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

func runSkills(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(skillsUsage)
		return nil
	}
	switch args[0] {
	case "update":
		return spec.Update(args[1:])
	default:
		return fmt.Errorf("unknown skills command: %s", args[0])
	}
}

// skill action flags (Shape 2 multi-skill host): --show / --install / --list only.
// Both orders are valid: skill --show <name> and skill <name> --show.
type skillAction string

const (
	skillActionShow    skillAction = "show"
	skillActionInstall skillAction = "install"
	skillActionList    skillAction = "list"
	skillActionHelp    skillAction = "help"
)

type parsedSkillArgs struct {
	Action skillAction
	Header bool
	Rest   []string
}

func parseSkillArgs(args []string) (parsedSkillArgs, error) {
	var (
		show    bool
		install bool
		list    bool
		header  bool
		help    bool
		rest    []string
	)
	for _, a := range args {
		switch a {
		case "--show":
			show = true
		case "--install":
			install = true
		case "--list", "-l":
			list = true
		case "--header":
			header = true
		case "-h", "--help":
			help = true
		default:
			rest = append(rest, a)
		}
	}

	// Install owns its own --help (targets, --global, etc.).
	if help && install && !show && !list {
		rest = append(rest, "--help")
		return parsedSkillArgs{Action: skillActionInstall, Rest: rest}, nil
	}
	// Skill-level help: bare --help, or --help with --show/--list.
	if help {
		return parsedSkillArgs{Action: skillActionHelp, Rest: rest}, nil
	}

	n := 0
	var action skillAction
	if show {
		n++
		action = skillActionShow
	}
	if install {
		n++
		action = skillActionInstall
	}
	if list {
		n++
		action = skillActionList
	}
	if n == 0 {
		return parsedSkillArgs{}, fmt.Errorf("expected one of --show, --install, or --list (try --help)")
	}
	if n > 1 {
		if show && install {
			return parsedSkillArgs{}, fmt.Errorf("cannot combine --show and --install")
		}
		return parsedSkillArgs{}, fmt.Errorf("expected exactly one of --show, --install, or --list (try --help)")
	}
	if header && action != skillActionShow {
		return parsedSkillArgs{}, fmt.Errorf("--header is only valid with --show")
	}
	return parsedSkillArgs{Action: action, Header: header, Rest: rest}, nil
}

func runSkill(args []string) error {
	if len(args) == 0 {
		fmt.Print(skillUsage)
		return nil
	}
	parsed, err := parseSkillArgs(args)
	if err != nil {
		return err
	}
	switch parsed.Action {
	case skillActionHelp:
		fmt.Print(skillUsage)
		return nil
	case skillActionList:
		for _, name := range []string{
			"doc-spec",
			"code-spec",
			"tdd",
			"tdd-cli-agent",
			"tdd-lite",
			"reproduce",
			"review",
			"review-perf",
			"output-assert",
			"implementer",
			"designer",
		} {
			fmt.Println(name)
		}
		return nil
	case skillActionShow:
		if len(parsed.Rest) == 0 {
			return fmt.Errorf("expected skill name for --show (try --help)")
		}
		if len(parsed.Rest) > 1 {
			return fmt.Errorf("unexpected arguments: %v", parsed.Rest[1:])
		}
		content, err := spec.Content(parsed.Rest[0])
		if err != nil {
			return err
		}
		if parsed.Header {
			out, err := skill_file.FormatHeaderWithDelimiters(content)
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		}
		fmt.Print(content)
		return nil
	case skillActionInstall:
		name, installArgs, err := splitSkillName(parsed.Rest)
		if err != nil {
			return fmt.Errorf("expected skill name for --install (try --help)")
		}
		return spec.Install(name, installArgs)
	default:
		return fmt.Errorf("unknown skill action: %s", parsed.Action)
	}
}

// splitSkillName takes the first non-flag token as the skill name; remaining
// tokens (including flags that appeared before the name) are install args.
func splitSkillName(rest []string) (name string, installArgs []string, err error) {
	for i, a := range rest {
		if strings.HasPrefix(a, "-") {
			continue
		}
		name = a
		installArgs = append(installArgs, rest[:i]...)
		installArgs = append(installArgs, rest[i+1:]...)
		return name, installArgs, nil
	}
	return "", nil, fmt.Errorf("missing skill name")
}

func runEdit(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Print(editUsage)
		return nil
	}
	var addLabel, addExplanation string
	remainArgs, err := lessflags.String("--add-label", &addLabel).
		String("--add-explanation", &addExplanation).
		Parse(args)
	if err != nil {
		return err
	}
	if len(remainArgs) != 1 {
		return fmt.Errorf("edit requires exactly one concrete leaf path")
	}
	return edit.Edit(remainArgs[0], addLabel, addExplanation)
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
