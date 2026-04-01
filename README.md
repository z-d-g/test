# md-cli

> Terminal markdown editor with live rendering. Written in Go.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-blue?style=flat-square)

## Features

- **Live rendering** — headings, bold, italic, code, links, tables, lists, images
- **Syntax-aware cursor** — jumps between rendered and raw source in code blocks, tables, and lists
- **Full editing** — selection, copy/cut/paste, undo/redo, word and line operations
- **Persistent cursor** — restores position per file across sessions
- **Adaptive theming** — respects terminal light/dark background
- **Print mode** — render markdown to stdout without the editor

## Install

```bash
go install github.com/z-d-g/md-cli/cmd/md-cli@latest
```

Or build from source:

```bash
git clone https://github.com/z-d-g/md-cli.git
cd md-cli
go build -o md-cli ./cmd/md-cli
```

## Usage

```bash
md-cli file.md              # open in editor
md-cli -p file.md           # print to stdout
cat file.md | md-cli -p     # pipe from stdin
```

### Keybindings

| Key     | Action |
|---------|--------|
| Ctrl+S  | Save   |
| Ctrl+Q  | Quit   |
| F1      | Help   |


Full reference: press `F1` in the editor.

## License

[MIT](LICENSE)
