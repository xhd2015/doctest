package version

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestVersionMatchesCmdDoctestVERSIONTxt(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	cmdVersionPath := filepath.Join(filepath.Dir(thisFile), "..", "cmd", "doctest", "VERSION.txt")
	data, err := os.ReadFile(cmdVersionPath)
	if err != nil {
		t.Fatalf("read cmd/doctest/VERSION.txt: %v", err)
	}
	if string(data) != string(versionBytes) {
		t.Fatalf("version/VERSION.txt and cmd/doctest/VERSION.txt differ")
	}
}