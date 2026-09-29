package service

import "testing"

func TestFormatRupiah(t *testing.T) {
	cases := []struct {
		amount float64
		want   string
	}{
		{0, "Rp0"},
		{999, "Rp999"},
		{57000, "Rp57.000"},
		{1000000, "Rp1.000.000"},
		{-500, "-Rp500"},
	}
	for _, c := range cases {
		if got := formatRupiah(c.amount); got != c.want {
			t.Errorf("formatRupiah(%v) = %q, want %q", c.amount, got, c.want)
		}
	}
}
