package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func ParseByteString(s string) int64 {
	re := regexp.MustCompile(`([\d\.]+)\s*(B|KB|MB|GB|TB)`)
	
	matches := re.FindAllStringSubmatch(strings.ToUpper(s), -1)
	if len(matches) == 0 {
		return 0
	}

	lastMatch := matches[len(matches)-1]

	val, err := strconv.ParseFloat(lastMatch[1], 64)
	if err != nil {
		return 0
	}

	multiplier := float64(1)
	switch lastMatch[2] {
	case "KB":
		multiplier = 1024
	case "MB":
		multiplier = 1024 * 1024
	case "GB":
		multiplier = 1024 * 1024 * 1024
	case "TB":
		multiplier = 1024 * 1024 * 1024 * 1024
	}

	return int64(val * multiplier)
}