package ui

import (
	"fmt"
	"strings"

	"github.com/born1337/hyperliquid-terminal/internal/style"
)

var TabNames = []string{
	"Market",
	"Positions",
	"Orders",
	"Fills",
	"Funding",
	"Portfolio",
	"Vaults",
}

func RenderTabs(activeIdx int, width int) string {
	var tabs []string
	for i, name := range TabNames {
		label := fmt.Sprintf("%d:%s", i, name)
		if i == activeIdx {
			tabs = append(tabs, style.ActiveTab.Render(label))
		} else {
			tabs = append(tabs, style.InactiveTab.Render(label))
		}
	}
	return strings.Join(tabs, " ")
}
