package main

import "testing"

func TestParseSize(t *testing.T) {
	cases := []struct {
		in   string
		want int64
	}{
		{"4.7GB", 4700000000},
		{"650MiB", 681574400},
		{"1024", 1024},
		{"0", 0},
		{"0B", 0},
		{"-5MB", -5000000},
		{"+5MB", 5000000},
		{"5gb", 5000000000},  // case-insensitive unit
		{"5 GB", 5000000000}, // space between number and unit
		{"1KiB", 1024},
		{"1PiB", 1 << 50},
		{"1.5B", 2},     // rounds half away from zero
		{"  10B  ", 10}, // leading/trailing whitespace
	}
	for _, c := range cases {
		got, err := ParseSize(c.in)
		if err != nil {
			t.Errorf("ParseSize(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseSize(%q) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestParseSizeErrors(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"GB",
		"5XB",
		"5.5.5GB",
		"-",
		"+",
		"MB5",
	}
	for _, in := range cases {
		if _, err := ParseSize(in); err == nil {
			t.Errorf("ParseSize(%q) expected error, got nil", in)
		}
	}
}

func TestFormatSize(t *testing.T) {
	cases := []struct {
		in   int64
		want string
	}{
		{0, "0 B"},
		{999, "999 B"},
		{1000, "1.00 KB"},
		{4700000000, "4.70 GB"},
		{-1000, "-1.00 KB"},
	}
	for _, c := range cases {
		got := FormatSize(c.in)
		if got != c.want {
			t.Errorf("FormatSize(%d) = %q, want %q", c.in, got, c.want)
		}
	}
}
