package nav

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Prefer keeping the full first paragraph under "## What is / 什么是".
// Only trim when that paragraph is unusually long (e.g. Anti-DDoS EN).
// History: max=450 caused AAD EN (459 runes) to collapse to a 92-rune first
// sentence while ZH (~160) kept the full paragraph — README looked "broken" on EN only.
const readmeBlurbMaxRunes = 600

// ReadmeNavFromIndex derives README heading + blurb from a service index.md.
// Blurb = first paragraph under "## What is …" / "## 什么是…", optionally
// trimmed to as many complete sentences as fit in readmeBlurbMaxRunes.
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
		lower := strings.ToLower(title)
		switch {
		case strings.HasPrefix(lower, "what is "):
			// slice by rune-safe prefix length of "what is " (always ASCII)
			return strings.TrimSpace(title[len("what is "):])
		case strings.HasPrefix(title, "什么是"):
			return strings.TrimSpace(strings.TrimPrefix(title, "什么是"))
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
	// Keep as many complete sentences as fit (not only the first — short EN
	// openers like "… Huawei Cloud." must not wipe the rest of the paragraph).
	if packed := packSentences(para, loc.SentenceEnd, readmeBlurbMaxRunes); packed != "" {
		return packed
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

// packSentences returns the longest prefix of complete sentences within maxRunes.
func packSentences(s, preferredEnd string, maxRunes int) string {
	ends := sentenceEndRunes(preferredEnd)
	runes := []rune(s)
	var lastEnd int // exclusive index after a sentence ender
	for i, r := range runes {
		if !isSentenceEnd(r, ends) {
			continue
		}
		if i+1 < len(runes) && !unicode.IsSpace(runes[i+1]) {
			continue
		}
		end := i + 1
		if end > maxRunes {
			break
		}
		lastEnd = end
	}
	if lastEnd == 0 {
		return ""
	}
	return strings.TrimSpace(string(runes[:lastEnd]))
}

func firstSentence(s, preferredEnd string) string {
	packed := packSentences(s, preferredEnd, utf8.RuneCountInString(s)+1)
	if packed == "" {
		return ""
	}
	// return only the first sentence from packed
	ends := sentenceEndRunes(preferredEnd)
	runes := []rune(packed)
	for i, r := range runes {
		if !isSentenceEnd(r, ends) {
			continue
		}
		if i+1 >= len(runes) || unicode.IsSpace(runes[i+1]) {
			return strings.TrimSpace(string(runes[:i+1]))
		}
	}
	return packed
}

func sentenceEndRunes(preferredEnd string) []rune {
	var ends []rune
	seen := map[rune]bool{}
	add := func(s string) {
		if s == "" {
			return
		}
		r, _ := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			return
		}
		if !seen[r] {
			seen[r] = true
			ends = append(ends, r)
		}
	}
	add(preferredEnd)
	for _, e := range []string{".", "。", "!", "？", "?", "！"} {
		add(e)
	}
	return ends
}

func isSentenceEnd(r rune, ends []rune) bool {
	for _, e := range ends {
		if r == e {
			return true
		}
	}
	return false
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
