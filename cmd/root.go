package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/born1337/hltui/internal/app"
	"github.com/born1337/hltui/internal/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var (
	testnet bool
	vault   bool
)

var rootCmd = &cobra.Command{
	Use:   "hltui <address>",
	Short: "Hyperliquid TUI Dashboard",
	Long:  "Real-time terminal dashboard for monitoring Hyperliquid positions, orders, fills, funding, portfolio, and vaults.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		address := args[0]

		// Validate address
		matched, _ := regexp.MatchString(`^0x[a-fA-F0-9]{40}$`, address)
		if !matched {
			return fmt.Errorf("invalid address: must be 0x followed by 40 hex characters")
		}

		cfg := config.New(address, testnet, vault)
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
