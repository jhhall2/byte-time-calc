package main

import (
	"fmt"
	"strconv"
	"time"
)

// durationUnitMultipliers covers everything time.ParseDuration does, plus
// "d" and "w" since backup windows and retention periods are routinely
// measured in days and weeks and stdlib stops at hours.
var durationUnitMultipliers = map[string]time.Duration{
	"ns": time.Nanosecond,
	"us": time.Microsecond,
	"µs": time.Microsecond,
	"μs": time.Microsecond,
	"ms": time.Millisecond,
	"s":  time.Second,
	"m":  time.Minute,
	"h":  time.Hour,
	"d":  24 * time.Hour,
	"w":  7 * 24 * time.Hour,
}

// ParseDuration parses the same syntax as time.ParseDuration - a sequence of
// decimal numbers each followed by a unit, e.g. "1h30m" - plus "d" and "w"
// units, so "1w2d" and "90s" both work.
func ParseDuration(s string) (time.Duration, error) {
	orig := s
	neg := false
	if s != "" && (s[0] == '+' || s[0] == '-') {
		neg = s[0] == '-'
		s = s[1:]
	}
	if s == "" {
		return 0, fmt.Errorf("invalid duration %q", orig)
	}

	var total time.Duration
	for s != "" {
		i := 0
		for i < len(s) && (isDigit(s[i]) || s[i] == '.') {
			i++
		}
		if i == 0 {
			return 0, fmt.Errorf("invalid duration %q", orig)
		}
		numStr := s[:i]
		s = s[i:]

		j := 0
		for j < len(s) && !isDigit(s[j]) && s[j] != '.' {
			j++
		}
		if j == 0 {
			return 0, fmt.Errorf("missing unit after %q in duration %q", numStr, orig)
		}
		unit := s[:j]
		s = s[j:]

		val, err := strconv.ParseFloat(numStr, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number %q in duration %q", numStr, orig)
		}
		mult, ok := durationUnitMultipliers[unit]
		if !ok {
			return 0, fmt.Errorf("unknown duration unit %q in %q", unit, orig)
		}
		total += time.Duration(val * float64(mult))
	}

	if neg {
		total = -total
	}
	return total, nil
}
