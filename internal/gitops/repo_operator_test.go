package gitops_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Lance52259/doc-draft/internal/config"
	"github.com/Lance52259/doc-draft/internal/gitops"
	"github.com/Lance52259/doc-draft/internal/model"
)

func TestApplyAndPushCreatesBranchWithoutDeleteNoise(t *testing.T) {
	dir := t.TempDir()
	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v (%s)", args, err, out)
		}
	}
	run("init", "-b", "master")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("base\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	run("add", ".")
	run("commit", "-m", "init")
	// bare remote for pull/fetch
	remote := t.TempDir()
	cmd := exec.Command("git", "init", "--bare", remote)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("bare: %v %s", err, out)
	}
	run("remote", "add", "origin", remote)
	run("push", "-u", "origin", "master")

	op := &gitops.RepoOperator{Settings: &config.Settings{}}
	result := &model.GenerateResult{
		PracticeID: "examples/ecs/basic",
		Files: []model.DocFileChange{
			{Path: "docs/zh-cn/best-practices/ecs/basic.md", Content: "# ok\n", Action: "create"},
		},
	}
	branch, err := op.ApplyAndPush(dir, "doc-craft/examples-ecs-basic", "master", result, "docs: test", true)
	if err != nil {
		t.Fatal(err)
	}
	if branch != "doc-craft/examples-ecs-basic" {
		t.Fatalf("branch=%s", branch)
	}
}
