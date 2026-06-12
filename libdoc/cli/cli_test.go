package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunHelpOutput(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		wants []string
	}{
		{name: "no args", args: nil, wants: []string{"Usage: doctest <command>", "agent", "validate", "build", "test", "skill"}},
		{name: "top long", args: []string{"--help"}, wants: []string{"Usage: doctest <command>", "Run doctest <command> --help"}},
		{name: "top short", args: []string{"-h"}, wants: []string{"Usage: doctest <command>", "agent fill-code"}},
		{name: "agent no args", args: []string{"agent"}, wants: []string{"Usage: doctest agent <command>", "generate", "fill-code"}},
		{name: "agent long help", args: []string{"agent", "--help"}, wants: []string{"Usage: doctest agent <command>", "generate", "fill-code"}},
		{name: "agent short help", args: []string{"agent", "-h"}, wants: []string{"Usage: doctest agent <command>", "generate", "fill-code"}},
		{name: "agent generate help", args: []string{"agent", "generate", "--help"}, wants: []string{"Usage: doctest agent generate", "<idea>", "-d", "--dir", "--agent-runner"}},
		{name: "agent generate short help", args: []string{"agent", "generate", "-h"}, wants: []string{"Usage: doctest agent generate", "fake-codex"}},
		{name: "agent fill-code help", args: []string{"agent", "fill-code", "--help"}, wants: []string{"Usage: doctest agent fill-code <target-dir>"}},
		{name: "validate help", args: []string{"validate", "--help"}, wants: []string{"Usage: doctest validate <dir>"}},
		{name: "validate short help", args: []string{"validate", "-h"}, wants: []string{"Usage: doctest validate <dir>"}},
		{name: "build help", args: []string{"build", "--help"}, wants: []string{"Usage: doctest build", "-v", "--verbose", "--rm", "--gen-dir", "-count"}},
		{name: "build short help", args: []string{"build", "-h"}, wants: []string{"Usage: doctest build", "--gen-dir"}},
		{name: "test help", args: []string{"test", "--help"}, wants: []string{"Usage: doctest test", "-v", "--verbose", "--rm", "-count"}},
		{name: "test short help", args: []string{"test", "-h"}, wants: []string{"Usage: doctest test", "-count"}},
		{name: "skill no args", args: []string{"skill"}, wants: []string{"Usage: doctest skill --list", "doc-spec", "code-spec"}},
		{name: "skill long help", args: []string{"skill", "--help"}, wants: []string{"Usage: doctest skill --list", "install"}},
		{name: "skill short help", args: []string{"skill", "-h"}, wants: []string{"Usage: doctest skill --list", "show"}},
		{name: "skill list", args: []string{"skill", "--list"}, wants: []string{"doc-spec", "code-spec"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, err := captureStdout(func() error {
				return Run(tt.args)
			})
			if err != nil {
				t.Fatalf("Run(%v): %v", tt.args, err)
			}
			for _, want := range tt.wants {
				if !strings.Contains(stdout, want) {
					t.Fatalf("stdout missing %q:\n%s", want, stdout)
				}
			}
		})
	}
}

func TestRunErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "unknown top command", args: []string{"nope"}, wantErr: "unknown command: nope"},
		{name: "unknown agent command", args: []string{"agent", "nope"}, wantErr: "unknown agent command: nope"},
		{name: "fill code no dir", args: []string{"agent", "fill-code"}, wantErr: "agent fill-code requires <target-dir>"},
		{name: "fill code extra arg", args: []string{"agent", "fill-code", "a", "b"}, wantErr: "agent fill-code requires <target-dir>"},
		{name: "generate no idea", args: []string{"agent", "generate"}, wantErr: "agent generate requires <idea>"},
		{name: "generate short dir missing value", args: []string{"agent", "generate", "idea", "-d"}, wantErr: "-d requires value"},
		{name: "generate long dir missing value", args: []string{"agent", "generate", "idea", "--dir"}, wantErr: "--dir requires value"},
		{name: "generate runner missing value", args: []string{"agent", "generate", "idea", "--agent-runner"}, wantErr: "--agent-runner requires value"},
		{name: "validate no dir", args: []string{"validate"}, wantErr: "validate requires <dir>"},
		{name: "validate too many dirs", args: []string{"validate", "a", "b"}, wantErr: "validate requires <dir>"},
		{name: "skill missing action", args: []string{"skill", "doc-spec"}, wantErr: "skill requires doc-spec, code-spec, tdd, or implementer plus show or install"},
		{name: "skill unknown action", args: []string{"skill", "doc-spec", "nope"}, wantErr: "unknown skill action: nope"},
		{name: "skill unknown name", args: []string{"skill", "unknown", "show"}, wantErr: "unknown skill: unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := captureStdout(func() error {
				return Run(tt.args)
			})
			if err == nil {
				t.Fatalf("Run(%v): expected error containing %q", tt.args, tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Run(%v): expected error containing %q, got %v", tt.args, tt.wantErr, err)
			}
		})
	}
}

func TestRunAgentGenerateParsesIdeaAndFlags(t *testing.T) {
	tmp := t.TempDir()
	outDir := tmp + string(os.PathSeparator) + "generated"
	stdout, err := captureStdout(func() error {
		return Run([]string{"agent", "generate", "invoice", "cli", "--dir", outDir, "--agent-runner", "fake-codex"})
	})
	if err != nil {
		t.Fatalf("Run generate: %v\nstdout:\n%s", err, stdout)
	}
	data, err := os.ReadFile(outDir + string(os.PathSeparator) + "SETUP.md")
	if err != nil {
		t.Fatalf("read generated SETUP.md: %v", err)
	}
	if !bytes.Contains(data, []byte("Generated from idea: invoice cli")) {
		t.Fatalf("generated setup missing joined idea:\n%s", string(data))
	}
}

func captureStdout(fn func() error) (string, error) {
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	err = fn()
	closeErr := w.Close()
	var buf bytes.Buffer
	_, copyErr := io.Copy(&buf, r)
	_ = r.Close()
	if err != nil {
		return buf.String(), err
	}
	if closeErr != nil {
		return buf.String(), closeErr
	}
	return buf.String(), copyErr
}
