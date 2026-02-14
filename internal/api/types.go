package api

import "encoding/json"

// clearinghouseState response
type ClearinghouseState struct {
	MarginSummary              MarginSummary   `json:"marginSummary"`
	CrossMaintenanceMarginUsed string          `json:"crossMaintenanceMarginUsed"`
	Withdrawable               string          `json:"withdrawable"`
	AssetPositions             []AssetPosition `json:"assetPositions"`
}

type MarginSummary struct {
	AccountValue    string `json:"accountValue"`
	TotalNtlPos     string `json:"totalNtlPos"`
	TotalRawUsd     string `json:"totalRawUsd"`
	TotalMarginUsed string `json:"totalMarginUsed"`
}

type AssetPosition struct {
	Type     string   `json:"type"`
	Position Position `json:"position"`
}

type Position struct {
	Coin           string      `json:"coin"`
	Szi            string      `json:"szi"`
	EntryPx        string      `json:"entryPx"`
	PositionValue  string      `json:"positionValue"`
	UnrealizedPnl  string      `json:"unrealizedPnl"`
	ReturnOnEquity string      `json:"returnOnEquity"`
	Leverage       Leverage    `json:"leverage"`
	LiquidationPx  *string     `json:"liquidationPx"` // nullable - null for cross-margin
	MarginUsed     string      `json:"marginUsed"`
	MaxLeverage    int         `json:"maxLeverage"`
	CumFunding     *CumFunding `json:"cumFunding,omitempty"`
}

type Leverage struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}

type CumFunding struct {
	AllTime     string `json:"allTime"`
	SinceOpen   string `json:"sinceOpen"`
	SinceChange string `json:"sinceChange"`
}

// allMids response: map[string]string
type AllMids map[string]string

// metaAndAssetCtxs response: [Meta, [AssetCtx...]]
type MetaAndAssetCtxs struct {
	Meta      Meta
	AssetCtxs []AssetCtx
}

type Meta struct {
	Universe []AssetMeta `json:"universe"`
}

type AssetMeta struct {
	Name         string `json:"name"`
	SzDecimals   int    `json:"szDecimals"`
	MaxLeverage  int    `json:"maxLeverage"`
	OnlyIsolated bool   `json:"onlyIsolated"`
}

type AssetCtx struct {
	Funding      string   `json:"funding"`
	OpenInterest string   `json:"openInterest"`
	PrevDayPx    string   `json:"prevDayPx"`
	DayNtlVlm   string   `json:"dayNtlVlm"`
	Premium      string   `json:"premium,omitempty"`
	OraclePx     string   `json:"oraclePx"`
	MarkPx       string   `json:"markPx"`
	MidPx        string   `json:"midPx,omitempty"`
	ImpactPxs    []string `json:"impactPxs,omitempty"`
}

func (m *MetaAndAssetCtxs) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if len(raw) < 2 {
		return nil
	}
	if err := json.Unmarshal(raw[0], &m.Meta); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[1], &m.AssetCtxs); err != nil {
		return err
	}
	return nil
}

// predictedFundings response
type PredictedFundings [][]interface{} // [[AssetMeta, PredictedFunding], ...]

type PredictedFunding struct {
	Venue   string `json:"venue"`
	Funding string `json:"funding"`
}

// frontendOpenOrders response
type OpenOrder struct {
	Coin             string `json:"coin"`
	Side             string `json:"side"`
	LimitPx          string `json:"limitPx"`
	Sz               string `json:"sz"`
	Oid              int64  `json:"oid"`
	Timestamp        int64  `json:"timestamp"`
	OrigSz           string `json:"origSz"`
	Cloid            string `json:"cloid,omitempty"`
	OrderType        string `json:"orderType"`
	TriggerPx        string `json:"triggerPx,omitempty"`
	IsTrigger        bool   `json:"isTrigger"`
	TriggerCondition string `json:"triggerCondition,omitempty"`
	ReduceOnly       bool   `json:"reduceOnly"`
	Children         []any  `json:"children,omitempty"`
}

// userFills response
type Fill struct {
	Coin          string `json:"coin"`
	Px            string `json:"px"`
	Sz            string `json:"sz"`
	Side          string `json:"side"`
	Time          int64  `json:"time"`
	StartPosition string `json:"startPosition"`
	Dir           string `json:"dir"`
	ClosedPnl     string `json:"closedPnl"`
	Hash          string `json:"hash"`
	Oid           int64  `json:"oid"`
	Crossed       bool   `json:"crossed"`
	Fee           string `json:"fee"`
	Tid           int64  `json:"tid"`
	FeeToken      string `json:"feeToken"`
}

// userFunding response - API returns {time, hash, delta: {type, coin, usdc, szi, fundingRate}}
type FundingPaymentRaw struct {
	Time  int64              `json:"time"`
	Hash  string             `json:"hash"`
	Delta FundingPaymentData `json:"delta"`
}

type FundingPaymentData struct {
	Type        string `json:"type"`
	Coin        string `json:"coin"`
	Usdc        string `json:"usdc"`
	Szi         string `json:"szi"`
	FundingRate string `json:"fundingRate"`
}

// FundingPayment is the flattened form used internally
type FundingPayment struct {
	Time        int64
	Coin        string
	Usdc        string
	Szi         string
	FundingRate string
}

// fundingHistory response
type FundingHistoryEntry struct {
	Coin        string `json:"coin"`
	FundingRate string `json:"fundingRate"`
	Premium     string `json:"premium"`
	Time        int64  `json:"time"`
}

// portfolio response: [["day", {accountValueHistory, pnlHistory, vlm}], ...]
type PortfolioPeriod struct {
	Name                string
	AccountValueHistory []TimeValue
	PnlHistory          []TimeValue
	Vlm                 string
}

type TimeValue struct {
	Time  int64
	Value string
}

func (tv *TimeValue) UnmarshalJSON(data []byte) error {
	var raw [2]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[0], &tv.Time); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[1], &tv.Value); err != nil {
		return err
	}
	return nil
}

type portfolioPeriodData struct {
	AccountValueHistory []TimeValue `json:"accountValueHistory"`
	PnlHistory          []TimeValue `json:"pnlHistory"`
	Vlm                 string      `json:"vlm"`
}

// ParsePortfolio parses the raw portfolio response which is [[name, data], ...]
func ParsePortfolio(data []byte) ([]PortfolioPeriod, error) {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	var periods []PortfolioPeriod
	for _, entry := range raw {
		var pair [2]json.RawMessage
		if err := json.Unmarshal(entry, &pair); err != nil {
			continue
		}

		var name string
		if err := json.Unmarshal(pair[0], &name); err != nil {
			continue
		}

		var pd portfolioPeriodData
		if err := json.Unmarshal(pair[1], &pd); err != nil {
			continue
		}

		periods = append(periods, PortfolioPeriod{
			Name:                name,
			AccountValueHistory: pd.AccountValueHistory,
			PnlHistory:          pd.PnlHistory,
			Vlm:                 pd.Vlm,
		})
	}
	return periods, nil
}

// userFees response
type UserFees struct {
	ActiveReferralDiscount string        `json:"activeReferralDiscount"`
	DailyUserVlm           []DailyVolume `json:"dailyUserVlm"`
	FeeSchedule            FeeSchedule   `json:"feeSchedule"`
	UserCrossRate          string        `json:"userCrossRate"`
	UserAddRate            string        `json:"userAddRate"`
}

type DailyVolume struct {
	Date  string `json:"date"`
	Cross string `json:"userCross"`
	Add   string `json:"userAdd"`
}

type FeeSchedule struct {
	Cross string `json:"cross"`
	Add   string `json:"add"`
}

// userVaultEquities response
type VaultEquity struct {
	VaultAddress         string `json:"vaultAddress"`
	Equity               string `json:"equity"`
	LockedUntilTimestamp int64  `json:"lockedUntilTimestamp"`
}

// vaultDetails response
type VaultDetails struct {
	Name             string         `json:"name"`
	VaultAddress     string         `json:"vaultAddress"`
	Leader           string         `json:"leader"`
	Description      string         `json:"description"`
	LeaderCommission float64        `json:"leaderCommission"`
	MaxDistributable float64        `json:"maxDistributable"`
	APR              float64        `json:"apr"`
	IsClosed         bool           `json:"isClosed"`
	AllowDeposits    bool           `json:"allowDeposits"`
	FollowerState    *FollowerState `json:"followerState,omitempty"`
}

type FollowerState struct {
	User           string `json:"user"`
	VaultEquity    string `json:"vaultEquity"`
	Pnl            string `json:"pnl"`
	AllTimePnl     string `json:"allTimePnl"`
	DaysFollowing  int    `json:"daysFollowing"`
	VaultEntryTime int64  `json:"vaultEntryTime"`
	LockupUntil    int64  `json:"lockupUntil"`
}
