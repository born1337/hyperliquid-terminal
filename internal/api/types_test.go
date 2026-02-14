package api

import (
	"encoding/json"
	"testing"
)

func TestClearinghouseStateUnmarshal(t *testing.T) {
	raw := `{
		"marginSummary": {
			"accountValue": "314096.11",
			"totalNtlPos": "1308154.22",
			"totalRawUsd": "52648.86",
			"totalMarginUsed": "261447.25"
		},
		"crossMaintenanceMarginUsed": "130361.27",
		"withdrawable": "34867.51",
		"assetPositions": [
			{
				"type": "oneWay",
				"position": {
					"coin": "BTC",
					"szi": "1.0982",
					"entryPx": "91001.0",
					"positionValue": "99789.95",
					"unrealizedPnl": "-148.25",
					"returnOnEquity": "-0.046",
					"leverage": {"type": "cross", "value": 31},
					"liquidationPx": null,
					"marginUsed": "3219.35",
					"maxLeverage": 50,
					"cumFunding": {
						"allTime": "100.50",
						"sinceOpen": "40.66",
						"sinceChange": "10.20"
					}
				}
			}
		]
	}`

	var state ClearinghouseState
	if err := json.Unmarshal([]byte(raw), &state); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if state.MarginSummary.AccountValue != "314096.11" {
		t.Errorf("AccountValue = %q, want %q", state.MarginSummary.AccountValue, "314096.11")
	}
	if state.Withdrawable != "34867.51" {
		t.Errorf("Withdrawable = %q, want %q", state.Withdrawable, "34867.51")
	}
	if len(state.AssetPositions) != 1 {
		t.Fatalf("AssetPositions len = %d, want 1", len(state.AssetPositions))
	}

	pos := state.AssetPositions[0].Position
	if pos.Coin != "BTC" {
		t.Errorf("Coin = %q, want %q", pos.Coin, "BTC")
	}
	if pos.Szi != "1.0982" {
		t.Errorf("Szi = %q, want %q", pos.Szi, "1.0982")
	}
	if pos.Leverage.Value != 31 {
		t.Errorf("Leverage.Value = %v, want 31", pos.Leverage.Value)
	}
	// liquidationPx should be nil for null
	if pos.LiquidationPx != nil {
		t.Errorf("LiquidationPx = %v, want nil", pos.LiquidationPx)
	}
	if pos.CumFunding == nil {
		t.Fatal("CumFunding is nil, want non-nil")
	}
	if pos.CumFunding.SinceOpen != "40.66" {
		t.Errorf("CumFunding.SinceOpen = %q, want %q", pos.CumFunding.SinceOpen, "40.66")
	}
}

func TestClearinghouseStateLiquidationPxPresent(t *testing.T) {
	raw := `{
		"marginSummary": {"accountValue": "100", "totalNtlPos": "100", "totalRawUsd": "0", "totalMarginUsed": "10"},
		"crossMaintenanceMarginUsed": "5",
		"withdrawable": "50",
		"assetPositions": [{
			"type": "oneWay",
			"position": {
				"coin": "ETH",
				"szi": "10",
				"entryPx": "3000",
				"positionValue": "30000",
				"unrealizedPnl": "500",
				"returnOnEquity": "0.1",
				"leverage": {"type": "isolated", "value": 5},
				"liquidationPx": "2500.00",
				"marginUsed": "6000",
				"maxLeverage": 25
			}
		}]
	}`

	var state ClearinghouseState
	if err := json.Unmarshal([]byte(raw), &state); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	pos := state.AssetPositions[0].Position
	if pos.LiquidationPx == nil {
		t.Fatal("LiquidationPx is nil, want non-nil")
	}
	if *pos.LiquidationPx != "2500.00" {
		t.Errorf("LiquidationPx = %q, want %q", *pos.LiquidationPx, "2500.00")
	}
}

func TestMetaAndAssetCtxsUnmarshal(t *testing.T) {
	raw := `[
		{"universe": [
			{"name": "BTC", "szDecimals": 5, "maxLeverage": 50, "onlyIsolated": false},
			{"name": "ETH", "szDecimals": 4, "maxLeverage": 50, "onlyIsolated": false}
		]},
		[
			{"funding": "0.0001", "openInterest": "1000", "prevDayPx": "90000", "dayNtlVlm": "5000000", "oraclePx": "91000", "markPx": "91001"},
			{"funding": "-0.00005", "openInterest": "500", "prevDayPx": "3000", "dayNtlVlm": "2000000", "oraclePx": "3100", "markPx": "3101"}
		]
	]`

	var result MetaAndAssetCtxs
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(result.Meta.Universe) != 2 {
		t.Fatalf("Universe len = %d, want 2", len(result.Meta.Universe))
	}
	if result.Meta.Universe[0].Name != "BTC" {
		t.Errorf("Universe[0].Name = %q, want BTC", result.Meta.Universe[0].Name)
	}
	if len(result.AssetCtxs) != 2 {
		t.Fatalf("AssetCtxs len = %d, want 2", len(result.AssetCtxs))
	}
	if result.AssetCtxs[0].Funding != "0.0001" {
		t.Errorf("AssetCtxs[0].Funding = %q, want 0.0001", result.AssetCtxs[0].Funding)
	}
	if result.AssetCtxs[1].Funding != "-0.00005" {
		t.Errorf("AssetCtxs[1].Funding = %q, want -0.00005", result.AssetCtxs[1].Funding)
	}
}

func TestFundingPaymentRawUnmarshal(t *testing.T) {
	raw := `{
		"time": 1770609600091,
		"hash": "0x0000000000000000000000000000000000000000000000000000000000000000",
		"delta": {
			"type": "funding",
			"coin": "kPEPE",
			"usdc": "-0.275994",
			"szi": "11770244.0",
			"fundingRate": "0.0000061016",
			"nSamples": null
		}
	}`

	var fp FundingPaymentRaw
	if err := json.Unmarshal([]byte(raw), &fp); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if fp.Time != 1770609600091 {
		t.Errorf("Time = %d, want 1770609600091", fp.Time)
	}
	if fp.Delta.Coin != "kPEPE" {
		t.Errorf("Delta.Coin = %q, want kPEPE", fp.Delta.Coin)
	}
	if fp.Delta.Usdc != "-0.275994" {
		t.Errorf("Delta.Usdc = %q, want -0.275994", fp.Delta.Usdc)
	}
	if fp.Delta.FundingRate != "0.0000061016" {
		t.Errorf("Delta.FundingRate = %q, want 0.0000061016", fp.Delta.FundingRate)
	}
	if fp.Delta.Szi != "11770244.0" {
		t.Errorf("Delta.Szi = %q, want 11770244.0", fp.Delta.Szi)
	}
}

func TestFundingPaymentArrayUnmarshal(t *testing.T) {
	raw := `[
		{"time": 100, "hash": "0x1", "delta": {"type": "funding", "coin": "BTC", "usdc": "-10.5", "szi": "1.0", "fundingRate": "0.0001"}},
		{"time": 200, "hash": "0x2", "delta": {"type": "funding", "coin": "ETH", "usdc": "5.25", "szi": "-2.0", "fundingRate": "-0.00005"}}
	]`

	var rawPayments []FundingPaymentRaw
	if err := json.Unmarshal([]byte(raw), &rawPayments); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}

	if len(rawPayments) != 2 {
		t.Fatalf("len = %d, want 2", len(rawPayments))
	}
	if rawPayments[0].Delta.Coin != "BTC" {
		t.Errorf("[0].Delta.Coin = %q, want BTC", rawPayments[0].Delta.Coin)
	}
	if rawPayments[1].Delta.Usdc != "5.25" {
		t.Errorf("[1].Delta.Usdc = %q, want 5.25", rawPayments[1].Delta.Usdc)
	}
}

func TestParsePortfolio(t *testing.T) {
	raw := `[
		["day", {
			"accountValueHistory": [[1770926342418, "100000.50"], [1770930000000, "100500.25"]],
			"pnlHistory": [[1770926342418, "0.0"], [1770930000000, "500.25"]],
			"vlm": "252695.71"
		}],
		["allTime", {
			"accountValueHistory": [[1763590320118, "50000.00"]],
			"pnlHistory": [[1763590320118, "0.0"]],
			"vlm": "76267102.44"
		}]
	]`

	periods, err := ParsePortfolio([]byte(raw))
	if err != nil {
		t.Fatalf("ParsePortfolio error: %v", err)
	}

	if len(periods) != 2 {
		t.Fatalf("len = %d, want 2", len(periods))
	}

	if periods[0].Name != "day" {
		t.Errorf("[0].Name = %q, want day", periods[0].Name)
	}
	if periods[0].Vlm != "252695.71" {
		t.Errorf("[0].Vlm = %q, want 252695.71", periods[0].Vlm)
	}
	if len(periods[0].PnlHistory) != 2 {
		t.Fatalf("[0].PnlHistory len = %d, want 2", len(periods[0].PnlHistory))
	}
	if periods[0].PnlHistory[0].Time != 1770926342418 {
		t.Errorf("[0].PnlHistory[0].Time = %d, want 1770926342418", periods[0].PnlHistory[0].Time)
	}
	if periods[0].PnlHistory[1].Value != "500.25" {
		t.Errorf("[0].PnlHistory[1].Value = %q, want 500.25", periods[0].PnlHistory[1].Value)
	}
	if len(periods[0].AccountValueHistory) != 2 {
		t.Fatalf("[0].AccountValueHistory len = %d, want 2", len(periods[0].AccountValueHistory))
	}

	if periods[1].Name != "allTime" {
		t.Errorf("[1].Name = %q, want allTime", periods[1].Name)
	}
	if periods[1].Vlm != "76267102.44" {
		t.Errorf("[1].Vlm = %q, want 76267102.44", periods[1].Vlm)
	}
}

func TestParsePortfolioEmpty(t *testing.T) {
	raw := `[]`
	periods, err := ParsePortfolio([]byte(raw))
	if err != nil {
		t.Fatalf("ParsePortfolio error: %v", err)
	}
	if len(periods) != 0 {
		t.Errorf("len = %d, want 0", len(periods))
	}
}

func TestAllMidsUnmarshal(t *testing.T) {
	raw := `{"BTC": "91234.50", "ETH": "3456.78", "SOL": "123.45"}`
	var mids AllMids
	if err := json.Unmarshal([]byte(raw), &mids); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if mids["BTC"] != "91234.50" {
		t.Errorf("BTC = %q, want 91234.50", mids["BTC"])
	}
	if len(mids) != 3 {
		t.Errorf("len = %d, want 3", len(mids))
	}
}

func TestOpenOrderUnmarshal(t *testing.T) {
	raw := `{
		"coin": "BTC",
		"side": "B",
		"limitPx": "85000.0",
		"sz": "0.5",
		"oid": 12345,
		"timestamp": 1770000000000,
		"origSz": "1.0",
		"orderType": "Limit",
		"isTrigger": false,
		"reduceOnly": true
	}`
	var order OpenOrder
	if err := json.Unmarshal([]byte(raw), &order); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if order.Coin != "BTC" {
		t.Errorf("Coin = %q, want BTC", order.Coin)
	}
	if order.Side != "B" {
		t.Errorf("Side = %q, want B", order.Side)
	}
	if !order.ReduceOnly {
		t.Error("ReduceOnly = false, want true")
	}
	if order.Oid != 12345 {
		t.Errorf("Oid = %d, want 12345", order.Oid)
	}
}

func TestFillUnmarshal(t *testing.T) {
	raw := `{
		"coin": "MON",
		"px": "0.023597",
		"sz": "406851.0",
		"side": "B",
		"time": 1771008838112,
		"startPosition": "1718276.0",
		"dir": "Open Long",
		"closedPnl": "0.0",
		"hash": "0x92b063e45d74fa7b",
		"oid": 320252381479,
		"crossed": true,
		"fee": "3.840185",
		"tid": 711604632708247,
		"feeToken": "USDC"
	}`
	var fill Fill
	if err := json.Unmarshal([]byte(raw), &fill); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if fill.Coin != "MON" {
		t.Errorf("Coin = %q, want MON", fill.Coin)
	}
	if fill.Fee != "3.840185" {
		t.Errorf("Fee = %q, want 3.840185", fill.Fee)
	}
	if !fill.Crossed {
		t.Error("Crossed = false, want true")
	}
}

func TestUserFeesUnmarshal(t *testing.T) {
	raw := `{
		"activeReferralDiscount": "0.0",
		"dailyUserVlm": [
			{"date": "2026-01-30", "userCross": "1860117.31", "userAdd": "43158.72", "exchange": "12587102816.56"}
		],
		"feeSchedule": {"cross": "0.00045", "add": "0.00015", "spotCross": "0.0006", "spotAdd": "0.0003"},
		"userCrossRate": "0.0004",
		"userAddRate": "0.00012"
	}`
	var fees UserFees
	if err := json.Unmarshal([]byte(raw), &fees); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if fees.UserCrossRate != "0.0004" {
		t.Errorf("UserCrossRate = %q, want 0.0004", fees.UserCrossRate)
	}
	if fees.UserAddRate != "0.00012" {
		t.Errorf("UserAddRate = %q, want 0.00012", fees.UserAddRate)
	}
	if fees.FeeSchedule.Cross != "0.00045" {
		t.Errorf("FeeSchedule.Cross = %q, want 0.00045", fees.FeeSchedule.Cross)
	}
	if len(fees.DailyUserVlm) != 1 {
		t.Fatalf("DailyUserVlm len = %d, want 1", len(fees.DailyUserVlm))
	}
	if fees.DailyUserVlm[0].Cross != "1860117.31" {
		t.Errorf("DailyUserVlm[0].Cross = %q, want 1860117.31", fees.DailyUserVlm[0].Cross)
	}
}

func TestVaultEquityUnmarshal(t *testing.T) {
	raw := `[
		{"vaultAddress": "0x010461c14e146ac35fe42271bdc1134ee31c703a", "equity": "80379985.9635280073", "lockedUntilTimestamp": 1767424986312}
	]`
	var equities []VaultEquity
	if err := json.Unmarshal([]byte(raw), &equities); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if len(equities) != 1 {
		t.Fatalf("len = %d, want 1", len(equities))
	}
	if equities[0].Equity != "80379985.9635280073" {
		t.Errorf("Equity = %q, want 80379985.9635280073", equities[0].Equity)
	}
}

func TestVaultDetailsUnmarshal(t *testing.T) {
	raw := `{
		"name": "HLP Strategy A",
		"vaultAddress": "0x010461c14e146ac35fe42271bdc1134ee31c703a",
		"leader": "0xdfc24b077bc1425ad1dea75bcb6f8158e10df303",
		"description": "A strategy",
		"leaderCommission": 0.1,
		"maxDistributable": 78558337.224312,
		"apr": 0.017889205660477817,
		"isClosed": false,
		"allowDeposits": true,
		"followerState": {
			"user": "0xdfc24b077bc1425ad1dea75bcb6f8158e10df303",
			"vaultEquity": "80380440.9702650607",
			"pnl": "1636917.2868910581",
			"allTimePnl": "5907351.8798060566",
			"daysFollowing": 982,
			"vaultEntryTime": 1683243598416,
			"lockupUntil": 1767424986312
		}
	}`
	var details VaultDetails
	if err := json.Unmarshal([]byte(raw), &details); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if details.Name != "HLP Strategy A" {
		t.Errorf("Name = %q, want HLP Strategy A", details.Name)
	}
	if details.APR != 0.017889205660477817 {
		t.Errorf("APR = %v, want 0.017889205660477817", details.APR)
	}
	if details.FollowerState == nil {
		t.Fatal("FollowerState is nil")
	}
	if details.FollowerState.Pnl != "1636917.2868910581" {
		t.Errorf("FollowerState.Pnl = %q, want 1636917.2868910581", details.FollowerState.Pnl)
	}
	if details.FollowerState.DaysFollowing != 982 {
		t.Errorf("DaysFollowing = %d, want 982", details.FollowerState.DaysFollowing)
	}
}

func TestTimeValueUnmarshal(t *testing.T) {
	raw := `[1770926342418, "100500.25"]`
	var tv TimeValue
	if err := json.Unmarshal([]byte(raw), &tv); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if tv.Time != 1770926342418 {
		t.Errorf("Time = %d, want 1770926342418", tv.Time)
	}
	if tv.Value != "100500.25" {
		t.Errorf("Value = %q, want 100500.25", tv.Value)
	}
}
