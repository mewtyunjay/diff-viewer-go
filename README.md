# diff-viewer-go

> WIP

A terminal UI for viewing diffs in a side-by-side format, built with Bubble Tea and Lipgloss

## Installation

```bash
go build -o diff-viewer-go
```

## Usage

```bash
./diff-viewer-go
```

## Keybindings

| Key | Action |
|-----|--------|
| `Tab` | Cycle focus between panels |
| `Shift+Tab` | Cycle focus backwards |
| `j` / `Down` | Move down / scroll down |
| `k` / `Up` | Move up / scroll up |
| `Ctrl+d` | Half page down |
| `Ctrl+u` | Half page up |
| `Ctrl+f` / `PgDn` | Page down |
| `Ctrl+b` / `PgUp` | Page up |
| `s` | Toggle synchronized scrolling |
| `q` / `Esc` | Quit |

## Requirements

- Go 1.21+
