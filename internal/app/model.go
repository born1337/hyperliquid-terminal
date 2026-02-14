package app

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/born1337/hyperliquid-terminal/internal/api"
	"github.com/born1337/hyperliquid-terminal/internal/config"
	"github.com/born1337/hyperliquid-terminal/internal/store"
	"github.com/born1337/hyperliquid-terminal/internal/ws"
	"github.com/born1337/hyperliquid-terminal/internal/views/fills"
	"github.com/born1337/hyperliquid-terminal/internal/views/funding"
	"github.com/born1337/hyperliquid-terminal/internal/views/market"
	"github.com/born1337/hyperliquid-terminal/internal/views/orders"
	"github.com/born1337/hyperliquid-terminal/internal/views/portfolio"
	"github.com/born1337/hyperliquid-terminal/internal/views/positions"
	"github.com/born1337/hyperliquid-terminal/internal/views/vaults"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
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

	// Wallet picker
	showWalletPicker bool
	walletCursor     int

	// Add wallet form
	walletFormActive  bool
	walletFormName    textinput.Model
	walletFormAddr    textinput.Model
	walletFormField   int // 0=name, 1=addr, 2=testnet, 3=vault
	walletFormTestnet bool
	walletFormVault   bool
	walletFormErr     string

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
		// Capture address to guard against stale writes after wallet switch
		addr := m.cfg.Address
		go func() {
			orders, err := m.api.GetOpenOrders(addr)
			if err == nil && m.cfg.Address == addr {
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

// switchWallet closes WS, switches config, clears store, and re-fetches data.
func (m *Model) switchWallet(idx int) tea.Cmd {
	if idx == m.cfg.ActiveWallet {
		return nil
	}

	// Close existing WS
	if m.ws != nil {
		m.ws.Close()
		m.ws = nil
	}
	m.wsConnected = false

	// Switch config (may change network)
	networkChanged := m.cfg.SwitchToWallet(idx)

	// Clear store
	if networkChanged {
		m.store.ClearAll()
		m.api = api.NewClient(m.cfg.InfoURL())
	} else {
		m.store.ClearUserData()
	}

	// New WS channel so stale goroutines drain harmlessly into the old one
	m.wsCh = make(chan ws.Message, 64)

	// Reset per-wallet views
	m.resetViewScrolls()

	m.loading = true
	m.errMsg = ""

	return tea.Batch(
		m.fetchInitialData(),
		m.connectWS(),
	)
}

// resetViewScrolls resets per-wallet view state (scroll positions, etc.)
func (m *Model) resetViewScrolls() {
	m.positions = positions.New(m.store)
	m.orders = orders.New(m.store)
	m.fills = fills.New(m.store)
	m.funding = funding.New(m.store)
	m.portfolio = portfolio.New(m.store)
	m.vaults = vaults.New(m.store)
	// Preserve market view state (sort, scroll, filter)
}

// initWalletForm sets up the add-wallet form text inputs.
func (m *Model) initWalletForm() {
	nameInput := textinput.New()
	nameInput.Placeholder = "My Wallet"
	nameInput.CharLimit = 20
	nameInput.Width = 30
	nameInput.Focus()

	addrInput := textinput.New()
	addrInput.Placeholder = "0x..."
	addrInput.CharLimit = 42
	addrInput.Width = 30

	m.walletFormName = nameInput
	m.walletFormAddr = addrInput
	m.walletFormField = 0
	m.walletFormTestnet = false
	m.walletFormVault = false
	m.walletFormErr = ""
	m.walletFormActive = true
}

// focusWalletFormField updates text input focus based on the current field.
func (m *Model) focusWalletFormField() {
	m.walletFormName.Blur()
	m.walletFormAddr.Blur()
	switch m.walletFormField {
	case 0:
		m.walletFormName.Focus()
	case 1:
		m.walletFormAddr.Focus()
	}
}

// deleteWalletAtCursor deletes the wallet at the current cursor position.
func (m *Model) deleteWalletAtCursor() {
	idx := m.walletCursor
	if idx < 0 || idx >= len(m.cfg.Wallets) {
		return
	}
	// Can't delete the currently active wallet
	if idx == m.cfg.ActiveWallet {
		return
	}
	// Can't delete CLI wallet
	if m.cfg.Wallets[idx].Name == "CLI" {
		return
	}
	// Must keep at least 1 wallet
	if len(m.cfg.Wallets) <= 1 {
		return
	}

	// Remove from slice
	m.cfg.Wallets = append(m.cfg.Wallets[:idx], m.cfg.Wallets[idx+1:]...)

	// Adjust active wallet index if needed
	if m.cfg.ActiveWallet > idx {
		m.cfg.ActiveWallet--
	}

	// Adjust cursor
	if m.walletCursor >= len(m.cfg.Wallets) {
		m.walletCursor = len(m.cfg.Wallets) - 1
	}

	// Save persisted wallets (skip CLI if present)
	persistStart := 0
	if len(m.cfg.Wallets) > 0 && m.cfg.Wallets[0].Name == "CLI" {
		persistStart = 1
	}
	_ = config.SaveWallets(m.cfg.Wallets[persistStart:])
}

// submitWalletForm validates and saves a new wallet.
func (m *Model) submitWalletForm() {
	addr := strings.TrimSpace(m.walletFormAddr.Value())
	name := strings.TrimSpace(m.walletFormName.Value())

	if addr == "" {
		m.walletFormErr = "Address is required"
		return
	}
	if !config.ValidateAddress(addr) {
		m.walletFormErr = "Invalid address (need 0x + 40 hex chars)"
		return
	}
	if name == "" {
		name = config.TruncateAddress(addr)
	}

	w := config.Wallet{
		Name:    name,
		Address: addr,
		Testnet: m.walletFormTestnet,
		Vault:   m.walletFormVault,
	}

	// Determine which wallets are persisted (skip CLI wallet at index 0 if present)
	persistStart := 0
	if len(m.cfg.Wallets) > 0 && m.cfg.Wallets[0].Name == "CLI" {
		persistStart = 1
	}
	persisted := m.cfg.Wallets[persistStart:]

	if err := config.SaveWallets(append(persisted, w)); err != nil {
		m.walletFormErr = "Save failed: " + err.Error()
		return
	}

	m.cfg.Wallets = append(m.cfg.Wallets, w)
	m.walletFormActive = false
	m.walletCursor = len(m.cfg.Wallets) - 1
}
