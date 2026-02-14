package api

import "encoding/json"

func (c *Client) GetMetaAndAssetCtxs() (*MetaAndAssetCtxs, error) {
	body, err := c.post(map[string]string{
		"type": "metaAndAssetCtxs",
	})
	if err != nil {
		return nil, err
	}
	var result MetaAndAssetCtxs
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetPredictedFundings() (PredictedFundings, error) {
	body, err := c.post(map[string]string{
		"type": "predictedFundings",
	})
	if err != nil {
		return nil, err
	}
	var result PredictedFundings
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}
