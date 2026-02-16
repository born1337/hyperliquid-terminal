package store

import (
	"sort"
	"sync"

	"github.com/born1337/hyperliquid-terminal/internal/api"
	"github.com/born1337/hyperliquid-terminal/internal/util"
	"github.com/born1337/hyperliquid-terminal/internal/ws"
)

type Store struct {
	mu sync.RWMutex

	// Account state
	ClearinghouseState *api.ClearinghouseState
	AllMids            api.AllMids
	MetaAndAssetCtxs   *api.MetaAndAssetCtxs

	// Per-view data
	OpenOrders      []api.OpenOrder
	Fills           []api.Fill
	FundingPayments []api.FundingPayment
	Portfolio       []api.PortfolioPeriod
	UserFees        *api.UserFees
	VaultEquities   []api.VaultEquity
	VaultDetails    map[string]*api.VaultDetails

	// Derived/cached
	FundingRates map[string]float64 // coin -> funding rate
}

func New() *Store {
	return &Store{
		AllMids:      make(api.AllMids),
		FundingRates: make(map[string]float64),
		VaultDetails: make(map[string]*api.VaultDetails),
	}
}

func (s *Store) Lock()    { s.mu.Lock() }
func (s *Store) Unlock()  { s.mu.Unlock() }
func (s *Store) RLock()   { s.mu.RLock() }
func (s *Store) RUnlock() { s.mu.RUnlock() }

// ClearUserData clears per-wallet data, preserving global market data (AllMids, Meta, FundingRates).
func (s *Store) ClearUserData() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ClearinghouseState = nil
	s.OpenOrders = nil
	s.Fills = nil
	s.FundingPayments = nil
	s.Portfolio = nil
	s.UserFees = nil
	s.VaultEquities = nil
	s.VaultDetails = make(map[string]*api.VaultDetails)
}

// ClearAll clears all data (used when switching networks).
func (s *Store) ClearAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ClearinghouseState = nil
	s.AllMids = make(api.AllMids)
	s.MetaAndAssetCtxs = nil
	s.OpenOrders = nil
	s.Fills = nil
	s.FundingPayments = nil
	s.Portfolio = nil
	s.UserFees = nil
	s.VaultEquities = nil
	s.VaultDetails = make(map[string]*api.VaultDetails)
	s.FundingRates = make(map[string]float64)
}

func (s *Store) UpdateMids(mids map[string]string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range mids {
		s.AllMids[k] = v
	}
}

func (s *Store) UpdateFundingRates() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.MetaAndAssetCtxs == nil {
		return
	}
	for i, asset := range s.MetaAndAssetCtxs.Meta.Universe {
		if i < len(s.MetaAndAssetCtxs.AssetCtxs) {
			s.FundingRates[asset.Name] = util.ParseFloat(s.MetaAndAssetCtxs.AssetCtxs[i].Funding)
		}
	}
}

// AccountValue returns the account value as float64.
func (s *Store) AccountValue() float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ClearinghouseState == nil {
		return 0
	}
	return util.ParseFloat(s.ClearinghouseState.MarginSummary.AccountValue)
}

// PositionsSorted returns a copy of positions sorted by unrealized PnL.
// The returned slice and associated data (mids, funding rates) are snapshot copies
// safe to use without holding the store lock.
func (s *Store) PositionsSorted(ascending bool) ([]api.AssetPosition, map[string]string, map[string]float64) {
	s.mu.RLock()
	if s.ClearinghouseState == nil {
		s.mu.RUnlock()
		return nil, nil, nil
	}
	positions := make([]api.AssetPosition, len(s.ClearinghouseState.AssetPositions))
	copy(positions, s.ClearinghouseState.AssetPositions)

	// Copy mids and funding rates so caller doesn't need to hold lock
	mids := make(map[string]string, len(s.AllMids))
	for k, v := range s.AllMids {
		mids[k] = v
	}
	fundingRates := make(map[string]float64, len(s.FundingRates))
	for k, v := range s.FundingRates {
		fundingRates[k] = v
	}
	s.mu.RUnlock()

	// Sort outside the lock
	sort.Slice(positions, func(i, j int) bool {
		pi := util.ParseFloat(positions[i].Position.UnrealizedPnl)
		pj := util.ParseFloat(positions[j].Position.UnrealizedPnl)
		if ascending {
			return pi < pj
		}
		return pi > pj
	})

	return positions, mids, fundingRates
}

// ApplyOrderUpdates applies incremental order updates from the WebSocket.
func (s *Store) ApplyOrderUpdates(updates []ws.OrderUpdate) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, u := range updates {
		switch u.Status {
		case "open", "replaced", "triggered":
			// Upsert: find by Oid or add new
			found := false
			for i, o := range s.OpenOrders {
				if o.Oid == u.Order.Oid {
					s.OpenOrders[i] = api.OpenOrder{
						Coin:      u.Order.Coin,
						Side:      u.Order.Side,
						LimitPx:   u.Order.LimitPx,
						Sz:        u.Order.Sz,
						Oid:       u.Order.Oid,
						Timestamp: u.Order.Timestamp,
						OrigSz:    u.Order.OrigSz,
						OrderType: u.Order.OrderType,
						ReduceOnly: u.Order.ReduceOnly,
					}
					found = true
					break
				}
			}
			if !found {
				s.OpenOrders = append(s.OpenOrders, api.OpenOrder{
					Coin:      u.Order.Coin,
					Side:      u.Order.Side,
					LimitPx:   u.Order.LimitPx,
					Sz:        u.Order.Sz,
					Oid:       u.Order.Oid,
					Timestamp: u.Order.Timestamp,
					OrigSz:    u.Order.OrigSz,
					OrderType: u.Order.OrderType,
					ReduceOnly: u.Order.ReduceOnly,
				})
			}
		case "filled", "canceled", "marginCanceled", "liquidatedCanceled":
			// Remove by Oid
			for i, o := range s.OpenOrders {
				if o.Oid == u.Order.Oid {
					s.OpenOrders = append(s.OpenOrders[:i], s.OpenOrders[i+1:]...)
					break
				}
			}
		}
	}
}

func (s *Store) MidPrice(coin string) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.AllMids[coin]; ok {
		return util.ParseFloat(v)
	}
	return 0
}

func (s *Store) FundingRate(coin string) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.FundingRates[coin]
}

// GetPortfolioPeriod returns a specific period by name (e.g. "allTime", "day").
func (s *Store) GetPortfolioPeriod(name string) *api.PortfolioPeriod {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := range s.Portfolio {
		if s.Portfolio[i].Name == name {
			return &s.Portfolio[i]
		}
	}
	return nil
}
