package debug

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseEmpty(t *testing.T) {
	s, err := Parse("")
	if err != nil {
		t.Fatal(err)
	}
	if s.BypassGoTest || s.CPUProfile != "" || s.MemProfile != "" || s.BlockProfile != "" {
		t.Fatal("expected zero settings")
	}
}

func TestParseBypassGoTest(t *testing.T) {
	s, err := Parse("bypass-go-test=1")
	if err != nil {
		t.Fatal(err)
	}
	if !s.BypassGoTest {
		t.Fatal("expected BypassGoTest")
	}
	s, err = Parse("bypass-go-test=0")
	if err != nil {
		t.Fatal(err)
	}
	if s.BypassGoTest {
		t.Fatal("expected off")
	}
}

func TestParseProfilesCombined(t *testing.T) {
	s, err := Parse("bypass-go-test=1,cpuprofile=/tmp/c.pprof,memprofile=/tmp/m.pprof,blockprofile=/tmp/b.pprof")
	if err != nil {
		t.Fatal(err)
	}
	if !s.BypassGoTest {
		t.Fatal("bypass")
	}
	if s.CPUProfile != "/tmp/c.pprof" || s.MemProfile != "/tmp/m.pprof" || s.BlockProfile != "/tmp/b.pprof" {
		t.Fatalf("profiles: %+v", s)
	}
}

func TestParseProfileEmptyPathErrors(t *testing.T) {
	for _, raw := range []string{"cpuprofile=", "memprofile=", "blockprofile="} {
		if _, err := Parse(raw); err == nil {
			t.Fatalf("expected error for %q", raw)
		}
	}
}

func TestParseUnknownKeyErrors(t *testing.T) {
	_, err := Parse("bypass-go-test=1,not-a-key=1")
	if err == nil {
		t.Fatal("expected error")
	}
	if want := `unknown key "not-a-key"`; !strings.Contains(err.Error(), want) {
		t.Fatalf("err=%v want substring %q", err, want)
	}
}

func TestParseMissingEquals(t *testing.T) {
	_, err := Parse("bypass-go-test")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestStartProfilesMemAndBlock(t *testing.T) {
	dir := t.TempDir()
	mem := filepath.Join(dir, "mem.pprof")
	block := filepath.Join(dir, "block.pprof")
	s := Settings{MemProfile: mem, BlockProfile: block}
	stop, err := StartProfiles(s)
	if err != nil {
		t.Fatal(err)
	}
	// Touch some allocs so heap profile is non-empty-ish.
	_ = make([]byte, 1<<20)
	stop()
	if st, err := os.Stat(mem); err != nil || st.Size() == 0 {
		t.Fatalf("mem profile: err=%v size=%v", err, st)
	}
	if _, err := os.Stat(block); err != nil {
		t.Fatalf("block profile: %v", err)
	}
}

func TestStartProfilesCPU(t *testing.T) {
	dir := t.TempDir()
	cpu := filepath.Join(dir, "cpu.pprof")
	stop, err := StartProfiles(Settings{CPUProfile: cpu})
	if err != nil {
		t.Fatal(err)
	}
	// Brief busy work so the profile has samples.
	for i := 0; i < 1e6; i++ {
		_ = i * i
	}
	stop()
	if st, err := os.Stat(cpu); err != nil || st.Size() == 0 {
		t.Fatalf("cpu profile: err=%v size=%v", err, st)
	}
}
