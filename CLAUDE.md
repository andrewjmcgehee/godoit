# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is "godo", a terminal-based todo application written in Go using the Bubble Tea TUI framework. The application provides an interactive interface for managing todos with priorities, completion states, and persistence via SQLite with modern tooling.

## Architecture

The application uses modern Go tooling and clean architecture:

- **main.go**: Application entry point and Bubble Tea program initialization
- **ui.go**: Bubble Tea model implementation with all UI state management and rendering
- **database.go**: Database connection setup, goose migrations, and domain model definitions
- **styles.go**: Enhanced Lipgloss styling with modern color palette and visual hierarchy
- **config.go**: Configuration utilities for data directory and database path
- **internal/db/**: Auto-generated sqlc code for type-safe SQL operations
- **migrations/**: Goose migration files for database schema versioning
- **sql/**: SQLite schema and query definitions for sqlc

Key architectural patterns:
- **Goose** for database migrations with embedded migration files
- **Sqlc** for type-safe SQL queries and automatic Go code generation
- Bubble Tea's Elm-like architecture (Model-View-Update)
- SQLite database stored in `~/.local/share/godo/todos.db`

## Commands

**Build and run:**
```bash
go build -o godo
./godo
```

**Run directly:**
```bash
go run .
```

**Dependencies:**
```bash
go mod tidy
```

**Database operations:**
```bash
# Generate sqlc code after modifying queries
sqlc generate

# Create new migration
goose -dir migrations create <migration_name> sql

# Check migration status
goose -dir migrations sqlite3 ~/.local/share/godo/todos.db status
```

## Data Model

The core Todo struct includes:
- ID (int): Primary key
- Content (string): Todo text  
- Priority (Priority): P0 (high/red), P1 (medium/yellow), P2 (low/green)
- Completed (bool): Completion status
- CreatedAt/UpdatedAt (time.Time): Timestamps

Priority is a custom string type with constants P0, P1, P2.

## Database Layer

- **Goose migrations**: Version-controlled schema changes in `migrations/`
- **Sqlc queries**: Type-safe SQL operations defined in `sql/queries/todos.sql`
- **Generated code**: Auto-generated repository code in `internal/db/`
- **Domain conversion**: Helper functions convert between sqlc types and domain types

## UI States

The application has three main UI states:
- **BrowsingState**: Default view for navigating todos
- **EditingState**: Editing existing todo content  
- **CreatingState**: Creating new todos

View modes:
- **ActiveView**: Shows incomplete todos
- **CompletedView**: Shows completed todos

## Styling

Enhanced visual design with:
- Modern color palette (soft purple, mint green, coral red, etc.)
- Priority badges with colored backgrounds and borders
- Improved typography with emojis and better spacing
- Visual hierarchy with borders, padding, and styling