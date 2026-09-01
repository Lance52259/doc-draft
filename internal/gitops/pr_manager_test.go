package gitops_test

import (
	"testing"

	"github.com/Lance52259/doc-draft/internal/gitops"
	"github.com/Lance52259/doc-draft/internal/model"
)

func TestParseOwnerRepo(t *testing.T) {
	o, r, err := gitops.ParseOwnerRepo("acme/docs")
	if err != nil || o != "acme" || r != "docs" {
		t.Fatalf("%s %s %v", o, r, err)
	}
	o, r, err = gitops.ParseOwnerRepo("https://github.com/acme/docs.git")
	if err != nil || o != "acme" || r != "docs" {
		t.Fatalf("%s %s %v", o, r, err)
	}
}

func TestSummarizeChanges(t *testing.T) {
	result := &model.GenerateResult{
		PracticeID: "examples/foo",
		Summary:    "新增 foo 文档",
		Files: []model.DocFileChange{
			{Path: "docs/best-practices/foo.md", Content: "# foo\n", Action: "create"},
		},
	}
	body := gitops.SummarizeChanges(result, "acme/examples", "abc")
	if !contains(body, "examples/foo") || !contains(body, "docs/best-practices/foo.md") {
		t.Fatalf("%s", body)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (func() bool {
		for i := 0; i+len(sub) <= len(s); i++ {
			if s[i:i+len(sub)] == sub {
				return true
			}
		}
		return false
	})()
}
