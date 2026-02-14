package ui

import (
	"github.com/born1337/hltui/internal/style"
	"github.com/charmbracelet/lipgloss"
)

func RenderHelp(width, height int) string {
	title := style.White.Render("HLTUI Keyboard Shortcuts")

	sections := []string{
		title,
		"",
		style.Cyan.Render("Navigation"),
		"  " + style.Yellow.Render("Tab / Shift+Tab") + "  Cycle views",
		"  " + style.Yellow.Render("←/→ or h/l") + "       Switch views",
		"  " + style.Yellow.Render("0-6") + "              Jump to view",
		"  " + style.Yellow.Render("j/k or ↑/↓") + "      Scroll up/down",
		"",
		style.Cyan.Render("Actions"),
		"  " + style.Yellow.Render("s") + "  Toggle sort direction",
		"  " + style.Yellow.Render("f") + "  Cycle OI filter (Market view)",
		"  " + style.Yellow.Render("r") + "  Refresh all data",
		"  " + style.Yellow.Render(";") + "  Toggle this help",
		"  " + style.Yellow.Render("q") + "  Quit",
		"",
		style.Cyan.Render("Views"),
		"  " + style.White.Render("0: Market") + "      All assets overview",
		"  " + style.White.Render("1: Positions") + "   Open positions with PnL",
		"  " + style.White.Render("2: Orders") + "      Open/pending orders",
		"  " + style.White.Render("3: Fills") + "       Recent trade history",
		"  " + style.White.Render("4: Funding") + "     Funding rates & payments",
		"  " + style.White.Render("5: Portfolio") + "   Performance & fees",
		"  " + style.White.Render("6: Vaults") + "      Vault investments",
		"",
		style.Dim.Render("Press ; or Esc to close"),
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 3).
		Width(50)

	return lipgloss.Place(width, height,
		lipgloss.Center, lipgloss.Center,
		box.Render(content),
	)
}
