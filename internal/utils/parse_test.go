package utils

import "testing"

func TestParseByteString(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want int64
	}{
		// basic units
		{"bytes", "100 B", 100},
		{"kilobytes", "1 KB", 1024},
		{"megabytes", "1 MB", 1024 * 1024},
		{"gigabytes", "1 GB", 1024 * 1024 * 1024},
		{"terabytes", "1 TB", 1024 * 1024 * 1024 * 1024},

		// decimals
		{"decimal KB", "1.5 KB", 1536},
		{"decimal GB", "2.5 GB", 2684354560}, // 2.5 * 1024^3

		// spacing variations
		{"no space", "500MB", 500 * 1024 * 1024},
		{"extra space", "500  MB", 500 * 1024 * 1024},

		// case insensitivity (input is uppercased internally)
		{"lowercase mb", "100 mb", 100 * 1024 * 1024},
		{"mixed case", "100 Mb", 100 * 1024 * 1024},

		// picks last match on a line (used for docker parsing)
		{"multiple matches", "Total: 500 MB Reclaimable: 200 MB", 200 * 1024 * 1024},

		// docker-style output with percentage
		{"with percentage", "1.2GB (50%)", 1288490188}, // 1.2 * 1024^3

		// edge cases
		{"empty string", "", 0},
		{"no match", "no bytes here", 0},
		{"zero bytes", "0 B", 0},
		{"zero KB", "0 KB", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseByteString(tt.s)
			if got != tt.want {
				t.Errorf("ParseByteString(%q) = %d, want %d", tt.s, got, tt.want)
			}
		})
	}
}
