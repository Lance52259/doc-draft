package gitops_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Lance52259/doc-draft/internal/config"
	"github.com/Lance52259/doc-draft/internal/gitops"
	"github.com/Lance52259/doc-draft/internal/model"
)

func TestPracticeBranch(t *testing.T) {
	got := gitops.PracticeBranch("examples/antiddos/basic")
	if got != "doc-craft/examples-antiddos-basic" {
		t.Fatalf("got %q", got)
	}
}

func TestFilterOpenPRs(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/Lance52259/hcbp-demo/pulls", func(w http.ResponseWriter, r *http.Request) {
		head := r.URL.Query().Get("head")
		if strings.Contains(head, "examples-ecs-basic") {
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"number":   42,
				"html_url": "https://github.com/Lance52259/hcbp-demo/pull/42",
				"title":    "docs(ecs): support new best practice for basic",
			}})
			return
		}
		_ = json.NewEncoder(w).Encode([]any{})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	prm, err := gitops.NewPRManager(&config.Settings{
		CRepo:      "Lance52259/hcbp-demo",
		CRepoToken: "test-token",
	})
	if err != nil {
		t.Fatal(err)
	}
	prm.APIBase = srv.URL
	prm.Client = srv.Client()

	practices := []model.Practice{
		{PracticeID: "examples/ecs/basic"},
		{PracticeID: "examples/vpc/basic"},
	}
	keep, skipped, err := prm.FilterOpenPRs(practices)
	if err != nil {
		t.Fatal(err)
	}
	if len(keep) != 1 || keep[0].PracticeID != "examples/vpc/basic" {
		t.Fatalf("keep=%+v", keep)
	}
	if len(skipped) != 1 || !strings.Contains(skipped[0], "examples/ecs/basic") || !strings.Contains(skipped[0], "#42") {
		t.Fatalf("skipped=%v", skipped)
	}
}
