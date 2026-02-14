package store

import (
	"sync"
	"testing"

	"github.com/born1337/hltui/internal/api"
)

func TestNewStore(t *testing.T) {
	s := New()
	if s.AllMids == nil {
		t.Error("AllMids is nil")
	}
	if s.FundingRates == nil {
		t.Error("FundingRates is nil")
	}
	if s.VaultDetails == nil {
		t.Error("VaultDetails is nil")
	}
	if s.ClearinghouseState != nil {
		t.Error("ClearinghouseState should be nil initially")
	}
}

func TestUpdateMids(t *testing.T) {
	s := New()
	s.UpdateMids(map[string]string{
		"BTC": "91000.00",
		"ETH": "3400.00",
	})

	if s.AllMids["BTC"] != "91000.00" {
		t.Errorf("BTC = %q, want 91000.00", s.AllMids["BTC"])
	}

	// Update should merge, not replace
	s.UpdateMids(map[string]string{
		"BTC": "92000.00",
		"SOL": "150.00",
	})

	if s.AllMids["BTC"] != "92000.00" {
		t.Errorf("BTC = %q, want 92000.00", s.AllMids["BTC"])
	}
	if s.AllMids["ETH"] != "3400.00" {
		t.Errorf("ETH = %q, want 3400.00 (should be preserved)", s.AllMids["ETH"])
	}
	if s.AllMids["SOL"] != "150.00" {
		t.Errorf("SOL = %q, want 150.00", s.AllMids["SOL"])
	}
}

func TestUpdateMidsConcurrency(t *testing.T) {
	s := New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.UpdateMids(map[string]string{"BTC": "91000"})
		}()
	}
	wg.Wait()
	if s.AllMids["BTC"] != "91000" {
		t.Errorf("BTC = %q after concurrent updates", s.AllMids["BTC"])
	}
}

func TestUpdateFundingRates(t *testing.T) {
	s := New()
	s.MetaAndAssetCtxs = &api.MetaAndAssetCtxs{
		Meta: api.Meta{
			Universe: []api.AssetMeta{
				{Name: "BTC"},
				{Name: "ETH"},
			},
		},
		AssetCtxs: []api.AssetCtx{
			{Funding: "0.0001"},
			{Funding: "-0.00005"},
		},
	}

	s.UpdateFundingRates()

	if s.FundingRates["BTC"] != 0.0001 {
		t.Errorf("BTC rate = %v, want 0.0001", s.FundingRates["BTC"])
	}
	if s.FundingRates["ETH"] != -0.00005 {
		t.Errorf("ETH rate = %v, want -0.00005", s.FundingRates["ETH"])
	}
}

func TestUpdateFundingRatesNilMeta(t *testing.T) {
	s := New()
	// Should not panic
	s.UpdateFundingRates()
}

func TestUpdateFundingRatesMismatchedLengths(t *testing.T) {
	s := New()
	s.MetaAndAssetCtxs = &api.MetaAndAssetCtxs{
		Meta: api.Meta{
			Universe: []api.AssetMeta{
				{Name: "BTC"},
				{Name: "ETH"},
				{Name: "SOL"},
			},
		},
		AssetCtxs: []api.AssetCtx{
			{Funding: "0.0001"},
		},
	}
	// Should not panic with mismatched lengths
	s.UpdateFundingRates()
	if s.FundingRates["BTC"] != 0.0001 {
		t.Errorf("BTC rate = %v, want 0.0001", s.FundingRates["BTC"])
	}
	// ETH and SOL should not have rates set
	if _, ok := s.FundingRates["ETH"]; ok {
		t.Error("ETH should not have a funding rate")
	}
}

func TestAccountValue(t *testing.T) {
	s := New()

	// Nil state
	if v := s.AccountValue(); v != 0 {
		t.Errorf("AccountValue() = %v, want 0 for nil state", v)
	}

	s.ClearinghouseState = &api.ClearinghouseState{
		MarginSummary: api.MarginSummary{
			AccountValue: "314096.11",
		},
	}
	if v := s.AccountValue(); v != 314096.11 {
		t.Errorf("AccountValue() = %v, want 314096.11", v)
	}
}

func TestPositionsSortedByPnl(t *testing.T) {
	s := New()
	s.AllMids = api.AllMids{"BTC": "91000", "ETH": "3400"}
	s.FundingRates = map[string]float64{"BTC": 0.0001, "ETH": -0.00005}

	s.ClearinghouseState = &api.ClearinghouseState{
		AssetPositions: []api.AssetPosition{
			{Position: api.Position{Coin: "BTC", UnrealizedPnl: "-148.25"}},
			{Position: api.Position{Coin: "ETH", UnrealizedPnl: "500.00"}},
			{Position: api.Position{Coin: "SOL", UnrealizedPnl: "100.00"}},
		},
	}

	positions, mids, rates := s.PositionsSortedByPnl()

	if len(positions) != 3 {
		t.Fatalf("len = %d, want 3", len(positions))
	}

	// Should be sorted: ETH(500) > SOL(100) > BTC(-148)
	if positions[0].Position.Coin != "ETH" {
		t.Errorf("[0].Coin = %q, want ETH", positions[0].Position.Coin)
	}
	if positions[1].Position.Coin != "SOL" {
		t.Errorf("[1].Coin = %q, want SOL", positions[1].Position.Coin)
	}
	if positions[2].Position.Coin != "BTC" {
		t.Errorf("[2].Coin = %q, want BTC", positions[2].Position.Coin)
	}

	// Verify mids are copied
	if mids["BTC"] != "91000" {
		t.Errorf("mids[BTC] = %q, want 91000", mids["BTC"])
	}

	// Verify rates are copied
	if rates["BTC"] != 0.0001 {
		t.Errorf("rates[BTC] = %v, want 0.0001", rates["BTC"])
	}
}

func TestPositionsSortedByPnlNilState(t *testing.T) {
	s := New()
	positions, mids, rates := s.PositionsSortedByPnl()
	if positions != nil {
		t.Error("positions should be nil for nil state")
	}
	if mids != nil {
		t.Error("mids should be nil for nil state")
	}
	if rates != nil {
		t.Error("rates should be nil for nil state")
	}
}

func TestPositionsSortedByPnlDoesNotMutateOriginal(t *testing.T) {
	s := New()
	s.ClearinghouseState = &api.ClearinghouseState{
		AssetPositions: []api.AssetPosition{
			{Position: api.Position{Coin: "BTC", UnrealizedPnl: "-100"}},
			{Position: api.Position{Coin: "ETH", UnrealizedPnl: "200"}},
		},
	}

	_, _, _ = s.PositionsSortedByPnl()

	// Original should be unchanged
	if s.ClearinghouseState.AssetPositions[0].Position.Coin != "BTC" {
		t.Error("original positions were mutated")
	}
}

func TestMidPrice(t *testing.T) {
	s := New()
	s.AllMids["BTC"] = "91000.50"

	if v := s.MidPrice("BTC"); v != 91000.50 {
		t.Errorf("MidPrice(BTC) = %v, want 91000.50", v)
	}
	if v := s.MidPrice("NONEXISTENT"); v != 0 {
		t.Errorf("MidPrice(NONEXISTENT) = %v, want 0", v)
	}
}

func TestFundingRate(t *testing.T) {
	s := New()
	s.FundingRates["BTC"] = 0.0001

	if v := s.FundingRate("BTC"); v != 0.0001 {
		t.Errorf("FundingRate(BTC) = %v, want 0.0001", v)
	}
	if v := s.FundingRate("NONEXISTENT"); v != 0 {
		t.Errorf("FundingRate(NONEXISTENT) = %v, want 0", v)
	}
}

func TestGetPortfolioPeriod(t *testing.T) {
	s := New()
	s.Portfolio = []api.PortfolioPeriod{
		{Name: "day", Vlm: "1000"},
		{Name: "allTime", Vlm: "50000"},
	}

	p := s.GetPortfolioPeriod("allTime")
	if p == nil {
		t.Fatal("GetPortfolioPeriod(allTime) returned nil")
	}
	if p.Vlm != "50000" {
		t.Errorf("Vlm = %q, want 50000", p.Vlm)
	}

	if s.GetPortfolioPeriod("nonexistent") != nil {
		t.Error("GetPortfolioPeriod(nonexistent) should return nil")
	}
}
