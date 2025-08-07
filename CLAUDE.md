# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Run
```bash
# Build the application
go build -o ai-mux

# Run the application
go run main.go

# Run with alternate screen and mouse support (already configured in main.go)
./ai-mux
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./component/modal
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies (remove unused, add missing)
go mod tidy

# Update dependencies
go get -u ./...
```

## Architecture

This is a TUI (Terminal User Interface) application built with the Charm libraries (Bubble Tea framework). The application follows the Elm architecture pattern with Model-Update-View.

### Core Structure

- **main.go**: Entry point containing the root `appModel` that orchestrates the entire application
  - Manages window sizing, keyboard input routing, and modal display
  - Uses Bubble Tea's alternate screen mode and mouse support

- **component/**: Reusable UI components following Bubble Tea patterns
  - **modal/**: Modal dialog system with overlay rendering
    - `ModalContent` interface for pluggable modal content
    - Advanced ANSI escape sequence parsing for proper overlay rendering
  - **pane/**: Main content panes
    - `ListPaneModel`: Displays work items with viewport scrolling
  - **workform/**: Form component for work item creation
    - Multi-step form using `huh` library
    - Currently incomplete (references undefined types)

- **theme/**: Centralized color scheme using Lipgloss
  - Catppuccin-inspired color palette

- **util/**: Shared utilities
  - Border styling functions

### Key Patterns

1. **Message Passing**: Components communicate via Tea messages
   - `ShowModalMsg`, `CloseMsg` for modal control
   - Components return commands that emit messages

2. **Component Interfaces**: 
   - `ModalContent` interface allows any component to be displayed in a modal
   - Components implement `Update()`, `View()`, and domain-specific methods

3. **Layout Management**:
   - Window resize events propagate to all components
   - Components calculate their dimensions based on terminal size

4. **Styling**:
   - Consistent use of Lipgloss for styling
   - Theme colors centralized in `theme/colors.go`

### Current State

The application appears to be a work-in-progress AI multiplexer/manager with:
- Basic modal system for adding work items
- List view for displaying work items (currently empty)
- Form infrastructure for creating work items (needs completion)

Key areas needing work:
- `workform.go` has compilation errors (missing imports, undefined types)
- No persistence layer implemented
- Work item data structures not yet defined

## Development Guidance

- Since this is a TUI, do not try to run the program to test things. I will test it myself.
- Don't create builds to test changes.