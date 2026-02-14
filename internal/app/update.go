package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		viewHeight := m.height - 12 // title + header + tabs + statusbar + blank lines
		m.market.SetHeight(viewHeight)
		m.positions.SetHeight(viewHeight)
		m.orders.SetHeight(viewHeight)
		m.fills.SetHeight(viewHeight)
		m.funding.SetHeight(viewHeight)
		m.portfolio.SetHeight(viewHeight)
		m.vaults.SetHeight(viewHeight)

	case InitialDataMsg:
		m.loading = false
		if msg.Err != nil {
			m.errMsg = "Error: " + msg.Err.Error()
			return m, refreshTick()
		}
		m.errMsg = ""

		m.store.Lock()
		m.store.ClearinghouseState = msg.State
		if msg.Mids != nil {
			m.store.AllMids = msg.Mids
		}
		m.store.MetaAndAssetCtxs = msg.Meta
		m.store.OpenOrders = msg.Orders
		m.store.Fills = msg.Fills
		m.store.FundingPayments = msg.Funding
		m.store.Portfolio = msg.Portfolio
		m.store.UserFees = msg.Fees
		m.store.VaultEquities = msg.Vaults
		m.store.Unlock()

		m.store.UpdateFundingRates()

		// Fetch vault details
		m.store.RLock()
		vaultEquities := m.store.VaultEquities
		m.store.RUnlock()
		for _, ve := range vaultEquities {
			cmds = append(cmds, m.fetchVaultDetails(ve.VaultAddress))
		}

		cmds = append(cmds, refreshTick())

	case wsConnectedMsg:
		m.ws = msg.client
		m.wsConnected = true
		cmds = append(cmds, waitForWS(m.wsCh))

	case WSMsg:
		m.handleWSMessage(msg.Msg)
		cmds = append(cmds, waitForWS(m.wsCh))

	case WSStatusMsg:
		m.wsConnected = msg.Connected

	case RefreshTickMsg:
		// Periodic refresh of account state
		cmds = append(cmds, m.fetchInitialData())
		if m.ws != nil {
			m.wsConnected = m.ws.Connected()
		}

	case VaultDetailsMsg:
		if msg.Err == nil && msg.Details != nil {
			m.store.Lock()
			m.store.VaultDetails[msg.Address] = msg.Details
			m.store.Unlock()
		}

	case tea.KeyMsg:
		// Add wallet form overlay captures all keys
		if m.walletFormActive {
			switch msg.String() {
			case "esc":
				m.walletFormActive = false
			case "tab":
				m.walletFormField = (m.walletFormField + 1) % 4
				m.focusWalletFormField()
			case "shift+tab":
				m.walletFormField = (m.walletFormField + 3) % 4
				m.focusWalletFormField()
			case "enter":
				if m.walletFormField == 2 {
					m.walletFormTestnet = !m.walletFormTestnet
				} else if m.walletFormField == 3 {
					m.walletFormVault = !m.walletFormVault
				} else {
					m.submitWalletForm()
				}
			case " ":
				if m.walletFormField == 2 {
					m.walletFormTestnet = !m.walletFormTestnet
				} else if m.walletFormField == 3 {
					m.walletFormVault = !m.walletFormVault
				}
			default:
				if m.walletFormField == 0 {
					m.walletFormName, _ = m.walletFormName.Update(msg)
				} else if m.walletFormField == 1 {
					m.walletFormAddr, _ = m.walletFormAddr.Update(msg)
				}
			}
			return m, tea.Batch(cmds...)
		}

		// Wallet picker overlay captures all keys
		if m.showWalletPicker {
			addIdx := len(m.cfg.Wallets) // index of "+ Add wallet" entry
			switch msg.String() {
			case "j", "down":
				if m.walletCursor < addIdx {
					m.walletCursor++
				}
			case "k", "up":
				if m.walletCursor > 0 {
					m.walletCursor--
				}
			case "enter":
				if m.walletCursor == addIdx {
					m.initWalletForm()
				} else {
					m.showWalletPicker = false
					cmd := m.switchWallet(m.walletCursor)
					if cmd != nil {
						cmds = append(cmds, cmd)
					}
				}
			case "d":
				m.deleteWalletAtCursor()
			case "esc", "w", "q":
				m.showWalletPicker = false
			}
			return m, tea.Batch(cmds...)
		}

		// Help overlay captures all keys
		if m.showHelp {
			if msg.String() == ";" || msg.String() == "esc" || msg.String() == "q" {
				m.showHelp = false
			}
			return m, nil
		}

		switch {
		case key.Matches(msg, Keys.Quit):
			if m.ws != nil {
				m.ws.Close()
			}
			return m, tea.Quit

		case key.Matches(msg, Keys.Help):
			m.showHelp = true

		case key.Matches(msg, Keys.Tab), key.Matches(msg, Keys.NextView):
			m.activeView = (m.activeView + 1) % 7

		case key.Matches(msg, Keys.ShiftTab), key.Matches(msg, Keys.PrevView):
			m.activeView = (m.activeView + 6) % 7

		case key.Matches(msg, Keys.View0):
			m.activeView = ViewMarket
		case key.Matches(msg, Keys.View1):
			m.activeView = ViewPositions
		case key.Matches(msg, Keys.View2):
			m.activeView = ViewOrders
		case key.Matches(msg, Keys.View3):
			m.activeView = ViewFills
		case key.Matches(msg, Keys.View4):
			m.activeView = ViewFunding
		case key.Matches(msg, Keys.View5):
			m.activeView = ViewPortfolio
		case key.Matches(msg, Keys.View6):
			m.activeView = ViewVaults

		case key.Matches(msg, Keys.WalletPicker):
			m.showWalletPicker = true
			m.walletCursor = m.cfg.ActiveWallet

		case key.Matches(msg, Keys.Refresh):
			m.loading = true
			cmds = append(cmds, m.fetchInitialData())

		default:
			// Dispatch to active sub-view
			switch m.activeView {
			case ViewMarket:
				var cmd tea.Cmd
				m.market, cmd = m.market.Update(msg)
				cmds = append(cmds, cmd)
			case ViewPositions:
				var cmd tea.Cmd
				m.positions, cmd = m.positions.Update(msg)
				cmds = append(cmds, cmd)
			case ViewOrders:
				var cmd tea.Cmd
				m.orders, cmd = m.orders.Update(msg)
				cmds = append(cmds, cmd)
			case ViewFills:
				var cmd tea.Cmd
				m.fills, cmd = m.fills.Update(msg)
				cmds = append(cmds, cmd)
			case ViewFunding:
				var cmd tea.Cmd
				m.funding, cmd = m.funding.Update(msg)
				cmds = append(cmds, cmd)
			case ViewPortfolio:
				var cmd tea.Cmd
				m.portfolio, cmd = m.portfolio.Update(msg)
				cmds = append(cmds, cmd)
			case ViewVaults:
				var cmd tea.Cmd
				m.vaults, cmd = m.vaults.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return m, tea.Batch(cmds...)
}
