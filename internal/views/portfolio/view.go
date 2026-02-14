package portfolio

import (
	"fmt"
	"math"
	"strings"

	"github.com/born1337/hltui/internal/api"
	"github.com/born1337/hltui/internal/style"
	"github.com/born1337/hltui/internal/util"
)

func (m Model) View() string {
	var b strings.Builder

	m.store.RLock()
	periods := m.store.Portfolio
	fees := m.store.UserFees
	m.store.RUnlock()

	// Performance summary
	b.WriteString(style.White.Render("Performance Summary"))
	b.WriteString("\n")
	b.WriteString(style.Dim.Render(strings.Repeat("─", 50)))
	b.WriteString("\n\n")

	if len(periods) > 0 {
		// Find relevant periods
		periodMap := make(map[string]*api.PortfolioPeriod)
		for i := range periods {
			periodMap[periods[i].Name] = &periods[i]
		}

		// Helper to compute total PnL from a period's pnlHistory
		calcTotalPnl := func(p *api.PortfolioPeriod) float64 {
			if p == nil || len(p.PnlHistory) == 0 {
				return 0
			}
			// Last entry in pnlHistory is cumulative
			return util.ParseFloat(p.PnlHistory[len(p.PnlHistory)-1].Value)
		}

		dayPnl := calcTotalPnl(periodMap["perpDay"])
		weekPnl := calcTotalPnl(periodMap["perpWeek"])
		monthPnl := calcTotalPnl(periodMap["perpMonth"])
		allTimePnl := calcTotalPnl(periodMap["perpAllTime"])

		dayStyle := style.PnlColor(dayPnl)
		weekStyle := style.PnlColor(weekPnl)
		monthStyle := style.PnlColor(monthPnl)
		allTimeStyle := style.PnlColor(allTimePnl)

		fmt.Fprintf(&b, "  %s  %s\n", style.SummaryLabel.Render("Today PnL:"), dayStyle.Render(util.FormatSignedUSD(dayPnl)))
		fmt.Fprintf(&b, "  %s  %s\n", style.SummaryLabel.Render("7-Day PnL:"), weekStyle.Render(util.FormatSignedUSD(weekPnl)))
		fmt.Fprintf(&b, "  %s  %s\n", style.SummaryLabel.Render("30-Day PnL:"), monthStyle.Render(util.FormatSignedUSD(monthPnl)))
		fmt.Fprintf(&b, "  %s  %s\n", style.SummaryLabel.Render("All-Time PnL:"), allTimeStyle.Render(util.FormatSignedUSD(allTimePnl)))

		// Volume info
		if allTime := periodMap["perpAllTime"]; allTime != nil && allTime.Vlm != "" {
			vlm := util.ParseFloat(allTime.Vlm)
			fmt.Fprintf(&b, "  %s  %s\n", style.SummaryLabel.Render("All-Time Vol:"), style.Cyan.Render(util.FormatUSD(vlm)))
		}

		// Sparkline for daily PnL
		if day := periodMap["perpDay"]; day != nil && len(day.PnlHistory) > 0 {
			b.WriteString("\n")
			b.WriteString(style.White.Render("PnL History (Today)"))
			b.WriteString("\n")
			b.WriteString(renderSparkline(day.PnlHistory))
			b.WriteString("\n")
		}

		// Account value chart using allTime
		if allTime := periodMap["perpAllTime"]; allTime != nil && len(allTime.AccountValueHistory) > 0 {
			b.WriteString("\n")
			b.WriteString(style.White.Render("Account Value History"))
			b.WriteString("\n")
			b.WriteString(renderSparkline(allTime.AccountValueHistory))
			b.WriteString("\n")
		}
	} else {
		b.WriteString(style.Dim.Render("  No portfolio data available"))
		b.WriteString("\n")
	}

	// Fee info
	b.WriteString("\n")
	b.WriteString(style.White.Render("Fee Schedule"))
	b.WriteString("\n")
	b.WriteString(style.Dim.Render(strings.Repeat("─", 50)))
	b.WriteString("\n\n")

	if fees != nil {
		crossRate := util.ParseFloat(fees.UserCrossRate) * 100
		addRate := util.ParseFloat(fees.UserAddRate) * 100
		fmt.Fprintf(&b, "  %s  %.4f%%\n", style.SummaryLabel.Render("Taker Rate:"), crossRate)
		fmt.Fprintf(&b, "  %s  %.4f%%\n", style.SummaryLabel.Render("Maker Rate:"), addRate)

		if fees.FeeSchedule.Cross != "" {
			schedCross := util.ParseFloat(fees.FeeSchedule.Cross) * 100
			schedAdd := util.ParseFloat(fees.FeeSchedule.Add) * 100
			fmt.Fprintf(&b, "  %s  %.4f%%\n", style.SummaryLabel.Render("Sched Taker:"), schedCross)
			fmt.Fprintf(&b, "  %s  %.4f%%\n", style.SummaryLabel.Render("Sched Maker:"), schedAdd)
		}
	} else {
		b.WriteString(style.Dim.Render("  Fee data not available"))
	}

	return b.String()
}

func renderSparkline(history []api.TimeValue) string {
	if len(history) == 0 {
		return ""
	}

	// Use last N entries
	maxBars := 60
	start := len(history) - maxBars
	if start < 0 {
		start = 0
	}
	recent := history[start:]

	vals := make([]float64, len(recent))
	for i, tv := range recent {
		vals[i] = util.ParseFloat(tv.Value)
	}

	bars := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	minVal, maxVal := vals[0], vals[0]
	for _, v := range vals {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	rng := maxVal - minVal
	if rng == 0 {
		rng = 1
	}

	var b strings.Builder
	b.WriteString("  ")
	for _, v := range vals {
		idx := int(math.Round((v - minVal) / rng * 7))
		if idx > 7 {
			idx = 7
		}
		if idx < 0 {
			idx = 0
		}
		if v >= 0 {
			b.WriteString(style.Green.Render(string(bars[idx])))
		} else {
			b.WriteString(style.Red.Render(string(bars[idx])))
		}
	}
	return b.String()
}
