package config

import "testing"

func TestNewMainnet(t *testing.T) {
	c := New("0xabc", false, false)
	if c.Address != "0xabc" {
		t.Errorf("Address = %q, want 0xabc", c.Address)
	}
	if c.IsTestnet {
		t.Error("IsTestnet should be false")
	}
	if c.APIBaseURL != "https://api.hyperliquid.xyz" {
		t.Errorf("APIBaseURL = %q", c.APIBaseURL)
	}
	if c.WSBaseURL != "wss://api.hyperliquid.xyz/ws" {
		t.Errorf("WSBaseURL = %q", c.WSBaseURL)
	}
	if c.InfoURL() != "https://api.hyperliquid.xyz/info" {
		t.Errorf("InfoURL() = %q", c.InfoURL())
	}
}

func TestNewTestnet(t *testing.T) {
	c := New("0xdef", true, true)
	if !c.IsTestnet {
		t.Error("IsTestnet should be true")
	}
	if !c.IsVault {
		t.Error("IsVault should be true")
	}
	if c.APIBaseURL != "https://api.hyperliquid-testnet.xyz" {
		t.Errorf("APIBaseURL = %q", c.APIBaseURL)
	}
	if c.WSBaseURL != "wss://api.hyperliquid-testnet.xyz/ws" {
		t.Errorf("WSBaseURL = %q", c.WSBaseURL)
	}
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}
