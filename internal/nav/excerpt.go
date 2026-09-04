package nav

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const readmeBlurbMaxRunes = 450

// ReadmeNavFromIndex derives README heading + blurb from a service index.md.
// Blurb = first paragraph under "## What is …" / "## 什么是…", truncated at a
// sentence boundary when longer than readmeBlurbMaxRunes (matches hcbp-demo style).
func ReadmeNavFromIndex(indexContent string, loc LocaleStrings, fallbackLabel string) (heading, blurb string) {
	name := ServiceNameFromIndex(indexContent, loc)
	if name == "" {
		name = fallbackLabel
	}
	if loc.ReadmeNavHeading == ZhCN.ReadmeNavHeading || loc.IntroLabel == ZhCN.IntroLabel {
		heading = name + "最佳实践"
	} else {
		heading = name + " Best Practices"
	}
	blurb = ExcerptWhatIsParagraph(indexContent, loc)
	return heading, blurb
}

// ServiceNameFromIndex returns the name in "## What is X" / "## 什么是X".
func ServiceNameFromIndex(content string, loc LocaleStrings) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	for _, line := range strings.Split(content, "\n") {
		trim := strings.TrimSpace(line)
		if !strings.HasPrefix(trim, "## ") {
			continue
		}
		title := strings.TrimSpace(strings.TrimPrefix(trim, "## "))
		switch {
		case strings.HasPrefix(strings.ToLower(title), "what is "):
			return strings.TrimSpace(title[len("What is "):])
		case strings.HasPrefix(title, "什么是"):
			name := strings.TrimSpace(strings.TrimPrefix(title, "什么是"))
			// drop trailing full-width/half-width description in parens already part of title — keep as-is
			return name
		}
	}
	_ = loc
	return ""
}

// ExcerptWhatIsParagraph returns the README-style intro excerpt from index.md.
func ExcerptWhatIsParagraph(content string, loc LocaleStrings) string {
	para := firstWhatIsParagraph(content)
	if para == "" {
		return ""
	}
	para = collapseSpace(para)
	if utf8.RuneCountInString(para) <= readmeBlurbMaxRunes {
		return para
	}
	// Prefer first sentence when the opening paragraph is long (e.g. Anti-DDoS).
	if sent := firstSentence(para, loc.SentenceEnd); sent != "" {
		return sent
	}
	return truncateRunes(para, readmeBlurbMaxRunes)
}

func firstWhatIsParagraph(content string) string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	lines := strings.Split(content, "\n")
	start := -1
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if !strings.HasPrefix(trim, "## ") {
			continue
		}
		title := strings.TrimSpace(strings.TrimPrefix(trim, "## "))
		lower := strings.ToLower(title)
		if strings.HasPrefix(lower, "what is ") || strings.HasPrefix(title, "什么是") {
			start = i + 1
			break
		}
	}
	if start < 0 {
		return ""
	}
	// skip blank lines
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	var b strings.Builder
	for i := start; i < len(lines); i++ {
		trim := strings.TrimSpace(lines[i])
		if trim == "" {
			if b.Len() > 0 {
				break
			}
			continue
		}
		if strings.HasPrefix(trim, "#") {
			break
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(trim)
	}
	return strings.TrimSpace(b.String())
}

func firstSentence(s, preferredEnd string) string {
	ends := []string{".", "。", "!", "？", "?"}
	if preferredEnd != "" {
		ends = append([]string{preferredEnd}, ends...)
	}
	runes := []rune(s)
	for i, r := range runes {
		for _, e := range ends {
			er, _ := utf8.DecodeRuneInString(e)
			if r != er {
				continue
			}
			// require end of sentence: end of text or whitespace after
			if i+1 >= len(runes) || unicode.IsSpace(runes[i+1]) {
				return strings.TrimSpace(string(runes[:i+1]))
			}
		}
	}
	return ""
}

func collapseSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return strings.TrimSpace(string(r[:n])) + "…"
}
