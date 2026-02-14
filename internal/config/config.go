package config

const Version = "0.1.0"

type Config struct {
	Address    string
	APIBaseURL string
	WSBaseURL  string
	IsTestnet  bool
	IsVault    bool
}

func New(address string, testnet, vault bool) *Config {
	apiBase := "https://api.hyperliquid.xyz"
	wsBase := "wss://api.hyperliquid.xyz/ws"
	if testnet {
		apiBase = "https://api.hyperliquid-testnet.xyz"
		wsBase = "wss://api.hyperliquid-testnet.xyz/ws"
	}
	return &Config{
		Address:    address,
		APIBaseURL: apiBase,
		WSBaseURL:  wsBase,
		IsTestnet:  testnet,
		IsVault:    vault,
	}
}

func (c *Config) InfoURL() string {
	return c.APIBaseURL + "/info"
}
