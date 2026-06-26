package assert

import "strings"

var namedAnsiColors = map[string]string{
	"bold":  "1",
	"red":   "31",
	"green": "32",
	"gray":  "90",
}

func resolveAnsiOpen(spec AnsiColorSpec) (string, error) {
	var b strings.Builder
	for _, tok := range spec.Tokens {
		params, err := ansiTokenParams(tok)
		if err != nil {
			return "", err
		}
		b.WriteString("\x1b[")
		b.WriteString(params)
		b.WriteString("m")
	}
	return b.String(), nil
}

func ansiTokenParams(tok string) (string, error) {
	if strings.HasPrefix(tok, "#") {
		return tok[1:], nil
	}
	params, ok := namedAnsiColors[tok]
	if !ok {
		return "", parseErr(-1, "unknown ansi color name %q; use # for raw SGR", tok)
	}
	return params, nil
}

func parseAnsiSpec(specText string) (AnsiColorSpec, error) {
	specText = strings.TrimSpace(specText)
	if specText == "" {
		return AnsiColorSpec{}, parseErr(-1, "empty ansi-color specifier")
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
		if _, ok := namedAnsiColors[name]; !ok {
			return AnsiColorSpec{}, parseErr(-1, "unknown ansi color name %q", name)
		}
		tokens = append(tokens, name)
		specText = specText[end:]
	}
	return AnsiColorSpec{Tokens: tokens}, nil
}

const ansiReset = "\x1b[0m"