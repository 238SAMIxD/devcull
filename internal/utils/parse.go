package utils

import (
	"regexp"
	"strconv"
	"strings"
)

var byteRe = regexp.MustCompile(`(?i)([\d.]+)\s*(B|KB|MB|GB|TB|KIB|MIB|GIB|TIB)`)

func ParseByteString(s string) int64 {
	matches := byteRe.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return 0
	}

	lastMatch := matches[len(matches)-1]

	val, err := strconv.ParseFloat(lastMatch[1], 64)
	if err != nil {
		return 0
	}

	multiplier := float64(1)
	switch strings.ToUpper(lastMatch[2]) {
	case "KB", "KIB":
		multiplier = 1024
	case "MB", "MIB":
		multiplier = 1024 * 1024
	case "GB", "GIB":
		multiplier = 1024 * 1024 * 1024
	case "TB", "TIB":
		multiplier = 1024 * 1024 * 1024 * 1024
	}

	return int64(val * multiplier)
}
