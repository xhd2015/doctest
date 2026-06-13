package implementer

import (
	"bytes"
	"context"
	"crypto/md5"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"

	"github.com/xhd2015/dot-pkgs/go-pkgs/logs"

	agentprovider "github.com/xhd2015/agent-pro/agent/cli/provider"
	"github.com/xhd2015/agent-pro/agent/cli/registry"
	agentexec "github.com/xhd2015/agent-pro/agent/exec"
	"github.com/xhd2015/agent-pro/agent_trace/events"
)

//go:embed PROMPT.md
var promptContent string

func PromptContent() string {
	return promptContent
}

func agentPrompt() string {
	s := promptContent
	if strings.HasPrefix(s, "---\n") {
		rest := s[4:]
		if idx := strings.Index(rest, "\n---\n"); idx >= 0 {
			s = rest[idx+5:]
			if strings.HasPrefix(s, "\n") {
				s = s[1:]
			}
		}
	}
	return s
}

type Options struct {
	Prompt      string
	AgentRunner string
	MockConfig  string
	SessionID   string
	Requirement string
	CatchUp     bool
	Status      bool
	ListSessions bool
}

func Run(opts Options) error {
	if opts.ListSessions {
		if opts.SessionID != "" {
			fmt.Fprintf(os.Stderr, "error: --list-sessions and --session-id are mutually exclusive\n")
			return nil
		}
		return ListSessions()
	}

	if opts.Status {
		if opts.SessionID == "" {
			fmt.Fprintf(os.Stderr, "error: --status requires --session-id\n")
			return nil
		}
		return ShowStatus(opts.SessionID)
	}

	if opts.CatchUp {
		if opts.SessionID == "" && os.Getenv("DOCTEST_AGENT_IMPLEMENTER_SESSION_ID") == "" && os.Getenv("CODEX_THREAD_ID") == "" {
			return fmt.Errorf("--trace requires --session-id")
		}
		return TraceSession(opts.SessionID)
	}

	prompt := strings.TrimSpace(opts.Prompt)
	if opts.Requirement != "" {
		data, err := os.ReadFile(opts.Requirement)
		if err != nil {
			return fmt.Errorf("read requirement file %s: %w", opts.Requirement, err)
		}
		reqContent := strings.TrimSpace(string(data))
		if prompt != "" {
			prompt = reqContent + "\n\n---\n\n" + prompt
		} else {
			prompt = reqContent
		}
	}
	if prompt == "" {
		return fmt.Errorf("agent implement requires <prompt>")
	}

	agentRunner := strings.TrimSpace(opts.AgentRunner)
	if agentRunner == "" {
		agentRunner = "opencode"
	}

	srcs, err := resolveSessionID(opts.SessionID)
	if err != nil {
		return err
	}
	srcs.agentRunner = agentRunner

	Logf("Session ID: %s (source: %s)\n", srcs.sessionID, sourceLabel(srcs))

	sessionDir, isNew, err := findOrCreateSession(srcs.sessionID, srcs)
	if err != nil {
		return fmt.Errorf("session: %w", err)
	}

	if isNew {
		if err := writeSessionPID(sessionDir); err != nil {
			return fmt.Errorf("write pid: %w", err)
		}
		defer removeSessionPID(sessionDir)
	}

	msgPath := filepath.Join(sessionDir, "messages.jsonl")
	msgEntry := map[string]string{
		"type":        "message",
		"content":     prompt,
		"create_time": time.Now().Format(time.RFC3339),
	}
	msgData, _ := json.Marshal(msgEntry)
	f, err := os.OpenFile(msgPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("write messages.jsonl: %w", err)
	}
	fmt.Fprintf(f, "%s\n", string(msgData))
	f.Close()

	tempDir, err := os.MkdirTemp("", "doctest-agent-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get executable: %w", err)
	}
	ypqPath := filepath.Join(tempDir, "yield-pending-questions")
	if out, err := exec.Command("cp", exe, ypqPath).CombinedOutput(); err != nil {
		return fmt.Errorf("copy yield-pending-questions: %w\n%s", err, string(out))
	}

	pathEntry := tempDir + string(filepath.ListSeparator)
	os.Setenv("PATH", pathEntry+os.Getenv("PATH"))

	questionsDir := filepath.Join(sessionDir, "questions")
	if err := os.MkdirAll(questionsDir, 0755); err != nil {
		return fmt.Errorf("create questions dir: %w", err)
	}
	questionFile := newQuestionsFile(questionsDir)
	os.Setenv("QUESTION_FIFO", questionFile)

	progressDir := filepath.Join(sessionDir, "progress")
	if err := os.MkdirAll(progressDir, 0755); err != nil {
		return fmt.Errorf("create progress dir: %w", err)
	}

	rpPath := filepath.Join(tempDir, "report-progress")
	if out, err := exec.Command("cp", exe, rpPath).CombinedOutput(); err != nil {
		return fmt.Errorf("copy report-progress: %w\n%s", err, string(out))
	}

	progressFile := filepath.Join(progressDir, time.Now().Format("20060102_150405")+"_progress_update.jsonl")
	os.Setenv("PROGRESS_FILE", progressFile)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logs.WatchLine(ctx, progressFile, logs.WatchLineOptions{}, func(line string) error {
			var entry map[string]any
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				return nil
			}
			desc, _ := entry["description"].(string)
			Logf("%s", desc)
			return nil
		})
	}()

	if opts.MockConfig != "" {
		os.Setenv("FAKE_CODEX_MOCK_CONFIG", opts.MockConfig)
	}

	var opencodeSessionID string
	if !isNew {
		opencodeSessionID = readOpencodeSessionID(sessionDir)
	}

	eventsPath := filepath.Join(sessionDir, "events.jsonl")
	eventsLogger, err := events.Open(eventsPath)
	if err != nil {
		return fmt.Errorf("open events.jsonl: %w", err)
	}
	defer eventsLogger.Close()

	capture := &sessionLogWriter{
		eventsFile: eventsLogger,
	}

	var fullPrompt string
	if isNew {
		fullPrompt = agentPrompt() + "\n\n---\n\n<user_request>\n" + prompt + "\n</user_request>\n"
	} else {
		fullPrompt = prompt
	}
	output, err := runAgent(agentRunner, fullPrompt, opencodeSessionID, capture)
	cancel()
	if err != nil {
		return fmt.Errorf("sub-agent failed: %w", err)
	}

	if isNew {
		sid := capture.sessionID
		if sid == "" {
			sid = srcs.sessionID
		}
		if updateErr := updateSessionMeta(sessionDir, sid, srcs); updateErr != nil {
			return fmt.Errorf("update session meta: %w", updateErr)
		}
	}

	fmt.Print(output)

	f, fErr := os.Open(questionFile)
	if fErr == nil {
		defer f.Close()
		var buf bytes.Buffer
		buf.ReadFrom(f)
		if buf.Len() > 0 {
			fmt.Print("\n\n---\nQUESTIONS\n---\n\n")
			fmt.Print(buf.String())
		}
	}

	return nil
}

func Logf(fmtStr string, args ...interface{}) {
	ts := time.Now().Format("2006-01-02T15:04:05")
	fmt.Println("[" + ts + "] " + fmt.Sprintf(fmtStr, args...))
}

type sessionIDSources struct {
	sessionID         string
	codexThreadID     string
	implSessionID     string
	explicitSessionID string
	agentRunner       string
}

func TraceSession(flagSessionID string) error {
	srcs, err := resolveSessionID(flagSessionID)
	if err != nil {
		return err
	}

	sessionDir, _, err := findOrCreateSession(srcs.sessionID, srcs)
	if err != nil {
		return fmt.Errorf("session: %w", err)
	}

	eventsPath := filepath.Join(sessionDir, "events.jsonl")

	printHeader := func(eventCount int) {
		fmt.Fprintf(os.Stdout, "\n═══════════════════════════════════════════════════════════════\n")
		fmt.Fprintf(os.Stdout, "  Session: %s\n", srcs.sessionID)
		fmt.Fprintf(os.Stdout, "  Events:  %d lines\n", eventCount)
		fmt.Fprintf(os.Stdout, "═══════════════════════════════════════════════════════════════\n\n")
	}

	data, err := os.ReadFile(eventsPath)
	if err != nil {
		if os.IsNotExist(err) {
			printHeader(0)
			fmt.Fprintf(os.Stdout, "  (no events yet)\n")
		} else {
			return fmt.Errorf("read events.jsonl: %w", err)
		}
	} else {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		eventCount := 0
		for _, line := range lines {
			if line != "" {
				eventCount++
			}
		}
		printHeader(eventCount)

		n := 0
		for _, line := range lines {
			if line == "" {
				continue
			}
			n++
			formatted := formatTraceEventLine(line)
			if formatted != "" {
				fmt.Fprintf(os.Stdout, "[%d]  %s\n", n, formatted)
			}
		}
	}

	sessionLive := isSessionLive(sessionDir)
	if !sessionLive {
		fmt.Fprintf(os.Stdout, "\n───────────────────────────────────────────────────────────────\n")
		fmt.Fprintf(os.Stdout, "  ✓ Done (session finished)\n")
		fmt.Fprintf(os.Stdout, "───────────────────────────────────────────────────────────────\n")
		return nil
	}

	fmt.Fprintf(os.Stdout, "\n  ⏳ Following new events (Ctrl+C to stop)...\n\n")

	var n int
	if data != nil {
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		for _, line := range lines {
			if line != "" {
				n++
			}
		}
	}
	n = n + 1

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var watchErr error
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		watchErr = logs.WatchLine(ctx, eventsPath, logs.WatchLineOptions{}, func(line string) error {
			formatted := formatTraceEventLine(line)
			if formatted != "" {
				fmt.Fprintf(os.Stdout, "[%d]  %s\n", n, formatted)
			}
			n++
			return nil
		})
	}()

	for {
		time.Sleep(2 * time.Second)
		if !isSessionLive(sessionDir) {
			cancel()
			break
		}
	}

	wg.Wait()

	if watchErr != nil && watchErr != context.Canceled {
		return watchErr
	}

	fmt.Fprintf(os.Stdout, "\n───────────────────────────────────────────────────────────────\n")
	fmt.Fprintf(os.Stdout, "  ✓ Session finished\n")
	fmt.Fprintf(os.Stdout, "───────────────────────────────────────────────────────────────\n")

	return nil
}

func ShowStatus(flagSessionID string) error {
	srcs, err := resolveSessionID(flagSessionID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil
	}

	base, err := sessionsBase()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return nil
	}

	sessionDir := findSession(base, srcs.sessionID, srcs)
	if sessionDir == "" {
		fmt.Fprintf(os.Stderr, "error: session not found: %s\n", srcs.sessionID)
		return nil
	}

	metaPath := filepath.Join(sessionDir, "meta.json")
	metaData, err := os.ReadFile(metaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read meta.json: %v\n", err)
		return nil
	}

	var metaMap map[string]any
	if err := json.Unmarshal(metaData, &metaMap); err != nil {
		fmt.Fprintf(os.Stderr, "error: invalid meta.json: %v\n", err)
		return nil
	}

	sesID, _ := metaMap["explicit_session_id"].(string)
	if sesID == "" {
		sesID = srcs.sessionID
	}

	runner, _ := metaMap["agent_runner"].(string)
	if runner == "" {
		runner = "opencode"
	}

	codex, _ := metaMap["main_agent_codex_thread_id"].(string)
	opencodeSID, _ := metaMap["opencode_session_id"].(string)

	createdAtStr, _ := metaMap["created_at"].(string)
	if t, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
		createdAtStr = t.Format("2006-01-02 15:04:05")
	}

	eventsPath := filepath.Join(sessionDir, "events.jsonl")
	var eventLines []string
	var lastTimestampMs int64
	eventsData, evErr := os.ReadFile(eventsPath)
	if evErr == nil {
		for _, line := range strings.Split(string(eventsData), "\n") {
			line = strings.TrimSpace(line)
			if line != "" && strings.HasPrefix(line, "{") {
				eventLines = append(eventLines, line)
			}
		}
		if len(eventLines) > 0 {
			lastTimestampMs = parseEventTimestamp(eventLines[len(eventLines)-1])
		}
	}

	eventCount := len(eventLines)
	lastRelative := "—"
	if lastTimestampMs > 0 {
		lastRelative = relativeTime(lastTimestampMs)
	}

	status := "finished"
	if isSessionLive(sessionDir) {
		pidData, _ := os.ReadFile(filepath.Join(sessionDir, "pid"))
		pid := strings.TrimSpace(string(pidData))
		status = fmt.Sprintf("running (PID %s)", pid)
	}

	fmt.Fprintf(os.Stdout, "\n═══════════════════════════════════════════════════════════════\n")
	fmt.Fprintf(os.Stdout, "  Session:  %s\n", sesID)
	fmt.Fprintf(os.Stdout, "  Status:   %s\n", status)
	fmt.Fprintf(os.Stdout, "  Runner:   %s\n", runner)
	fmt.Fprintf(os.Stdout, "  Created:  %s\n", createdAtStr)

	codexDisplay := codex
	if codexDisplay == "" {
		codexDisplay = "—"
	}
	fmt.Fprintf(os.Stdout, "  Codex:    %s\n", codexDisplay)
	opencodeDisplay := opencodeSID
	if opencodeDisplay == "" {
		opencodeDisplay = "—"
	}
	fmt.Fprintf(os.Stdout, "  Opencode: %s\n", opencodeDisplay)
	fmt.Fprintf(os.Stdout, "  Events:   %d lines (last: %s)\n", eventCount, lastRelative)
	fmt.Fprintf(os.Stdout, "═══════════════════════════════════════════════════════════════\n\n")

	if eventCount == 0 {
		fmt.Fprintf(os.Stdout, "No events yet\n")
	} else {
		for i, line := range eventLines {
			formatted := formatTraceEventLine(line)
			if formatted == "" {
				formatted = line
			}
			ts := parseEventTimestamp(line)
			rel := "—"
			if ts > 0 {
				rel = relativeTime(ts)
			}
			lines := strings.Split(formatted, "\n")
			fmt.Fprintf(os.Stdout, "  [%d] %s — %s\n", i+1, lines[0], rel)
			for _, l := range lines[1:] {
				fmt.Fprintf(os.Stdout, "       %s\n", l)
			}
		}
	}

	return nil
}

func ListSessions() error {
	base, err := sessionsBase()
	if err != nil {
		return err
	}

	type sessionInfo struct {
		ID        string
		Runner    string
		CreatedAt time.Time
	}

	var sessions []sessionInfo
	today := time.Now()

	for i := 0; i < 7; i++ {
		dateDir := today.AddDate(0, 0, -i).Format("2006/01/02")
		dayPath := filepath.Join(base, dateDir)
		entries, err := os.ReadDir(dayPath)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "sess_") {
				continue
			}
			sessDir := filepath.Join(dayPath, entry.Name())
			metaPath := filepath.Join(sessDir, "meta.json")
			data, err := os.ReadFile(metaPath)
			if err != nil {
				continue
			}
			var m map[string]any
			if err := json.Unmarshal(data, &m); err != nil {
				continue
			}

			sid, _ := m["explicit_session_id"].(string)
			if sid == "" {
				sid = entry.Name()
			}
			runner, _ := m["agent_runner"].(string)
			if runner == "" {
				runner = "opencode"
			}

			var createdAt time.Time
			if ts, ok := m["created_at"].(string); ok {
				createdAt, _ = time.Parse(time.RFC3339, ts)
			}

			sessions = append(sessions, sessionInfo{
				ID:        sid,
				Runner:    runner,
				CreatedAt: createdAt,
			})
		}
	}

	if len(sessions) == 0 {
		fmt.Println("No sessions found")
		return nil
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	fmt.Printf("\nSessions (%d):\n\n", len(sessions))
	for _, s := range sessions {
		timeStr := s.CreatedAt.Format("2006-01-02 15:04:05")
		fmt.Printf("%-15s %-10s %s\n", s.ID, s.Runner, timeStr)
	}

	return nil
}

func writeSessionPID(sessionDir string) error {
	pid := os.Getpid()
	return os.WriteFile(filepath.Join(sessionDir, "pid"), []byte(strconv.Itoa(pid)), 0644)
}

func removeSessionPID(sessionDir string) {
	os.Remove(filepath.Join(sessionDir, "pid"))
}

func isSessionLive(sessionDir string) bool {
	data, err := os.ReadFile(filepath.Join(sessionDir, "pid"))
	if err != nil {
		return false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return false
	}
	return processExists(pid)
}

func formatTraceEventLine(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || !strings.HasPrefix(trimmed, "{") {
		return ""
	}
	var event struct {
		Type      string              `json:"type"`
		SessionID string              `json:"sessionID,omitempty"`
		Part      *traceEventPart     `json:"part,omitempty"`
		Error     map[string]any      `json:"error,omitempty"`
	}
	if err := json.Unmarshal([]byte(trimmed), &event); err != nil {
		return ""
	}

	switch event.Type {
	case "text":
		if event.Part != nil {
			text := strings.TrimSpace(event.Part.Text)
			if text != "" {
				return "💬   ASSISTANT\n     " + strings.ReplaceAll(text, "\n", "\n     ")
			}
		}
	case "tool_use":
		if event.Part != nil {
			toolType := friendlyToolName(event.Part.Tool)
			status := ""
			if event.Part.State != nil {
				status = event.Part.State.Status
			}
			switch status {
			case "completed":
				summary := ""
				if event.Part.State != nil {
					if event.Part.State.Title != "" {
						summary = event.Part.State.Title
					}
					if event.Part.State.Output != "" {
						if summary != "" {
							summary = summary + "\n" + event.Part.State.Output
						} else {
							summary = event.Part.State.Output
						}
					}
				}
				if summary == "" && event.Part.State != nil && event.Part.State.Input != nil {
					summary = toolInputSummary(event.Part.Tool, event.Part.State.Input)
				}
				line := "⚡  " + toolType + " (done)"
				if summary != "" {
					line += "\n     " + strings.ReplaceAll(summary, "\n", "\n     ")
				}
				return line
			case "error", "failed":
				errStr := ""
				if event.Part.State != nil && event.Part.State.Error != "" {
					errStr = event.Part.State.Error
				}
				line := "⚡  " + toolType + " (failed)"
				if errStr != "" {
					line += "\n     " + strings.ReplaceAll(errStr, "\n", "\n     ")
				}
				return line
			case "in_progress", "running":
				return "⚡  " + toolType + " (running)"
			case "pending":
				return "⚡  " + toolType + " (pending)"
			default:
				return "⚡  " + toolType
			}
		}
	case "step_start":
		return "─── Step started ───"
	case "step_finish":
		reason := ""
		if event.Part != nil {
			reason = event.Part.Reason
		}
		line := "─── Step finished"
		if reason != "" {
			line += " (" + reason + ")"
		}
		line += " ───"
		return line
	case "error":
		errMsg := ""
		if event.Error != nil {
			if name, ok := event.Error["name"].(string); ok {
				errMsg = name
			}
			if data, ok := event.Error["data"].(map[string]any); ok {
				if msg, ok := data["message"].(string); ok && msg != "" {
					errMsg = msg
				}
			}
		}
		if errMsg == "" {
			errMsg = "unknown error"
		}
		return "❌  ERROR\n     " + strings.ReplaceAll(errMsg, "\n", "\n     ")
	case "reasoning":
		if event.Part != nil {
			text := strings.TrimSpace(event.Part.Text)
			if text != "" {
				return "🧠  REASONING\n     " + strings.ReplaceAll(text, "\n", "\n     ")
			}
		}
	}
	return ""
}

type traceEventPart struct {
	ID     string              `json:"id"`
	Type   string              `json:"type"`
	Tool   string              `json:"tool,omitempty"`
	Text   string              `json:"text,omitempty"`
	Reason string              `json:"reason,omitempty"`
	State  *traceEventPartState `json:"state,omitempty"`
}

type traceEventPartState struct {
	Status string         `json:"status"`
	Error  string         `json:"error,omitempty"`
	Title  string         `json:"title,omitempty"`
	Output string         `json:"output,omitempty"`
	Input  map[string]any `json:"input,omitempty"`
}

func friendlyToolName(tool string) string {
	switch tool {
	case "bash":
		return "Shell"
	case "read", "Read":
		return "Read"
	case "edit", "Edit":
		return "Edit"
	case "write", "Write":
		return "Write"
	case "glob", "Glob":
		return "Glob"
	case "grep", "Grep":
		return "Grep"
	case "task", "Task":
		return "SubAgent"
	case "todowrite", "TodoWrite":
		return "Plan"
	case "skill", "Skill":
		return "Skill"
	case "webfetch", "WebFetch":
		return "WebFetch"
	default:
		return tool
	}
}

func toolInputSummary(tool string, input map[string]any) string {
	switch tool {
	case "bash":
		if cmd, ok := input["command"].(string); ok {
			return cmd
		}
		if desc, ok := input["description"].(string); ok {
			return desc
		}
	case "task":
		if desc, ok := input["description"].(string); ok {
			if prompt, ok := input["prompt"].(string); ok {
				return desc + ": " + ellipsize(prompt, 100)
			}
			return desc
		}
	}
	for _, key := range []string{"command", "description", "pattern", "query", "path", "filePath", "question"} {
		if v, ok := input[key].(string); ok {
			return ellipsize(v, 120)
		}
	}
	return ""
}

func ellipsize(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func resolveSessionID(flagSessionID string) (*sessionIDSources, error) {
	if flagSessionID != "" {
		codexID := os.Getenv("CODEX_THREAD_ID")
		return &sessionIDSources{
			sessionID:         flagSessionID,
			codexThreadID:     codexID,
			explicitSessionID: flagSessionID,
		}, nil
	}
	if v := os.Getenv("DOCTEST_AGENT_IMPLEMENTER_SESSION_ID"); v != "" {
		codexID := os.Getenv("CODEX_THREAD_ID")
		return &sessionIDSources{
			sessionID:     v,
			codexThreadID: codexID,
			implSessionID: v,
		}, nil
	}
	if v := os.Getenv("CODEX_THREAD_ID"); v != "" {
		return &sessionIDSources{
			sessionID:     v,
			codexThreadID: v,
		}, nil
	}
	genID := generateSessionID()
	return nil, fmt.Errorf("cannot detect session id, if you're running inside opencode, try again with: `doctest agent implement --session-id %s <prompt>`, and use the same session id in subsequent followups, don't generate your session id, use the provided session id %s explicity.", genID, genID)
}

func generateSessionID() string {
	return "gen_" + fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
}

func sourceLabel(srcs *sessionIDSources) string {
	if srcs.explicitSessionID != "" {
		return "--session-id"
	}
	if srcs.implSessionID != "" {
		return "DOCTEST_AGENT_IMPLEMENTER_SESSION_ID"
	}
	if srcs.codexThreadID != "" {
		return "CODEX_THREAD_ID"
	}
	return "generated"
}

func runAgent(agentRunner, prompt, sessionID string, rawLog *sessionLogWriter) (string, error) {
	env := agentexec.NewEnv(&agentexec.PathsConfig{
		RootDirName: ".agent-pro",
		DataDirName: "data",
		BinDirName:  "bin",
	}, "AGENT_PRO_CONFIG_HOME")

	runner, err := agentprovider.Build(registry.AgentRunnerID(agentRunner), "", ".", env)
	if err != nil {
		return "", err
	}

	opts := &registry.AskOptions{
		Workspace: ".",
		SessionID: sessionID,
		RawLog:    rawLog,
	}

	output, err := runner.Agent.Ask(context.Background(), prompt, opts, func(delta string) {
		fmt.Print(delta)
	})
	if err != nil {
		return output, err
	}
	return output, nil
}

func HandleYieldPendingQuestions(args []string) error {
	questionFifo := os.Getenv("QUESTION_FIFO")
	if questionFifo == "" {
		return fmt.Errorf("QUESTION_FIFO must be set")
	}

	if len(args) == 0 {
		return fmt.Errorf("usage: yield-pending-questions '<json>' '<json>' ...")
	}

	f, err := os.OpenFile(questionFifo, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open questions file: %w", err)
	}
	defer f.Close()

	for i, arg := range args {
		var input struct {
			ID       string `json:"id"`
			Question string `json:"question"`
			Options  []struct {
				Option      string `json:"option"`
				Explanation string `json:"explanation"`
			} `json:"options"`
		}
		if err := json.Unmarshal([]byte(arg), &input); err != nil || input.Question == "" {
			continue
		}
		id := input.ID
		if id == "" {
			id = fmt.Sprintf("%d", i+1)
		}
		entry := map[string]any{
			"type":     "question",
			"id":       id,
			"question": input.Question,
		}
		if len(input.Options) > 0 {
			entry["options"] = input.Options
		}
		data, _ := json.Marshal(entry)
		fmt.Fprintf(f, "%s\n", string(data))
	}
	return nil
}

func HandleReportProgress(args []string) error {
	progressFile := os.Getenv("PROGRESS_FILE")
	if progressFile == "" {
		return fmt.Errorf("PROGRESS_FILE must be set")
	}

	if len(args) == 0 {
		return fmt.Errorf("usage: report-progress <description>")
	}

	description := strings.Join(args, " ")

	entry := map[string]string{
		"type":        "progress",
		"description": description,
		"timestamp":   time.Now().Format(time.RFC3339),
	}
	data, _ := json.Marshal(entry)

	f, err := os.OpenFile(progressFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open progress file: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", string(data))
	return nil
}

func sessionsBase() (string, error) {
	if v := os.Getenv("DOCTEST_DEBUG_SESSION_HOME"); v != "" {
		return v, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".doctest", "implementer", "sessions"), nil
}

type meta struct {
	ExplicitSessionID                string    `json:"explicit_session_id,omitempty"`
	DoctestAgentImplementerSessionID string    `json:"doctest_agent_implementer_session_id,omitempty"`
	MainAgentCodexThreadID           string    `json:"main_agent_codex_thread_id,omitempty"`
	OpencodeSessionID                string    `json:"opencode_session_id,omitempty"`
	AgentRunner                      string    `json:"agent_runner,omitempty"`
	CreatedAt                        time.Time `json:"created_at"`
}

func findOrCreateSession(threadID string, srcs *sessionIDSources) (dir string, isNew bool, err error) {
	base, err := sessionsBase()
	if err != nil {
		return "", false, err
	}

	dir = findSession(base, threadID, srcs)
	if dir != "" {
		return dir, false, nil
	}

	dir, err = createSession(base, threadID, srcs)
	if err != nil {
		return "", false, err
	}
	return dir, true, nil
}

func findSession(base, threadID string, srcs *sessionIDSources) string {
	matchField := sessionMatchField(srcs)
	today := time.Now()
	for i := 0; i < 7; i++ {
		dateDir := today.AddDate(0, 0, -i).Format("2006/01/02")
		dayPath := filepath.Join(base, dateDir)
		entries, err := os.ReadDir(dayPath)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "sess_") {
				continue
			}
			sessDir := filepath.Join(dayPath, entry.Name())
			metaPath := filepath.Join(sessDir, "meta.json")
			data, err := os.ReadFile(metaPath)
			if err != nil {
				continue
			}
			var m map[string]any
			if err := json.Unmarshal(data, &m); err != nil {
				continue
			}
			if v, ok := m[matchField]; ok {
				if s, ok := v.(string); ok && s == threadID {
					return sessDir
				}
			}
		}
	}
	return ""
}

func sessionMatchField(srcs *sessionIDSources) string {
	if srcs.explicitSessionID != "" {
		return "explicit_session_id"
	}
	if srcs.implSessionID != "" {
		return "doctest_agent_implementer_session_id"
	}
	return "main_agent_codex_thread_id"
}

func createSession(base, threadID string, srcs *sessionIDSources) (string, error) {
	now := time.Now()
	dateDir := now.Format("2006/01/02")
	sessID := fmt.Sprintf("sess_%s_%d", now.Format("150405"), now.UnixNano())
	sessDir := filepath.Join(base, dateDir, sessID)

	if err := os.MkdirAll(sessDir, 0755); err != nil {
		return "", fmt.Errorf("create session dir: %w", err)
	}

	m := meta{
		ExplicitSessionID:                srcs.explicitSessionID,
		DoctestAgentImplementerSessionID: srcs.implSessionID,
		MainAgentCodexThreadID:           srcs.codexThreadID,
		AgentRunner:                      srcs.agentRunner,
		CreatedAt:                        now,
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal meta: %w", err)
	}
	metaPath := filepath.Join(sessDir, "meta.json")
	if err := os.WriteFile(metaPath, append(data, '\n'), 0644); err != nil {
		return "", fmt.Errorf("write meta.json: %w", err)
	}

	return sessDir, nil
}

func newQuestionsFile(dir string) string {
	base := time.Now().Format("2006_01_02_15_04_05")
	name := base + ".json"
	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.WriteFile(path, nil, 0644)
		return path
	}
	for n := 1; ; n++ {
		name := base + "_" + strconv.Itoa(n) + ".json"
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.WriteFile(path, nil, 0644)
			return path
		}
	}
}

type sessionLogWriter struct {
	mu         sync.Mutex
	sessionID  string
	eventsFile events.Logger
}

func (w *sessionLogWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.eventsFile != nil {
		_ = w.eventsFile.Append(p)
	}

	if w.sessionID == "" {
		line := strings.TrimSpace(string(p))
		if line != "" && line[0] == '{' {
			var event struct {
				SessionID string `json:"sessionID,omitempty"`
			}
			if json.Unmarshal([]byte(line), &event) == nil && event.SessionID != "" {
				w.sessionID = event.SessionID
			}
		}
	}

	return len(p), nil
}

func (w *sessionLogWriter) Close() error {
	if w.eventsFile != nil {
		return w.eventsFile.Close()
	}
	return nil
}

func readOpencodeSessionID(sessionDir string) string {
	metaPath := filepath.Join(sessionDir, "meta.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return ""
	}
	var m meta
	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}
	return m.OpencodeSessionID
}

func updateSessionMeta(sessionDir, innerSessionID string, srcs *sessionIDSources) error {
	metaPath := filepath.Join(sessionDir, "meta.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return err
	}
	var m meta
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	m.OpencodeSessionID = innerSessionID
	if srcs != nil {
		m.MainAgentCodexThreadID = srcs.codexThreadID
	}
	newData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(metaPath, append(newData, '\n'), 0644)
}

func processExists(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func parseEventTimestamp(line string) int64 {
	var event struct {
		Timestamp int64 `json:"timestamp"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(line)), &event); err != nil {
		return 0
	}
	return event.Timestamp
}

func relativeTime(timestampMs int64) string {
	if timestampMs <= 0 {
		return "—"
	}
	now := time.Now().UnixMilli()
	diff := now - timestampMs
	if diff < 0 {
		diff = 0
	}
	seconds := diff / 1000
	if seconds < 60 {
		return fmt.Sprintf("%ds ago", seconds)
	}
	minutes := seconds / 60
	if minutes < 60 {
		return fmt.Sprintf("%dm ago", minutes)
	}
	hours := minutes / 60
	if hours < 24 {
		return fmt.Sprintf("%dh ago", hours)
	}
	days := hours / 24
	return fmt.Sprintf("%dd ago", days)
}
