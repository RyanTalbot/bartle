package lint

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/RyanTalbot/bartle/internal/config"
)

type Result struct {
	Valid  bool
	Errors []string
}

func ValidateMessage(msg string, cfg config.Config) Result {
	var res Result

	// TODO: handle body/footer rules.
	firstLine := strings.Split(msg, "\n")[0]
	firstLine = strings.ReplaceAll(firstLine, "\r", "") // normalize CRLF
	firstLine = strings.TrimSpace(firstLine)
	if firstLine == "" {
		res.Errors = append(res.Errors, Errorf("empty commit message"))
		return finish(res)
	}

	switch strings.ToLower(cfg.Style) {
	case "conventional", "":
		res = validateConventional(firstLine, cfg.Rules)
	case "jira":
		res = validateJIRA(firstLine, cfg.Rules)
	default:
		res = validateConventional(firstLine, cfg.Rules)
	}

	return finish(res)
}

func finish(r Result) Result {
	r.Valid = len(r.Errors) == 0
	return r
}

func validateConventional(line string, rules config.Rules) Result {
	var out Result

	parsed, ok := ParseConventionalLine(line)
	if !ok {
		colonIdx := strings.Index(line, ":")
		if colonIdx < 0 {
			out.Errors = append(out.Errors, Errorf("missing ':' separator (e.g., %s)", FormatExample(rules.ScopeRequired)))
			return finish(out)
		}

		subject := strings.TrimSpace(line[colonIdx+1:])
		if subject == "" {
			out.Errors = append(out.Errors, Errorf("empty subject after ':'"))
		}

		hasOpen := strings.Contains(line, "(")
		hasClose := strings.Contains(line, ")")
		if hasOpen && !hasClose {
			out.Errors = append(out.Errors, Errorf("unclosed scope '(' — expected ')': e.g., %s", FormatExample(true)))
		}
		if rules.ScopeRequired && !hasOpen {
			out.Errors = append(out.Errors, Errorf("missing scope (e.g., type(scope): subject)"))
		}

		if len(out.Errors) == 0 {
			out.Errors = append(out.Errors, Errorf("not conventional format (e.g., %s)", FormatExample(rules.ScopeRequired)))
		}
		return finish(out)
	}

	if parsed.Type != strings.ToLower(parsed.Type) {
		out.Errors = append(out.Errors, Errorf("type must be lowercase (got %q)", parsed.Type))
	}

	if !inStringSet(rules.Types, parsed.Type) {
		out.Errors = append(out.Errors, Errorf("type %q not allowed (choose one of: %s)", parsed.Type, strings.Join(rules.Types, ", ")))
	}

	if rules.ScopeRequired && parsed.Scope == "" {
		out.Errors = append(out.Errors, Errorf("scope required (e.g., %s)", "type(scope): subject"))
	}

	if rules.MaxLineLength > 0 && utf8.RuneCountInString(line) > rules.MaxLineLength {
		out.Errors = append(out.Errors, Errorf("first line too long (%d > %d)",
			utf8.RuneCountInString(line), rules.MaxLineLength))
	}

	if rules.LowercaseStart && len(parsed.Subject) > 0 && isUpper(rune(parsed.Subject[0])) {
		out.Errors = append(out.Errors, Errorf("subject should start lowercase"))
	}

	return finish(out)
}

func validateJIRA(line string, rules config.Rules) Result {
	var out Result

	colon := strings.Index(line, ":")
	if colon <= 0 {
		out.Errors = append(out.Errors, Errorf("missing ':' separator (e.g., ABC-123: summary)"))
		return finish(out)
	}

	prefix := strings.TrimSpace(line[:colon])
	subject := strings.TrimSpace(line[colon+1:])

	if subject == "" {
		out.Errors = append(out.Errors, Errorf("empty subject after ':'"))
	}

	if !looksLikeTicket(prefix) {
		out.Errors = append(out.Errors, Errorf("prefix %q doesn't look like a ticket (e.g., ABC-123)", prefix))
	}

	if rules.MaxLineLength > 0 && utf8.RuneCountInString(line) > rules.MaxLineLength {
		out.Errors = append(out.Errors, Errorf("first line too long (%d > %d)",
			utf8.RuneCountInString(line), rules.MaxLineLength))
	}

	return finish(out)
}

func looksLikeTicket(s string) bool {
	if len(s) < 5 {
		return false
	}
	dash := strings.Index(s, "-")
	if dash < 2 {
		return false
	}
	left, right := s[:dash], s[dash+1:]
	if left == "" || right == "" {
		return false
	}
	for _, r := range left {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	for _, r := range right {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func inStringSet(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func isUpper(r rune) bool { return unicode.IsUpper(r) }
