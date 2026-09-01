package gitops

import (
	"fmt"
	"strings"

	"github.com/Lance52259/doc-draft/internal/model"
)

// PracticeBranch returns the doc-craft head branch for a practice_id.
// Example: examples/antiddos/basic → doc-craft/examples-antiddos-basic
func PracticeBranch(practiceID string) string {
	slug := strings.ReplaceAll(practiceID, "/", "-")
	name := "doc-craft/" + slug
	if len(name) > 200 {
		return name[:200]
	}
	return name
}

// FilterOpenPRs drops practices that already have an open PR from PracticeBranch.
// skipped entries are human-readable reasons for pipeline.Skipped / logs.
func (m *PRManager) FilterOpenPRs(practices []model.Practice) (keep []model.Practice, skipped []string, err error) {
	for _, p := range practices {
		branch := PracticeBranch(p.PracticeID)
		pr, err := m.FindOpenPR(branch)
		if err != nil {
			return nil, nil, fmt.Errorf("%s: find open PR for %s: %w", p.PracticeID, branch, err)
		}
		if pr != nil {
			msg := fmt.Sprintf("%s: skip, open PR #%d %s (head %s)", p.PracticeID, pr.Number, pr.URL, branch)
			skipped = append(skipped, msg)
			continue
		}
		keep = append(keep, p)
	}
	return keep, skipped, nil
}
