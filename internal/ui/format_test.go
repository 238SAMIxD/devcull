package ui

import "testing"

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name string
		b    int64
		want string
	}{
		{"zero", 0, "0 B"},
		{"one byte", 1, "1 B"},
		{"below KB", 1023, "1023 B"},
		{"exact KB", 1024, "1.0 KB"},
		{"1.5 KB", 1536, "1.5 KB"},
		{"below MB", 1024*1024 - 1, "1024.0 KB"},
		{"exact MB", 1024 * 1024, "1.0 MB"},
		{"2.5 MB", 2.5 * 1024 * 1024, "2.5 MB"},
		{"exact GB", 1024 * 1024 * 1024, "1.0 GB"},
		{"1.5 GB", 1.5 * 1024 * 1024 * 1024, "1.5 GB"},
		{"exact TB", 1024 * 1024 * 1024 * 1024, "1.0 TB"},
		{"exact PB", 1024 * 1024 * 1024 * 1024 * 1024, "1.0 PB"},
		{"exact EB", 1024 * 1024 * 1024 * 1024 * 1024 * 1024, "1.0 EB"},
		{"negative", -100, "-100 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatBytes(tt.b)
			if got != tt.want {
				t.Errorf("FormatBytes(%d) = %q, want %q", tt.b, got, tt.want)
			}
		})
	}
}
