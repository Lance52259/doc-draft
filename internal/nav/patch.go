package nav

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var (
	summaryServiceRe = regexp.MustCompile(`^\s{2}\* \[([^\]]+)\]\(best-practices/([^)/]+)/\)\s*$`)
	summaryPracticeRe = regexp.MustCompile(`^\s{4}\* \[([^\]]+)\]\(best-practices/([^)/]+)/([^)]+\.md)\)\s*$`)
	indexListRe      = regexp.MustCompile(`^\* \[([^\]]+)\]\(([^)]+\.md)\)\s*-\s*(.+)$`)
	readmeNavRe      = regexp.MustCompile(`(?m)^### \[([^\]]+)\]\(([^)/]+)/index\.md\)\s*$`)
)

// PatchSUMMARY inserts a practice link under the service section.
// If the service section is missing, a new service block is inserted in
// alphabetical order by service directory name among sibling services.
// Existing lines are preserved; only minimal lines are added.
func PatchSUMMARY(content, service, serviceLabel, practiceSlug, practiceTitle string) (string, error) {
	service = strings.TrimSpace(service)
	practiceSlug = strings.TrimSuffix(strings.TrimSpace(practiceSlug), ".md")
	if service == "" || practiceSlug == "" {
		return "", fmt.Errorf("service and practice slug are required")
	}
	if serviceLabel == "" {
		serviceLabel = strings.ToUpper(service)
	}
	if practiceTitle == "" {
		practiceTitle = practiceSlug
	}

	content = strings.ReplaceAll(content, "\r\n", "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	lines := strings.Split(content, "\n")
	// Split keeps trailing empty from final newline; drop last empty for processing then restore
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	practicePath := fmt.Sprintf("best-practices/%s/%s.md", service, practiceSlug)
	practiceLine := fmt.Sprintf("    * [%s](%s)", practiceTitle, practicePath)
	introLine := fmt.Sprintf("    * [简介](best-practices/%s/index.md)", service)
	serviceLine := fmt.Sprintf("  * [%s](best-practices/%s/)", serviceLabel, service)

	// Already present?
	for _, line := range lines {
		if strings.Contains(line, "("+practicePath+")") {
			return ensureTrailingNewline(strings.Join(lines, "\n")), nil
		}
	}

	svcStart, svcEnd, found := findServiceBlock(lines, service)
	if found {
		insertAt := practiceInsertIndex(lines, svcStart, svcEnd, practiceSlug)
		out := make([]string, 0, len(lines)+1)
		out = append(out, lines[:insertAt]...)
		out = append(out, practiceLine)
		out = append(out, lines[insertAt:]...)
		return ensureTrailingNewline(strings.Join(out, "\n")), nil
	}

	// New service block under 最佳实践
	block := []string{serviceLine, introLine, practiceLine}
	insertAt := newServiceInsertIndex(lines, service)
	out := make([]string, 0, len(lines)+len(block))
	out = append(out, lines[:insertAt]...)
	out = append(out, block...)
	out = append(out, lines[insertAt:]...)
	return ensureTrailingNewline(strings.Join(out, "\n")), nil
}

func findServiceBlock(lines []string, service string) (start, end int, found bool) {
	marker := fmt.Sprintf("(best-practices/%s/)", service)
	start = -1
	for i, line := range lines {
		if m := summaryServiceRe.FindStringSubmatch(line); m != nil && m[2] == service {
			start = i
			break
		}
		// fallback: contain service dir link at 2-space indent
		if strings.HasPrefix(line, "  * ") && strings.Contains(line, marker) {
			start = i
			break
		}
	}
	if start < 0 {
		return 0, 0, false
	}
	end = len(lines)
	for i := start + 1; i < len(lines); i++ {
		line := lines[i]
		if summaryServiceRe.MatchString(line) {
			end = i
			break
		}
		// next top-level or another 2-space service under 最佳实践
		if strings.HasPrefix(line, "* ") {
			end = i
			break
		}
		if strings.HasPrefix(line, "  * ") && !strings.HasPrefix(line, "    * ") {
			end = i
			break
		}
	}
	return start, end, true
}

func practiceInsertIndex(lines []string, svcStart, svcEnd int, practiceSlug string) int {
	target := practiceSlug + ".md"
	// Prefer after 简介, among practice lines sorted by filename
	insertAt := svcStart + 1
	for i := svcStart + 1; i < svcEnd; i++ {
		line := lines[i]
		m := summaryPracticeRe.FindStringSubmatch(line)
		if m == nil {
			// keep intro / non-practice before practices
			if strings.Contains(line, "/index.md)") {
				insertAt = i + 1
			}
			continue
		}
		file := m[3]
		if file > target {
			return i
		}
		insertAt = i + 1
	}
	return insertAt
}

func newServiceInsertIndex(lines []string, service string) int {
	// Find 最佳实践 section services and insert by service dir name
	bestIdx := -1
	for i, line := range lines {
		if strings.Contains(line, "](best-practices/)") || strings.Contains(line, "](best-practices)") {
			bestIdx = i
			break
		}
	}
	if bestIdx < 0 {
		return len(lines)
	}
	insertAt := bestIdx + 1
	for i := bestIdx + 1; i < len(lines); i++ {
		line := lines[i]
		if strings.HasPrefix(line, "* ") {
			return i
		}
		m := summaryServiceRe.FindStringSubmatch(line)
		if m == nil {
			// skip 简介 under 最佳实践
			if strings.HasPrefix(line, "  * ") {
				insertAt = i + 1
			}
			continue
		}
		if m[2] > service {
			return i
		}
		// skip whole block — find end of this service
		_, end, ok := findServiceBlock(lines, m[2])
		if ok {
			insertAt = end
			i = end - 1
		} else {
			insertAt = i + 1
		}
	}
	return insertAt
}

// PatchServiceIndex inserts one list item under ## 最佳实践列表, sorted by filename.
func PatchServiceIndex(content, practiceSlug, practiceTitle, oneLiner string) (string, error) {
	practiceSlug = strings.TrimSuffix(strings.TrimSpace(practiceSlug), ".md")
	if practiceSlug == "" {
		return "", fmt.Errorf("practice slug required")
	}
	if practiceTitle == "" {
		practiceTitle = practiceSlug
	}
	if oneLiner == "" {
		oneLiner = "介绍如何使用Terraform完成本实践的自动化部署。"
	}
	oneLiner = strings.TrimSuffix(strings.TrimSpace(oneLiner), "。") + "。"

	content = strings.ReplaceAll(content, "\r\n", "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	link := practiceSlug + ".md"
	item := fmt.Sprintf("* [%s](%s) - %s", practiceTitle, link, oneLiner)

	if strings.Contains(content, "]("+link+")") {
		return content, nil
	}

	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	listStart := -1
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "## 最佳实践列表") {
			listStart = i
			break
		}
	}
	if listStart < 0 {
		// append section
		lines = append(lines, "", "## 最佳实践列表", "", "本章节包含以下最佳实践：", "", item)
		return ensureTrailingNewline(strings.Join(lines, "\n")), nil
	}

	listEnd := len(lines)
	for i := listStart + 1; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), "## ") {
			listEnd = i
			break
		}
	}

	// collect existing list items in range
	type entry struct {
		file string
		line string
		idx  int
	}
	var entries []entry
	firstItem, lastItem := -1, -1
	for i := listStart + 1; i < listEnd; i++ {
		m := indexListRe.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		if firstItem < 0 {
			firstItem = i
		}
		lastItem = i
		entries = append(entries, entry{file: m[2], line: lines[i], idx: i})
	}

	entries = append(entries, entry{file: link, line: item, idx: -1})
	sort.SliceStable(entries, func(i, j int) bool { return entries[i].file < entries[j].file })

	newLines := make([]string, 0, len(entries))
	for _, e := range entries {
		newLines = append(newLines, e.line)
	}

	out := make([]string, 0, len(lines)+1)
	if firstItem < 0 {
		// no items yet — insert before listEnd after blank/intro lines
		out = append(out, lines[:listEnd]...)
		if listEnd > 0 && lines[listEnd-1] != "" {
			out = append(out, "")
		}
		out = append(out, newLines...)
		out = append(out, lines[listEnd:]...)
	} else {
		out = append(out, lines[:firstItem]...)
		out = append(out, newLines...)
		out = append(out, lines[lastItem+1:]...)
	}
	return ensureTrailingNewline(strings.Join(out, "\n")), nil
}

// PatchBestPracticesREADME inserts a service navigation block sorted by link path.
func PatchBestPracticesREADME(content, service, heading, blurb string) (string, error) {
	service = strings.TrimSpace(service)
	if service == "" {
		return "", fmt.Errorf("service required")
	}
	if heading == "" {
		heading = fmt.Sprintf("%s最佳实践", strings.ToUpper(service))
	}
	if blurb == "" {
		blurb = fmt.Sprintf("%s 相关最佳实践。", strings.ToUpper(service))
	}

	content = strings.ReplaceAll(content, "\r\n", "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	link := service + "/index.md"
	if strings.Contains(content, "]("+link+")") {
		return content, nil
	}

	block := fmt.Sprintf("### [%s](%s)\n\n%s\n", heading, link, strings.TrimSpace(blurb))

	navIdx := strings.Index(content, "## 文档导航")
	if navIdx < 0 {
		return ensureTrailingNewline(content + "\n## 文档导航\n\n" + block), nil
	}

	// Split into before nav body / nav sections / after
	afterNavHeader := content[navIdx:]
	headerEnd := strings.Index(afterNavHeader, "\n")
	if headerEnd < 0 {
		return ensureTrailingNewline(content + "\n\n" + block), nil
	}
	prefix := content[:navIdx+headerEnd+1]
	rest := afterNavHeader[headerEnd+1:]

	// rest may start with blank lines then ### sections; stop at next ##
	nextH2 := regexp.MustCompile(`(?m)^## `).FindStringIndex(rest)
	navBody, suffix := rest, ""
	if nextH2 != nil {
		navBody = rest[:nextH2[0]]
		suffix = rest[nextH2[0]:]
	}

	type sec struct {
		key  string
		text string
	}
	var sections []sec
	parts := readmeNavRe.FindAllStringIndex(navBody, -1)
	if len(parts) == 0 {
		return ensureTrailingNewline(prefix + "\n" + block + "\n" + suffix), nil
	}
	for i, loc := range parts {
		end := len(navBody)
		if i+1 < len(parts) {
			end = parts[i+1][0]
		}
		chunk := strings.TrimRight(navBody[loc[0]:end], "\n") + "\n"
		m := readmeNavRe.FindStringSubmatch(navBody[loc[0]:loc[1]])
		key := ""
		if m != nil {
			key = m[2]
		}
		sections = append(sections, sec{key: key, text: chunk})
	}
	sections = append(sections, sec{key: service, text: block})
	sort.SliceStable(sections, func(i, j int) bool { return sections[i].key < sections[j].key })

	var b strings.Builder
	b.WriteString(prefix)
	if !strings.HasSuffix(prefix, "\n") {
		b.WriteByte('\n')
	}
	b.WriteByte('\n')
	for _, s := range sections {
		b.WriteString(strings.TrimRight(s.text, "\n"))
		b.WriteString("\n\n")
	}
	b.WriteString(suffix)
	return ensureTrailingNewline(b.String()), nil
}

// IsDestructiveUpdate reports whether updated content likely dropped too much of baseline.
func IsDestructiveUpdate(baseline, updated string) bool {
	baseLines := nonEmptyLines(baseline)
	upLines := nonEmptyLines(updated)
	if len(baseLines) == 0 {
		return false
	}
	if len(upLines) < len(baseLines)*8/10 { // lost >20% of lines
		return true
	}
	// key anchors for SUMMARY
	if strings.Contains(baseline, "# Summary") && !strings.Contains(updated, "# Summary") {
		return true
	}
	if strings.Contains(baseline, "best-practices/anti-ddos/") && !strings.Contains(updated, "best-practices/anti-ddos/") {
		return true
	}
	return false
}

func nonEmptyLines(s string) []string {
	var out []string
	for _, line := range strings.Split(s, "\n") {
		if strings.TrimSpace(line) != "" {
			out = append(out, line)
		}
	}
	return out
}

func ensureTrailingNewline(s string) string {
	if !strings.HasSuffix(s, "\n") {
		return s + "\n"
	}
	return s
}

// ServiceLabelFromSUMMARY finds existing display label for a service, if any.
func ServiceLabelFromSUMMARY(content, service string) string {
	for _, line := range strings.Split(content, "\n") {
		m := summaryServiceRe.FindStringSubmatch(line)
		if m != nil && m[2] == service {
			return m[1]
		}
	}
	return ""
}
