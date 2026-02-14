package api

import "encoding/json"

func (c *Client) GetUserFills(user string) ([]Fill, error) {
	body, err := c.post(map[string]string{
		"type": "userFills",
		"user": user,
	})
	if err != nil {
		return nil, err
	}
	var fills []Fill
	if err := json.Unmarshal(body, &fills); err != nil {
		return nil, err
	}
	return fills, nil
}

func (c *Client) GetUserFillsByTime(user string, startTime int64) ([]Fill, error) {
	body, err := c.post(map[string]interface{}{
		"type":      "userFillsByTime",
		"user":      user,
		"startTime": startTime,
	})
	if err != nil {
		return nil, err
	}
	var fills []Fill
	if err := json.Unmarshal(body, &fills); err != nil {
		return nil, err
	}
	return fills, nil
}
