package ai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lance52259/doc-draft/internal/ai/provider"
	"github.com/Lance52259/doc-draft/internal/model"
)

const outputSchemaHint = `你必须只输出一个 JSON 对象（不要 Markdown 围栏），格式如下：
{
  "practice_id": "<与输入相同的 practice_id>",
  "summary": "<一句话中文摘要，用于 PR 说明>",
  "files": [
    {
      "path": "<相对 C 仓库根目录的路径>",
      "action": "create",
      "content": "<完整 Markdown 文件内容>"
    }
  ]
}
约束：
1. path 不得包含 ..，且必须落在允许的文档目录下
2. content 为完整文档正文
3. 不要在 JSON 外输出任何文字`

// PromptTemplates assembles chat messages for generation.
type PromptTemplates struct {
	TemplatesDir string
}

func (p *PromptTemplates) LoadTemplate(name string) (string, error) {
	path := filepath.Join(p.TemplatesDir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("template not found: %s: %w", path, err)
	}
	return string(data), nil
}

func (p *PromptTemplates) BuildMessages(skill *Skill, practice model.Practice, sourceContext, targetPath, templateName, docsRoot string) ([]provider.ChatMessage, error) {
	template, err := p.LoadTemplate(templateName)
	if err != nil {
		return nil, err
	}
	refs, _ := skill.ReadReferences(8)
	var refParts []string
	for _, r := range refs {
		refParts = append(refParts, "### reference:"+r[0]+"\n"+r[1])
	}
	refText := strings.Join(refParts, "\n\n")
	if refText == "" {
		refText = "(无)"
	}

	system := "你是文档工程师，根据最佳实践源码与 Skill 约束，为文档仓库生成变更。\n" + outputSchemaHint
	user := fmt.Sprintf(`## Skill: %s
%s

## 目标
- practice_id: %s
- 建议写入路径: %s
- C 仓文档根: %s

## 文档模板（请在结构上对齐，可按源内容充实）
%s

## 参考资料
%s

## 源最佳实践上下文
%s
`, skill.Title, skill.Body, practice.PracticeID, targetPath, docsRoot, template, refText, sourceContext)

	return []provider.ChatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: user},
	}, nil
}
