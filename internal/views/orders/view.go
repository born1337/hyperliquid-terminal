package orders

import (
	"fmt"
	"strings"

	"github.com/born1337/hltui/internal/style"
	"github.com/born1337/hltui/internal/util"
)

func (m Model) View() string {
	m.store.RLock()
	orders := m.store.OpenOrders
	m.store.RUnlock()

	if len(orders) == 0 {
		return style.Dim.Render("  No open orders")
	}

	var b strings.Builder

	header := fmt.Sprintf("%-9s %-6s %-8s %12s %12s %12s %-7s %-16s",
		"COIN", "SIDE", "TYPE", "SIZE", "PRICE", "TRIGGER", "REDUCE", "TIME",
	)
	b.WriteString(style.TableHeader.Render(header))
	b.WriteString("\n")

	visibleRows := m.height - 3
	if visibleRows < 1 {
		visibleRows = len(orders)
	}
	start := m.scroll
	if start >= len(orders) {
		start = len(orders) - 1
	}
	if start < 0 {
		start = 0
	}
	end := start + visibleRows
	if end > len(orders) {
		end = len(orders)
	}

	for _, o := range orders[start:end] {
		sideStyle := style.Green
		if o.Side == "A" || o.Side == "sell" {
			sideStyle = style.Red
		}
		side := "BUY"
		if o.Side == "A" || o.Side == "sell" {
			side = "SELL"
		}

		trigger := "-"
		if o.TriggerPx != "" && o.TriggerPx != "0" {
			trigger = util.FormatPrice(util.ParseFloat(o.TriggerPx))
		}

		reduce := " "
		if o.ReduceOnly {
			reduce = "yes"
		}

		orderType := o.OrderType
		if o.IsTrigger {
			orderType = "trigger"
		}

		row := fmt.Sprintf("%-9s %s %-8s %12s %12s %12s %-7s %-16s",
			style.White.Render(o.Coin),
			sideStyle.Render(fmt.Sprintf("%-6s", side)),
			orderType,
			util.FormatSize(util.ParseFloat(o.Sz)),
			util.FormatPrice(util.ParseFloat(o.LimitPx)),
			trigger,
			reduce,
			util.FormatTime(o.Timestamp),
		)
		b.WriteString(row)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(style.Dim.Render(fmt.Sprintf("  %d open orders", len(orders))))

	return b.String()
}
