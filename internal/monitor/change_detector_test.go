package monitor_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Lance52259/doc-draft/internal/config"
	"github.com/Lance52259/doc-draft/internal/model"
	"github.com/Lance52259/doc-draft/internal/monitor"
)

func makeExamples(t *testing.T, root string) {
	t.Helper()
	_ = os.MkdirAll(filepath.Join(root, "examples", "alpha"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "examples", "alpha", "README.md"), []byte("# Alpha Practice\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(root, "examples", "beta"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "examples", "beta", "main.sh"), []byte("echo ok\n"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "examples", "README.md"), []byte("ignore me\n"), 0o644)
}

func TestEnumeratePractices(t *testing.T) {
	root := t.TempDir()
	makeExamples(t, root)
	practices, err := monitor.EnumeratePractices(filepath.Join(root, "examples"), "examples", []string{"README.md", "README", ".gitkeep"}, "directory")
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, p := range practices {
		ids[p.PracticeID] = true
	}
	if !ids["examples/alpha"] || !ids["examples/beta"] || len(practices) != 2 {
		t.Fatalf("unexpected practices: %+v", practices)
	}
	for _, p := range practices {
		if p.PracticeID == "examples/alpha" && p.Title != "Alpha Practice" {
			t.Fatalf("title=%q", p.Title)
		}
	}
}

func TestEnumerateNestedPractices(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "examples", "ecs", "basic"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "examples", "ecs", "basic", "main.tf"), []byte("resource \"huaweicloud_vpc\" \"t\" {}\n"), 0o644)
	_ = os.MkdirAll(filepath.Join(root, "examples", "vpc", "peering"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "examples", "vpc", "peering", "main.tf"), []byte(""), 0o644)

	practices, err := monitor.EnumeratePractices(filepath.Join(root, "examples"), "examples", []string{"README.md"}, "nested_directory")
	if err != nil {
		t.Fatal(err)
	}
	ids := map[string]bool{}
	for _, p := range practices {
		ids[p.PracticeID] = true
	}
	if !ids["examples/ecs/basic"] || !ids["examples/vpc/peering"] {
		t.Fatalf("got %+v", practices)
	}
}

func TestDetectHybrid(t *testing.T) {
	tmp := t.TempDir()
	b := filepath.Join(tmp, "b")
	c := filepath.Join(tmp, "c")
	makeExamples(t, b)
	_ = os.MkdirAll(filepath.Join(c, "docs", "best-practices"), 0o755)
	_ = os.WriteFile(filepath.Join(c, "docs", "best-practices", "alpha.md"), []byte("# alpha\n"), 0o644)
	manifest, _ := json.Marshal(map[string]any{"practices": []string{"examples/beta"}})
	_ = os.WriteFile(filepath.Join(c, "synced-practices.json"), manifest, 0o644)

	s := &config.Settings{
		BRepo:            "o/b",
		CRepo:            "o/c",
		SyncedStrategy:   "hybrid",
		CDocsRoot:        "docs/best-practices",
		CSyncedManifest:  "synced-practices.json",
		BExamplesPath:    "examples",
		IgnoreNames:      []string{"README.md", "README", ".gitkeep"},
		Granularity:      "directory",
	}
	ctx := &monitor.RepoContext{
		B: model.RepoRef{Name: "o/b", LocalPath: b, CommitSHA: "b1"},
		C: model.RepoRef{Name: "o/c", LocalPath: c, CommitSHA: "c1"},
	}
	result, err := (&monitor.ChangeDetector{Settings: s}).Detect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.NewPractices) != 0 {
		t.Fatalf("expected no new, got %+v", result.NewPractices)
	}
}

func TestDetectNewOnly(t *testing.T) {
	tmp := t.TempDir()
	b := filepath.Join(tmp, "b")
	c := filepath.Join(tmp, "c")
	makeExamples(t, b)
	_ = os.MkdirAll(filepath.Join(c, "docs", "best-practices"), 0o755)

	s := &config.Settings{
		BRepo:          "o/b",
		CRepo:          "o/c",
		SyncedStrategy: "path_infer",
		CDocsRoot:      "docs/best-practices",
		BExamplesPath:  "examples",
		IgnoreNames:    []string{"README.md", "README", ".gitkeep"},
		Granularity:    "directory",
	}
	ctx := &monitor.RepoContext{
		B: model.RepoRef{LocalPath: b},
		C: model.RepoRef{LocalPath: c},
	}
	result, err := (&monitor.ChangeDetector{Settings: s}).Detect(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.NewPractices) != 2 {
		t.Fatalf("got %d new", len(result.NewPractices))
	}
}
