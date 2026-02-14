package market

import (
	"github.com/born1337/hltui/internal/store"
	tea "github.com/charmbracelet/bubbletea"
)

// OI filter thresholds in USD. Pressing 'f' cycles through these.
var oiThresholds = []float64{
	0,           // no filter
	100_000,     // $100K
	500_000,     // $500K
	1_000_000,   // $1M (default)
	5_000_000,   // $5M
	10_000_000,  // $10M
	50_000_000,  // $50M
	100_000_000, // $100M
}

const defaultOIIndex = 3 // $1M

type Model struct {
	store      *store.Store
	scroll     int
	height     int
	sortAsc    bool
	oiFilterIdx int // index into oiThresholds
}

func New(s *store.Store) Model {
	return Model{store: s, oiFilterIdx: defaultOIIndex}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) OIThreshold() float64 {
	return oiThresholds[m.oiFilterIdx]
}

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
		case "f":
			m.oiFilterIdx = (m.oiFilterIdx + 1) % len(oiThresholds)
			m.scroll = 0
		}
	}
	return m, nil
}

func (m *Model) SetHeight(h int) {
	m.height = h
}
