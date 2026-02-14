package api

import "encoding/json"

func (c *Client) GetClearinghouseState(user string) (*ClearinghouseState, error) {
	body, err := c.post(map[string]string{
		"type": "clearinghouseState",
		"user": user,
	})
	if err != nil {
		return nil, err
	}
	var state ClearinghouseState
	if err := json.Unmarshal(body, &state); err != nil {
		return nil, err
	}
	return &state, nil
}
