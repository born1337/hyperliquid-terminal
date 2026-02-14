package api

import "encoding/json"

func (c *Client) GetUserFunding(user string, startTime int64) ([]FundingPayment, error) {
	body, err := c.post(map[string]interface{}{
		"type":      "userFunding",
		"user":      user,
		"startTime": startTime,
	})
	if err != nil {
		return nil, err
	}

	// API returns [{time, hash, delta: {coin, usdc, szi, fundingRate}}, ...]
	var rawPayments []FundingPaymentRaw
	if err := json.Unmarshal(body, &rawPayments); err != nil {
		return nil, err
	}

	// Flatten into our internal type
	payments := make([]FundingPayment, len(rawPayments))
	for i, rp := range rawPayments {
		payments[i] = FundingPayment{
			Time:        rp.Time,
			Coin:        rp.Delta.Coin,
			Usdc:        rp.Delta.Usdc,
			Szi:         rp.Delta.Szi,
			FundingRate: rp.Delta.FundingRate,
		}
	}
	return payments, nil
}

func (c *Client) GetFundingHistory(coin string, startTime int64) ([]FundingHistoryEntry, error) {
	body, err := c.post(map[string]interface{}{
		"type":      "fundingHistory",
		"coin":      coin,
		"startTime": startTime,
	})
	if err != nil {
		return nil, err
	}
	var entries []FundingHistoryEntry
	if err := json.Unmarshal(body, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
