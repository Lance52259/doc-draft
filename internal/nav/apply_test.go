package nav_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Lance52259/doc-draft/internal/model"
	"github.com/Lance52259/doc-draft/internal/nav"
)

func writeTree(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		_ = os.MkdirAll(filepath.Dir(path), 0o755)
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestApplyToFilesExistingServiceBilingual(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"docs/zh-cn/SUMMARY.md": `# Summary

* [最佳实践](best-practices/)
  * [简介](best-practices/README.md)
  * [ECS](best-practices/ecs/)
    * [简介](best-practices/ecs/index.md)
    * [部署基础实例](best-practices/ecs/simple_instance.md)
`,
		"docs/en-us/SUMMARY.md": `# Summary

* [Best Practices](best-practices/)
  * [Introduction](best-practices/README.md)
  * [ECS](best-practices/ecs/)
    * [Introduction](best-practices/ecs/index.md)
    * [Deploy Basic Instance](best-practices/ecs/simple_instance.md)
`,
		"docs/zh-cn/best-practices/ecs/index.md": `# 简介

## 最佳实践列表

* [部署基础实例](simple_instance.md) - 基础实例。
`,
		"docs/en-us/best-practices/ecs/index.md": `# Introduction

## Best Practices List

* [Deploy Basic Instance](simple_instance.md) - Basic.
`,
		"docs/zh-cn/best-practices/README.md": "# 中心\n\n## 文档导航\n\n### [ECS最佳实践](ecs/index.md)\n\nECS。\n",
		"docs/en-us/best-practices/README.md": "# Center\n\n## Documentation Navigation\n\n### [ECS Best Practices](ecs/index.md)\n\nECS.\n",
	})

	files := []model.DocFileChange{
		{Path: "docs/zh-cn/best-practices/ecs/prepaid_instance.md", Action: "create", Content: "# 部署包周期实例\n"},
		{Path: "docs/en-us/best-practices/ecs/prepaid_instance.md", Action: "create", Content: "# Deploy PrePaid Instance\n"},
		{Path: "docs/zh-cn/SUMMARY.md", Action: "update", Content: "# 目录\n"},
		{Path: "docs/en-us/SUMMARY.md", Action: "update", Content: "# wiped\n"},
	}
	out, err := nav.ApplyToFiles(files, nav.ApplyOptions{
		CRepoRoot: root,
		Service:   "ecs",
		Slug:      "prepaid_instance",
	})
	if err != nil {
		t.Fatal(err)
	}

	got := map[string]string{}
	for _, f := range out {
		got[f.Path] = f.Content
	}
	if _, ok := got["docs/zh-cn/best-practices/README.md"]; ok {
		t.Fatal("existing service must not touch zh README")
	}
	if _, ok := got["docs/en-us/best-practices/README.md"]; ok {
		t.Fatal("existing service must not touch en README")
	}
	for _, p := range []string{"docs/zh-cn/SUMMARY.md", "docs/en-us/SUMMARY.md"} {
		if !strings.Contains(got[p], "simple_instance.md") || !strings.Contains(got[p], "prepaid_instance.md") {
			t.Fatalf("%s:\n%s", p, got[p])
		}
		if strings.Contains(got[p], "# 目录") || strings.Contains(got[p], "# wiped") {
			t.Fatalf("destructive %s", p)
		}
	}
	if !strings.Contains(got["docs/en-us/best-practices/ecs/index.md"], "prepaid_instance.md") {
		t.Fatalf("en index: %s", got["docs/en-us/best-practices/ecs/index.md"])
	}
	if !strings.Contains(got["docs/zh-cn/best-practices/ecs/index.md"], "prepaid_instance.md") {
		t.Fatalf("zh index: %s", got["docs/zh-cn/best-practices/ecs/index.md"])
	}
}

func TestApplyToFilesNewServiceBilingual(t *testing.T) {
	root := t.TempDir()
	writeTree(t, root, map[string]string{
		"docs/zh-cn/SUMMARY.md": `# Summary

* [最佳实践](best-practices/)
  * [简介](best-practices/README.md)
  * [ECS](best-practices/ecs/)
    * [简介](best-practices/ecs/index.md)
`,
		"docs/en-us/SUMMARY.md": `# Summary

* [Best Practices](best-practices/)
  * [Introduction](best-practices/README.md)
  * [ECS](best-practices/ecs/)
    * [Introduction](best-practices/ecs/index.md)
`,
		"docs/zh-cn/best-practices/README.md": "# 中心\n\n## 文档导航\n\n### [ECS最佳实践](ecs/index.md)\n\nECS。\n",
		"docs/en-us/best-practices/README.md": "# Center\n\n## Documentation Navigation\n\n### [ECS Best Practices](ecs/index.md)\n\nECS.\n",
	})

	files := []model.DocFileChange{
		{Path: "docs/zh-cn/best-practices/aad/black_white_lists.md", Action: "create", Content: "# 部署黑白名单防护\n"},
		{Path: "docs/en-us/best-practices/aad/black_white_lists.md", Action: "create", Content: "# Deploy Black White Lists\n"},
		{Path: "docs/zh-cn/best-practices/aad/index.md", Action: "create", Content: "# 简介\n\n## 最佳实践列表\n\n"},
		{Path: "docs/en-us/best-practices/aad/index.md", Action: "create", Content: "# Introduction\n\n## Best Practices List\n\n"},
	}
	out, err := nav.ApplyToFiles(files, nav.ApplyOptions{
		CRepoRoot: root,
		Service:   "aad",
		Slug:      "black_white_lists",
	})
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, f := range out {
		got[f.Path] = f.Content
	}
	for _, p := range []string{"docs/en-us/SUMMARY.md", "docs/zh-cn/SUMMARY.md"} {
		if !strings.Contains(got[p], "best-practices/aad/") || !strings.Contains(got[p], "black_white_lists.md") {
			t.Fatalf("%s:\n%s", p, got[p])
		}
	}
	if !strings.Contains(got["docs/en-us/best-practices/README.md"], "aad/index.md") {
		t.Fatalf("en readme: %s", got["docs/en-us/best-practices/README.md"])
	}
	if !strings.Contains(got["docs/zh-cn/best-practices/README.md"], "aad/index.md") {
		t.Fatalf("zh readme: %s", got["docs/zh-cn/best-practices/README.md"])
	}
}

func TestApplyRequiresBothSUMMARY(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "docs", "zh-cn"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "docs", "zh-cn", "SUMMARY.md"), []byte("# Summary\n"), 0o644)
	_, err := nav.ApplyToFiles(nil, nav.ApplyOptions{CRepoRoot: root, Service: "ecs", Slug: "x"})
	if err == nil || !strings.Contains(err.Error(), "docs/en-us/SUMMARY.md") {
		t.Fatalf("err=%v", err)
	}
}
