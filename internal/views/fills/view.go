package fills

import (
	"fmt"
	"strings"

	"github.com/born1337/hyperliquid-terminal/internal/style"
	"github.com/born1337/hyperliquid-terminal/internal/util"
)

var separator80 = strings.Repeat("─", 80)

func (m Model) View() string {
	m.store.RLock()
	allFills := m.store.Fills
	m.store.RUnlock()

	if len(allFills) == 0 {
		return style.Dim.Render("  No recent fills")
	}

	var b strings.Builder

	header := fmt.Sprintf("%-16s %-9s %-6s %12s %12s %14s %12s",
		"TIME", "COIN", "SIDE", "SIZE", "PRICE", "REALIZED PNL", "FEE",
	)
	b.WriteString(style.TableHeader.Render(header))
	b.WriteString("\n")

	visibleRows := m.height - 5
	if visibleRows < 1 {
		visibleRows = len(allFills)
	}
	start := m.scroll
	if start >= len(allFills) {
		start = len(allFills) - 1
	}
	if start < 0 {
		start = 0
	}
	end := start + visibleRows
	if end > len(allFills) {
		end = len(allFills)
	}

	var totalRealizedPnl, totalFees float64
	for _, f := range allFills {
		totalRealizedPnl += util.ParseFloat(f.ClosedPnl)
		totalFees += util.ParseFloat(f.Fee)
	}

	for _, f := range allFills[start:end] {
		sideStyle := style.Green
		side := "BUY"
		if f.Side == "A" || f.Side == "sell" {
			sideStyle = style.Red
			side = "SELL"
		}

		closedPnl := util.ParseFloat(f.ClosedPnl)
		pnlStyle := style.PnlColor(closedPnl)
		pnlStr := "-"
		if closedPnl != 0 {
			pnlStr = util.FormatSignedUSD(closedPnl)
		}

		fee := util.ParseFloat(f.Fee)

		row := fmt.Sprintf("%-16s %-9s %s %12s %12s %14s %12s",
			util.FormatTimeFull(f.Time),
			style.White.Render(f.Coin),
			sideStyle.Render(fmt.Sprintf("%-6s", side)),
			util.FormatSize(util.ParseFloat(f.Sz)),
			util.FormatPrice(util.ParseFloat(f.Px)),
			pnlStyle.Render(pnlStr),
			style.Red.Render(util.FormatUSD(fee)),
		)
		b.WriteString(row)
		b.WriteString("\n")
	}

	// Summary
	b.WriteString("\n")
	b.WriteString(style.Dim.Render(separator80))
	b.WriteString("\n")
	pnlStyle := style.PnlColor(totalRealizedPnl)
	summary := fmt.Sprintf("  %s %s   %s %s   %s",
		style.White.Render("Realized PnL:"),
		pnlStyle.Render(util.FormatSignedUSD(totalRealizedPnl)),
		style.White.Render("Total Fees:"),
		style.Red.Render(util.FormatUSD(totalFees)),
		style.Dim.Render(fmt.Sprintf("(%d fills)", len(allFills))),
	)
	b.WriteString(summary)

	return b.String()
}
