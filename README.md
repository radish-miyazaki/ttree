# ttree

A terminal-based interactive tree structure editor. Create, edit, and organize hierarchical data with vim-style navigation, then copy the ASCII tree output to your clipboard.

## Demo

<img src="vhs/demo.gif" alt="ttree demo" width="600">

## Features

- Interactive tree editing in the terminal
- Real-time ASCII tree preview
- Vim-style navigation (`j`/`k`)
- Indent/unindent nodes with Tab
- Copy rendered tree to clipboard

## Installation

```bash
go install github.com/radish-miyazaki/ttree@latest
```

Or build from source:

```bash
git clone https://github.com/radish-miyazaki/ttree.git
cd ttree
go build -o ttree .
```

## Usage

```bash
./ttree
```

### Key Bindings

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Tab` | Indent node (make child of previous sibling) |
| `Shift+Tab` | Unindent node (make sibling of parent) |
| `Enter` | Create new sibling node |
| `Ctrl+D` | Delete current node |
| `Ctrl+C` | Copy tree to clipboard |
| `Ctrl+Q` / `Esc` | Quit |

### Example Output

```
├── src
│   ├── components
│   │   ├── Header.tsx
│   │   └── Footer.tsx
│   └── utils
│       └── helpers.ts
├── tests
└── README.md
```

## License

MIT
