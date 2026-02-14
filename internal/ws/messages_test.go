package ws

import (
	"encoding/json"
	"testing"
)

func TestAllMidsDataUnmarshal(t *testing.T) {
	raw := `{"mids": {"BTC": "91000.50", "ETH": "3400.00"}}`
	var data AllMidsData
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if data.Mids["BTC"] != "91000.50" {
		t.Errorf("BTC = %q, want 91000.50", data.Mids["BTC"])
	}
	if len(data.Mids) != 2 {
		t.Errorf("len = %d, want 2", len(data.Mids))
	}
}

func TestOrderUpdateUnmarshal(t *testing.T) {
	raw := `{
		"order": {
			"coin": "BTC",
			"side": "B",
			"limitPx": "85000.0",
			"sz": "0.5",
			"oid": 12345,
			"timestamp": 1770000000000,
			"origSz": "1.0",
			"orderType": "Limit",
			"reduceOnly": false
		},
		"status": "filled",
		"statusTimestamp": 1770000001000
	}`
	var update OrderUpdate
	if err := json.Unmarshal([]byte(raw), &update); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if update.Order.Coin != "BTC" {
		t.Errorf("Coin = %q, want BTC", update.Order.Coin)
	}
	if update.Status != "filled" {
		t.Errorf("Status = %q, want filled", update.Status)
	}
}

func TestUserEventFillsUnmarshal(t *testing.T) {
	raw := `{
		"fills": [
			{
				"coin": "ETH",
				"px": "3400.00",
				"sz": "1.0",
				"side": "B",
				"time": 1770000000000,
				"closedPnl": "50.00",
				"hash": "0xabc",
				"fee": "1.36",
				"tid": 123456,
				"dir": "Open Long"
			}
		]
	}`
	var event UserEvent
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if len(event.Fills) != 1 {
		t.Fatalf("Fills len = %d, want 1", len(event.Fills))
	}
	if event.Fills[0].Coin != "ETH" {
		t.Errorf("Fills[0].Coin = %q, want ETH", event.Fills[0].Coin)
	}
	if event.Fills[0].ClosedPnl != "50.00" {
		t.Errorf("Fills[0].ClosedPnl = %q, want 50.00", event.Fills[0].ClosedPnl)
	}
}

func TestUserEventFundingsUnmarshal(t *testing.T) {
	raw := `{
		"fundings": [
			{
				"time": 1770000000000,
				"coin": "BTC",
				"usdc": "-10.50",
				"szi": "1.0",
				"fundingRate": "0.0001"
			}
		]
	}`
	var event UserEvent
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if len(event.Fundings) != 1 {
		t.Fatalf("Fundings len = %d, want 1", len(event.Fundings))
	}
	if event.Fundings[0].Coin != "BTC" {
		t.Errorf("Fundings[0].Coin = %q, want BTC", event.Fundings[0].Coin)
	}
}

func TestMessageUnmarshal(t *testing.T) {
	raw := `{"channel": "allMids", "data": {"mids": {"BTC": "91000"}}}`
	var msg Message
	if err := json.Unmarshal([]byte(raw), &msg); err != nil {
		t.Fatalf("Unmarshal error: %v", err)
	}
	if msg.Channel != "allMids" {
		t.Errorf("Channel = %q, want allMids", msg.Channel)
	}
	if msg.Data == nil {
		t.Fatal("Data is nil")
	}
}

func TestSubRequestMarshal(t *testing.T) {
	tests := []struct {
		name string
		sub  SubRequest
		want string
	}{
		{
			name: "allMids",
			sub:  SubAllMids(),
			want: `{"method":"subscribe","subscription":{"type":"allMids"}}`,
		},
		{
			name: "userFills",
			sub:  SubUserFills("0xabc"),
			want: `{"method":"subscribe","subscription":{"type":"userFills","user":"0xabc"}}`,
		},
	}

	for _, tt := range tests {
		data, err := json.Marshal(tt.sub)
		if err != nil {
			t.Errorf("%s: Marshal error: %v", tt.name, err)
			continue
		}
		// Compare as maps to avoid key-order issues
		var got, want map[string]interface{}
		json.Unmarshal(data, &got)
		json.Unmarshal([]byte(tt.want), &want)
		gotJ, _ := json.Marshal(got)
		wantJ, _ := json.Marshal(want)
		if string(gotJ) != string(wantJ) {
			t.Errorf("%s: got %s, want %s", tt.name, string(data), tt.want)
		}
	}
}
