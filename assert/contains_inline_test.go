package assert

import "testing"

func TestContainsInlineAnyOf(t *testing.T) {
	p := MustParse(`<contains>
default
<any-of><expect>not configured</expect><expect>profiles</expect><expect>no profiles</expect></any-of>
</contains>`)

	if err := Match(p, "default\nnot configured\n", Contains()); err != nil {
		t.Fatalf("expected inline any-of in contains to match: %v", err)
	}
}

func TestContainsInlineAnyOfRejectsMissingBranches(t *testing.T) {
	p := MustParse(`<contains>
<any-of><expect>not configured</expect><expect>profiles</expect><expect>no profiles</expect></any-of>
</contains>`)

	if err := Match(p, "default\nempty\n", Contains()); err == nil {
		t.Fatal("expected inline any-of in contains to reject non-matching output")
	}
}
