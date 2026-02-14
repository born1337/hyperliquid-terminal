package api

func (c *Client) GetPortfolio(user string) ([]PortfolioPeriod, error) {
	body, err := c.post(map[string]string{
		"type": "portfolio",
		"user": user,
	})
	if err != nil {
		return nil, err
	}
	return ParsePortfolio(body)
}
