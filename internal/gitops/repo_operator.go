package gitops

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Lance52259/doc-draft/internal/config"
	"github.com/Lance52259/doc-draft/internal/model"
	"github.com/Lance52259/doc-draft/internal/monitor"
)

// RepoOperator applies generated files and pushes a branch.
type RepoOperator struct {
	Settings *config.Settings
}

func (o *RepoOperator) ApplyAndPush(repoPath, branchName, baseBranch string, result *model.GenerateResult, commitMessage string, dryRun bool) (string, error) {
	if err := run(repoPath, "git", "fetch", "origin"); err != nil {
		return "", err
	}
	if err := run(repoPath, "git", "checkout", baseBranch); err != nil {
		return "", err
	}
	_ = run(repoPath, "git", "pull", "origin", baseBranch)

	_ = run(repoPath, "git", "branch", "-D", branchName)
	if err := run(repoPath, "git", "checkout", "-b", branchName); err != nil {
		return "", err
	}

	for _, change := range result.Files {
		target := filepath.Join(repoPath, filepath.FromSlash(change.Path))
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return "", err
		}
		if err := os.WriteFile(target, []byte(change.Content), 0o644); err != nil {
			return "", err
		}
	}

	if err := run(repoPath, "git", "add", "-A"); err != nil {
		return "", err
	}
	dirty, err := runOut(repoPath, "git", "status", "--porcelain")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(dirty) == "" {
		return "", nil
	}

	env := append(os.Environ(),
		"GIT_AUTHOR_NAME=doc-craft",
		"GIT_AUTHOR_EMAIL=doc-craft@users.noreply.github.com",
		"GIT_COMMITTER_NAME=doc-craft",
		"GIT_COMMITTER_EMAIL=doc-craft@users.noreply.github.com",
	)
	cmd := exec.Command("git", "commit", "-m", commitMessage)
	cmd.Dir = repoPath
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	if dryRun {
		return branchName, nil
	}

	if token := o.Settings.CRepoToken; token != "" {
		url := monitor.ToHTTPSURL(o.Settings.CRepo, token)
		if err := run(repoPath, "git", "remote", "set-url", "origin", url); err != nil {
			return "", err
		}
	}
	if err := run(repoPath, "git", "push", "--set-upstream", "origin", branchName, "--force"); err != nil {
		return "", err
	}
	return branchName, nil
}

func run(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runOut(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	return string(out), err
}

// SummarizeChanges builds a PR body.
func SummarizeChanges(result *model.GenerateResult, bRepo, bSHA string) string {
	if bSHA == "" {
		bSHA = "unknown"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "## Summary\n")
	fmt.Fprintf(&b, "- 来源：`%s`@%s / `%s`\n", bRepo, bSHA, result.PracticeID)
	fmt.Fprintf(&b, "- 类型：最佳实践文档同步（doc-craft / DeepSeek 自动生成）\n\n")
	fmt.Fprintf(&b, "## 变更说明\n")
	summary := result.Summary
	if summary == "" {
		summary = "(无摘要)"
	}
	fmt.Fprintf(&b, "%s\n\n## 文件\n", summary)
	for _, f := range result.Files {
		fmt.Fprintf(&b, "- `%s` `%s`\n", f.Action, f.Path)
	}
	fmt.Fprintf(&b, "\n## 人工检查清单\n")
	fmt.Fprintf(&b, "- [ ] 路径与 C 仓信息架构一致\n")
	fmt.Fprintf(&b, "- [ ] 代码块/命令可运行或已标明环境\n")
	fmt.Fprintf(&b, "- [ ] 无密钥或敏感信息\n")
	return b.String()
}
