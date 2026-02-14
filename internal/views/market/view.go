package market

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/born1337/hyperliquid-terminal/internal/style"
	"github.com/born1337/hyperliquid-terminal/internal/util"
	"github.com/charmbracelet/lipgloss"
)

const (
	colRank   = 4
	colAsset  = 10
	colPrice  = 14
	colChg    = 10
	colVol    = 14
	colFund   = 10
	colOI     = 14
	colOracle = 14
)

type assetRow struct {
	name      string
	midPx     float64
	prevDayPx float64
	chg24h    float64
	volume    float64
	funding   float64
	openInt   float64
	oraclePx  float64
}

func (m Model) View() string {
	m.store.RLock()
	meta := m.store.MetaAndAssetCtxs
	mids := m.store.AllMids
	m.store.RUnlock()

	if meta == nil || len(meta.AssetCtxs) == 0 {
		return style.Dim.Render("  Loading market data...")
	}

	// Build rows
	var rows []assetRow
	for i, asset := range meta.Meta.Universe {
		if i >= len(meta.AssetCtxs) {
			break
		}
		ctx := meta.AssetCtxs[i]

		midPx := 0.0
		if v, ok := mids[asset.Name]; ok {
			midPx = util.ParseFloat(v)
		}

		prevDay := util.ParseFloat(ctx.PrevDayPx)
		chg := 0.0
		if prevDay > 0 {
			chg = ((midPx - prevDay) / prevDay) * 100
		}

		rows = append(rows, assetRow{
			name:      asset.Name,
			midPx:     midPx,
			prevDayPx: prevDay,
			chg24h:    chg,
			volume:    util.ParseFloat(ctx.DayNtlVlm),
			funding:   util.ParseFloat(ctx.Funding),
			openInt:   util.ParseFloat(ctx.OpenInterest),
			oraclePx:  util.ParseFloat(ctx.OraclePx),
		})
	}

	// Filter by OI threshold
	threshold := m.OIThreshold()
	if threshold > 0 {
		filtered := rows[:0]
		for _, r := range rows {
			if r.openInt*r.midPx >= threshold {
				filtered = append(filtered, r)
			}
		}
		rows = filtered
	}

	// Sort by 24H % change
	sort.Slice(rows, func(i, j int) bool {
		if m.sortAsc {
			return rows[i].chg24h < rows[j].chg24h
		}
		return rows[i].chg24h > rows[j].chg24h
	})

	arrow := " ▼"
	if m.sortAsc {
		arrow = " ▲"
	}

	var b strings.Builder

	// Header
	header := padRight("#", colRank) +
		padRight("ASSET", colAsset) +
		padLeft("PRICE", colPrice) +
		"  " + padLeft("24H %"+arrow, colChg) +
		"  " + padLeft("24H VOL", colVol) +
		"  " + padLeft("FUND/1H", colFund) +
		"  " + padLeft("OPEN INT", colOI) +
		"  " + padLeft("ORACLE", colOracle)
	b.WriteString(style.TableHeader.Render(header))
	b.WriteString("\n")

	// Visible range
	visibleRows := m.height - 3
	if visibleRows < 1 {
		visibleRows = len(rows)
	}
	start := m.scroll
	if start >= len(rows) {
		start = len(rows) - 1
	}
	if start < 0 {
		start = 0
	}
	end := start + visibleRows
	if end > len(rows) {
		end = len(rows)
	}

	for i, r := range rows[start:end] {
		rank := start + i + 1
		chgStyle := style.PnlColor(r.chg24h)
		fundStyle := style.PnlColor(r.funding)

		cells := []string{
			style.Dim.Render(padRight(fmt.Sprintf("%d", rank), colRank)),
			style.White.Render(padRight(r.name, colAsset)),
			padLeft(formatPrice(r.midPx), colPrice),
			"  ",
			chgStyle.Render(padLeft(formatChg(r.chg24h), colChg)),
			"  ",
			style.Cyan.Render(padLeft(formatCompact(r.volume), colVol)),
			"  ",
			fundStyle.Render(padLeft(util.FormatFundingRate(r.funding), colFund)),
			"  ",
			padLeft(formatCompact(r.openInt), colOI),
			"  ",
			style.Dim.Render(padLeft(formatPrice(r.oraclePx), colOracle)),
		}
		b.WriteString(strings.Join(cells, ""))
		b.WriteString("\n")
	}

	// Footer
	b.WriteString("\n")
	filterLabel := "OFF"
	if threshold > 0 {
		filterLabel = "≥" + formatCompact(threshold)
	}
	b.WriteString(style.Dim.Render(fmt.Sprintf("  %d assets  (sorted by 24h %%)  ", len(rows))))
	b.WriteString(style.Yellow.Render(fmt.Sprintf("[f] OI filter: %s", filterLabel)))

	return b.String()
}

func formatPrice(val float64) string {
	if val == 0 {
		return "-"
	}
	if val >= 10000 {
		return fmt.Sprintf("$%.2f", val)
	}
	if val >= 1 {
		return fmt.Sprintf("$%.2f", val)
	}
	if val >= 0.01 {
		return fmt.Sprintf("$%.4f", val)
	}
	// Find significant digits
	if val >= 0.0001 {
		return fmt.Sprintf("$%.6f", val)
	}
	return fmt.Sprintf("$%.8f", val)
}

func formatChg(val float64) string {
	if math.IsNaN(val) || math.IsInf(val, 0) {
		return "-"
	}
	if val >= 0 {
		return fmt.Sprintf("+%.2f%%", val)
	}
	return fmt.Sprintf("%.2f%%", val)
}

func formatCompact(val float64) string {
	if val == 0 {
		return "-"
	}
	abs := math.Abs(val)
	if abs >= 1_000_000_000 {
		return fmt.Sprintf("$%.2fB", val/1_000_000_000)
	}
	if abs >= 1_000_000 {
		return fmt.Sprintf("$%.2fM", val/1_000_000)
	}
	if abs >= 1_000 {
		return fmt.Sprintf("$%.1fK", val/1_000)
	}
	return fmt.Sprintf("$%.0f", val)
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
