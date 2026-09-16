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
		in        int64
		precision int
		want      string
	}{
		{0, 2, "0 B"},
		{999, 2, "999 B"},
		{1000, 2, "1.00 KB"},
		{4700000000, 2, "4.70 GB"},
		{-1000, 2, "-1.00 KB"},
		{4700000000, 0, "5 GB"},
		{4700000000, 4, "4.7000 GB"},
	}
	for _, c := range cases {
		got := FormatSize(c.in, c.precision)
		if got != c.want {
			t.Errorf("FormatSize(%d, %d) = %q, want %q", c.in, c.precision, got, c.want)
		}
	}
}
