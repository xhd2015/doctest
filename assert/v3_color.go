package assert

import "strings"

var v3NamedColors = map[string]string{
	"bold":  "1",
	"red":   "31",
	"green": "32",
	"gray":  "90",
}

const (
	v3ColorReset          = "\x1b[0m"
	v3AnsiColorOpenPrefix = "<ansi-color"
	v3AnsiColorCloseTag   = "</ansi-color>"
	v3AnsiColorTagName    = "ansi-color"
)

type v3ColorSpec struct {
	Tokens []string
}

func v3ResolveColorOpen(spec v3ColorSpec) (string, error) {
	var b strings.Builder
	for _, tok := range spec.Tokens {
		params, err := v3ColorTokenParams(tok)
		if err != nil {
			return "", err
		}
		b.WriteString("\x1b[")
		b.WriteString(params)
		b.WriteString("m")
	}
	return b.String(), nil
}

func v3ColorTokenParams(tok string) (string, error) {
	if strings.HasPrefix(tok, "#") {
		return tok[1:], nil
	}
	params, ok := v3NamedColors[tok]
	if !ok {
		return "", parseErr(-1, "unknown color name %q; use # for raw SGR", tok)
	}
	return params, nil
}

func v3ParseColorSpec(specText string) (v3ColorSpec, error) {
	specText = strings.TrimSpace(specText)
	if specText == "" {
		return v3ColorSpec{}, parseErr(-1, "empty ansi-color specifier")
	}
	var tokens []string
	for specText != "" {
		specText = strings.TrimLeft(specText, " \t")
		if specText == "" {
			break
		}
		if specText[0] == '#' {
			end := 1
			for end < len(specText) && specText[end] != ' ' && specText[end] != '\t' {
				end++
			}
			tokens = append(tokens, specText[:end])
			specText = specText[end:]
			continue
		}
		end := 0
		for end < len(specText) && specText[end] != ' ' && specText[end] != '\t' {
			end++
		}
		name := specText[:end]
		if _, ok := v3NamedColors[name]; !ok {
			return v3ColorSpec{}, parseErr(-1, "unknown color name %q", name)
		}
		tokens = append(tokens, name)
		specText = specText[end:]
	}
	return v3ColorSpec{Tokens: tokens}, nil
}
