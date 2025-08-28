package main

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

var asciiArt = []string{
	" ██████╗  ██████╗     ██████╗  ██████╗     ██╗████████╗",
	"██╔════╝ ██╔═══██╗    ██╔══██╗██╔═══██╗    ██║╚══██╔══╝",
	"██║  ███╗██║   ██║    ██║  ██║██║   ██║    ██║   ██║   ",
	"██║   ██║██║   ██║    ██║  ██║██║   ██║    ██║   ██║   ",
	"╚██████╔╝╚██████╔╝    ██████╔╝╚██████╔╝    ██║   ██║   ",
	" ╚═════╝  ╚═════╝     ╚═════╝  ╚═════╝     ╚═╝   ╚═╝   ",
}

var (
	magenta   = lipgloss.Color("5")
	green     = lipgloss.Color("2")
	yellow    = lipgloss.Color("3")
	red       = lipgloss.Color("1")
	gray      = lipgloss.Color("15")
	lightGray = lipgloss.Color("15")
	blue      = lipgloss.Color("9")
)

var (
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
	tabGap    = tab.
			BorderTop(false).
			BorderLeft(false).
			BorderRight(false)
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(gray).
				MarginRight(1)
	itemStyle = lipgloss.NewStyle().
			Foreground(gray).
			MarginRight(1)
	completedItemStyle = lipgloss.NewStyle().
				Foreground(gray).
				Strikethrough(true).
				MarginRight(1)
	selectedCompletedItemStyle = lipgloss.NewStyle().
					Foreground(gray).
					MarginRight(1).
					Strikethrough(true)
	priorityP0Style = lipgloss.NewStyle().
			Foreground(red)
	priorityP1Style = lipgloss.NewStyle().
			Foreground(yellow)
	priorityP2Style = lipgloss.NewStyle().
			Foreground(green)
	messageStyle = lipgloss.NewStyle().
			Foreground(red).
			MarginTop(1).
			Padding(0, 1).
			Background(lipgloss.Color("#2D1B1B"))
	emptyStyle = lipgloss.NewStyle().
			Foreground(lightGray).
			Italic(true).
			Padding(1, 2).
			Align(lipgloss.Center)
	cursorStyle = lipgloss.NewStyle().
			Foreground(magenta).
			Padding(0, 1)
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
			xRatio := float64(x) / float64(xSteps)
			yRatio := float64(y) / float64(ySteps)
			topColor := x0y0.BlendLuv(x1y0, xRatio)
			bottomColor := x0y1.BlendLuv(x1y1, xRatio)
			finalColor := topColor.BlendLuv(bottomColor, yRatio)
			grid[y][x] = finalColor.Hex()
		}
	}
	return grid
}

func (s State) renderTitle() string {
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
				Background(bgColor)
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
			cursor := cursorStyle.Render("   ")
			if i == s.cursor {
				cursor = cursorStyle.Render(">>>")
			}
			var content string
			if i == s.cursor {
				if todo.Completed {
					content = fmt.Sprintf("%s: %s", todo.Priority, todo.Content)
					content = selectedCompletedItemStyle.Render(content)
				} else {
					priorityText := s.renderPriority(Priority(todo.Priority))
					content = fmt.Sprintf("%s: %s", priorityText, todo.Content)
					content = selectedItemStyle.Render(content)
				}
			} else {
				if todo.Completed {
					content = fmt.Sprintf("%s: %s", todo.Priority, todo.Content)
					content = completedItemStyle.Render(content)
				} else {
					priorityText := s.renderPriority(Priority(todo.Priority))
					content = fmt.Sprintf("%s: %s", priorityText, todo.Content)
					content = itemStyle.Render(content)
				}
			}
			b.WriteString(fmt.Sprintf("%s %s\n", cursor, content))
		}
	}
	mainContent := b.String()
	help := s.renderHelp()
	lines := strings.Split(mainContent, "\n")
	helpLines := strings.Split(help, "\n")
	maxLines := max(len(lines), len(helpLines))
	result := make([]string, maxLines)
	for i := range maxLines {
		if i < len(lines) && i < len(helpLines) {
			mainLine := lines[i]
			helpLine := helpLines[i]
			if strings.TrimSpace(helpLine) != "" {
				result[i] = helpLine
			} else {
				result[i] = mainLine
			}
		} else if i < len(lines) {
			result[i] = lines[i]
		} else if i < len(helpLines) {
			result[i] = helpLines[i]
		}
	}
	return strings.Join(result, "\n")
}

func (s State) renderCreateView() string {
	formBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(magenta).
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center)
	titleStyle := lipgloss.NewStyle().
		Foreground(magenta).
		MarginBottom(1).
		Align(lipgloss.Center)
	inputFieldStyle := lipgloss.NewStyle().
		Foreground(gray).
		Padding(0, 2).
		Width(50).
		Border(lipgloss.NormalBorder()).
		BorderForeground(yellow)
	keymapBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(magenta).
		Padding(0, 2).
		MarginRight(1)
	keyStyle := lipgloss.NewStyle().
		Foreground(yellow).
		Width(10)
	descStyle := lipgloss.NewStyle().
		Foreground(lightGray)

	var content []string
	content = append(content, titleStyle.Render("create new todo"))
	content = append(content, inputFieldStyle.Render(s.editingText+"█"))

	var keymaps []string
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("enter"), descStyle.Render("save todo")))
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("esc"), descStyle.Render("cancel")))
	keymapContent := strings.Join(keymaps, "\n")

	content = append(content, keymapBoxStyle.Render(keymapContent))
	formContent := strings.Join(content, "\n")
	form := formBoxStyle.Render(formContent)
	if s.windowWidth > 0 && s.windowHeight > 0 {
		availableHeight := s.windowHeight - len(asciiArt) - 4
		form = lipgloss.Place(
			s.windowWidth,
			availableHeight,
			lipgloss.Center,
			lipgloss.Center,
			form,
		)
	}
	return form
}

func (s State) renderEditView() string {
	formBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(magenta).
		Padding(1, 2).
		Width(60).
		Align(lipgloss.Center)
	titleStyle := lipgloss.NewStyle().
		Foreground(magenta).
		MarginBottom(1).
		Align(lipgloss.Center)
	inputFieldStyle := lipgloss.NewStyle().
		Foreground(gray).
		Padding(0, 2).
		Width(50).
		Border(lipgloss.NormalBorder()).
		BorderForeground(yellow)
	keymapBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(magenta).
		Padding(0, 2).
		MarginRight(1)
	keyStyle := lipgloss.NewStyle().
		Foreground(yellow).
		Width(10)
	descStyle := lipgloss.NewStyle().
		Foreground(lightGray)

	var content []string
	content = append(content, titleStyle.Render("edit todo"))
	content = append(content, inputFieldStyle.Render(s.editingText+"█"))

	var keymaps []string
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("enter"), descStyle.Render("save changes")))
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("esc"), descStyle.Render("cancel")))
	keymapContent := strings.Join(keymaps, "\n")

	content = append(content, keymapBoxStyle.Render(keymapContent))
	formContent := strings.Join(content, "\n")
	form := formBoxStyle.Render(formContent)
	if s.windowWidth > 0 && s.windowHeight > 0 {
		availableHeight := s.windowHeight - len(asciiArt) - 4
		form = lipgloss.Place(
			s.windowWidth,
			availableHeight,
			lipgloss.Center,
			lipgloss.Center,
			form,
		)
	}
	return form
}

func (s State) renderPriority(priority Priority) string {
	switch priority {
	case P0:
		return priorityP0Style.Render(string(P0))
	case P1:
		return priorityP1Style.Render(string(P1))
	case P2:
		return priorityP2Style.Render(string(P2))
	default:
		return priorityP2Style.Render(string(P2))
	}
}

func (s State) renderHelp() string {
	var keymaps string
	if s.showHelp {
		keymaps = s.renderKeymaps()
	} else {
		keymapBoxStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(magenta).
			Padding(0, 2).
			MarginRight(1)
		titleStyle := lipgloss.NewStyle().
			Foreground(magenta).
			Align(lipgloss.Center)
		keymaps = keymapBoxStyle.Render(titleStyle.Render("? keymaps"))
	}
	availableHeight := s.windowHeight - len(asciiArt) - 4
	return lipgloss.Place(
		s.windowWidth,
		availableHeight,
		lipgloss.Right,
		lipgloss.Bottom,
		keymaps,
	)
}

func (s State) renderKeymaps() string {
	keymapBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(magenta).
		Padding(0, 2).
		MarginRight(1)
	titleStyle := lipgloss.NewStyle().
		Foreground(magenta).
		MarginBottom(1).
		Align(lipgloss.Center)
	keyStyle := lipgloss.NewStyle().
		Foreground(yellow).
		Width(15)
	descStyle := lipgloss.NewStyle().
		Foreground(gray)
	var keymaps []string
	keymaps = append(keymaps, titleStyle.Render("? keymaps"))
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("↑/k"), descStyle.Render("move up")))
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("↓/j"), descStyle.Render("move down")))
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("tab"), descStyle.Render("cycle tabs")))
	if s.viewMode == ActiveView {
		keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("n"), descStyle.Render("new todo")))
		keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("e"), descStyle.Render("edit todo")))
		keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("space"), descStyle.Render("mark done")))
		keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("p"), descStyle.Render("cycle priority")))
	} else {
		keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("space"), descStyle.Render("mark not done")))
	}
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("d"), descStyle.Render("delete todo")))
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("?"), descStyle.Render("toggle help")))
	keymaps = append(keymaps, lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render("q / ctrl+c"), descStyle.Render("quit")))
	content := strings.Join(keymaps, "\n")
	return keymapBoxStyle.Render(content)
}
