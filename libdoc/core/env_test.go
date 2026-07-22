package core

import (
	"strings"
	"testing"
)

func TestChildEnvKeyReplace(t *testing.T) {
	base := []string{
		"PATH=/bin",
		"GOCACHE=/old",
		"DOCTEST_SESSION_ID=old-sid",
		"KEEP=1",
	}
	got := ChildEnv(base,
		"GOCACHE=/new",
		"DOCTEST_SESSION_ID=new-sid",
	)
	var gocache, sid, keep, path string
	gocacheN, sidN := 0, 0
	for _, e := range got {
		k, v, _ := strings.Cut(e, "=")
		switch k {
		case "GOCACHE":
			gocache = v
			gocacheN++
		case "DOCTEST_SESSION_ID":
			sid = v
			sidN++
		case "KEEP":
			keep = v
		case "PATH":
			path = v
		}
	}
	if gocache != "/new" || gocacheN != 1 {
		t.Fatalf("GOCACHE=%q count=%d, want /new once", gocache, gocacheN)
	}
	if sid != "new-sid" || sidN != 1 {
		t.Fatalf("SESSION_ID=%q count=%d, want new-sid once", sid, sidN)
	}
	if keep != "1" || path != "/bin" {
		t.Fatalf("KEEP=%q PATH=%q", keep, path)
	}
}

func TestChildEnvLastOverrideWins(t *testing.T) {
	got := ChildEnv([]string{"K=1"}, "K=2", "K=3")
	n := 0
	var v string
	for _, e := range got {
		k, val, _ := strings.Cut(e, "=")
		if k == "K" {
			n++
			v = val
		}
	}
	if n != 1 || v != "3" {
		t.Fatalf("K count=%d val=%q, want once 3", n, v)
	}
}
