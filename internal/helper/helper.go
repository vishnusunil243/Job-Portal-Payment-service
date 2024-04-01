package helper

import (
	"strconv"
	"strings"
)

func IsValidDuration(s string) bool {
	if !strings.HasSuffix(s, "year") && !strings.HasSuffix(s, "month") {
		return false
	}
	_, err := strconv.Atoi(string(s[0]))
	if err != nil {
		return false
	}
	return true
}
func IsNegative(v float64) bool {
	if v < 0 {
		return true
	}
	return false
}
