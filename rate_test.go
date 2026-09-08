package main

import "testing"

func TestParseRate(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"25MB/s", 25000000},
		{"1.5GiB/s", 1.5 * (1 << 30)},
		{"100Mbps", 12500000},
		{"56kbps", 7000},
		{"56KBPS", 7000}, // case-insensitive bit-rate unit
		{"10bps", 1.25},
		{"0bps", 0},
		{"0MB/s", 0},
	}
	for _, c := range cases {
		got, err := ParseRate(c.in)
		if err != nil {
			t.Errorf("ParseRate(%q) returned error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseRate(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestParseRateErrors(t *testing.T) {
	cases := []string{
		"",
		"   ",
		"Mbps",
		"MB/s",
		"10Xbps",
		"25MB",  // missing /s
		"garbage",
	}
	for _, in := range cases {
		if _, err := ParseRate(in); err == nil {
			t.Errorf("ParseRate(%q) expected error, got nil", in)
		}
	}
}

func TestFormatRate(t *testing.T) {
	cases := []struct {
		in   float64
		want string
	}{
		{0, "0 B/s"},
		{25000000, "25.00 MB/s"},
		{12500000, "12.50 MB/s"},
	}
	for _, c := range cases {
		got := FormatRate(c.in)
		if got != c.want {
			t.Errorf("FormatRate(%v) = %q, want %q", c.in, got, c.want)
		}
	}
}
