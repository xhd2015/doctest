package build

import (
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
)

func TestLeafTimingsFromSubtests(t *testing.T) {
	cases := []core.TreeCase{
		{Path: "flags/default-metrics-on"},
		{Path: "recording/enabled-writes-run-start-end"},
		{Path: "warn-predicate/no-fire-total-zero"},
	}
	result := goTestJSONResult{
		testElapsedNs: map[string]int64{
			"TestDoctestSuite/flags/default-metrics-on":                 12 * int64(time.Millisecond),
			"TestDoctestSuite/recording/enabled-writes-run-start-end": 1500 * int64(time.Millisecond),
			"TestDoctestSuite/warn-predicate/no-fire-total-zero":         3 * int64(time.Millisecond),
			"TestDoctestSuite": 2 * int64(time.Second), // parent — ignored
		},
	}
	got := leafTimingsFromSubtests(cases, result, 2*time.Second)
	if len(got) != 3 {
		t.Fatalf("len=%d want 3", len(got))
	}
	byPath := map[string]int64{}
	for _, lt := range got {
		byPath[lt.Path] = lt.ElapsedNs
	}
	if byPath["flags/default-metrics-on"] != 12*int64(time.Millisecond) {
		t.Fatalf("pure leaf elapsed=%d", byPath["flags/default-metrics-on"])
	}
	if byPath["recording/enabled-writes-run-start-end"] != 1500*int64(time.Millisecond) {
		t.Fatalf("recording leaf elapsed=%d", byPath["recording/enabled-writes-run-start-end"])
	}
	if byPath["warn-predicate/no-fire-total-zero"] != 3*int64(time.Millisecond) {
		t.Fatalf("warn leaf elapsed=%d", byPath["warn-predicate/no-fire-total-zero"])
	}
}

func TestLeafTimingsFromSubtestsSingleLeafFallback(t *testing.T) {
	cases := []core.TreeCase{{Path: "only"}}
	result := goTestJSONResult{
		pkgElapsedNs: map[string]int64{"testcase/suite": 400 * int64(time.Millisecond)},
	}
	got := leafTimingsFromSubtests(cases, result, time.Second)
	if len(got) != 1 || got[0].ElapsedNs != 400*int64(time.Millisecond) {
		t.Fatalf("got=%+v", got)
	}
}
