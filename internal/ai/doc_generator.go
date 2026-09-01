package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Lance52259/doc-draft/internal/ai/provider"
	"github.com/Lance52259/doc-draft/internal/config"
	"github.com/Lance52259/doc-draft/internal/mapping"
	"github.com/Lance52259/doc-draft/internal/model"
)

var jsonBlock = regexp.MustCompile("(?s)```(?:json)?\\s*(\\{.*?\\})\\s*```")

// DocGenerator orchestrates Skill + DeepSeek generation.
type DocGenerator struct {
	Settings *config.Settings
	Provider provider.Provider
	Skills   *SkillLoader
	Prompts  *PromptTemplates
}

func NewDocGenerator(settings *config.Settings, p provider.Provider) *DocGenerator {
	return &DocGenerator{
		Settings: settings,
		Provider: p,
		Skills:   &SkillLoader{SkillsRoot: filepath.Join(settings.RepoRoot, settings.SkillRoot)},
		Prompts:  &PromptTemplates{TemplatesDir: settings.TemplatesDir()},
	}
}

// ResolveTargetPath returns target path, skill id, template name.
func ResolveTargetPath(settings *config.Settings, practice model.Practice) (target, skillID, template string) {
	defaults := settings.Mapping.Defaults
	skillID = defaults.SkillID
	if skillID == "" {
		skillID = settings.SkillID
	}
	template = defaults.Template
	if template == "" {
		template = "best_practice_template.md"
	}
	pattern := defaults.TargetPathPattern
	if pattern == "" {
		pattern = "docs/zh-cn/best-practices/{service}/{practice_slug}.md"
	}
	for _, rule := range settings.Mapping.Rules {
		match := strings.TrimRight(rule.Match, "*")
		match = strings.TrimRight(match, "/")
		if match != "" && strings.HasPrefix(practice.PracticeID, match) {
			if rule.SkillID != "" {
				skillID = rule.SkillID
			}
			if rule.Template != "" {
				template = rule.Template
			}
			if rule.TargetPathPattern != "" {
				pattern = rule.TargetPathPattern
			}
			break
		}
	}
	resolver := mapping.NewResolver(settings.Mapping, settings.CDocsRoot)
	doc := resolver.Resolve(practice)
	service := doc.Service
	slug := doc.Slug
	if service == "" {
		service = practice.Slug()
	}
	if slug == "" {
		slug = practice.Slug()
	}
	target = strings.ReplaceAll(pattern, "{practice_slug}", slug)
	target = strings.ReplaceAll(target, "{service}", service)
	target = strings.ReplaceAll(target, "{practice_id}", practice.PracticeID)
	return target, skillID, template
}

// PackSourceContext packs example files into a prompt budget.
func PackSourceContext(practiceDir string, maxChars int) (string, error) {
	st, err := os.Stat(practiceDir)
	if err != nil || !st.IsDir() {
		return fmt.Sprintf("(missing directory: %s)", practiceDir), nil
	}

	priority := map[string]int{"README.md": 0, "README": 0, "readme.md": 0, "SKILL.md": 0, "doc.md": 0}
	type item struct {
		path string
		rel  string
		rank int
	}
	var files []item
	_ = filepath.WalkDir(practiceDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		ext := strings.ToLower(filepath.Ext(d.Name()))
		switch ext {
		case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".zip", ".tar", ".gz":
			return nil
		}
		rel, _ := filepath.Rel(practiceDir, path)
		rank := 1
		if r, ok := priority[d.Name()]; ok {
			rank = r
		}
		files = append(files, item{path: path, rel: filepath.ToSlash(rel), rank: rank})
		return nil
	})

	// simple sort: rank, then path
	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			if files[j].rank < files[i].rank || (files[j].rank == files[i].rank && files[j].rel < files[i].rel) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}

	var parts []string
	used := 0
	for _, f := range files {
		data, err := os.ReadFile(f.path)
		if err != nil {
			continue
		}
		chunk := "### file:" + f.rel + "\n" + string(data) + "\n"
		if used+len(chunk) > maxChars {
			remain := maxChars - used
			if remain < 64 {
				break
			}
			chunk = chunk[:remain] + "\n...[truncated]...\n"
			parts = append(parts, chunk)
			break
		}
		parts = append(parts, chunk)
		used += len(chunk)
	}
	if len(parts) == 0 {
		return "(empty practice directory)", nil
	}
	return strings.Join(parts, "\n"), nil
}

func ExtractJSON(text string) (map[string]any, error) {
	text = strings.TrimSpace(text)
	var out map[string]any
	if err := json.Unmarshal([]byte(text), &out); err == nil {
		return out, nil
	}
	if m := jsonBlock.FindStringSubmatch(text); m != nil {
		if err := json.Unmarshal([]byte(m[1]), &out); err == nil {
			return out, nil
		}
	}
	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start >= 0 && end > start {
		if err := json.Unmarshal([]byte(text[start:end+1]), &out); err == nil {
			return out, nil
		}
	}
	return nil, fmt.Errorf("unable to parse JSON from model output")
}

func ValidatePaths(files []model.DocFileChange, allowlist []string) error {
	for _, f := range files {
		path := filepath.ToSlash(strings.TrimLeft(f.Path, "/"))
		if path == "" || strings.Contains(path, "..") {
			return fmt.Errorf("unsafe path: %q", f.Path)
		}
		ok := false
		for _, prefix := range allowlist {
			if strings.HasPrefix(path, prefix) {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("path not in allowlist: %s", path)
		}
	}
	return nil
}

func (g *DocGenerator) Generate(ctx context.Context, practice model.Practice, practiceDir string) (*model.GenerateResult, error) {
	targetPath, skillID, templateName := ResolveTargetPath(g.Settings, practice)
	skill, err := g.Skills.Load(skillID)
	if err != nil {
		return nil, err
	}
	sourceContext, err := PackSourceContext(practiceDir, g.Settings.MaxContextChars)
	if err != nil {
		return nil, err
	}
	messages, err := g.Prompts.BuildMessages(skill, practice, sourceContext, targetPath, templateName, g.Settings.CDocsRoot)
	if err != nil {
		return nil, err
	}

	var lastErr error
	var raw string
	for attempt := 0; attempt <= g.Settings.AIMaxRetries; attempt++ {
		completion, err := g.Provider.Complete(ctx, messages, 0.2, g.Settings.ResponseFormatJSON)
		if err != nil {
			lastErr = err
			continue
		}
		raw = completion.Content
		data, err := ExtractJSON(raw)
		if err != nil {
			lastErr = err
			continue
		}
		files, err := parseFiles(data, targetPath)
		if err != nil {
			lastErr = err
			continue
		}
		if err := ValidatePaths(files, g.Settings.PathAllowlist); err != nil {
			lastErr = err
			continue
		}
		summary, _ := data["summary"].(string)
		return &model.GenerateResult{
			PracticeID:  practice.PracticeID,
			Files:       files,
			Summary:     summary,
			RawResponse: raw,
		}, nil
	}
	return nil, fmt.Errorf("doc generation failed for %s: %w", practice.PracticeID, lastErr)
}

func parseFiles(data map[string]any, fallbackPath string) ([]model.DocFileChange, error) {
	rawFiles, ok := data["files"].([]any)
	if ok && len(rawFiles) > 0 {
		var files []model.DocFileChange
		for _, item := range rawFiles {
			m, ok := item.(map[string]any)
			if !ok {
				continue
			}
			path, _ := m["path"].(string)
			content, _ := m["content"].(string)
			action, _ := m["action"].(string)
			if action == "" {
				action = "create"
			}
			if path == "" || content == "" {
				continue
			}
			files = append(files, model.DocFileChange{Path: filepath.ToSlash(path), Content: content, Action: action})
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("AI JSON missing usable files[]")
		}
		return files, nil
	}
	if content, ok := data["content"].(string); ok && content != "" {
		return []model.DocFileChange{{Path: fallbackPath, Content: content, Action: "create"}}, nil
	}
	return nil, fmt.Errorf("AI JSON missing files[]")
}
