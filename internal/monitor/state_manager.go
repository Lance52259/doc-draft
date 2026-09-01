package monitor

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RunState persists last run metadata.
type RunState struct {
	BCommit             string            `json:"b_commit,omitempty"`
	CCommit             string            `json:"c_commit,omitempty"`
	ProcessedPractices  []string          `json:"processed_practices"`
	OpenPRs             map[string]string `json:"open_prs"`
}

// StateManager loads/saves RunState JSON.
type StateManager struct {
	Path string
}

func (m *StateManager) Load() (*RunState, error) {
	data, err := os.ReadFile(m.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return &RunState{OpenPRs: map[string]string{}}, nil
		}
		return nil, err
	}
	var st RunState
	if err := json.Unmarshal(data, &st); err != nil {
		return nil, err
	}
	if st.OpenPRs == nil {
		st.OpenPRs = map[string]string{}
	}
	return &st, nil
}

func (m *StateManager) Save(st *RunState) error {
	if err := os.MkdirAll(filepath.Dir(m.Path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(st, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.Path, append(data, '\n'), 0o644)
}
