package lint

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Parsed struct {
	Type     string
	Scope    string
	Subject  string
	Breaking bool // optional future use
}

func ParseConventionalLine(line string) (Parsed, bool) {
	message := strings.TrimSpace(line)
	if message == "" {
		return Parsed{}, false
	}

	var parsed Parsed
	cursor := 0
	length := len(message)

	typeStart := cursor
	for cursor < length {
		char := message[cursor]
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
			cursor++
			continue
		}
		break
	}
	if cursor == typeStart {
		return Parsed{}, false
	}
	parsed.Type = message[typeStart:cursor]

	if cursor < length && message[cursor] == '(' {
		cursor++ // skip '('
		scopeStart := cursor

		for cursor < length && message[cursor] != ')' {
			cursor++
		}
		if cursor >= length || message[cursor] != ')' {
			// Missing closing parenthesis
			return Parsed{}, false
		}

		parsed.Scope = message[scopeStart:cursor]
		cursor++ // skip ')'
	}

	if cursor < length && message[cursor] == '!' {
		parsed.Breaking = true
		cursor++
	}

	if cursor >= length || message[cursor] != ':' {
		return Parsed{}, false
	}
	cursor++ // skip ':'

	for cursor < length {
		runeValue, size := utf8.DecodeRuneInString(message[cursor:])
		if runeValue == utf8.RuneError && size == 1 {
			// Invalid UTF-8 sequence, treat as non-space
			break
		}
		if !unicode.IsSpace(runeValue) {
			break
		}
		cursor += size
	}

	if cursor >= length {
		// Nothing after ':'
		return Parsed{}, false
	}
	parsed.Subject = message[cursor:]

	return parsed, true
}

func FormatExample(requireScope bool) string {
	if requireScope {
		return "type(scope): subject"
	}
	return "type: subject"
}

func Errorf(format string, args ...any) string {
	return " - " + fmt.Sprintf(format, args...)
}
