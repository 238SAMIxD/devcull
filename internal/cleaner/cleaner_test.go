package cleaner

import "testing"

type mockCleaner struct {
	name     string
	category Category
}

func (m *mockCleaner) Name() string           { return m.name }
func (m *mockCleaner) Category() Category     { return m.category }
func (m *mockCleaner) IsInstalled() bool      { return true }
func (m *mockCleaner) EstimateReclaimable() (int64, error) { return 0, nil }
func (m *mockCleaner) Clean(dryRun bool) (int64, error) { return 0, nil }

func TestMatchesArg(t *testing.T) {
	tests := []struct {
		name    string
		cleaner Cleaner
		arg     string
		want    bool
	}{
		// 1. Name matches
		{"exact name", &mockCleaner{"Pip", CategoryPython}, "pip", true},
		{"name case insensitive", &mockCleaner{"Homebrew", CategorySystem}, "HOMEBREW", true},
		{"name with spaces", &mockCleaner{"Docker", CategorySystem}, "  docker  ", true},
		
		// 2. Category matches
		{"exact category", &mockCleaner{"Pip", CategoryPython}, "python", true},
		{"category partial (prefix)", &mockCleaner{"NPM", CategoryNode}, "node", true},
		{"category partial (substring)", &mockCleaner{"NPM", CategoryNode}, "& js", true},
		{"category case insensitive", &mockCleaner{"Cargo", CategoryRust}, "RUST", true},

		// 3. Explicit aliases
		{"brew -> homebrew", &mockCleaner{"Homebrew", CategorySystem}, "brew", true},
		{"macos -> homebrew", &mockCleaner{"Homebrew", CategorySystem}, "macos", true},
		{"py -> python category", &mockCleaner{"Poetry", CategoryPython}, "py", true},
		{"python3 -> python category", &mockCleaner{"Uv", CategoryPython}, "python3", true},
		{"node -> node category", &mockCleaner{"Yarn", CategoryNode}, "node", true},
		{"ts -> node category", &mockCleaner{"Bun", CategoryNode}, "ts", true},
		{"java -> java category", &mockCleaner{"Maven", CategoryJava}, "java", true},
		{"jvm -> java category", &mockCleaner{"Gradle", CategoryJava}, "jvm", true},
		{"c# -> dotnet category", &mockCleaner{"Dotnet", CategoryCSharp}, "c#", true},
		{"nuget -> dotnet category", &mockCleaner{"Dotnet", CategoryCSharp}, "nuget", true},
		{"golang -> go category", &mockCleaner{"Go", CategoryGo}, "golang", true},
		
		// Apple specific
		{"ios -> apple category", &mockCleaner{"Xcode", CategoryApple}, "ios", true},
		{"mac -> cocoapods", &mockCleaner{"CocoaPods", CategorySystem}, "mac", true},

		// 4. Non-matches
		{"wrong name", &mockCleaner{"Pip", CategoryPython}, "npm", false},
		{"wrong category", &mockCleaner{"Pip", CategoryPython}, "rust", false},
		{"unrelated alias", &mockCleaner{"Cargo", CategoryRust}, "brew", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MatchesArg(tt.cleaner, tt.arg)
			if got != tt.want {
				t.Errorf("MatchesArg(cleaner={%s, %s}, %q) = %v, want %v", tt.cleaner.Name(), tt.cleaner.Category(), tt.arg, got, tt.want)
			}
		})
	}
}
