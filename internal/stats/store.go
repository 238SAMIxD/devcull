package stats

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

type ToolStats struct {
	TotalReclaimed int64   `json:"total_reclaimed"`
	RecentRuns     []int64 `json:"recent_runs"`
}

type State struct {
	AllTimeTotal int64                 `json:"all_time_total"`
	Tools        map[string]*ToolStats `json:"tools"`

	lock *flock.Flock
}

func getStateFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := filepath.Join(configDir, "devcull")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(dir, "stats.json"), nil
}

func Load() (*State, error) {
	path, err := getStateFilePath()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &State{Tools: make(map[string]*ToolStats)}, nil
	} else if err != nil {
		return nil, err
	}

	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}

	if s.Tools == nil {
		s.Tools = make(map[string]*ToolStats)
	}

	return &s, nil
}

func LoadAndLock() (*State, error) {
	path, err := getStateFilePath()
	if err != nil {
		return nil, err
	}

	lock := flock.New(path + ".lock")
	if err := lock.Lock(); err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &State{Tools: make(map[string]*ToolStats), lock: lock}, nil
	} else if err != nil {
		lock.Unlock()
		return nil, err
	}

	var s State
	if err := json.Unmarshal(b, &s); err != nil {
		lock.Unlock()
		return nil, err
	}

	if s.Tools == nil {
		s.Tools = make(map[string]*ToolStats)
	}
	s.lock = lock

	return &s, nil
}

func (s *State) Unlock() error {
	if s.lock != nil {
		err := s.lock.Unlock()
		s.lock = nil
		return err
	}
	return nil
}

func (s *State) AddRun(tool string, reclaimed int64) {
	if reclaimed <= 0 {
		return
	}

	s.AllTimeTotal += reclaimed

	ts, exists := s.Tools[tool]
	if !exists {
		ts = &ToolStats{}
		s.Tools[tool] = ts
	}

	ts.TotalReclaimed += reclaimed

	ts.RecentRuns = append([]int64{reclaimed}, ts.RecentRuns...)
	if len(ts.RecentRuns) > 3 {
		ts.RecentRuns = ts.RecentRuns[:3]
	}
}

func (s *State) Save() error {
	path, err := getStateFilePath()
	if err != nil {
		return err
	}

	if s.lock != nil {
		defer s.Unlock()
	} else {
		lock := flock.New(path + ".lock")
		if err := lock.Lock(); err != nil {
			return err
		}
		defer lock.Unlock()
	}

	b, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}
