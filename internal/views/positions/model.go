package positions

import (
	"github.com/born1337/hyperliquid-terminal/internal/store"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	store   *store.Store
	scroll  int
	height  int
	sortAsc bool
}

func New(s *store.Store) Model {
	return Model{store: s}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.scroll++
		case "k", "up":
			if m.scroll > 0 {
				m.scroll--
			}
		case "g":
			m.scroll = 0
		case "s":
			m.sortAsc = !m.sortAsc
			m.scroll = 0
		}
	}
	return m, nil
}

func (m *Model) SetHeight(h int) {
	m.height = h
}
