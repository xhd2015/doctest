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
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	agentprovider "github.com/xhd2015/agent-pro/agent/cli/provider"
	"github.com/xhd2015/agent-pro/agent/cli/registry"
	agentexec "github.com/xhd2015/agent-pro/agent/exec"
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
}

func Run(opts Options) error {
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

	sessionDir, isNew, err := findOrCreateSession(srcs.sessionID, srcs)
	if err != nil {
		return fmt.Errorf("session: %w", err)
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

	if opts.MockConfig != "" {
		os.Setenv("FAKE_CODEX_MOCK_CONFIG", opts.MockConfig)
	}

	var opencodeSessionID string
	if !isNew {
		opencodeSessionID = readOpencodeSessionID(sessionDir)
	}

	capture := &sessionIDCapture{}

	var fullPrompt string
	if isNew {
		fullPrompt = agentPrompt() + "\n\n---\n\n<user_request>\n" + prompt + "\n</user_request>\n"
	} else {
		fullPrompt = prompt
	}
	output, err := runAgent(agentRunner, fullPrompt, opencodeSessionID, capture)
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

type sessionIDSources struct {
	sessionID         string
	codexThreadID     string
	implSessionID     string
	explicitSessionID string
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
	return nil, fmt.Errorf("cannot detect session id, if you're running inside opencode, run with: doctest agent implement --session-id %s <prompt>, and use the same session id in subsequent followups; or set DOCTEST_AGENT_IMPLEMENTER_SESSION_ID or CODEX_THREAD_ID", genID)
}

func generateSessionID() string {
	return "gen_" + fmt.Sprintf("%x", md5.Sum([]byte(uuid.New().String())))
}

func runAgent(agentRunner, prompt, sessionID string, rawLog *sessionIDCapture) (string, error) {
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

	output, err := runner.Agent.Ask(context.Background(), prompt, opts, func(delta string) {})
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

func sessionsBase() (string, error) {
	if v := os.Getenv("DOCTEST_DEBUG_SESSION_HOME"); v != "" {
		return v, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".agent-pro", "dedicated-agents", "doctest-agent-implementer", "sessions"), nil
}

type meta struct {
	ExplicitSessionID                  string    `json:"explicit_session_id,omitempty"`
	DoctestAgentImplementerSessionID   string    `json:"doctest_agent_implementer_session_id,omitempty"`
	MainAgentCodexThreadID             string    `json:"main_agent_codex_thread_id,omitempty"`
	OpencodeSessionID                  string    `json:"opencode_session_id,omitempty"`
	CreatedAt                          time.Time `json:"created_at"`
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

type sessionIDCapture struct {
	sessionID string
}

func (c *sessionIDCapture) Write(p []byte) (int, error) {
	if c.sessionID != "" {
		return len(p), nil
	}
	line := strings.TrimSpace(string(p))
	if line == "" || line[0] != '{' {
		return len(p), nil
	}
	var event struct {
		SessionID string `json:"sessionID,omitempty"`
	}
	if json.Unmarshal([]byte(line), &event) == nil && event.SessionID != "" {
		c.sessionID = event.SessionID
	}
	return len(p), nil
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
