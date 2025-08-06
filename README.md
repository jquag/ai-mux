# AI Mux

A terminal-based multiplexer for managing AI-assisted development workflows. AI Mux provides a seamless interface for organizing work items, managing git worktrees, and integrating with Claude Code for AI-powered development assistance.

## Features

- **Work Item Management**: Create and track development tasks with descriptions, branch names, and planning modes
- **Git Worktree Integration**: Automatically manages git worktrees for isolated development environments
- **Claude Code Integration**: Deep integration with Claude Code including custom hooks and settings management
- **Beautiful TUI**: Built with the Charm libraries (Bubble Tea framework) featuring a modern, keyboard-driven interface
- **Modal System**: Intuitive modal dialogs for creating and managing work items
- **Tmux Integration**: Seamlessly works with tmux for enhanced terminal multiplexing

## Prerequisites

AI Mux requires the following tools to be installed:

- **git**: For version control and worktree management
- **claude**: Claude Code CLI for AI assistance
- **tmux**: Terminal multiplexer for session management
- **Go 1.23.2+**: For building from source

## Installation

```bash
# Clone the repository
git clone https://github.com/jquag/ai-mux.git
cd ai-mux

# Download dependencies
go mod download

# Build the application
go build -o ai-mux

# Or run directly
go run main.go
```

## Usage

### Basic Commands

```bash
# Run the application
./ai-mux

# Handle Claude events (used by hooks)
./ai-mux --event < event.json
```

### Keyboard Shortcuts

- **`n`**: Create a new work item
- **`↑/↓`**: Navigate through work items
- **`Enter`**: Select/activate a work item
- **`Esc`**: Close modals
- **`q`**: Quit the application (when not in a modal)
- **`Ctrl+C`**: Force quit

## Configuration

AI Mux creates a `.ai-mux` directory in your home folder for configuration and state management:

```
~/.ai-mux/
├── claude-settings.json    # Claude Code settings and hooks
├── worktrees/              # Git worktrees for work items
└── work-items.json         # Persisted work items (if implemented)
```

### Claude Code Integration

AI Mux automatically configures Claude Code with custom hooks for seamless integration. The `claude-settings.json` file is created on first run with predefined hooks that notify AI Mux of Claude Code events.

You can also set a custom directory using the `AI_MUX_DIR` environment variable:

```bash
export AI_MUX_DIR=/path/to/custom/dir
./ai-mux
```

## Architecture

AI Mux follows the Elm architecture pattern (Model-Update-View) using the Bubble Tea framework:

### Components

- **`app`**: Root application model that orchestrates the UI
- **`worklist`**: Displays and manages work items with scrollable viewport
- **`workform`**: Multi-step form for creating new work items
- **`modal`**: Reusable modal dialog system with overlay rendering
- **`footer`**: Status bar and keyboard hints
- **`theme`**: Catppuccin-inspired color scheme using Lipgloss

### Data Flow

1. User input triggers Tea messages
2. Components update their state and return commands
3. Commands emit new messages for state changes
4. Views are re-rendered with updated state

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./component/modal
```

### Code Structure

```
ai-mux/
├── main.go                 # Entry point and system checks
├── component/              # UI components
│   ├── app/               # Root application model
│   ├── modal/             # Modal dialog system
│   ├── worklist/          # Work item list view
│   ├── workform/          # Work item creation form
│   └── footer/            # Footer component
├── data/                   # Data structures
│   └── workitem.go        # Work item model
├── util/                   # Utilities
│   ├── claude.go          # Claude Code integration
│   ├── git.go             # Git operations
│   └── borders.go         # UI styling utilities
├── theme/                  # Visual theming
│   └── colors.go          # Color palette
└── claudeSettings.json     # Embedded Claude settings
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

### Development Notes

- The application uses embedded files (see `//go:embed` directive in main.go)
- Modal rendering uses advanced ANSI escape sequences for proper overlay display
- The project follows standard Go conventions and uses the Charm libraries idiomatically

## License

[Add your license information here]

## Acknowledgments

Built with the excellent [Charm](https://charm.sh/) libraries:
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Huh](https://github.com/charmbracelet/huh) - Terminal forms