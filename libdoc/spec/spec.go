package spec

import (
	"fmt"
	"sort"

	"github.com/xhd2015/doctest/doc"
	"github.com/xhd2015/doctest/libdoc/designer"
	"github.com/xhd2015/doctest/libdoc/implementer"
	"github.com/xhd2015/skills/install"
)

type entry struct {
	SkillName   string
	FileName    string
	ContentFunc func() (string, error)
}

var entries = map[string]entry{
	"doc-spec":    {SkillName: "doc-style-test-specification", FileName: "DOC_STYLE_TEST_SPECIFICATION.md"},
	"code-spec":   {SkillName: "doc-style-test-code-specification", FileName: "DOC_STYLE_TEST_CODE_SPECIFICATION.md"},
	"tdd":                    {SkillName: "doctest-tdd", FileName: "DOCTEST_TDD.md"},
	"tdd-cli-agent": {SkillName: "tdd-cli-agent", FileName: "DOCTEST_TDD_CLI_AGENT.md"},
	"tdd-lite":               {SkillName: "doctest-tdd-lite", FileName: "DOCTEST_TDD_LITE.md"},
	"reproduce":   {SkillName: "doctest-reproduce", FileName: "DOCTEST_REPRODUCE.md"},
	"review":        {SkillName: "doctest-review", FileName: "DOCTEST_REVIEW.md"},
	"output-assert": {SkillName: "doctest-output-assert", FileName: "DOCTEST_OUTPUT_ASSERT.md"},
	"implementer": {SkillName: "doctest-implementer", ContentFunc: func() (string, error) { return implementer.PromptContent(), nil }},
	"designer":    {SkillName: "doctest-designer", ContentFunc: func() (string, error) { return designer.PromptContent(), nil }},
}

func Content(name string) (string, error) {
	ent, ok := entries[name]
	if !ok {
		return "", fmt.Errorf("unknown skill: %s", name)
	}
	if ent.ContentFunc != nil {
		return ent.ContentFunc()
	}
	return doc.Content(ent.FileName)
}

func Install(name string, args []string) error {
	ent, ok := entries[name]
	if !ok {
		return fmt.Errorf("unknown skill: %s", name)
	}
	content, err := Content(name)
	if err != nil {
		return err
	}
	return install.HandleInstall(install.InstallOptions{
		SkillDirName: ent.SkillName,
		SkillContent: content,
		Usage:        "doctest skill --install " + name,
	}, args)
}

// AllUpdateSkills returns every registry skill for batch update.
func AllUpdateSkills() ([]install.UpdateSkill, error) {
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)

	skills := make([]install.UpdateSkill, 0, len(names))
	for _, name := range names {
		ent := entries[name]
		content, err := Content(name)
		if err != nil {
			return nil, err
		}
		skills = append(skills, install.UpdateSkill{
			InstallOptions: install.InstallOptions{
				SkillDirName: ent.SkillName,
				SkillContent: content,
				Usage:        "doctest skills update",
			},
			Name: name,
		})
	}
	return skills, nil
}

// Update runs batch update for all registry skills.
func Update(args []string) error {
	skills, err := AllUpdateSkills()
	if err != nil {
		return err
	}
	return install.HandleUpdateMany(skills, args)
}
