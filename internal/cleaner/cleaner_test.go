package cleaner

import (
	"context"
	"testing"
)

type mockCleaner struct {
	name     string
	category Category
}

func (m *mockCleaner) Name() string       { return m.name }
func (m *mockCleaner) Category() Category { return m.category }
func (m *mockCleaner) Aliases() []string  { return nil }

func (m *mockCleaner) IsInstalled(ctx context.Context) bool                   { return true }
func (m *mockCleaner) EstimateReclaimable(ctx context.Context) (int64, error) { return 0, nil }
func (m *mockCleaner) Clean(ctx context.Context, dryRun bool) (int64, error)  { return 0, nil }

func TestMatchesArg(t *testing.T) {
	tests := []struct {
		name    string
		cleaner Cleaner
		arg     string
		want    bool
	}{
		{"exact name", &mockCleaner{"Pip", CategoryPython}, "pip", true},
		{"name case insensitive", &mockCleaner{"Homebrew", CategorySystem}, "HOMEBREW", true},
		{"name with spaces", &mockCleaner{"Docker", CategorySystem}, "  docker  ", true},

		{"exact category", &mockCleaner{"Pip", CategoryPython}, "python", true},
		{"category partial (prefix)", &mockCleaner{"NPM", CategoryNode}, "node", true},
		{"category alias (devops)", &mockCleaner{"Docker", CategorySystem}, "devops", true},
		{"category case insensitive", &mockCleaner{"Cargo", CategoryRust}, "RUST", true},

		{"brew -> homebrew", &mockCleaner{"Homebrew", CategorySystem}, "brew", true},
		{"macos -> apple category", &mockCleaner{"Xcode", CategoryApple}, "macos", true},
		{"py -> python category", &mockCleaner{"Poetry", CategoryPython}, "py", true},
		{"python3 -> python category", &mockCleaner{"Uv", CategoryPython}, "python3", true},
		{"node -> node category", &mockCleaner{"Yarn", CategoryNode}, "node", true},
		{"ts -> node category", &mockCleaner{"Bun", CategoryNode}, "ts", true},
		{"java -> java category", &mockCleaner{"Maven", CategoryJava}, "java", true},
		{"jvm -> java category", &mockCleaner{"Gradle", CategoryJava}, "jvm", true},
		{"c# -> dotnet category", &mockCleaner{"Dotnet", CategoryCSharp}, "c#", true},
		{"nuget -> dotnet category", &mockCleaner{"Dotnet", CategoryCSharp}, "nuget", true},
		{"golang -> go category", &mockCleaner{"Go", CategoryGo}, "golang", true},

		{"ios -> apple category", &mockCleaner{"Xcode", CategoryApple}, "ios", true},
		{"mac -> cocoapods", &mockCleaner{"CocoaPods", CategoryApple}, "mac", true},

		{"wrong name", &mockCleaner{"Pip", CategoryPython}, "npm", false},
		{"wrong category", &mockCleaner{"Pip", CategoryPython}, "rust", false},
		{"unrelated alias", &mockCleaner{"Cargo", CategoryRust}, "brew", false},

		{"c must not match apple", &mockCleaner{"Xcode", CategoryApple}, "c", false},
		{"c must not match c#", &mockCleaner{"Dotnet", CategoryCSharp}, "c", false},
		{"o must not match python", &mockCleaner{"Pip", CategoryPython}, "o", false},
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
