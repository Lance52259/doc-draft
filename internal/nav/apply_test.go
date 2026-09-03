package nav_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Lance52259/doc-draft/internal/model"
	"github.com/Lance52259/doc-draft/internal/nav"
)

func TestApplyToFilesExistingService(t *testing.T) {
	root := t.TempDir()
	summaryRel := "docs/zh-cn/SUMMARY.md"
	docs := filepath.Join(root, "docs", "zh-cn", "best-practices", "ecs")
	_ = os.MkdirAll(docs, 0o755)
	base := `# Summary

* [最佳实践](best-practices/)
  * [简介](best-practices/README.md)
  * [ECS](best-practices/ecs/)
    * [简介](best-practices/ecs/index.md)
    * [部署基础实例](best-practices/ecs/simple_instance.md)
`
	_ = os.WriteFile(filepath.Join(root, filepath.FromSlash(summaryRel)), []byte(base), 0o644)
	_ = os.WriteFile(filepath.Join(docs, "index.md"), []byte(`# 简介

## 最佳实践列表

* [部署基础实例](simple_instance.md) - 基础实例。
`), 0o644)
	_ = os.WriteFile(filepath.Join(root, "docs", "zh-cn", "best-practices", "README.md"), []byte("# 中心\n\n## 文档导航\n\n### [ECS最佳实践](ecs/index.md)\n\nECS。\n"), 0o644)

	files := []model.DocFileChange{
		{Path: "docs/zh-cn/best-practices/ecs/prepaid_instance.md", Action: "create", Content: "# 部署包周期实例\n\n## 应用场景\n"},
		{Path: summaryRel, Action: "update", Content: "# 目录\n\n* [ECS](best-practices/ecs/)\n"},
		{Path: "docs/zh-cn/best-practices/ecs/index.md", Action: "update", Content: "# wiped\n"},
	}
	out, err := nav.ApplyToFiles(files, nav.ApplyOptions{
		CRepoRoot: root,
		DocsRoot:  "docs/zh-cn/best-practices",
		Service:   "ecs",
		Slug:      "prepaid_instance",
	})
	if err != nil {
		t.Fatal(err)
	}

	var summary, index, readme string
	for _, f := range out {
		switch f.Path {
		case summaryRel:
			summary = f.Content
		case "docs/zh-cn/best-practices/ecs/index.md":
			index = f.Content
		case "docs/zh-cn/best-practices/README.md":
			readme = f.Content
		}
	}
	if readme != "" {
		t.Fatal("existing service must not touch README")
	}
	if strings.Contains(summary, "# 目录") || !strings.Contains(summary, "simple_instance.md") {
		t.Fatalf("summary broken:\n%s", summary)
	}
	if !strings.Contains(summary, "prepaid_instance.md") {
		t.Fatalf("summary missing practice:\n%s", summary)
	}
	if !strings.Contains(index, "simple_instance.md") || !strings.Contains(index, "prepaid_instance.md") {
		t.Fatalf("index:\n%s", index)
	}
}

func TestApplyToFilesNewService(t *testing.T) {
	root := t.TempDir()
	summaryRel := "docs/zh-cn/SUMMARY.md"
	_ = os.MkdirAll(filepath.Join(root, "docs", "zh-cn", "best-practices"), 0o755)
	base := `# Summary

* [最佳实践](best-practices/)
  * [简介](best-practices/README.md)
  * [ECS](best-practices/ecs/)
    * [简介](best-practices/ecs/index.md)
    * [部署基础实例](best-practices/ecs/simple_instance.md)
`
	_ = os.WriteFile(filepath.Join(root, filepath.FromSlash(summaryRel)), []byte(base), 0o644)
	_ = os.WriteFile(filepath.Join(root, "docs", "zh-cn", "best-practices", "README.md"), []byte(`# 中心

## 文档导航

### [ECS最佳实践](ecs/index.md)

ECS。
`), 0o644)

	files := []model.DocFileChange{
		{Path: "docs/zh-cn/best-practices/aad/black_white_lists.md", Action: "create", Content: "# 部署黑白名单防护\n"},
		{Path: "docs/zh-cn/best-practices/aad/index.md", Action: "create", Content: "# 简介\n\n## 最佳实践列表\n\n"},
		{Path: summaryRel, Action: "update", Content: "# 目录\n"},
	}
	out, err := nav.ApplyToFiles(files, nav.ApplyOptions{
		CRepoRoot: root,
		DocsRoot:  "docs/zh-cn/best-practices",
		Service:   "aad",
		Slug:      "black_white_lists",
	})
	if err != nil {
		t.Fatal(err)
	}
	var summary, index, readme string
	for _, f := range out {
		switch f.Path {
		case summaryRel:
			summary = f.Content
		case "docs/zh-cn/best-practices/aad/index.md":
			index = f.Content
		case "docs/zh-cn/best-practices/README.md":
			readme = f.Content
		}
	}
	if !strings.Contains(summary, "# Summary") || !strings.Contains(summary, "ecs/simple_instance.md") {
		t.Fatalf("summary:\n%s", summary)
	}
	if !strings.Contains(summary, "best-practices/aad/") {
		t.Fatalf("missing aad service:\n%s", summary)
	}
	if index == "" || !strings.Contains(index, "black_white_lists.md") {
		t.Fatalf("index create:\n%s", index)
	}
	if !strings.Contains(readme, "aad/index.md") || !strings.Contains(readme, "ecs/index.md") {
		t.Fatalf("readme:\n%s", readme)
	}
}

func TestApplyRequiresSUMMARY(t *testing.T) {
	root := t.TempDir()
	_, err := nav.ApplyToFiles(nil, nav.ApplyOptions{
		CRepoRoot: root,
		DocsRoot:  "docs/zh-cn/best-practices",
		Service:   "ecs",
		Slug:      "x",
	})
	if err == nil || !strings.Contains(err.Error(), "SUMMARY.md must exist") {
		t.Fatalf("err=%v", err)
	}
}
