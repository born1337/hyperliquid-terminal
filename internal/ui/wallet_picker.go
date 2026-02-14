package ui

import (
	"fmt"

	"github.com/born1337/hyperliquid-terminal/internal/config"
	"github.com/born1337/hyperliquid-terminal/internal/style"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func RenderWalletPicker(wallets []config.Wallet, activeIdx, cursorIdx, width, height int) string {
	title := style.White.Render("Switch Wallet")

	lines := []string{title, ""}

	for i, w := range wallets {
		truncAddr := config.TruncateAddress(w.Address)
		label := fmt.Sprintf("%s  %s", w.Name, style.Dim.Render(truncAddr))

		if w.Testnet {
			label += "  " + style.Yellow.Render("[T]")
		}
		if w.Vault {
			label += "  " + style.Magenta.Render("[V]")
		}

		activeMarker := "  "
		if i == activeIdx {
			activeMarker = style.Green.Render("* ")
		}

		if i == cursorIdx {
			line := style.Cyan.Render("▸ ") + activeMarker + label
			lines = append(lines, line)
		} else {
			line := "  " + activeMarker + label
			lines = append(lines, line)
		}
	}

	// "+ Add wallet" entry
	lines = append(lines, "")
	addLabel := style.Green.Render("+ Add wallet")
	if cursorIdx == len(wallets) {
		lines = append(lines, style.Cyan.Render("▸ ")+"  "+addLabel)
	} else {
		lines = append(lines, "    "+addLabel)
	}

	lines = append(lines, "")
	lines = append(lines, style.Dim.Render("j/k: navigate  enter: select"))
	lines = append(lines, style.Dim.Render("d: delete  esc: cancel"))

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

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

func RenderAddWalletForm(name, addr textinput.Model, focusField int, testnet, vault bool, errMsg string, width, height int) string {
	title := style.White.Render("Add Wallet")

	lines := []string{title, ""}

	// Name field
	nameLabel := "Name:    "
	if focusField == 0 {
		nameLabel = style.Cyan.Render(nameLabel)
	} else {
		nameLabel = style.Dim.Render(nameLabel)
	}
	lines = append(lines, nameLabel+name.View())

	// Address field
	addrLabel := "Address: "
	if focusField == 1 {
		addrLabel = style.Cyan.Render(addrLabel)
	} else {
		addrLabel = style.Dim.Render(addrLabel)
	}
	lines = append(lines, addrLabel+addr.View())

	lines = append(lines, "")

	// Testnet toggle
	testnetCheck := "[ ]"
	if testnet {
		testnetCheck = "[x]"
	}
	testnetStr := testnetCheck + " Testnet"
	if focusField == 2 {
		testnetStr = style.Cyan.Render(testnetStr)
	}

	// Vault toggle
	vaultCheck := "[ ]"
	if vault {
		vaultCheck = "[x]"
	}
	vaultStr := vaultCheck + " Vault"
	if focusField == 3 {
		vaultStr = style.Cyan.Render(vaultStr)
	}

	lines = append(lines, testnetStr+"    "+vaultStr)

	// Error message
	if errMsg != "" {
		lines = append(lines, "")
		lines = append(lines, style.Red.Render(errMsg))
	}

	lines = append(lines, "")
	lines = append(lines, style.Dim.Render("tab: next field  enter: save  esc: cancel"))

	content := lipgloss.JoinVertical(lipgloss.Left, lines...)

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
