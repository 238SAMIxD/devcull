package stats

import "testing"

func TestAddRun(t *testing.T) {
	t.Run("first run adds tool and updates totals", func(t *testing.T) {
		s := &State{Tools: make(map[string]*ToolStats)}
		s.AddRun("npm", 1000)

		if s.AllTimeTotal != 1000 {
			t.Errorf("AllTimeTotal = %d, want 1000", s.AllTimeTotal)
		}
		ts := s.Tools["npm"]
		if ts == nil {
			t.Fatal("expected npm tool stats to exist")
		}
		if ts.TotalReclaimed != 1000 {
			t.Errorf("TotalReclaimed = %d, want 1000", ts.TotalReclaimed)
		}
		if len(ts.RecentRuns) != 1 || ts.RecentRuns[0] != 1000 {
			t.Errorf("RecentRuns = %v, want [1000]", ts.RecentRuns)
		}
	})

	t.Run("multiple runs accumulate", func(t *testing.T) {
		s := &State{Tools: make(map[string]*ToolStats)}
		s.AddRun("go", 500)
		s.AddRun("go", 300)

		if s.AllTimeTotal != 800 {
			t.Errorf("AllTimeTotal = %d, want 800", s.AllTimeTotal)
		}
		ts := s.Tools["go"]
		if ts.TotalReclaimed != 800 {
			t.Errorf("TotalReclaimed = %d, want 800", ts.TotalReclaimed)
		}
		if len(ts.RecentRuns) != 2 {
			t.Fatalf("len(RecentRuns) = %d, want 2", len(ts.RecentRuns))
		}
		// most recent first
		if ts.RecentRuns[0] != 300 || ts.RecentRuns[1] != 500 {
			t.Errorf("RecentRuns = %v, want [300, 500]", ts.RecentRuns)
		}
	})

	t.Run("caps at 3 recent runs", func(t *testing.T) {
		s := &State{Tools: make(map[string]*ToolStats)}
		s.AddRun("pip", 100)
		s.AddRun("pip", 200)
		s.AddRun("pip", 300)
		s.AddRun("pip", 400)

		ts := s.Tools["pip"]
		if len(ts.RecentRuns) != 3 {
			t.Fatalf("len(RecentRuns) = %d, want 3", len(ts.RecentRuns))
		}
		// should keep the 3 most recent: 400, 300, 200
		if ts.RecentRuns[0] != 400 || ts.RecentRuns[1] != 300 || ts.RecentRuns[2] != 200 {
			t.Errorf("RecentRuns = %v, want [400, 300, 200]", ts.RecentRuns)
		}
		if ts.TotalReclaimed != 1000 {
			t.Errorf("TotalReclaimed = %d, want 1000", ts.TotalReclaimed)
		}
	})

	t.Run("zero reclaimed is ignored", func(t *testing.T) {
		s := &State{Tools: make(map[string]*ToolStats)}
		s.AddRun("uv", 0)

		if s.AllTimeTotal != 0 {
			t.Errorf("AllTimeTotal = %d, want 0", s.AllTimeTotal)
		}
		if _, exists := s.Tools["uv"]; exists {
			t.Error("expected uv tool stats to not exist after zero reclaimed")
		}
	})

	t.Run("negative reclaimed is ignored", func(t *testing.T) {
		s := &State{Tools: make(map[string]*ToolStats)}
		s.AddRun("cargo", -500)

		if s.AllTimeTotal != 0 {
			t.Errorf("AllTimeTotal = %d, want 0", s.AllTimeTotal)
		}
		if _, exists := s.Tools["cargo"]; exists {
			t.Error("expected cargo tool stats to not exist after negative reclaimed")
		}
	})

	t.Run("multiple tools tracked independently", func(t *testing.T) {
		s := &State{Tools: make(map[string]*ToolStats)}
		s.AddRun("npm", 100)
		s.AddRun("go", 200)
		s.AddRun("npm", 300)

		if s.AllTimeTotal != 600 {
			t.Errorf("AllTimeTotal = %d, want 600", s.AllTimeTotal)
		}
		if s.Tools["npm"].TotalReclaimed != 400 {
			t.Errorf("npm TotalReclaimed = %d, want 400", s.Tools["npm"].TotalReclaimed)
		}
		if s.Tools["go"].TotalReclaimed != 200 {
			t.Errorf("go TotalReclaimed = %d, want 200", s.Tools["go"].TotalReclaimed)
		}
	})
}
