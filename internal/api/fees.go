package api

import "encoding/json"

func (c *Client) GetUserFees(user string) (*UserFees, error) {
	body, err := c.post(map[string]string{
		"type": "userFees",
		"user": user,
	})
	if err != nil {
		return nil, err
	}
	var fees UserFees
	if err := json.Unmarshal(body, &fees); err != nil {
		return nil, err
	}
	return &fees, nil
}
