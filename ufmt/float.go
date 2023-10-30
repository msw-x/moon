package ufmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Float64(v float64, precision int) (s string) {
	f := floatFmt(precision)
	if f == "" {
		return strconv.Itoa(int(v))
	}
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf(f, v), "0"), ".")
}

func floatFmt(precision int) string {
	var f string
	switch precision {
	case 1:
		f = "%.1f"
	case 2:
		f = "%.2f"
	case 3:
		f = "%.3f"
	case 4:
		f = "%.4f"
	case 5:
		f = "%.5f"
	case 6:
		f = "%.6f"
	case 7:
		f = "%.7f"
	case 8:
		f = "%.8f"
	case 9:
		f = "%.9f"
	case 10:
		f = "%.10f"
	case 11:
		f = "%.11f"
	case 12:
		f = "%.12f"
	case 13:
		f = "%.13f"
	case 14:
		f = "%.14f"
	case 15:
		f = "%.15f"
	case 16:
		f = "%.16f"
	case 17:
		f = "%.17f"
	case 18:
		f = "%.18f"
	case 19:
		f = "%.19f"
	case 20:
		f = "%.20f"
	}
	return f
}
