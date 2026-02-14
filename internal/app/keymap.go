package app

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Quit     key.Binding
	Tab      key.Binding
	ShiftTab key.Binding
	NextView key.Binding
	PrevView key.Binding
	View0    key.Binding
	View1    key.Binding
	View2    key.Binding
	View3    key.Binding
	View4    key.Binding
	View5    key.Binding
	View6    key.Binding
	Up       key.Binding
	Down     key.Binding
	Refresh  key.Binding
	Help     key.Binding
}

var Keys = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next view"),
	),
	ShiftTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev view"),
	),
	NextView: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "next view"),
	),
	PrevView: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "prev view"),
	),
	View0: key.NewBinding(key.WithKeys("0"), key.WithHelp("0", "market")),
	View1: key.NewBinding(key.WithKeys("1"), key.WithHelp("1", "positions")),
	View2: key.NewBinding(key.WithKeys("2"), key.WithHelp("2", "orders")),
	View3: key.NewBinding(key.WithKeys("3"), key.WithHelp("3", "fills")),
	View4: key.NewBinding(key.WithKeys("4"), key.WithHelp("4", "funding")),
	View5: key.NewBinding(key.WithKeys("5"), key.WithHelp("5", "portfolio")),
	View6: key.NewBinding(key.WithKeys("6"), key.WithHelp("6", "vaults")),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("k/up", "scroll up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("j/down", "scroll down"),
	),
	Refresh: key.NewBinding(
		key.WithKeys("r"),
		key.WithHelp("r", "refresh"),
	),
	Help: key.NewBinding(
		key.WithKeys(";"),
		key.WithHelp(";", "help"),
	),
}
