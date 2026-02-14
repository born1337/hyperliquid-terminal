package positions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/born1337/hltui/internal/style"
	"github.com/born1337/hltui/internal/util"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	positions, mids, fundingRates := m.store.PositionsSortedByPnl()
	if len(positions) == 0 {
		return style.Dim.Render("  No open positions")
	}

	// Sort by PNL
	sort.Slice(positions, func(i, j int) bool {
		pi := util.ParseFloat(positions[i].Position.UnrealizedPnl)
		pj := util.ParseFloat(positions[j].Position.UnrealizedPnl)
		if m.sortAsc {
			return pi < pj
		}
		return pi > pj
	})

	arrow := " ▼"
	if m.sortAsc {
		arrow = " ▲"
	}

	var b strings.Builder

	// Header
	header := fmt.Sprintf("%-9s %-6s %-5s %14s %11s %14s %16s %12s %12s %12s %12s",
		"COIN", "SIDE", "LEV", "VALUE", "FUND RATE", "FUND FEE", "PNL"+arrow, "ROE", "ENTRY", "CURRENT", "LIQ",
	)
	b.WriteString(style.TableHeader.Render(header))
	b.WriteString("\n")

	// Track totals for summary
	var totalPnl, winners, losers float64

	// Determine visible range
	visibleRows := m.height - 6
	if visibleRows < 1 {
		visibleRows = len(positions)
	}
	start := m.scroll
	if start >= len(positions) {
		start = len(positions) - 1
	}
	if start < 0 {
		start = 0
	}
	end := start + visibleRows
	if end > len(positions) {
		end = len(positions)
	}

	for i, ap := range positions {
		p := ap.Position
		szi := util.ParseFloat(p.Szi)
		pnl := util.ParseFloat(p.UnrealizedPnl)
		roe := util.ParseFloat(p.ReturnOnEquity) * 100
		entryPx := util.ParseFloat(p.EntryPx)
		posValue := util.ParseFloat(p.PositionValue)
		lev := p.Leverage.Value

		fundingFee := 0.0
		if p.CumFunding != nil {
			fundingFee = util.ParseFloat(p.CumFunding.SinceOpen)
		}

		currentPx := util.ParseFloat(mids[p.Coin])
		fundRate := fundingRates[p.Coin]

		totalPnl += pnl
		if pnl > 0 {
			winners += pnl
		} else {
			losers += pnl
		}

		if i < start || i >= end {
			continue
		}

		side := "LONG"
		sideStyle := style.Green
		if szi < 0 {
			side = "SHORT"
			sideStyle = style.Red
		}

		pnlStyle := style.PnlColor(pnl)
		roeStyle := style.PnlColor(roe)
		fundFeeStyle := style.PnlColor(fundingFee)
		fundRateStyle := style.PnlColor(fundRate)

		// Handle nullable liquidation price
		liqStr := "-"
		if p.LiquidationPx != nil && *p.LiquidationPx != "" {
			liqStr = util.FormatPrice(util.ParseFloat(*p.LiquidationPx))
		}

		// Build each cell with fixed width, then join
		// Use padRight for left-aligned, padLeft for right-aligned
		cells := []string{
			style.White.Render(padRight(p.Coin, 9)),
			sideStyle.Render(padRight(side, 6)),
			style.Dim.Render(padRight(util.FormatLeverage(lev), 5)),
			padLeft(util.FormatUSD(posValue), 14),
			fundRateStyle.Render(padLeft(util.FormatFundingRate(fundRate), 11)),
			fundFeeStyle.Render(padLeft(util.FormatSignedUSD(fundingFee), 14)),
			pnlStyle.Render(padLeft(util.FormatSignedUSD(pnl), 16)),
			roeStyle.Render(padLeft(util.FormatPercent(roe), 12)),
			padLeft(util.FormatPrice(entryPx), 12),
			padLeft(util.FormatPrice(currentPx), 12),
			style.Dim.Render(padLeft(liqStr, 12)),
		}
		b.WriteString(strings.Join(cells, " "))
		b.WriteString("\n")
	}

	// Summary
	b.WriteString("\n")
	b.WriteString(style.Dim.Render(strings.Repeat("─", 80)))
	b.WriteString("\n")

	summaryLine := fmt.Sprintf("  %s %s   %s %s   %s %s",
		style.White.Render("Winners:"),
		style.Green.Render(util.FormatSignedUSD(winners)),
		style.White.Render("Losers:"),
		style.Red.Render(util.FormatUSD(losers)),
		style.White.Render("Net PnL:"),
		pnlSummaryStyle(totalPnl).Render(util.FormatSignedUSD(totalPnl)),
	)
	b.WriteString(summaryLine)

	return b.String()
}

func pnlSummaryStyle(val float64) lipgloss.Style {
	s := style.PnlColor(val)
	return s.Bold(true)
}

func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

func padLeft(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return strings.Repeat(" ", width-w) + s
}
