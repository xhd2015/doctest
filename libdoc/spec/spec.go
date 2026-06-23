package spec

import (
	"fmt"

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
	"tdd":         {SkillName: "doc-style-test-based-tdd", FileName: "DOC_STYLE_TEST_BASED_TDD.md"},
	"tdd-lite":    {SkillName: "doc-style-test-based-tdd-lite", FileName: "DOC_STYLE_TEST_BASED_TDD_LITE.md"},
	"implementer": {SkillName: "doc-style-test-based-tdd-implementer", ContentFunc: func() (string, error) { return implementer.PromptContent(), nil }},
	"designer":    {SkillName: "doc-style-test-based-tdd-designer", ContentFunc: func() (string, error) { return designer.PromptContent(), nil }},
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
		Usage:        "doctest skill " + name + " install",
	}, args)
}
