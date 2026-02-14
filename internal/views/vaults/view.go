package vaults

import (
	"fmt"
	"strings"

	"github.com/born1337/hyperliquid-terminal/internal/style"
	"github.com/born1337/hyperliquid-terminal/internal/util"
)

const (
	colName    = 28
	colEquity  = 16
	colPnl     = 16
	colAllTime = 16
	colAPR     = 8
)

func (m Model) View() string {
	m.store.RLock()
	equities := m.store.VaultEquities
	details := m.store.VaultDetails
	m.store.RUnlock()

	if len(equities) == 0 {
		return style.Dim.Render("  No vault investments")
	}

	var b strings.Builder

	header := padRight("VAULT", colName) + "  " +
		padLeft("YOUR EQUITY", colEquity) + "  " +
		padLeft("YOUR PNL", colPnl) + "  " +
		padLeft("ALL-TIME PNL", colAllTime) + "  " +
		padLeft("APR", colAPR)
	b.WriteString(style.TableHeader.Render(header))
	b.WriteString("\n")

	var totalEquity float64

	for _, ve := range equities {
		equity := util.ParseFloat(ve.Equity)
		totalEquity += equity

		name := truncAddr(ve.VaultAddress)
		pnlCell := style.Dim.Render(padLeft("-", colPnl))
		allTimeCell := style.Dim.Render(padLeft("-", colAllTime))
		aprCell := style.Dim.Render(padLeft("-", colAPR))

		if d, ok := details[ve.VaultAddress]; ok {
			name = truncName(d.Name, colName)

			if d.FollowerState != nil {
				pnl := util.ParseFloat(d.FollowerState.Pnl)
				allTimePnl := util.ParseFloat(d.FollowerState.AllTimePnl)
				pnlCell = style.PnlColor(pnl).Render(padLeft(util.FormatSignedUSD(pnl), colPnl))
				allTimeCell = style.PnlColor(allTimePnl).Render(padLeft(util.FormatSignedUSD(allTimePnl), colAllTime))

				if d.FollowerState.VaultEquity != "" {
					equity = util.ParseFloat(d.FollowerState.VaultEquity)
				}
			}

			if d.APR != 0 {
				aprCell = style.PnlColor(d.APR).Render(padLeft(fmt.Sprintf("%.1f%%", d.APR*100), colAPR))
			}
		}

		cells := []string{
			style.White.Render(padRight(name, colName)),
			"  ",
			style.Green.Render(padLeft(util.FormatUSD(equity), colEquity)),
			"  ",
			pnlCell,
			"  ",
			allTimeCell,
			"  ",
			aprCell,
		}
		b.WriteString(strings.Join(cells, ""))
		b.WriteString("\n")
	}

	// Total
	totalWidth := colName + colEquity + colPnl + colAllTime + colAPR + 8
	b.WriteString("\n")
	b.WriteString(style.Dim.Render(strings.Repeat("─", totalWidth)))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s %s  %s",
		style.White.Render("Total Vault Equity:"),
		style.Green.Render(util.FormatUSD(totalEquity)),
		style.Dim.Render(fmt.Sprintf("(%d vaults)", len(equities))),
	))

	return b.String()
}

func truncName(name string, maxLen int) string {
	if len(name) > maxLen {
		return name[:maxLen-3] + "..."
	}
	return name
}

func truncAddr(addr string) string {
	if len(addr) < 16 {
		return addr
	}
	return addr[:10] + "..." + addr[len(addr)-4:]
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
