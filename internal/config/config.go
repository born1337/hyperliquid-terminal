package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const Version = "0.1.0"

// Wallet represents a named wallet entry from the config file.
type Wallet struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Testnet bool   `json:"testnet,omitempty"`
	Vault   bool   `json:"vault,omitempty"`
}

type walletsFile struct {
	Wallets []Wallet `json:"wallets"`
}

type Config struct {
	Address    string
	APIBaseURL string
	WSBaseURL  string
	IsTestnet  bool
	IsVault    bool

	Wallets      []Wallet
	ActiveWallet int
	WalletName   string
}

func New(address string, testnet, vault bool) *Config {
	cfg := &Config{
		Address:   address,
		IsTestnet: testnet,
		IsVault:   vault,
	}
	cfg.setURLs()
	return cfg
}

// NewWithWallets creates a Config from a wallet list.
func NewWithWallets(wallets []Wallet, initialIdx int, cliTestnet, cliVault bool) *Config {
	w := wallets[initialIdx]
	testnet := w.Testnet || cliTestnet
	vault := w.Vault || cliVault

	cfg := &Config{
		Address:      w.Address,
		IsTestnet:    testnet,
		IsVault:      vault,
		Wallets:      wallets,
		ActiveWallet: initialIdx,
		WalletName:   w.Name,
	}
	cfg.setURLs()
	return cfg
}

// SwitchToWallet changes the active wallet and recalculates URLs.
// Returns true if the network changed (requiring full reconnect).
func (c *Config) SwitchToWallet(idx int) (networkChanged bool) {
	if idx < 0 || idx >= len(c.Wallets) {
		return false
	}
	w := c.Wallets[idx]
	oldTestnet := c.IsTestnet

	c.Address = w.Address
	c.IsTestnet = w.Testnet
	c.IsVault = w.Vault
	c.ActiveWallet = idx
	c.WalletName = w.Name
	c.setURLs()

	return c.IsTestnet != oldTestnet
}

// TruncatedAddress returns address as 0x1234...5678.
func (c *Config) TruncatedAddress() string {
	return TruncateAddress(c.Address)
}

// TruncateAddress returns an address as 0x1234...5678.
func TruncateAddress(addr string) string {
	if len(addr) <= 10 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}

func (c *Config) setURLs() {
	if c.IsTestnet {
		c.APIBaseURL = "https://api.hyperliquid-testnet.xyz"
		c.WSBaseURL = "wss://api.hyperliquid-testnet.xyz/ws"
	} else {
		c.APIBaseURL = "https://api.hyperliquid.xyz"
		c.WSBaseURL = "wss://api.hyperliquid.xyz/ws"
	}
}

func (c *Config) InfoURL() string {
	return c.APIBaseURL + "/info"
}

// ValidateAddress checks if an address is a valid 0x-prefixed 40-hex-char string.
func ValidateAddress(addr string) bool {
	re := regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	return re.MatchString(strings.TrimSpace(addr))
}

// SaveWallets writes wallets to ~/.config/hltui/wallets.json.
func SaveWallets(wallets []Wallet) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(home, ".config", "hltui")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	wf := walletsFile{Wallets: wallets}
	data, err := json.MarshalIndent(wf, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "wallets.json"), data, 0o644)
}

// LoadWallets reads the wallets config file from ~/.config/hltui/wallets.json.
func LoadWallets() ([]Wallet, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".config", "hltui", "wallets.json")

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var wf walletsFile
	if err := json.Unmarshal(data, &wf); err != nil {
		return nil, fmt.Errorf("parse wallets.json: %w", err)
	}

	addrRe := regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)
	for i, w := range wf.Wallets {
		w.Address = strings.TrimSpace(w.Address)
		if !addrRe.MatchString(w.Address) {
			return nil, fmt.Errorf("wallet %d (%s): invalid address %q", i, w.Name, w.Address)
		}
		if w.Name == "" {
			w.Name = TruncateAddress(w.Address)
		}
		wf.Wallets[i] = w
	}

	if len(wf.Wallets) == 0 {
		return nil, fmt.Errorf("wallets.json contains no wallets")
	}

	return wf.Wallets, nil
}
