package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// decimal and binary byte units, longest suffix checked last so "B" doesn't
// shadow "KiB" when matching.
var sizeUnits = []struct {
	name  string
	bytes float64
}{
	{"PiB", 1 << 50}, {"TiB", 1 << 40}, {"GiB", 1 << 30}, {"MiB", 1 << 20}, {"KiB", 1 << 10},
	{"PB", 1e15}, {"TB", 1e12}, {"GB", 1e9}, {"MB", 1e6}, {"KB", 1e3},
	{"B", 1},
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

// ParseSize turns a string like "4.7GB", "650MiB", or a bare "1024" (bytes)
// into a byte count.
func ParseSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty size")
	}

	i := 0
	if i < len(s) && (s[i] == '+' || s[i] == '-') {
		i++
	}
	for i < len(s) && (isDigit(s[i]) || s[i] == '.') {
		i++
	}
	numPart := s[:i]
	unitPart := strings.TrimSpace(s[i:])

	if numPart == "" || numPart == "+" || numPart == "-" {
		return 0, fmt.Errorf("no numeric value in %q", s)
	}
	val, err := strconv.ParseFloat(numPart, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q", numPart)
	}

	mult, ok := unitMultiplier(unitPart)
	if !ok {
		return 0, fmt.Errorf("unknown size unit %q", unitPart)
	}
	return int64(math.Round(val * mult)), nil
}

func unitMultiplier(unit string) (float64, bool) {
	if unit == "" {
		return 1, true
	}
	for _, u := range sizeUnits {
		if strings.EqualFold(u.name, unit) {
			return u.bytes, true
		}
	}
	return 0, false
}

// FormatSize renders a byte count the way "du -h" style tools do: decimal
// steps of 1000, two decimal places once we're past raw bytes.
func FormatSize(n int64) string {
	if n < 0 {
		return "-" + FormatSize(-n)
	}
	if n < 1000 {
		return fmt.Sprintf("%d B", n)
	}

	units := []string{"KB", "MB", "GB", "TB", "PB"}
	f := float64(n)
	i := -1
	for f >= 1000 && i < len(units)-1 {
		f /= 1000
		i++
	}
	return fmt.Sprintf("%.2f %s", f, units[i])
}
