package ai_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Lance52259/doc-draft/internal/ai"
	"github.com/Lance52259/doc-draft/internal/config"
	"github.com/Lance52259/doc-draft/internal/model"
)

func TestExtractJSON(t *testing.T) {
	m, err := ai.ExtractJSON(`{"a": 1}`)
	if err != nil || int(m["a"].(float64)) != 1 {
		t.Fatalf("%v %v", m, err)
	}
	m, err = ai.ExtractJSON("```json\n{\"a\": 2}\n```")
	if err != nil || int(m["a"].(float64)) != 2 {
		t.Fatalf("%v %v", m, err)
	}
}

func TestValidatePaths(t *testing.T) {
	err := ai.ValidatePaths([]model.DocFileChange{{Path: "docs/best-practices/x.md", Content: "hi", Action: "create"}}, []string{"docs/"})
	if err != nil {
		t.Fatal(err)
	}
	err = ai.ValidatePaths([]model.DocFileChange{{Path: "etc/passwd", Content: "x", Action: "create"}}, []string{"docs/"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveTargetPath(t *testing.T) {
	s := &config.Settings{
		SkillID: "best-practice-doc",
		Mapping: config.MappingConfig{Defaults: config.MappingDefaults{
			TargetPathPattern: "docs/zh-cn/best-practices/{service}/{practice_slug}.md",
			SkillID:           "best-practice-doc",
			Template:          "best_practice_template.md",
		}},
	}
	p := model.Practice{PracticeID: "examples/ecs/foo-bar", SourcePath: "examples/ecs/foo-bar"}
	path, skill, template := ai.ResolveTargetPath(s, p)
	if path != "docs/zh-cn/best-practices/ecs/foo-bar.md" || skill != "best-practice-doc" || template == "" {
		t.Fatalf("%s %s %s", path, skill, template)
	}
}

func TestPackSourceContext(t *testing.T) {
	d := filepath.Join(t.TempDir(), "examples", "foo")
	_ = os.MkdirAll(d, 0o755)
	_ = os.WriteFile(filepath.Join(d, "README.md"), []byte("# Hello\n"), 0o644)
	_ = os.WriteFile(filepath.Join(d, "run.sh"), []byte("echo 1\n"), 0o644)
	text, err := ai.PackSourceContext(d, 10000)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(text, "file:README.md") || !contains(text, "file:run.sh") {
		t.Fatalf("%s", text)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || (len(s) > 0 && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()))
}
