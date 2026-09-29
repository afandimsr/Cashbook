package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLooksLikeTransaction(t *testing.T) {
	cases := []struct {
		text string
		want bool
	}{
		{"halo", false},
		{"makasih ya", false},
		{"ok", false},
		{"kok bisa?", false},
		{"Selamat pagi!", false},
		{"beli kopi 25rb", true},
		{"Gaji 5jt", true},
		{"bensin 50 rb", true},
		{"makan dua puluh ribu", true},
		{"parkir seribu", true},
		{"Terima Rp 200.000", true},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, looksLikeTransaction(c.text), c.text)
	}
}
