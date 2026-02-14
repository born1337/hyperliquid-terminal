package app

import (
	"encoding/json"
	"time"

	"github.com/born1337/hltui/internal/api"
	"github.com/born1337/hltui/internal/config"
	"github.com/born1337/hltui/internal/store"
	"github.com/born1337/hltui/internal/ws"
	"github.com/born1337/hltui/internal/views/fills"
	"github.com/born1337/hltui/internal/views/funding"
	"github.com/born1337/hltui/internal/views/market"
	"github.com/born1337/hltui/internal/views/orders"
	"github.com/born1337/hltui/internal/views/portfolio"
	"github.com/born1337/hltui/internal/views/positions"
	"github.com/born1337/hltui/internal/views/vaults"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	ViewMarket = iota
	ViewPositions
	ViewOrders
	ViewFills
	ViewFunding
	ViewPortfolio
	ViewVaults
)

type Model struct {
	cfg    *config.Config
	store  *store.Store
	api    *api.Client
	ws     *ws.Client
	wsCh   chan ws.Message

	activeView  int
	showHelp    bool
	width       int
	height      int
	errMsg      string
	wsConnected bool
	loading     bool

	// Sub-models
	market    market.Model
	positions positions.Model
	orders    orders.Model
	fills     fills.Model
	funding   funding.Model
	portfolio portfolio.Model
	vaults    vaults.Model
}

func NewModel(cfg *config.Config) Model {
	s := store.New()
	wsCh := make(chan ws.Message, 64)

	return Model{
		cfg:    cfg,
		store:  s,
		api:    api.NewClient(cfg.InfoURL()),
		wsCh:   wsCh,
		loading: true,

		market:    market.New(s),
		positions: positions.New(s),
		orders:    orders.New(s),
		fills:     fills.New(s),
		funding:   funding.New(s),
		portfolio: portfolio.New(s),
		vaults:    vaults.New(s),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchInitialData(),
		m.connectWS(),
	)
}

func (m Model) fetchInitialData() tea.Cmd {
	return func() tea.Msg {
		addr := m.cfg.Address

		state, err := m.api.GetClearinghouseState(addr)
		if err != nil {
			return InitialDataMsg{Err: err}
		}

		mids, _ := m.api.GetAllMids()
		meta, _ := m.api.GetMetaAndAssetCtxs()
		openOrders, _ := m.api.GetOpenOrders(addr)
		userFills, _ := m.api.GetUserFills(addr)

		// Funding from last 7 days
		weekAgo := time.Now().Add(-7 * 24 * time.Hour).UnixMilli()
		fundingPayments, _ := m.api.GetUserFunding(addr, weekAgo)

		portfolioEntries, _ := m.api.GetPortfolio(addr)
		fees, _ := m.api.GetUserFees(addr)
		vaultEquities, _ := m.api.GetUserVaultEquities(addr)

		return InitialDataMsg{
			State:     state,
			Mids:      mids,
			Meta:      meta,
			Orders:    openOrders,
			Fills:     userFills,
			Funding:   fundingPayments,
			Portfolio: portfolioEntries,
			Fees:      fees,
			Vaults:    vaultEquities,
		}
	}
}

func (m Model) connectWS() tea.Cmd {
	return func() tea.Msg {
		client := ws.NewClient(m.cfg.WSBaseURL, m.wsCh)
		if err := client.Connect(); err != nil {
			return WSStatusMsg{Connected: false}
		}

		// Subscribe
		client.Subscribe(ws.SubAllMids())
		client.Subscribe(ws.SubUserFills(m.cfg.Address))
		client.Subscribe(ws.SubUserFundings(m.cfg.Address))
		client.Subscribe(ws.SubOrderUpdates(m.cfg.Address))

		return wsConnectedMsg{client: client}
	}
}

type wsConnectedMsg struct {
	client *ws.Client
}

func waitForWS(ch chan ws.Message) tea.Cmd {
	return func() tea.Msg {
		msg := <-ch
		return WSMsg{Msg: msg}
	}
}

func refreshTick() tea.Cmd {
	return tea.Tick(30*time.Second, func(time.Time) tea.Msg {
		return RefreshTickMsg{}
	})
}

func (m *Model) handleWSMessage(msg ws.Message) {
	switch msg.Channel {
	case "allMids":
		var data ws.AllMidsData
		if err := json.Unmarshal(msg.Data, &data); err == nil {
			m.store.UpdateMids(data.Mids)
		}
	case "orderUpdates":
		// Refetch orders on any order update
		go func() {
			orders, err := m.api.GetOpenOrders(m.cfg.Address)
			if err == nil {
				m.store.Lock()
				m.store.OpenOrders = orders
				m.store.Unlock()
			}
		}()
	case "user":
		var event ws.UserEvent
		if err := json.Unmarshal(msg.Data, &event); err == nil {
			if len(event.Fills) > 0 {
				// Prepend new fills
				m.store.Lock()
				newFills := make([]api.Fill, len(event.Fills))
				for i, f := range event.Fills {
					newFills[i] = api.Fill{
						Coin:      f.Coin,
						Px:        f.Px,
						Sz:        f.Sz,
						Side:      f.Side,
						Time:      f.Time,
						ClosedPnl: f.ClosedPnl,
						Hash:      f.Hash,
						Fee:       f.Fee,
						Tid:       f.Tid,
						Dir:       f.Dir,
					}
				}
				m.store.Fills = append(newFills, m.store.Fills...)
				m.store.Unlock()
			}
		}
	}
}

func (m Model) fetchVaultDetails(addr string) tea.Cmd {
	return func() tea.Msg {
		details, err := m.api.GetVaultDetails(addr, m.cfg.Address)
		return VaultDetailsMsg{Address: addr, Details: details, Err: err}
	}
}
