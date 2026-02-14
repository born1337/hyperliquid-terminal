package ws

import "encoding/json"

// Subscription request
type SubRequest struct {
	Method       string      `json:"method"`
	Subscription interface{} `json:"subscription"`
}

// Incoming WS message envelope
type Message struct {
	Channel string          `json:"channel"`
	Data    json.RawMessage `json:"data"`
}

// allMids channel data
type AllMidsData struct {
	Mids map[string]string `json:"mids"`
}

// orderUpdates channel data
type OrderUpdate struct {
	Order  OrderInfo `json:"order"`
	Status string    `json:"status"`
	StatusTimestamp int64 `json:"statusTimestamp"`
}

type OrderInfo struct {
	Coin       string `json:"coin"`
	Side       string `json:"side"`
	LimitPx    string `json:"limitPx"`
	Sz         string `json:"sz"`
	Oid        int64  `json:"oid"`
	Timestamp  int64  `json:"timestamp"`
	OrigSz     string `json:"origSz"`
	OrderType  string `json:"orderType"`
	ReduceOnly bool   `json:"reduceOnly"`
}

// user channel data: fills, funding, etc. vary by subscription type
type UserEvent struct {
	Fills    []UserFillWs    `json:"fills,omitempty"`
	Fundings []UserFundingWs `json:"fundings,omitempty"`
}

type UserFillWs struct {
	Coin      string `json:"coin"`
	Px        string `json:"px"`
	Sz        string `json:"sz"`
	Side      string `json:"side"`
	Time      int64  `json:"time"`
	ClosedPnl string `json:"closedPnl"`
	Hash      string `json:"hash"`
	Fee       string `json:"fee"`
	Tid       int64  `json:"tid"`
	Dir       string `json:"dir"`
}

type UserFundingWs struct {
	Time        int64  `json:"time"`
	Coin        string `json:"coin"`
	Usdc        string `json:"usdc"`
	Szi         string `json:"szi"`
	FundingRate string `json:"fundingRate"`
}
