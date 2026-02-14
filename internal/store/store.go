package store

import (
	"sort"
	"sync"

	"github.com/born1337/hltui/internal/api"
	"github.com/born1337/hltui/internal/util"
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

// PositionsSortedByPnl returns a copy of positions sorted by unrealized PnL descending.
// The returned slice and associated data (mids, funding rates) are snapshot copies
// safe to use without holding the store lock.
func (s *Store) PositionsSortedByPnl() ([]api.AssetPosition, map[string]string, map[string]float64) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ClearinghouseState == nil {
		return nil, nil, nil
	}
	positions := make([]api.AssetPosition, len(s.ClearinghouseState.AssetPositions))
	copy(positions, s.ClearinghouseState.AssetPositions)
	sort.Slice(positions, func(i, j int) bool {
		pi := util.ParseFloat(positions[i].Position.UnrealizedPnl)
		pj := util.ParseFloat(positions[j].Position.UnrealizedPnl)
		return pi > pj
	})

	// Copy mids and funding rates so caller doesn't need to hold lock
	mids := make(map[string]string, len(s.AllMids))
	for k, v := range s.AllMids {
		mids[k] = v
	}
	fundingRates := make(map[string]float64, len(s.FundingRates))
	for k, v := range s.FundingRates {
		fundingRates[k] = v
	}

	return positions, mids, fundingRates
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
