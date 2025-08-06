# AI Mux

A terminal-based multiplexer for managing AI-assisted development workflows. AI Mux provides an interface for organizing work items, managing git worktrees, and integrating with Claude Code for AI-powered development assistance.

## Features

- **Work Item Management**: Create and implement work tasks in parallel
- **Git Worktree Integration**: Automatically manages git worktrees for isolated development environments
- **Claude Code Integration**: Integration with Claude Code for status updates
- **Tmux Integration**: Works with tmux for enhanced terminal multiplexing

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

## State Management

AI Mux creates a `.ai-mux` directory in the folder where you run it (a git repo workspace) for state management:

```
~/.ai-mux/
├── claude-settings.json    # Claude Code settings and hooks
╰─┬ [UUID]/                 # State files for open work items
  ├── item.json             # details about the item
  └── state-log.txt         # state log updated by claude, used for showing the status of the item
```

### Claude Code Integration

AI Mux automatically configures Claude Code with custom hooks for integration. The `claude-settings.json` file is created on first run with predefined hooks that notify AI Mux of Claude Code events.


## Development

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.


## License

[Add your license information here]
