# ttree

Interactive terminal-based tree structure editor written in Go using the Bubble Tea framework.

## Quick Reference

```bash
# Build
go build -o ttree .

# Run
./ttree

# Test
go test ./...

# Test with coverage
go test -cover ./...
```

## Project Structure

```
├── main.go              # Entry point, initializes Bubble Tea program
├── internal/
│   ├── tree/            # Tree data structure (Node, Tree)
│   ├── ui/              # UI model, view, update, keybindings
│   └── render/          # ASCII tree rendering
```

## Architecture

- **Bubble Tea (TUI)**: Uses the Elm architecture - Model, View, Update pattern
- **Model** (`internal/ui/model.go`): Application state including tree, cursor, mode
- **Update** (`internal/ui/update.go`): Handles keyboard input and state transitions
- **View** (`internal/ui/view.go`): Renders the current state to terminal
- **Tree** (`internal/tree/tree.go`): Core tree data structure with Node operations
- **Renderer** (`internal/render/ascii.go`): Converts tree to ASCII art output

## Key Bindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Tab` | Indent node |
| `Shift+Tab` | Unindent node |
| `Enter` | Create new sibling |
| `Ctrl+D` / `Ctrl+Backspace` | Delete node |
| `Ctrl+C` | Copy tree to clipboard |
| `Ctrl+Q` / `Esc` | Quit |

## Code Conventions

- Follow standard Go formatting (`go fmt`)
- Keep packages focused: tree logic in `tree/`, UI in `ui/`, rendering in `render/`
- Use Bubble Tea patterns for state management
- Tests are co-located with source files (`*_test.go`)
- Write all documentation and code comments in English

## Commit Conventions

Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification.

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Changes that don't affect code meaning (whitespace, formatting)
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `test`: Adding or updating tests
- `chore`: Build process or tooling changes

### Examples

```
feat(tree): add node collapse/expand toggle
fix(ui): cursor position after node deletion
docs: update keybindings table
refactor(render): simplify prefix calculation
test(tree): add indent/unindent edge cases
```
