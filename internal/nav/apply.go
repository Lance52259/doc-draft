package nav

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Lance52259/doc-draft/internal/model"
)

var h1Re = regexp.MustCompile(`(?m)^#\s+(.+?)\s*$`)

// ApplyOptions controls surgical navigation updates after AI generation.
type ApplyOptions struct {
	CRepoRoot     string
	DocsRoot      string // e.g. docs/zh-cn/best-practices
	Service       string // C service dir
	Slug          string // practice file stem
	ServiceLabel  string
	PracticeTitle string
	OneLiner      string
	READMEHeading string
	READMEBlurb   string
}

// ApplyToFiles enforces navigation rules:
//   - SUMMARY.md always exists → always surgical insert at the correct place
//   - service dir first time → keep/require index.md create; patch README + SUMMARY service block
//   - service dir already exists → surgical insert into index list + SUMMARY practice line only
func ApplyToFiles(files []model.DocFileChange, opt ApplyOptions) ([]model.DocFileChange, error) {
	if opt.CRepoRoot == "" || opt.Service == "" || opt.Slug == "" {
		return files, nil
	}
	docsRoot := opt.DocsRoot
	if docsRoot == "" {
		docsRoot = "docs/zh-cn/best-practices"
	}

	title := opt.PracticeTitle
	if title == "" {
		title = TitleFromFiles(files, docsRoot, opt.Service, opt.Slug)
	}
	oneLiner := opt.OneLiner
	if oneLiner == "" {
		oneLiner = fmt.Sprintf("介绍如何使用Terraform自动化完成「%s」。", title)
	}
	label := opt.ServiceLabel
	if label == "" {
		label = strings.ToUpper(opt.Service)
	}

	summaryPath := filepath.ToSlash(filepath.Join(filepath.Dir(docsRoot), "SUMMARY.md"))
	indexPath := filepath.ToSlash(filepath.Join(docsRoot, opt.Service, "index.md"))
	readmePath := filepath.ToSlash(filepath.Join(docsRoot, "README.md"))
	serviceDir := filepath.Join(opt.CRepoRoot, filepath.FromSlash(docsRoot), opt.Service)

	summaryBase, summaryOK := readFile(opt.CRepoRoot, summaryPath)
	if !summaryOK || strings.TrimSpace(summaryBase) == "" {
		return nil, fmt.Errorf("SUMMARY.md must exist at %s", summaryPath)
	}
	if existing := ServiceLabelFromSUMMARY(summaryBase, opt.Service); existing != "" {
		label = existing
	}

	indexBase, indexExists := readFile(opt.CRepoRoot, indexPath)
	readmeBase, readmeOK := readFile(opt.CRepoRoot, readmePath)
	newService := !dirExists(serviceDir) && !indexExists

	out := make([]model.DocFileChange, 0, len(files)+3)
	var aiIndex *model.DocFileChange
	for _, f := range files {
		p := filepath.ToSlash(f.Path)
		switch p {
		case summaryPath:
			// always rebuilt surgically below
			continue
		case readmePath:
			continue
		case indexPath:
			if !newService {
				// existing service: ignore AI rewrite; patch list below
				continue
			}
			cp := f
			aiIndex = &cp
		default:
			out = append(out, f)
		}
	}

	// SUMMARY.md: always append/insert only
	patchedSummary, err := PatchSUMMARY(summaryBase, opt.Service, label, opt.Slug, title)
	if err != nil {
		return nil, fmt.Errorf("patch SUMMARY: %w", err)
	}
	if IsDestructiveUpdate(summaryBase, patchedSummary) {
		return nil, fmt.Errorf("refusing destructive SUMMARY patch")
	}
	out = append(out, model.DocFileChange{Path: summaryPath, Action: "update", Content: patchedSummary})

	if newService {
		// index.md does not exist → must create (prefer AI body; ensure list contains this practice)
		if aiIndex != nil && strings.TrimSpace(aiIndex.Content) != "" {
			content := aiIndex.Content
			if !strings.Contains(content, "]("+opt.Slug+".md)") && !strings.Contains(content, "]("+opt.Slug+")") {
				if patched, err := PatchServiceIndex(content, opt.Slug, title, oneLiner); err == nil {
					content = patched
				}
			}
			out = append(out, model.DocFileChange{Path: indexPath, Action: "create", Content: ensureTrailingNewline(content)})
		} else {
			return nil, fmt.Errorf("new service %s requires create of %s (missing from AI output)", opt.Service, indexPath)
		}
		if readmeOK {
			heading := opt.READMEHeading
			if heading == "" {
				heading = fmt.Sprintf("%s最佳实践", label)
			}
			blurb := opt.READMEBlurb
			if blurb == "" {
				blurb = fmt.Sprintf("%s 相关 Terraform 最佳实践。", label)
			}
			patched, err := PatchBestPracticesREADME(readmeBase, opt.Service, heading, blurb)
			if err != nil {
				return nil, fmt.Errorf("patch README: %w", err)
			}
			if IsDestructiveUpdate(readmeBase, patched) {
				return nil, fmt.Errorf("refusing destructive README patch")
			}
			out = append(out, model.DocFileChange{Path: readmePath, Action: "update", Content: patched})
		}
		return out, nil
	}

	// Existing service directory: only append list item in index.md
	patchedIndex, err := PatchServiceIndex(indexBase, opt.Slug, title, oneLiner)
	if err != nil {
		return nil, fmt.Errorf("patch index: %w", err)
	}
	if IsDestructiveUpdate(indexBase, patchedIndex) {
		return nil, fmt.Errorf("refusing destructive index patch")
	}
	out = append(out, model.DocFileChange{Path: indexPath, Action: "update", Content: patchedIndex})
	return out, nil
}

// TitleFromFiles extracts `# title` from the practice markdown in files.
func TitleFromFiles(files []model.DocFileChange, docsRoot, service, slug string) string {
	want := filepath.ToSlash(filepath.Join(docsRoot, service, slug+".md"))
	for _, f := range files {
		if filepath.ToSlash(f.Path) != want {
			continue
		}
		if m := h1Re.FindStringSubmatch(f.Content); m != nil {
			return strings.TrimSpace(m[1])
		}
	}
	return slug
}

// LoadBaselines returns current navigation files for prompt injection.
// SUMMARY.md is always expected; index/README only when present.
func LoadBaselines(cRepoRoot, docsRoot, service string) map[string]string {
	if docsRoot == "" {
		docsRoot = "docs/zh-cn/best-practices"
	}
	out := map[string]string{}
	summaryPath := filepath.ToSlash(filepath.Join(filepath.Dir(docsRoot), "SUMMARY.md"))
	if s, ok := readFile(cRepoRoot, summaryPath); ok {
		out[summaryPath] = s
	}
	indexPath := filepath.ToSlash(filepath.Join(docsRoot, service, "index.md"))
	if s, ok := readFile(cRepoRoot, indexPath); ok {
		out[indexPath] = s
	}
	readmePath := filepath.ToSlash(filepath.Join(docsRoot, "README.md"))
	if s, ok := readFile(cRepoRoot, readmePath); ok {
		out[readmePath] = s
	}
	return out
}

func readFile(root, rel string) (string, bool) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(rel)))
	if err != nil {
		return "", false
	}
	return string(data), true
}

func dirExists(path string) bool {
	st, err := os.Stat(path)
	return err == nil && st.IsDir()
}
