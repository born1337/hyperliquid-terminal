package style

import "github.com/charmbracelet/lipgloss"

var (
	Green   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	Red     = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	Yellow  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	Cyan    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	Magenta = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	White   = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Bold(true)
	Dim     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	Bold    = lipgloss.NewStyle().Bold(true)

	HeaderBar = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Padding(0, 1)

	ActiveTab = lipgloss.NewStyle().
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Padding(0, 1)

	InactiveTab = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Padding(0, 1)

	StatusBar = lipgloss.NewStyle().
			Background(lipgloss.Color("236")).
			Foreground(lipgloss.Color("8")).
			Padding(0, 1)

	StatusConnected = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)

	StatusDisconnected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("9")).
				Bold(true)

	SummaryLabel = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Width(16)

	TableHeader = Dim.Underline(true)

	Border = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("62"))
)

func PnlColor(val float64) lipgloss.Style {
	if val < 0 {
		return Red
	}
	return Green
}

func SideColor(szi float64) lipgloss.Style {
	if szi < 0 {
		return Red
	}
	return Green
}
