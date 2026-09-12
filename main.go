// bytetime answers one question: given two of {size, rate, duration}, what
// is the third? It's the calculation behind "how long will this transfer
// take" and "what throughput do I need to hit this deadline."
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math"
	"os"
	"time"
)

type result struct {
	SizeBytes       int64   `json:"size_bytes"`
	SizeHuman       string  `json:"size_human"`
	RateBytesPerSec float64 `json:"rate_bytes_per_sec"`
	RateHuman       string  `json:"rate_human"`
	DurationSeconds float64 `json:"duration_seconds"`
	DurationHuman   string  `json:"duration_human"`
}

func main() {
	sizeFlag := flag.String("size", "", "data size, e.g. 4.7GB, 650MiB")
	rateFlag := flag.String("rate", "", "transfer rate, e.g. 25MB/s, 100Mbps")
	durationFlag := flag.String("duration", "", "duration, e.g. 90s, 1h30m, 2d, 1w")
	jsonOut := flag.Bool("json", false, "print the result as JSON instead of text")
	flag.Usage = usage
	flag.Parse()

	given := 0
	for _, v := range []string{*sizeFlag, *rateFlag, *durationFlag} {
		if v != "" {
			given++
		}
	}
	if given != 2 {
		fmt.Fprintln(os.Stderr, "bytetime: give exactly two of --size, --rate, --duration to solve for the third")
		flag.Usage()
		os.Exit(2)
	}

	res, err := solve(*sizeFlag, *rateFlag, *durationFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bytetime:", err)
		os.Exit(1)
	}

	if *jsonOut {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(res); err != nil {
			fmt.Fprintln(os.Stderr, "bytetime:", err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("%s at %s takes %s\n", res.SizeHuman, res.RateHuman, res.DurationHuman)
}

func solve(sizeStr, rateStr, durationStr string) (*result, error) {
	var (
		sizeBytes int64
		rateBps   float64
		dur       time.Duration
		err       error
	)

	haveSize := sizeStr != ""
	haveRate := rateStr != ""
	haveDuration := durationStr != ""

	if haveSize {
		sizeBytes, err = ParseSize(sizeStr)
		if err != nil {
			return nil, fmt.Errorf("size: %w", err)
		}
	}
	if haveRate {
		rateBps, err = ParseRate(rateStr)
		if err != nil {
			return nil, fmt.Errorf("rate: %w", err)
		}
	}
	if haveDuration {
		dur, err = ParseDuration(durationStr)
		if err != nil {
			return nil, fmt.Errorf("duration: %w", err)
		}
	}

	switch {
	case !haveDuration:
		if rateBps == 0 {
			return nil, errors.New("rate cannot be zero when solving for duration")
		}
		dur = time.Duration(float64(sizeBytes) / rateBps * float64(time.Second))
	case !haveRate:
		if dur <= 0 {
			return nil, errors.New("duration must be positive when solving for rate")
		}
		rateBps = float64(sizeBytes) / dur.Seconds()
	case !haveSize:
		sizeBytes = int64(math.Round(rateBps * dur.Seconds()))
	}

	return &result{
		SizeBytes:       sizeBytes,
		SizeHuman:       FormatSize(sizeBytes),
		RateBytesPerSec: rateBps,
		RateHuman:       FormatRate(rateBps),
		DurationSeconds: dur.Seconds(),
		DurationHuman:   dur.String(),
	}, nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "bytetime solves for whichever of size, rate, or duration you leave out.")
	fmt.Fprintln(os.Stderr, "\nUsage:")
	fmt.Fprintln(os.Stderr, "  bytetime --size 4.7GB --rate 25MB/s")
	fmt.Fprintln(os.Stderr, "  bytetime --size 4.7GB --duration 3m20s")
	fmt.Fprintln(os.Stderr, "  bytetime --rate 100Mbps --duration 1h")
	fmt.Fprintln(os.Stderr, "\nFlags:")
	flag.PrintDefaults()
}
