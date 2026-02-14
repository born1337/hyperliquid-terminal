package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/born1337/hyperliquid-terminal/internal/app"
	"github.com/born1337/hyperliquid-terminal/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	testnet bool
	vault   bool
)

var rootCmd = &cobra.Command{
	Use:   "hltui [address]",
	Short: "Hyperliquid TUI Dashboard",
	Long:  "Real-time terminal dashboard for monitoring Hyperliquid positions, orders, fills, funding, portfolio, and vaults.\n\nRun with an address argument or configure wallets in ~/.config/hltui/wallets.json",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *config.Config

		// Try loading wallets from config
		wallets, walletsErr := config.LoadWallets()

		if len(args) == 1 {
			address := args[0]

			// Validate address
			matched, _ := regexp.MatchString(`^0x[a-fA-F0-9]{40}$`, address)
			if !matched {
				return fmt.Errorf("invalid address: must be 0x followed by 40 hex characters")
			}

			if walletsErr == nil && len(wallets) > 0 {
				// Prepend CLI address as first wallet, merge with config wallets
				cliWallet := config.Wallet{Name: "CLI", Address: address, Testnet: testnet, Vault: vault}
				allWallets := append([]config.Wallet{cliWallet}, wallets...)
				cfg = config.NewWithWallets(allWallets, 0, testnet, vault)
			} else {
				cfg = config.New(address, testnet, vault)
			}
		} else {
			// No address arg — require wallets.json
			if walletsErr != nil {
				return fmt.Errorf("no address provided and wallets not configured: %w\n\nUsage: hltui <address>\n   or: create ~/.config/hltui/wallets.json", walletsErr)
			}
			cfg = config.NewWithWallets(wallets, 0, testnet, vault)
		}

		m := app.NewModel(cfg)

		p := tea.NewProgram(m, tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return err
		}
		return nil
	},
	Version: config.Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&testnet, "testnet", "t", false, "Use testnet API")
	rootCmd.Flags().BoolVarP(&vault, "vault", "V", false, "Treat address as vault")
}
