package app

import (
	"github.com/born1337/hltui/internal/ui"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Help overlay
	if m.showHelp {
		return ui.RenderHelp(m.width, m.height)
	}

	// Title bar
	titleBar := ui.RenderTitleBar(m.width, m.cfg.IsTestnet, m.wsConnected)

	// Account header
	header := ui.RenderHeader(m.store, m.width)

	// Tabs
	tabs := ui.RenderTabs(m.activeView, m.width)

	// Active view
	var viewContent string
	if m.loading {
		viewContent = "  Loading..."
	} else {
		switch m.activeView {
		case ViewMarket:
			viewContent = m.market.View()
		case ViewPositions:
			viewContent = m.positions.View()
		case ViewOrders:
			viewContent = m.orders.View()
		case ViewFills:
			viewContent = m.fills.View()
		case ViewFunding:
			viewContent = m.funding.View()
		case ViewPortfolio:
			viewContent = m.portfolio.View()
		case ViewVaults:
			viewContent = m.vaults.View()
		}
	}

	// Status bar
	statusBar := ui.RenderStatusBar(m.width, m.errMsg)

	return lipgloss.JoinVertical(lipgloss.Left,
		titleBar,
		"",
		header,
		"",
		tabs,
		"",
		viewContent,
		"",
		statusBar,
	)
}
