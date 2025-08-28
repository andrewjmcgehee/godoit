package main

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
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
	// Tab styling
	activeTabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      " ",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┘",
		BottomRight: "└",
	}

	tabBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "┴",
		BottomRight: "┴",
	}

	tab = lipgloss.NewStyle().
		Border(tabBorder, true).
		BorderForeground(magenta).
		Padding(0, 1)

	activeTab = tab.Border(activeTabBorder, true)

	tabGap = tab.
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false)

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

func colorGrid(xSteps, ySteps int) [][]string {
	x0y0, _ := colorful.Hex("#F25D94")
	x1y0, _ := colorful.Hex("#EDFF82")
	x0y1, _ := colorful.Hex("#643AFF")
	x1y1, _ := colorful.Hex("#14F9D5")

	grid := make([][]string, ySteps)
	for y := range ySteps {
		grid[y] = make([]string, xSteps)
		for x := range xSteps {
			// Calculate normalized coordinates
			xRatio := float64(x) / float64(xSteps-1)
			yRatio := float64(y) / float64(ySteps-1)

			// Blend horizontally for top and bottom rows
			topColor := x0y0.BlendLuv(x1y0, xRatio)
			bottomColor := x0y1.BlendLuv(x1y1, xRatio)

			// Blend vertically between top and bottom colors
			finalColor := topColor.BlendLuv(bottomColor, yRatio)

			grid[y][x] = finalColor.Hex()
		}
	}
	return grid
}

func (s State) renderTitle() string {
	asciiArt := []string{
		" ██████╗  ██████╗     ██████╗  ██████╗     ██╗████████╗",
		"██╔════╝ ██╔═══██╗    ██╔══██╗██╔═══██╗    ██║╚══██╔══╝",
		"██║  ███╗██║   ██║    ██║  ██║██║   ██║    ██║   ██║   ",
		"██║   ██║██║   ██║    ██║  ██║██║   ██║    ██║   ██║   ",
		"╚██████╔╝╚██████╔╝    ██████╔╝╚██████╔╝    ██║   ██║   ",
		" ╚═════╝  ╚═════╝     ╚═════╝  ╚═════╝     ╚═╝   ╚═╝   ",
	}
	rows := len(asciiArt)
	cols := utf8.RuneCountInString(asciiArt[0])
	widthMinusTitle := s.windowWidth - cols
	leftPad := widthMinusTitle / 2
	rightPad := leftPad
	if widthMinusTitle%2 != 0 {
		rightPad += 1
	}
	colorized := [][]string{}
	colorized = append(colorized, strings.Split(strings.Repeat(" ", s.windowWidth), ""))
	for row := range rows {
		line := strings.Repeat(" ", leftPad)
		line += asciiArt[row]
		line += strings.Repeat(" ", rightPad)
		colorized = append(colorized, strings.Split(line, ""))
	}
	colorized = append(colorized, strings.Split(strings.Repeat(" ", s.windowWidth), ""))
	subtitle := "seriously, just do the thing already..."
	colors := colorGrid(s.windowWidth, len(colorized))
	for r := range colorized {
		for c, char := range colorized[r] {
			styledChar := lipgloss.NewStyle()
			bgColor := lipgloss.Color(colors[r][c])
			if char == " " {
				styledChar = styledChar.Foreground(bgColor)
			} else {
				styledChar = styledChar.Foreground(lipgloss.Color("0"))
			}
			styledChar = styledChar.
				Background(bgColor).
				Bold(true)
			colorized[r][c] = styledChar.Render(char)
		}
	}
	lines := []string{}
	for _, row := range colorized {
		lines = append(lines, strings.Join(row, ""))
	}
	asciiTitle := strings.Join(lines, "\n")
	styledSubtitle := lipgloss.NewStyle().
		Foreground(magenta).
		Bold(true).
		Italic(true).
		Align(lipgloss.Center).
		Render(subtitle)
	titleContent := lipgloss.JoinVertical(lipgloss.Center, asciiTitle, "", styledSubtitle)
	titleBox := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Render(titleContent)
	if s.windowWidth > 0 {
		titleBox = lipgloss.Place(s.windowWidth, lipgloss.Height(titleBox),
			lipgloss.Center, lipgloss.Top, titleBox)
	}
	return titleBox
}

func (s State) View() string {
	if s.windowHeight == 0 {
		return "Loading..."
	}
	var b strings.Builder
	title := s.renderTitle()
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

func (s State) renderTabs() string {
	activeTabText := fmt.Sprintf("📋 Active (%s)", func() string {
		ctx := context.Background()
		cnt, err := s.database.Queries.CountActiveTodos(ctx)
		if err != nil {
			return "?"
		}
		return fmt.Sprintf("%d", cnt)
	}())

	completedTabText := fmt.Sprintf("✅ Completed (%s)", func() string {
		ctx := context.Background()
		cnt, err := s.database.Queries.CountCompletedTodos(ctx)
		if err != nil {
			return "?"
		}
		return fmt.Sprintf("%d", cnt)
	}())

	var activeTabRendered, completedTabRendered string
	if s.viewMode == ActiveView {
		activeTabRendered = activeTab.Render(activeTabText)
		completedTabRendered = tab.Render(completedTabText)
	} else {
		activeTabRendered = tab.Render(activeTabText)
		completedTabRendered = activeTab.Render(completedTabText)
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		activeTabRendered,
		completedTabRendered,
	)

	gap := tabGap.Render(strings.Repeat(" ", max(0, s.windowWidth-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)

	return row
}

func (s State) renderBrowseView() string {
	var b strings.Builder

	// Render tabs instead of header
	tabs := s.renderTabs()
	b.WriteString(tabs + "\n\n")
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
			"↑/k: up", "↓/j: down", "space: mark complete", "n: new", "e: edit",
			"d: delete", "p: cycle priority", "tab: next tab", "q: quit",
		}
	case CompletedView:
		helps = []string{
			"↑/k: up", "↓/j: down", "space: mark incomplete", "d: delete",
			"tab: next tab", "q: quit",
		}
	}
	helpText := strings.Join(helps, " • ")
	return helpStyle.Render("💡 " + helpText)
}
