package assert_test

import (
	"testing"
	"github.com/xhd2015/doctest/assert"
)

func TestProbeMixed(t *testing.T) {
	actual := "\nSkill is up to date: /tmp/.agents/skills/skill-alpha\nskill not installed: skill-beta"
	assert.Output(t, actual, `
<contains>
Skill is up to date
skill-alpha
skill not installed: skill-beta
</contains>`)
}
