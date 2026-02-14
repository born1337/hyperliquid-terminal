package app

import (
	"github.com/born1337/hyperliquid-terminal/internal/api"
	"github.com/born1337/hyperliquid-terminal/internal/ws"
)

// Initial data loaded from HTTP
type InitialDataMsg struct {
	State     *api.ClearinghouseState
	Mids      api.AllMids
	Meta      *api.MetaAndAssetCtxs
	Orders    []api.OpenOrder
	Fills     []api.Fill
	Funding   []api.FundingPayment
	Portfolio []api.PortfolioPeriod
	Fees      *api.UserFees
	Vaults    []api.VaultEquity
	Err       error
}

// WebSocket message received
type WSMsg struct {
	Msg ws.Message
}

// Periodic refresh tick
type RefreshTickMsg struct{}

// Vault details loaded
type VaultDetailsMsg struct {
	Address string
	Details *api.VaultDetails
	Err     error
}

// Error message
type ErrMsg struct {
	Err error
}

// WS connection status changed
type WSStatusMsg struct {
	Connected bool
}
