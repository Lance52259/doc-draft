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
  "summary": "<one-sentence English summary for the PR description>",
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
3. 不要在 JSON 外输出任何文字
4. 优先只输出实践正文的 create；SUMMARY.md / 已有 index.md / README.md 导航由编排层手术式补丁，除非你能在提供的基线上仅插入新行且不删改任何已有条目
5. 严禁把 SUMMARY.md 重写成短目录或改掉「# Summary」标题`

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

// BuildMessagesInput carries generation prompt inputs.
type BuildMessagesInput struct {
	Skill         *Skill
	Practice      model.Practice
	SourceContext string
	TargetPath    string
	TemplateName  string
	DocsRoot      string
	// NavBaselines maps relative path → current C-repo file content.
	NavBaselines map[string]string
}

func (p *PromptTemplates) BuildMessages(in BuildMessagesInput) ([]provider.ChatMessage, error) {
	template, err := p.LoadTemplate(in.TemplateName)
	if err != nil {
		return nil, err
	}
	refs, _ := in.Skill.ReadReferences(8)
	var refParts []string
	for _, r := range refs {
		refParts = append(refParts, "### reference:"+r[0]+"\n"+r[1])
	}
	refText := strings.Join(refParts, "\n\n")
	if refText == "" {
		refText = "(无)"
	}

	baselineText := formatNavBaselines(in.NavBaselines)

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

## C 仓导航基线（必须保留全部已有内容；禁止整文件重写）
%s

## 源最佳实践上下文
%s
`, in.Skill.Title, in.Skill.Body, in.Practice.PracticeID, in.TargetPath, in.DocsRoot, template, refText, baselineText, in.SourceContext)

	return []provider.ChatMessage{
		{Role: "system", Content: system},
		{Role: "user", Content: user},
	}, nil
}

func formatNavBaselines(baselines map[string]string) string {
	if len(baselines) == 0 {
		return "(未提供基线：请勿输出 SUMMARY.md / README.md / 已有 index.md 的 update；仅输出实践正文 create。)"
	}
	var keys []string
	for k := range baselines {
		keys = append(keys, k)
	}
	// stable-ish order
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[j] < keys[i] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	var b strings.Builder
	b.WriteString("以下为 C 仓**当前真实文件**。若你输出对应 path 的 update，必须以这些全文为起点，只插入本实践相关新行，不得删除任何已有行。\n")
	const maxPerFile = 120000
	for _, k := range keys {
		content := baselines[k]
		if len(content) > maxPerFile {
			content = content[:maxPerFile] + "\n…(truncated)…\n"
		}
		b.WriteString("\n### baseline:")
		b.WriteString(k)
		b.WriteString("\n```markdown\n")
		b.WriteString(content)
		if !strings.HasSuffix(content, "\n") {
			b.WriteByte('\n')
		}
		b.WriteString("```\n")
	}
	return b.String()
}
