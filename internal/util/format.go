package util

import (
	"fmt"
	"strings"
	"time"
)

func FormatUSD(val float64) string {
	negative := val < 0
	if negative {
		val = -val
	}
	s := fmt.Sprintf("%.2f", val)
	parts := strings.Split(s, ".")
	intPart := parts[0]
	decPart := parts[1]

	// Add commas
	n := len(intPart)
	if n > 3 {
		var b strings.Builder
		for i, c := range intPart {
			if i > 0 && (n-i)%3 == 0 {
				b.WriteByte(',')
			}
			b.WriteRune(c)
		}
		intPart = b.String()
	}

	if negative {
		return "-$" + intPart + "." + decPart
	}
	return "$" + intPart + "." + decPart
}

func FormatSignedUSD(val float64) string {
	if val >= 0 {
		return "+" + FormatUSD(val)
	}
	return FormatUSD(val)
}

func FormatPercent(val float64) string {
	if val >= 0 {
		return fmt.Sprintf("+%.2f%%", val)
	}
	return fmt.Sprintf("%.2f%%", val)
}

func FormatFundingRate(val float64) string {
	return fmt.Sprintf("%.4f%%", val*100*24)
}

func FormatLeverage(val float64) string {
	if val == float64(int(val)) {
		return fmt.Sprintf("%.0fx", val)
	}
	return fmt.Sprintf("%.2fx", val)
}

func FormatPrice(val float64) string {
	if val >= 1000 {
		s := fmt.Sprintf("%.2f", val)
		parts := strings.Split(s, ".")
		intPart := parts[0]
		decPart := parts[1]
		n := len(intPart)
		if n > 3 {
			var b strings.Builder
			for i, c := range intPart {
				if i > 0 && (n-i)%3 == 0 {
					b.WriteByte(',')
				}
				b.WriteRune(c)
			}
			intPart = b.String()
		}
		return "$" + intPart + "." + decPart
	}
	if val >= 1 {
		return fmt.Sprintf("$%.2f", val)
	}
	return fmt.Sprintf("$%.4f", val)
}

func FormatTime(ts int64) string {
	t := time.UnixMilli(ts)
	return t.Format("Jan 02 15:04")
}

func FormatTimeFull(ts int64) string {
	t := time.UnixMilli(ts)
	return t.Format("2006-01-02 15:04:05")
}

func FormatSize(val float64) string {
	if val < 0 {
		val = -val
	}
	if val >= 1 {
		return fmt.Sprintf("%.4f", val)
	}
	return fmt.Sprintf("%.6f", val)
}
