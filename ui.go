package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/andrewjmcgehee/godo/internal/orm"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewMode int

const (
	ActiveView ViewMode = iota
	CompletedView
)

type UIState int

const (
	BrowsingState UIState = iota
	EditingState
	CreatingState
)

type State struct {
	database     *Database
	todos        []orm.Todo
	cursor       int
	viewMode     ViewMode
	uiState      UIState
	editingTodo  *orm.Todo
	editingText  string
	message      string
	windowWidth  int
	windowHeight int
}

type todoLoadedMsg struct {
	todos []orm.Todo
}

type todoCreatedMsg struct {
	todo *orm.Todo
}

type todoUpdatedMsg struct {
	success bool
}

type todoDeletedMsg struct {
	success bool
}

func InitialState(database *Database) State {
	return State{
		database: database,
		todos:    []orm.Todo{},
		cursor:   0,
		viewMode: ActiveView,
		uiState:  BrowsingState,
	}
}

func (s State) Init() tea.Cmd {
	return s.loadTodos()
}

func (s State) loadTodos() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		var todos []orm.Todo
		var err error
		if s.viewMode == ActiveView {
			todos, err = s.database.Queries.GetActiveTodos(ctx)
		} else {
			todos, err = s.database.Queries.GetCompletedTodos(ctx)
		}
		if err != nil {
			return tea.Msg(fmt.Sprintf("Error loading todos: %v", err))
		}
		return todoLoadedMsg{todos: todos}
	})
}

func (s State) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.windowWidth = msg.Width
		s.windowHeight = msg.Height
	case todoLoadedMsg:
		s.todos = msg.todos
		if s.cursor >= len(s.todos) && len(s.todos) > 0 {
			s.cursor = len(s.todos) - 1
		} else if len(s.todos) == 0 {
			s.cursor = 0
		}
		s.message = ""
	case todoCreatedMsg:
		s.uiState = BrowsingState
		s.editingText = ""
		return s, s.loadTodos()
	case todoUpdatedMsg:
		s.uiState = BrowsingState
		s.editingTodo = nil
		s.editingText = ""
		return s, s.loadTodos()
	case todoDeletedMsg:
		return s, s.loadTodos()
	case tea.KeyMsg:
		return s.handleKeyPress(msg)
	}
	return s, nil
}

func (s State) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch s.uiState {
	case BrowsingState:
		return s.handleBrowsingKeys(msg)
	case EditingState, CreatingState:
		return s.handleEditingKeys(msg)
	}
	return s, nil
}

func (s State) handleBrowsingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return s, tea.Quit
	case "up", "k":
		if s.cursor > 0 {
			s.cursor--
		}
	case "down", "j":
		if s.cursor < len(s.todos)-1 {
			s.cursor++
		}
	case "n":
		s.uiState = CreatingState
		s.editingText = ""
	case "c":
		if len(s.todos) > 0 && s.cursor < len(s.todos) {
			s.uiState = EditingState
			s.editingTodo = &s.todos[s.cursor]
			s.editingText = s.editingTodo.Content
		}
	case " ":
		if len(s.todos) > 0 && s.cursor < len(s.todos) {
			return s, s.toggleTodo(s.todos[s.cursor].ID)
		}
	case "d":
		if len(s.todos) > 0 && s.cursor < len(s.todos) {
			return s, s.deleteTodo(s.todos[s.cursor].ID)
		}
	case "p":
		if len(s.todos) > 0 && s.cursor < len(s.todos) {
			return s, s.cyclePriority(s.todos[s.cursor].ID, Priority(s.todos[s.cursor].Priority))
		}
	case "a":
		if s.viewMode != ActiveView {
			s.viewMode = ActiveView
			s.cursor = 0
			return s, s.loadTodos()
		}
	case "l":
		if s.viewMode != CompletedView {
			s.viewMode = CompletedView
			s.cursor = 0
			return s, s.loadTodos()
		}
	}
	return s, nil
}

func (s State) handleEditingKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		s.uiState = BrowsingState
		s.editingTodo = nil
		s.editingText = ""
	case "enter":
		if strings.TrimSpace(s.editingText) == "" {
			return s, nil
		}
		if s.uiState == CreatingState {
			return s, s.createTodo(s.editingText)
		} else if s.uiState == EditingState && s.editingTodo != nil {
			return s, s.updateTodo(s.editingTodo.ID, s.editingText)
		}
	case "backspace":
		if len(s.editingText) > 0 {
			s.editingText = s.editingText[:len(s.editingText)-1]
		}
	default:
		if len(msg.String()) == 1 {
			s.editingText += msg.String()
		}
	}
	return s, nil
}

func (s State) createTodo(content string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		now := time.Now()
		todo, err := s.database.Queries.CreateTodo(ctx, orm.CreateTodoParams{
			Content:   content,
			Priority:  string(P2),
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return tea.Msg(fmt.Sprintf("Error creating todo: %v", err))
		}
		return todoCreatedMsg{todo: &todo}
	})
}

func (s State) updateTodo(id int, content string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		now := time.Now()
		err := s.database.Queries.UpdateTodoContent(ctx, orm.UpdateTodoContentParams{
			ID:        id,
			Content:   content,
			UpdatedAt: now,
		})
		if err != nil {
			return tea.Msg(fmt.Sprintf("Error updating todo: %v", err))
		}
		return todoUpdatedMsg{success: true}
	})
}

func (s State) toggleTodo(id int) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		now := time.Now()
		err := s.database.Queries.ToggleTodoCompleted(ctx, orm.ToggleTodoCompletedParams{
			ID:        id,
			UpdatedAt: now,
		})
		if err != nil {
			return tea.Msg(fmt.Sprintf("Error toggling todo: %v", err))
		}
		return todoUpdatedMsg{success: true}
	})
}

func (s State) deleteTodo(id int) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		ctx := context.Background()
		err := s.database.Queries.DeleteTodo(ctx, id)
		if err != nil {
			return tea.Msg(fmt.Sprintf("Error deleting todo: %v", err))
		}
		return todoDeletedMsg{success: true}
	})
}

func (s State) cyclePriority(id int, prev Priority) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		var next Priority
		switch prev {
		case P2:
			next = P1
		case P1:
			next = P0
		case P0:
			next = P2
		default:
			next = P2
		}
		ctx := context.Background()
		now := time.Now()
		err := s.database.Queries.UpdateTodoPriority(ctx, orm.UpdateTodoPriorityParams{
			ID:        id,
			Priority:  string(next),
			UpdatedAt: now,
		})
		if err != nil {
			return tea.Msg(fmt.Sprintf("Error updating priority: %v", err))
		}
		return todoUpdatedMsg{success: true}
	})
}
