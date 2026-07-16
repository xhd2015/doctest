package legacy_v2

import "strings"

var namedColors = map[string]string{
	"bold":  "1",
	"red":   "31",
	"green": "32",
	"gray":  "90",
}

const (
	colorReset          = "\x1b[0m"
	ansiColorOpenPrefix = "<ansi-color"
	ansiColorCloseTag   = "</ansi-color>"
	ansiColorTagName    = "ansi-color"
)

func resolveColorOpen(spec colorSpec) (string, error) {
	var b strings.Builder
	for _, tok := range spec.Tokens {
		params, err := colorTokenParams(tok)
		if err != nil {
			return "", err
		}
		b.WriteString("\x1b[")
		b.WriteString(params)
		b.WriteString("m")
	}
	return b.String(), nil
}

func colorTokenParams(tok string) (string, error) {
	if strings.HasPrefix(tok, "#") {
		return tok[1:], nil
	}
	params, ok := namedColors[tok]
	if !ok {
		return "", parseErr(-1, "unknown color name %q; use # for raw SGR", tok)
	}
	return params, nil
}

func parseColorSpec(specText string) (colorSpec, error) {
	specText = strings.TrimSpace(specText)
	if specText == "" {
		return colorSpec{}, parseErr(-1, "empty ansi-color specifier")
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
		if _, ok := namedColors[name]; !ok {
			return colorSpec{}, parseErr(-1, "unknown color name %q", name)
		}
		tokens = append(tokens, name)
		specText = specText[end:]
	}
	return colorSpec{Tokens: tokens}, nil
}
