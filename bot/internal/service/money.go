package service

import (
	"math"
	"strconv"
	"strings"
)

// formatRupiah renders an amount the way the rest of CashBook does (see
// frontend's formatIDR): dot-separated thousands, no decimals, "Rp" prefix
// with no space — e.g. 57000 -> "Rp57.000", 1000000 -> "Rp1.000.000".
func formatRupiah(amount float64) string {
	n := int64(math.Round(amount))
	neg := n < 0
	if neg {
		n = -n
	}
	digits := strconv.FormatInt(n, 10)
	var groups []string
	for len(digits) > 3 {
		groups = append([]string{digits[len(digits)-3:]}, groups...)
		digits = digits[:len(digits)-3]
	}
	groups = append([]string{digits}, groups...)
	sign := ""
	if neg {
		sign = "-"
	}
	return sign + "Rp" + strings.Join(groups, ".")
}
