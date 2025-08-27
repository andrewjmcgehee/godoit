package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Color palette inspired by modern terminal themes
var (
	// Base colors
	magenta   = lipgloss.Color("5")
	green     = lipgloss.Color("2")
	yellow    = lipgloss.Color("3")
	red       = lipgloss.Color("1")
	gray      = lipgloss.Color("15")
	lightGray = lipgloss.Color("15")
	blue      = lipgloss.Color("9")
	white     = lipgloss.Color("0")
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(blue).
			Padding(0, 2).
			Bold(true).
			Border(lipgloss.NormalBorder()).
			Align(lipgloss.Center)

	// Header styles with subtle borders
	headerStyle = lipgloss.NewStyle().
			Foreground(gray).
			Bold(true).
			MarginTop(1).
			MarginBottom(1).
			PaddingLeft(1).
			Border(lipgloss.Border{
			Left: "▎",
		}).
		BorderForeground(magenta)

	// Selected item with rounded corners and subtle shadow effect
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(gray).
				Background(magenta).
				Padding(0, 2).
				MarginRight(1).
				Bold(true)

	itemStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(gray).
			MarginRight(1)

	// Completed items with strike-through and muted colors
	completedItemStyle = lipgloss.NewStyle().
				Padding(0, 2).
				Foreground(lightGray).
				Strikethrough(true).
				MarginRight(1)

	selectedCompletedItemStyle = lipgloss.NewStyle().
					Foreground(gray).
					Background(lightGray).
					Padding(0, 2).
					MarginRight(1).
					Strikethrough(true).
					Bold(true)

	// Priority styles with better visual hierarchy
	priorityP0Style = lipgloss.NewStyle().
			Foreground(red).
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("#2D1B1B")). // Dark red background
			Border(lipgloss.NormalBorder()).
			BorderForeground(red)

	priorityP1Style = lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("#2D2B1B")). // Dark yellow background
			Border(lipgloss.NormalBorder()).
			BorderForeground(yellow)

	priorityP2Style = lipgloss.NewStyle().
			Foreground(green).
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("#1B2D26")). // Dark green background
			Border(lipgloss.NormalBorder()).
			BorderForeground(green)

	// Help text with better spacing and subtle styling
	helpStyle = lipgloss.NewStyle().
			Foreground(lightGray).
			MarginTop(2).
			PaddingTop(1).
			Border(lipgloss.Border{
			Top: "─",
		}).
		BorderForeground(white).
		Italic(true)

	// Input field with modern styling
	inputStyle = lipgloss.NewStyle().
			Foreground(gray).
			Background(blue).
			Padding(0, 2).
			MarginTop(1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(magenta)

	// Message/error styles
	messageStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true).
			MarginTop(1).
			Padding(0, 1).
			Background(lipgloss.Color("#2D1B1B"))

	// Empty state message
	emptyStyle = lipgloss.NewStyle().
			Foreground(lightGray).
			Italic(true).
			Padding(1, 2).
			Align(lipgloss.Center)

	// Cursor/arrow indicator
	cursorStyle = lipgloss.NewStyle().
			Foreground(magenta).
			Bold(true)
)

func (s State) View() string {
	if s.windowHeight == 0 {
		return "Loading..."
	}
	var b strings.Builder
	title := titleStyle.Render("godo - go do what you're procrastinating doing")
	b.WriteString(title + "\n")
	switch s.uiState {
	case CreatingState:
		b.WriteString(s.renderCreateView())
	case EditingState:
		b.WriteString(s.renderEditView())
	default:
		b.WriteString(s.renderBrowseView())
	}
	if s.message != "" {
		b.WriteString("\n" + messageStyle.Render("⚠ "+s.message))
	}
	return b.String()
}

func (s State) renderBrowseView() string {
	var b strings.Builder
	viewIcon := "📋"
	viewTitle := "Active Todos"
	if s.viewMode == CompletedView {
		viewIcon = "✅"
		viewTitle = "Completed Todos"
	}
	header := headerStyle.Render(fmt.Sprintf("%s %s (%d)", viewIcon, viewTitle, len(s.todos)))
	b.WriteString(header + "\n")
	if len(s.todos) == 0 {
		emptyMsg := "😌 Nothing here!"
		if s.viewMode == ActiveView {
			emptyMsg = "😌 No active todos! Press 'n' to create your first one."
		}
		b.WriteString(emptyStyle.Render(emptyMsg) + "\n")
	} else {
		for i, todo := range s.todos {
			cursor := "  "
			if i == s.cursor {
				cursor = cursorStyle.Render("> ")
			} else {
				cursor = "  "
			}
			priorityText := s.renderPriority(Priority(todo.Priority))
			content := fmt.Sprintf("%s%s: %s", cursor, priorityText, todo.Content)

			if i == s.cursor {
				if todo.Completed {
					content = selectedCompletedItemStyle.Render(content)
				} else {
					content = selectedItemStyle.Render(content)
				}
			} else {
				if todo.Completed {
					content = completedItemStyle.Render(content)
				} else {
					content = itemStyle.Render(content)
				}
			}

			b.WriteString(content + "\n")
		}
	}

	b.WriteString(s.renderHelp())

	return b.String()
}

func (s State) renderCreateView() string {
	var b strings.Builder

	header := headerStyle.Render("📝 Create New Todo")
	b.WriteString(header + "\n")

	prompt := itemStyle.Render("Enter your todo:")
	b.WriteString(prompt + "\n")

	input := inputStyle.Render(s.editingText + "│")
	b.WriteString(input + "\n")

	help := helpStyle.Render("↵ Enter: save • Esc: cancel")
	b.WriteString(help)

	return b.String()
}

func (s State) renderEditView() string {
	var b strings.Builder

	header := headerStyle.Render("✏️  Edit Todo")
	b.WriteString(header + "\n")

	prompt := itemStyle.Render("Update content:")
	b.WriteString(prompt + "\n")

	input := inputStyle.Render(s.editingText + "│")
	b.WriteString(input + "\n")

	help := helpStyle.Render("↵ Enter: save • Esc: cancel")
	b.WriteString(help)

	return b.String()
}

func (s State) renderPriority(priority Priority) string {
	switch priority {
	case P0:
		return priorityP0Style.Render(" P0 ")
	case P1:
		return priorityP1Style.Render(" P1 ")
	case P2:
		return priorityP2Style.Render(" P2 ")
	default:
		return priorityP2Style.Render(" P2 ")
	}
}

func (s State) renderHelp() string {
	var helps []string

	switch s.viewMode {
	case ActiveView:
		helps = []string{
			"↑/k up", "↓/j down", "space toggle", "c edit", "n new",
			"d delete", "p priority", "l completed", "q quit",
		}
	case CompletedView:
		helps = []string{
			"↑/k up", "↓/j down", "space toggle", "d delete",
			"a active", "q quit",
		}
	}

	helpText := strings.Join(helps, " • ")
	return helpStyle.Render("💡 " + helpText)
}
