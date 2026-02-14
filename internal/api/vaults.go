package api

import "encoding/json"

func (c *Client) GetUserVaultEquities(user string) ([]VaultEquity, error) {
	body, err := c.post(map[string]string{
		"type": "userVaultEquities",
		"user": user,
	})
	if err != nil {
		return nil, err
	}
	var equities []VaultEquity
	if err := json.Unmarshal(body, &equities); err != nil {
		return nil, err
	}
	return equities, nil
}

func (c *Client) GetVaultDetails(vaultAddress, user string) (*VaultDetails, error) {
	body, err := c.post(map[string]string{
		"type":         "vaultDetails",
		"vaultAddress": vaultAddress,
		"user":         user,
	})
	if err != nil {
		return nil, err
	}
	var details VaultDetails
	if err := json.Unmarshal(body, &details); err != nil {
		return nil, err
	}
	return &details, nil
}
