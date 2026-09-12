package main

import (
	"testing"
	"time"
)

func TestParseDuration(t *testing.T) {
	cases := []struct {
		in   string
		want time.Duration
	}{
		{"90s", 90 * time.Second},
		{"3m20s", 3*time.Minute + 20*time.Second},
		{"1h30m", time.Hour + 30*time.Minute},
		{"1d", 24 * time.Hour},
		{"2d", 48 * time.Hour},
		{"1w", 7 * 24 * time.Hour},
		{"1w3d12h", 7*24*time.Hour + 3*24*time.Hour + 12*time.Hour},
		{"-2d", -48 * time.Hour},
		{"+2d", 48 * time.Hour},
		{"1.5d", 36 * time.Hour},
		{"500ms", 500 * time.Millisecond},
		{"0d", 0},
	}
	for _, c := range cases {
		got, err := ParseDuration(c.in)
		if err != nil {
			t.Errorf("ParseDuration(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseDuration(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseDurationErrors(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"d",
		"5",
		"5x",
		"1w2",
		"-",
		"+",
	}
	for _, in := range cases {
		if _, err := ParseDuration(in); err == nil {
			t.Errorf("ParseDuration(%q) expected error, got nil", in)
		}
	}
}
