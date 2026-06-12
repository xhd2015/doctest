package agent

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GenerateOptions struct {
	Idea        string
	Dir         string
	AgentRunner string
}

func Generate(opts GenerateOptions) error {
	dir := strings.TrimSpace(opts.Dir)
	if dir == "" {
		dir = "doctest-test-cases"
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# Generated Doc Tests\n"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte("## Preconditions\n- Generated from idea: "+opts.Idea+"\n"), 0644); err != nil {
		return err
	}
	output, err := runFakeCodex(opts.Idea)
	if output != "" {
		fmt.Println(output)
	}
	return err
}

func FillCode(dir string) error {
	if strings.TrimSpace(dir) == "" {
		return fmt.Errorf("agent fill-code requires <target-dir>")
	}
	info, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", dir)
	}
	return nil
}

func runFakeCodex(prompt string) (string, error) {
	path := strings.TrimSpace(os.Getenv("AGENT_RUNNER_FAKE_CODEX_PATH"))
	if path == "" {
		path = strings.TrimSpace(os.Getenv("AGENT_RUNNER_CODEX_PATH"))
	}
	if path == "" {
		return "", nil
	}
	cmd := exec.CommandContext(context.Background(), path, "exec", "--json", prompt)
	cmd.Env = os.Environ()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}
	var out strings.Builder
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		var event struct {
			Item *struct {
				Text string `json:"text"`
			} `json:"item"`
		}
		if json.Unmarshal([]byte(scanner.Text()), &event) == nil && event.Item != nil && event.Item.Text != "" {
			out.WriteString(event.Item.Text)
		}
	}
	if scanErr := scanner.Err(); scanErr != nil {
		return out.String(), scanErr
	}
	return out.String(), cmd.Wait()
}
