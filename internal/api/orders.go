package api

import "encoding/json"

func (c *Client) GetOpenOrders(user string) ([]OpenOrder, error) {
	body, err := c.post(map[string]string{
		"type": "frontendOpenOrders",
		"user": user,
	})
	if err != nil {
		return nil, err
	}
	var orders []OpenOrder
	if err := json.Unmarshal(body, &orders); err != nil {
		return nil, err
	}
	return orders, nil
}
