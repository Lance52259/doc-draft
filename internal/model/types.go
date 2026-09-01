package model

// Practice is one best-practice unit under B repo examples/.
type Practice struct {
	PracticeID string            `json:"practice_id"`
	SourcePath string            `json:"source_path"`
	Title      string            `json:"title"`
	Files      []string          `json:"files"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Slug returns the basename of practice_id.
func (p Practice) Slug() string {
	for i := len(p.PracticeID) - 1; i >= 0; i-- {
		if p.PracticeID[i] == '/' {
			return p.PracticeID[i+1:]
		}
	}
	return p.PracticeID
}

// Service returns the service segment for paths like examples/ecs/basic → ecs.
func (p Practice) Service() string {
	parts := splitPath(p.PracticeID)
	if len(parts) >= 3 && parts[0] == "examples" {
		return parts[1]
	}
	if len(parts) >= 2 && parts[0] != "examples" {
		return parts[0]
	}
	return ""
}

func splitPath(p string) []string {
	var out []string
	start := 0
	for i := 0; i < len(p); i++ {
		if p[i] == '/' {
			if i > start {
				out = append(out, p[start:i])
			}
			start = i + 1
		}
	}
	if start < len(p) {
		out = append(out, p[start:])
	}
	return out
}

// DocFileChange is one generated file write.
type DocFileChange struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Action  string `json:"action"`
}

// GenerateResult is AI output for one practice.
type GenerateResult struct {
	PracticeID  string          `json:"practice_id"`
	Files       []DocFileChange `json:"files"`
	Summary     string          `json:"summary"`
	RawResponse string          `json:"-"`
}

// DetectionResult is the outcome of B vs C comparison.
type DetectionResult struct {
	NewPractices []Practice `json:"new_practices"`
	SyncedIDs    []string   `json:"synced_ids"`
	BCommit      string     `json:"b_commit,omitempty"`
	CCommit      string     `json:"c_commit,omitempty"`
}

// PipelineResult summarizes a full run.
type PipelineResult struct {
	Detected  DetectionResult  `json:"detected"`
	Generated []GenerateResult `json:"generated"`
	PRURLs    []string         `json:"pr_urls"`
	Skipped   []string         `json:"skipped"`
	Errors    []string         `json:"errors"`
	DryRun    bool             `json:"dry_run"`
}

// RepoRef describes a prepared working tree.
type RepoRef struct {
	Name      string
	Branch    string
	Token     string
	LocalPath string
	CommitSHA string
}
