package api

import "encoding/json"

func (c *Client) GetAllMids() (AllMids, error) {
	body, err := c.post(map[string]string{
		"type": "allMids",
	})
	if err != nil {
		return nil, err
	}
	var mids AllMids
	if err := json.Unmarshal(body, &mids); err != nil {
		return nil, err
	}
	return mids, nil
}
