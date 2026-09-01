package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var frontMatter = regexp.MustCompile(`(?s)^---\s*\n(.*?)\n---\s*\n(.*)$`)

// Skill is a loaded skills/<id>/SKILL.md definition.
type Skill struct {
	ID       string
	Root     string
	Title    string
	Body     string
	Metadata map[string]string
}

// SkillLoader loads Skill files from disk.
type SkillLoader struct {
	SkillsRoot string
}

func (l *SkillLoader) Load(skillID string) (*Skill, error) {
	root := filepath.Join(l.SkillsRoot, skillID)
	skillFile := filepath.Join(root, "SKILL.md")
	raw, err := os.ReadFile(skillFile)
	if err != nil {
		return nil, fmt.Errorf("skill not found: %s: %w", skillFile, err)
	}
	meta, body := parseFrontMatter(string(raw))
	title := meta["name"]
	if title == "" {
		title = meta["title"]
	}
	if title == "" {
		title = skillID
	}
	return &Skill{ID: skillID, Root: root, Title: title, Body: body, Metadata: meta}, nil
}

func parseFrontMatter(text string) (map[string]string, string) {
	m := frontMatter.FindStringSubmatch(text)
	if m == nil {
		return map[string]string{}, text
	}
	meta := map[string]string{}
	for _, line := range strings.Split(m[1], "\n") {
		if !strings.Contains(line, ":") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		val = strings.Trim(val, `"'`)
		meta[key] = val
	}
	return meta, strings.TrimSpace(m[2])
}

// ReadReferences loads optional references under the skill.
func (s *Skill) ReadReferences(limit int) ([][2]string, error) {
	refDir := filepath.Join(s.Root, "references")
	entries, err := os.ReadDir(refDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var out [][2]string
	for _, e := range entries {
		if e.IsDir() || len(out) >= limit {
			continue
		}
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".md" && ext != ".txt" && ext != ".yaml" && ext != ".yml" {
			continue
		}
		path := filepath.Join(refDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		rel := filepath.ToSlash(filepath.Join("references", name))
		out = append(out, [2]string{rel, string(data)})
	}
	return out, nil
}
