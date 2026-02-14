package api

import (
	"testing"
	"time"
)

// Integration tests that hit the real Hyperliquid API.
// Run with: go test -v -run TestIntegration -tags integration ./internal/api/
// These tests use a known public address for verification.

const testAddr = "0x4cb5f4d145cd16460932bbb9b871bb6fd5db97e3"
const testAPIURL = "https://api.hyperliquid.xyz/info"

func newTestClient() *Client {
	return NewClient(testAPIURL)
}

func TestIntegrationClearinghouseState(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := newTestClient()
	state, err := c.GetClearinghouseState(testAddr)
	if err != nil {
		t.Fatalf("GetClearinghouseState error: %v", err)
	}

	if state.MarginSummary.AccountValue == "" {
		t.Error("AccountValue is empty")
	}
	if state.MarginSummary.TotalNtlPos == "" {
		t.Error("TotalNtlPos is empty")
	}
	if state.Withdrawable == "" {
		t.Error("Withdrawable is empty")
	}

	t.Logf("Account value: %s", state.MarginSummary.AccountValue)
	t.Logf("Positions: %d", len(state.AssetPositions))

	if len(state.AssetPositions) > 0 {
		pos := state.AssetPositions[0].Position
		if pos.Coin == "" {
			t.Error("First position Coin is empty")
		}
		if pos.Szi == "" {
			t.Error("First position Szi is empty")
		}
		t.Logf("First position: %s size=%s pnl=%s", pos.Coin, pos.Szi, pos.UnrealizedPnl)

		// Verify liquidationPx handling (null for cross positions)
		t.Logf("First position liqPx is nil: %v", pos.LiquidationPx == nil)
	}
}

func TestIntegrationAllMids(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := newTestClient()
	mids, err := c.GetAllMids()
	if err != nil {
		t.Fatalf("GetAllMids error: %v", err)
	}

	if len(mids) == 0 {
		t.Error("AllMids is empty")
	}

	// BTC and ETH should always exist
	if _, ok := mids["BTC"]; !ok {
		t.Error("BTC not found in AllMids")
	}
	if _, ok := mids["ETH"]; !ok {
		t.Error("ETH not found in AllMids")
	}

	t.Logf("Number of assets: %d, BTC=%s, ETH=%s", len(mids), mids["BTC"], mids["ETH"])
}

func TestIntegrationMetaAndAssetCtxs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := newTestClient()
	result, err := c.GetMetaAndAssetCtxs()
	if err != nil {
		t.Fatalf("GetMetaAndAssetCtxs error: %v", err)
	}

	if len(result.Meta.Universe) == 0 {
		t.Error("Universe is empty")
	}
	if len(result.AssetCtxs) == 0 {
		t.Error("AssetCtxs is empty")
	}
	if len(result.Meta.Universe) != len(result.AssetCtxs) {
		t.Errorf("Universe len (%d) != AssetCtxs len (%d)",
			len(result.Meta.Universe), len(result.AssetCtxs))
	}

	// First asset should be BTC
	if result.Meta.Universe[0].Name != "BTC" {
		t.Errorf("First asset = %q, want BTC", result.Meta.Universe[0].Name)
	}
	if result.AssetCtxs[0].Funding == "" {
		t.Error("First AssetCtx Funding is empty")
	}

	t.Logf("Assets: %d, BTC funding=%s", len(result.Meta.Universe), result.AssetCtxs[0].Funding)
}

func TestIntegrationOpenOrders(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := newTestClient()
	orders, err := c.GetOpenOrders(testAddr)
	if err != nil {
		t.Fatalf("GetOpenOrders error: %v", err)
	}

	// May be empty, that's fine
	t.Logf("Open orders: %d", len(orders))
	for _, o := range orders {
		if o.Coin == "" {
			t.Error("Order Coin is empty")
		}
		t.Logf("  %s %s %s @ %s", o.Coin, o.Side, o.Sz, o.LimitPx)
	}
}

func TestIntegrationUserFills(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := newTestClient()
	fills, err := c.GetUserFills(testAddr)
	if err != nil {
		t.Fatalf("GetUserFills error: %v", err)
	}

	if len(fills) == 0 {
		t.Skip("No fills for test address")
	}

	fill := fills[0]
	if fill.Coin == "" {
		t.Error("Fill Coin is empty")
	}
	if fill.Px == "" {
		t.Error("Fill Px is empty")
	}
	if fill.Fee == "" {
		t.Error("Fill Fee is empty")
	}
	if fill.Time == 0 {
		t.Error("Fill Time is zero")
	}

	t.Logf("Fills: %d, latest: %s %s %s @ %s, fee=%s", len(fills), fill.Coin, fill.Side, fill.Sz, fill.Px, fill.Fee)
}

func TestIntegrationUserFunding(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := newTestClient()
	weekAgo := time.Now().Add(-7 * 24 * time.Hour).UnixMilli()
	payments, err := c.GetUserFunding(testAddr, weekAgo)
	if err != nil {
		t.Fatalf("GetUserFunding error: %v", err)
	}

	if len(payments) == 0 {
		t.Skip("No funding payments for test address")
	}

	fp := payments[0]
	// Verify the delta fields are properly flattened
	if fp.Coin == "" {
		t.Error("FundingPayment Coin is empty (delta flattening likely broken)")
	}
	if fp.Usdc == "" {
		t.Error("FundingPayment Usdc is empty (delta flattening likely broken)")
	}
	if fp.FundingRate == "" {
		t.Error("FundingPayment FundingRate is empty")
	}
	if fp.Time == 0 {
		t.Error("FundingPayment Time is zero")
	}

	t.Logf("Funding payments: %d, first: coin=%s usdc=%s rate=%s", len(payments), fp.Coin, fp.Usdc, fp.FundingRate)
}

func TestIntegrationPortfolio(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := newTestClient()
	periods, err := c.GetPortfolio(testAddr)
	if err != nil {
		t.Fatalf("GetPortfolio error: %v", err)
	}

	if len(periods) == 0 {
		t.Fatal("Portfolio returned 0 periods")
	}

	// Should have day, week, month, allTime, perpDay, perpWeek, perpMonth, perpAllTime
	periodNames := make(map[string]bool)
	for _, p := range periods {
		periodNames[p.Name] = true
		t.Logf("Period %q: pnlHistory=%d entries, vlm=%s", p.Name, len(p.PnlHistory), p.Vlm)
	}

	for _, expected := range []string{"day", "allTime", "perpDay", "perpAllTime"} {
		if !periodNames[expected] {
			t.Errorf("Missing period %q", expected)
		}
	}

	// Verify PnlHistory has data
	for _, p := range periods {
		if p.Name == "perpAllTime" && len(p.PnlHistory) == 0 {
			t.Error("perpAllTime pnlHistory is empty")
		}
		if p.Name == "perpAllTime" && len(p.PnlHistory) > 0 {
			tv := p.PnlHistory[0]
			if tv.Time == 0 {
				t.Error("PnlHistory[0].Time is zero")
			}
			t.Logf("perpAllTime first pnl entry: time=%d value=%s", tv.Time, tv.Value)
		}
	}
}

func TestIntegrationUserFees(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := newTestClient()
	fees, err := c.GetUserFees(testAddr)
	if err != nil {
		t.Fatalf("GetUserFees error: %v", err)
	}

	if fees.UserCrossRate == "" {
		t.Error("UserCrossRate is empty")
	}
	if fees.UserAddRate == "" {
		t.Error("UserAddRate is empty")
	}
	if fees.FeeSchedule.Cross == "" {
		t.Error("FeeSchedule.Cross is empty")
	}

	t.Logf("Fees: taker=%s, maker=%s, schedule.cross=%s", fees.UserCrossRate, fees.UserAddRate, fees.FeeSchedule.Cross)
}

func TestIntegrationUserVaultEquities(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	c := newTestClient()
	equities, err := c.GetUserVaultEquities(testAddr)
	if err != nil {
		t.Fatalf("GetUserVaultEquities error: %v", err)
	}

	t.Logf("Vault equities: %d", len(equities))
	for _, ve := range equities {
		if ve.VaultAddress == "" {
			t.Error("VaultAddress is empty")
		}
		if ve.Equity == "" {
			t.Error("Equity is empty")
		}
		t.Logf("  vault=%s equity=%s", ve.VaultAddress[:10]+"...", ve.Equity)
	}
}
