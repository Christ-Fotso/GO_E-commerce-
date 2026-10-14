package admin

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			MarginBottom(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("39")).
				Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			MarginTop(1)

	badgeAdminStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208")).
			Bold(true)

	badgeCustomerStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("245"))

	badgePaidStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("42"))

	badgePendingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("214"))

	badgeShippingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("39"))

	badgeDeliveredStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("49"))

	badgeCancelledStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("196"))
)
