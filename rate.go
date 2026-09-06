package main

import (
	"fmt"
	"strconv"
	"strings"
)

// decimal prefixes for bit rates (bps, kbps, Mbps, ...). Network speeds are
// conventionally decimal, not binary, so there's no Kibps/Mibps here.
var bitPrefixMultipliers = map[string]float64{
	"":  1,
	"k": 1e3,
	"m": 1e6,
	"g": 1e9,
	"t": 1e12,
}

// ParseRate accepts two families of input:
//
//	"<size>/s"   bytes per second, e.g. "25MB/s", "1.5GiB/s"
//	"<n><p>bps"  bits per second, e.g. "100Mbps", "56kbps", "10bps"
//
// The second form is divided by 8 to normalize everything to bytes/sec
// internally, since that's what ParseSize and FormatSize deal in.
func ParseRate(s string) (float64, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return 0, fmt.Errorf("empty rate")
	}
	lower := strings.ToLower(trimmed)

	if strings.HasSuffix(lower, "bps") {
		numPart := trimmed[:len(trimmed)-3]
		i := 0
		for i < len(numPart) && (isDigit(numPart[i]) || numPart[i] == '.') {
			i++
		}
		if i == 0 {
			return 0, fmt.Errorf("no numeric value in rate %q", s)
		}
		val, err := strconv.ParseFloat(numPart[:i], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number in rate %q", s)
		}
		prefix := strings.ToLower(strings.TrimSpace(numPart[i:]))
		mult, ok := bitPrefixMultipliers[prefix]
		if !ok {
			return 0, fmt.Errorf("unknown bit-rate prefix %q in %q", prefix, s)
		}
		return (val * mult) / 8, nil
	}

	if strings.HasSuffix(trimmed, "/s") {
		sizePart := trimmed[:len(trimmed)-2]
		bytes, err := ParseSize(sizePart)
		if err != nil {
			return 0, err
		}
		return float64(bytes), nil
	}

	return 0, fmt.Errorf("rate %q must end in a byte unit plus /s (e.g. MB/s) or a bit unit plus bps (e.g. Mbps)", s)
}

// FormatRate renders bytes/sec using the same scale as FormatSize.
func FormatRate(bytesPerSec float64) string {
	return FormatSize(int64(bytesPerSec)) + "/s"
}
