package cleaner

import "testing"

func TestResolveAlias(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		// explicit aliases
		{"brew", "brew", "homebrew"},
		{"python", "python", "pip"},
		{"python3", "python3", "pip"},
		{"pip3", "pip3", "pip"},
		{"node", "node", "npm"},
		{"nodejs", "nodejs", "npm"},
		{"rust", "rust", "cargo"},
		{"golang", "golang", "go"},
		{"cs", "cs", "dotnet"},
		{"csharp", "csharp", "dotnet"},
		{"c#", "c#", "dotnet"},
		{".net", ".net", "dotnet"},
		{"nuget", "nuget", "dotnet"},
		{"mac", "mac", "cocoapods"},

		// case insensitivity
		{"BREW uppercase", "BREW", "homebrew"},
		{"Python mixed", "Python", "pip"},
		{"NODE uppercase", "NODE", "npm"},

		// identity (no alias)
		{"npm passthrough", "npm", "npm"},
		{"docker passthrough", "docker", "docker"},
		{"go passthrough", "go", "go"},
		{"pip passthrough", "pip", "pip"},
		{"cargo passthrough", "cargo", "cargo"},
		{"yarn passthrough", "yarn", "yarn"},
		{"pnpm passthrough", "pnpm", "pnpm"},
		{"uv passthrough", "uv", "uv"},
		{"bun passthrough", "bun", "bun"},
		{"deno passthrough", "deno", "deno"},

		// unknown tool lowercased
		{"unknown tool", "SomeTool", "sometool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveAlias(tt.input)
			if got != tt.want {
				t.Errorf("ResolveAlias(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
