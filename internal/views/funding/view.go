package funding

import (
	"fmt"
	"strings"

	"github.com/born1337/hyperliquid-terminal/internal/style"
	"github.com/born1337/hyperliquid-terminal/internal/util"
)

const (
	colTime    = 20
	colCoin    = 10
	colPayment = 14
	colRate    = 12
	colPos     = 16
)

func (m Model) View() string {
	m.store.RLock()
	payments := m.store.FundingPayments
	m.store.RUnlock()

	if len(payments) == 0 {
		return style.Dim.Render("  No funding payments")
	}

	var b strings.Builder

	// Header
	header := padRight("TIME", colTime) + "  " +
		padRight("COIN", colCoin) + "  " +
		padLeft("PAYMENT", colPayment) + "  " +
		padLeft("RATE", colRate) + "  " +
		padLeft("POSITION", colPos)
	b.WriteString(style.TableHeader.Render(header))
	b.WriteString("\n")

	// Visible range
	visibleRows := m.height - 5
	if visibleRows < 1 {
		visibleRows = len(payments)
	}
	start := m.scroll
	if start >= len(payments) {
		start = len(payments) - 1
	}
	if start < 0 {
		start = 0
	}
	end := start + visibleRows
	if end > len(payments) {
		end = len(payments)
	}

	var totalPayment float64
	for _, fp := range payments {
		totalPayment += util.ParseFloat(fp.Usdc)
	}

	for _, fp := range payments[start:end] {
		payment := util.ParseFloat(fp.Usdc)
		payStyle := style.PnlColor(payment)
		rate := util.ParseFloat(fp.FundingRate)
		rateStyle := style.PnlColor(rate)

		cells := []string{
			style.Dim.Render(padRight(util.FormatTimeFull(fp.Time), colTime)),
			"  ",
			style.White.Render(padRight(fp.Coin, colCoin)),
			"  ",
			payStyle.Render(padLeft(util.FormatSignedUSD(payment), colPayment)),
			"  ",
			rateStyle.Render(padLeft(util.FormatFundingRate(rate), colRate)),
			"  ",
			padLeft(formatPosition(util.ParseFloat(fp.Szi)), colPos),
		}
		b.WriteString(strings.Join(cells, ""))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	payStyle := style.PnlColor(totalPayment)
	b.WriteString(fmt.Sprintf("  %s %s  %s",
		style.White.Render("Total Funding:"),
		payStyle.Render(util.FormatSignedUSD(totalPayment)),
		style.Dim.Render(fmt.Sprintf("(%d payments)", len(payments))),
	))

	return b.String()
}

// formatPosition formats position size in a compact way for large numbers.
func formatPosition(val float64) string {
	if val < 0 {
		val = -val
	}
	if val >= 1_000_000 {
		return fmt.Sprintf("%.2fM", val/1_000_000)
	}
	if val >= 1_000 {
		return fmt.Sprintf("%.2fK", val/1_000)
	}
	if val >= 1 {
		return fmt.Sprintf("%.4f", val)
	}
	return fmt.Sprintf("%.6f", val)
}

func padRight(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padLeft(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return strings.Repeat(" ", width-len(s)) + s
}
