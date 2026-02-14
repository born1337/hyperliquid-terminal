package ui

import (
	"fmt"
	"strings"

	"github.com/born1337/hyperliquid-terminal/internal/store"
	"github.com/born1337/hyperliquid-terminal/internal/style"
	"github.com/born1337/hyperliquid-terminal/internal/util"
	"github.com/charmbracelet/lipgloss"
)

func RenderHeader(s *store.Store, width int) string {
	s.RLock()
	defer s.RUnlock()

	if s.ClearinghouseState == nil {
		return style.Dim.Render("Loading account data...")
	}

	ms := s.ClearinghouseState.MarginSummary
	acctVal := util.ParseFloat(ms.AccountValue)
	posVal := util.ParseFloat(ms.TotalNtlPos)
	marginUsed := util.ParseFloat(ms.TotalMarginUsed)
	maintMargin := util.ParseFloat(s.ClearinghouseState.CrossMaintenanceMarginUsed)
	withdrawable := util.ParseFloat(s.ClearinghouseState.Withdrawable)

	leverage := 0.0
	marginRatio := 0.0
	if acctVal > 0 {
		leverage = posVal / acctVal
		marginRatio = (maintMargin / acctVal) * 100
	}

	line1 := fmt.Sprintf("Acct: %s  Pos: %s  Margin: %s  Lev: %s",
		style.Green.Render(util.FormatUSD(acctVal)),
		style.White.Render(util.FormatUSD(posVal)),
		style.White.Render(util.FormatUSD(marginUsed)),
		style.Magenta.Render(util.FormatLeverage(leverage)),
	)

	mrStyle := style.Green
	if marginRatio > 80 {
		mrStyle = style.Red
	} else if marginRatio > 50 {
		mrStyle = style.Yellow
	}

	line2 := fmt.Sprintf("Withdrawable: %s            Maint: %s    MR: %s",
		style.Green.Render(util.FormatUSD(withdrawable)),
		style.White.Render(util.FormatUSD(maintMargin)),
		mrStyle.Render(fmt.Sprintf("%.2f%%", marginRatio)),
	)

	content := lipgloss.JoinVertical(lipgloss.Left, line1, line2)
	return lipgloss.NewStyle().Width(width).Render(content)
}

func RenderTitleBar(width int, isTestnet, wsConnected bool, walletName, truncAddr string) string {
	title := style.White.Render("HLTUI v0.1.0")
	if walletName != "" {
		title += "  " + style.Cyan.Render(walletName) + " " + style.Dim.Render(truncAddr)
	}

	var indicators []string
	if isTestnet {
		indicators = append(indicators, style.Yellow.Render("[TESTNET]"))
	} else {
		indicators = append(indicators, style.Green.Render("[MAINNET]"))
	}

	if wsConnected {
		indicators = append(indicators, style.StatusConnected.Render("[WS: CONNECTED]"))
	} else {
		indicators = append(indicators, style.StatusDisconnected.Render("[WS: DISCONNECTED]"))
	}

	right := strings.Join(indicators, " ")
	gap := width - lipgloss.Width(title) - lipgloss.Width(right)
	if gap < 1 {
		gap = 1
	}

	return title + strings.Repeat(" ", gap) + right
}
