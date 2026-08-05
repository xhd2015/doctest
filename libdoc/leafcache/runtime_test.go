package leafcache

import "testing"

func TestSkipEnabled(t *testing.T) {
	if !SkipEnabled(0, false, false) {
		t.Fatal("default should enable skip")
	}
	if SkipEnabled(1, false, false) {
		t.Fatal("-count disables skip")
	}
	if SkipEnabled(0, true, false) {
		t.Fatal("-a disables skip")
	}
	if SkipEnabled(0, false, true) {
		t.Fatal("measure (cover/race) disables skip")
	}
}
