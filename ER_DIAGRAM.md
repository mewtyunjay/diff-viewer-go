# ER Diagram: diff-viewer-go Project

## Entity-Relationship Model
Treating **Files as Tables** and **Functions/Attributes/Types as Entities**

```
┌─────────────────────────────────────────────────────────────────────┐
│                         PROJECT STRUCTURE                           │
└─────────────────────────────────────────────────────────────────────┘

                              CORE LAYER
                         ┌─────────────────┐
                         │ diff/types.go   │
                         ├─────────────────┤
                         │ • LineType      │
                         │ • Line          │
                         │ • FileDiff      │
                         │ • Result        │
                         └────────┬────────┘
                                  │ (used by all)
                    ┌─────────────┼─────────────┐
                    │             │             │
        ┌───────────▼────┐ ┌─────▼──────┐ ┌──▼─────────────┐
        │ parser/errors  │ │  parser/   │ │     tui/       │
        │     .go        │ │  options.go│ │    model.go    │
        ├────────────────┤ ├────────────┤ ├────────────────┤
        │ • ErrNotGitRepo│ │ • Options  │ │ • FocusedPanel │
        │ • ErrEmptyDiff │ │ • Option   │ │ • Model        │
        │ • ErrInvalidDif│ │ • Default │ │ • Init()       │
        │ • GitError     │ │   Options()│ │ • Update()     │
        │ • ParseError   │ │ • WithGit  │ │ • View()       │
        │ • Error()      │ │   Path()   │ │ • scroll*()    │
        │ • Unwrap()     │ │ • WithWork │ │ • render*()    │
        │                │ │   Dir()    │ │ • update*()    │
        └────────────────┘ └────────────┘ └────────────────┘
                                  │               │
                                  │               │
        ┌─────────────────────────▼────┐         │
        │   parser/git.go              │         │
        ├──────────────────────────────┤         │
        │ • GitRunner                  │         │
        │   - gitPath                  │         │
        │   - workDir                  │         │
        │ • NewGitRunner()             │         │
        │ • IsGitRepository()          │         │
        │ • RunDiff()                  │         │
        │ • FindGitRoot()              │         │
        └────────────┬─────────────────┘         │
                     │                           │
        ┌────────────▼──────────────┐            │
        │   parser/parser.go        │            │
        ├───────────────────────────┤            │
        │ • Parser                  │            │
        │   - opts                  │            │
        │   - git                   │            │
        │ • New()                   │            │
        │ • ParseString()           │            │
        │ • ParseReader()           │────────────┼──┐
        │ • ParseGitDiff()          │            │  │
        │ • IsGitRepository()       │            │  │
        └────────────┬──────────────┘            │  │
                     │                           │  │
        ┌────────────▼──────────────┐            │  │
        │  parser/unified.go        │            │  │
        ├───────────────────────────┤            │  │
        │ • hunk (internal)         │            │  │
        │   - oldStart/oldCount     │            │  │
        │   - newStart/newCount     │            │  │
        │   - lines                 │            │  │
        │ • parseUnified()          │            │  │
        │ • parseFileDiff()         │            │  │
        │ • parseHunk()             │────────┐   │  │
        │ • diffGitRE (regex)       │        │   │  │
        │ • hunkHeaderRE (regex)    │        │   │  │
        │ • binaryFileRE (regex)    │        │   │  │
        └─────────────────────────────────────┤───┼──┘
                                              │   │
        ┌─────────────────────────────────────▼──┐│
        │    parser/aligner.go                  ││
        ├───────────────────────────────────────┤│
        │ • alignHunks()                        ││
        │ • alignHunk()                         ││
        │ • alignBlock()                        ││
        │ (Transforms to left/right line pairs)││
        └─────────────────────────────────────┬─┘│
                                              │  │
                                    ┌─────────┘  │
                                    │            │
                     ┌──────────────▼────┐  ┌───▼────────────┐
                     │  tui/keys.go      │  │  tui/styles.go │
                     ├───────────────────┤  ├────────────────┤
                     │ • KeyMap          │  │ • PanelStyle   │
                     │   - Up            │  │ • TitleStyle   │
                     │   - Down          │  │ • AddLineStyle │
                     │   - Left/Right    │  │ • DeleteLine   │
                     │   - Tab/ShiftTab  │  │   Style        │
                     │   - PageUp/Down   │  │ • ContextLine  │
                     │   - Enter/Quit    │  │   Style        │
                     │   - SyncToggle    │  │ • FileItemStyle│
                     │ • Default KeyMap  │  │ • LineNumStyle │
                     │ • ShortHelp()     │  │ • Placeholder  │
                     │ • FullHelp()      │  │   Style        │
                     └───────────────────┘  │ • HelpStyle    │
                                            └────────────────┘
                                                    │
                                                    │
                                            ┌───────▼────────┐
                                            │  main.go       │
                                            ├────────────────┤
                                            │ • main()       │
                                            │ • runTUI()     │
                                            │ • handleError()│
                                            │ (Orchestrates) │
                                            └────────────────┘
```

---

## Detailed Entity Relationships

### **Table 1: diff/types.go** (Core Data Models)
```
ENTITIES:
  LineType (enum)
    ├─ Context
    ├─ Add
    ├─ Delete
    └─ Placeholder

  Line (struct)
    ├─ Type: LineType
    └─ Content: string

  FileDiff (struct)
    ├─ Name: string
    ├─ OldPath: string
    ├─ NewPath: string
    ├─ IsNew: bool
    ├─ IsDeleted: bool
    ├─ IsBinary: bool
    ├─ AddCount: int
    ├─ DelCount: int
    ├─ LeftLines: []Line
    └─ RightLines: []Line

  Result (struct)
    └─ Files: []FileDiff

RELATIONSHIPS:
  Line        → used by FileDiff.LeftLines, FileDiff.RightLines
  LineType    → used by Line.Type
  FileDiff    → collected in Result.Files
```

### **Table 2: parser/errors.go** (Error Handling)
```
ENTITIES:
  ErrNotGitRepo (error)
  ErrEmptyDiff (error)
  ErrInvalidDiff (error)

  GitError (struct)
    ├─ Args: []string
    ├─ Stderr: string
    ├─ Err: error
    ├─ Error(): string
    └─ Unwrap(): error

  ParseError (struct)
    ├─ Line: int
    ├─ Message: string
    ├─ Cause: error
    ├─ Error(): string
    └─ Unwrap(): error

RELATIONSHIPS:
  GitError    → returned by parser/git.go operations
  ParseError  → returned by parser/unified.go operations
```

### **Table 3: parser/options.go** (Configuration)
```
ENTITIES:
  Options (struct)
    ├─ GitPath: string (default: "git")
    └─ WorkDir: string

  Option (func type)
    └─ signature: func(*Options)

  Functions:
    ├─ DefaultOptions() → Options
    ├─ WithGitPath(path) → Option
    └─ WithWorkDir(dir) → Option

RELATIONSHIPS:
  Option      → applied to Options via functional pattern
  Options     → used by parser/parser.go.Parser
```

### **Table 4: parser/git.go** (Git Integration)
```
ENTITIES:
  GitRunner (struct)
    ├─ gitPath: string
    ├─ workDir: string
    ├─ NewGitRunner(path, dir) → GitRunner
    ├─ IsGitRepository(ctx) → bool
    ├─ RunDiff(ctx, args...) → string, error
    └─ FindGitRoot(ctx) → string, error

RELATIONSHIPS:
  GitRunner   → stored in parser/parser.go.Parser.git
  Options     → provides GitPath and WorkDir
  GitError    → returned by RunDiff/FindGitRoot
```

### **Table 5: parser/parser.go** (Main Parser Orchestrator)
```
ENTITIES:
  Parser (struct)
    ├─ opts: Options
    ├─ git: *GitRunner
    ├─ New(...Option) → Parser
    ├─ ParseString(input) → Result, error
    ├─ ParseReader(r) → Result, error
    ├─ ParseGitDiff(ctx, args...) → Result, error
    └─ IsGitRepository(ctx) → bool, error

RELATIONSHIPS:
  Parser      ← created by main.go
  Options     → stored in Parser.opts
  GitRunner   → stored in Parser.git
  Result      → returned by Parse* methods
  ParseError  → returned by Parse* methods
  (delegates to parseUnified)
```

### **Table 6: parser/unified.go** (Unified Diff Parsing)
```
ENTITIES:
  hunk (internal struct)
    ├─ oldStart: int
    ├─ oldCount: int
    ├─ newStart: int
    ├─ newCount: int
    └─ lines: []diff.Line

  Functions:
    ├─ parseUnified(input) → Result, error
    ├─ parseFileDiff(lines, pos) → FileDiff, int, error
    ├─ parseHunk(lines, pos) → hunk, int, error
    ├─ diffGitRE (regex pattern)
    ├─ hunkHeaderRE (regex pattern)
    └─ binaryFileRE (regex pattern)

RELATIONSHIPS:
  hunk        → internal structure for parsing
  FileDiff    → built from parsed file hunks
  Line        → populated from parsed hunk content
  alignHunks  → called to convert hunks to FileDiff.LeftLines/RightLines
  ParseError  → returned on parsing errors
```

### **Table 7: parser/aligner.go** (Side-by-Side Alignment)
```
ENTITIES:
  Functions:
    ├─ alignHunks([]hunk) → []Line, []Line
    ├─ alignHunk(hunk) → []Line, []Line
    └─ alignBlock(deletes, adds, left, right) → void

  Logic:
    ├─ Context lines → appear on both sides
    ├─ Delete lines → left side with placeholder on right
    ├─ Add lines → right side with placeholder on left
    └─ Paired Delete+Add → same row

RELATIONSHIPS:
  alignHunks  ← called by parser/unified.go.parseUnified
  Line        → generates aligned Line pairs with Placeholder type
  FileDiff    → result populates FileDiff.LeftLines/RightLines
```

### **Table 8: tui/model.go** (TUI Application State)
```
ENTITIES:
  FocusedPanel (enum)
    ├─ FocusFileList
    ├─ FocusLeftDiff
    └─ FocusRightDiff

  Model (struct) - implements tea.Model
    ├─ width, height: int
    ├─ focused: FocusedPanel
    ├─ keys: KeyMap
    ├─ files: []diff.FileDiff
    ├─ selectedIdx: int
    ├─ leftViewport, rightViewport: viewport.Model
    ├─ syncScroll: bool
    ├─ ready: bool
    ├─ NewModel(files) → Model
    ├─ Init() → tea.Cmd
    ├─ Update(msg) → (Model, tea.Cmd)
    ├─ View() → string
    ├─ renderFileListPanel() → string
    ├─ renderDiffPanel() → string
    ├─ renderDiffLines(lines) → string
    ├─ scrollUp/Down() → void
    ├─ updateDiffContent() → void
    └─ updateViewportSizes() → void

RELATIONSHIPS:
  Model       ← created by main.go
  FileDiff    → stored in Model.files
  KeyMap      → stored in Model.keys
  FocusedPanel → used in Model.focused
  styles.*    → used by render* methods
```

### **Table 9: tui/keys.go** (Keyboard Bindings)
```
ENTITIES:
  KeyMap (struct)
    ├─ Up, Down, Left, Right: key.Binding
    ├─ Tab, ShiftTab: key.Binding
    ├─ PageUp, PageDown: key.Binding
    ├─ HalfPageUp, HalfPageDown: key.Binding
    ├─ Enter, Quit: key.Binding
    ├─ Help, SyncToggle: key.Binding
    ├─ DefaultKeyMap: KeyMap
    ├─ ShortHelp() → [][]key.Binding
    └─ FullHelp() → [][]key.Binding

RELATIONSHIPS:
  KeyMap      ← stored in tui/model.go.Model.keys
  (no direct dependencies)
```

### **Table 10: tui/styles.go** (Styling/Theming)
```
ENTITIES:
  Style variables:
    ├─ PanelStyle
    ├─ FocusedPanelStyle
    ├─ TitleStyle
    ├─ TitleInactiveStyle
    ├─ AddLineStyle
    ├─ DeleteLineStyle
    ├─ ContextLineStyle
    ├─ FileItemStyle
    ├─ FileItemSelectedStyle
    ├─ AddCountStyle
    ├─ DelCountStyle
    ├─ LineNumStyle
    ├─ PlaceholderStyle
    ├─ EmptyLineStyle
    └─ HelpStyle

RELATIONSHIPS:
  All styles → used by tui/model.go render* methods
  (no direct dependencies on other modules)
```

### **Table 11: main.go** (Application Entry Point)
```
ENTITIES:
  Functions:
    ├─ main() → orchestrates application startup
    ├─ runTUI(model) → error
    └─ handleError(err) → void

DEPENDENCIES:
  ├─ parser.New()          → creates Parser
  ├─ parser.ParseGitDiff() → gets Result
  ├─ tui.NewModel()        → creates Model
  ├─ tea.NewProgram()      → creates TUI program
  └─ handles all error types

RELATIONSHIPS:
  main() orchestrates:
    parser/parser.go → for parsing git diff
    tui/model.go     → for TUI rendering
    diff/types.go    → uses Result, FileDiff
    parser/errors.go → handles error types
```

---

## Dependency Graph (Call Flow)

```
ENTRY POINT
│
└─ main()
   ├─ parser.New(opts...)           [parser/parser.go]
   │  └─ NewGitRunner()             [parser/git.go]
   │
   ├─ parser.ParseGitDiff()         [parser/parser.go]
   │  ├─ git.RunDiff()              [parser/git.go]
   │  ├─ parseUnified()             [parser/unified.go]
   │  │  ├─ parseFileDiff()         [parser/unified.go]
   │  │  │  └─ parseHunk()          [parser/unified.go]
   │  │  │     └─ alignHunks()      [parser/aligner.go]
   │  │  └─ alignBlock()            [parser/aligner.go]
   │  └─ returns Result             [diff/types.go]
   │
   ├─ tui.NewModel(files)           [tui/model.go]
   │  └─ uses KeyMap                [tui/keys.go]
   │  └─ uses styles                [tui/styles.go]
   │
   ├─ tea.NewProgram(model)         [bubbletea]
   │  └─ model.Update()             [tui/model.go]
   │     └─ model.View()            [tui/model.go]
   │        ├─ renderFileListPanel()
   │        ├─ renderDiffPanel()
   │        └─ renderDiffLines()
   │
   └─ handleError(err)              [main.go]
      └─ matches error types        [parser/errors.go]

DATA FLOW
│
└─ FileDiff (from diff/types.go)
   ├─ populated by parseUnified()
   ├─ aligned by alignHunk()
   ├─ stored in Result.Files
   ├─ passed to Model.files
   └─ rendered in View()
```

---

## Entity Cardinality & Relationships

```
main.go              1:1  parser/parser.go        (creates one Parser)
main.go              1:1  tui/model.go            (creates one Model)

parser/parser.go     1:1  parser/git.go           (contains one GitRunner)
parser/parser.go     1:1  parser/options.go       (contains one Options)
parser/parser.go     1:N  diff/types.go::Result   (returns Result with N FileDiff)

tui/model.go         1:1  tui/keys.go::KeyMap     (contains one KeyMap)
tui/model.go         1:N  diff/types.go::FileDiff (stores multiple FileDiff)
tui/model.go         1:2  viewport.Model          (left and right viewports)

parser/unified.go    N:1  parser/aligner.go       (aligns N hunks)
parser/aligner.go    N:2  diff/types.go::Line     (produces left & right line arrays)

diff/types.go::Result         1:N  diff/types.go::FileDiff
diff/types.go::FileDiff       1:N  diff/types.go::Line
diff/types.go::Line           1:1  diff/types.go::LineType
```

---

## Cross-Module Imports Summary

```
main.go
  ├─ imports: parser, tui, diff
  └─ external: github.com/charmbracelet/bubbletea

parser/parser.go
  ├─ imports: parser/git, parser/options, parser/unified, diff
  └─ external: standard library

parser/git.go
  ├─ imports: parser/errors
  └─ external: os/exec, context

parser/unified.go
  ├─ imports: diff, parser/aligner, parser/errors
  └─ external: regexp, strings

parser/aligner.go
  ├─ imports: diff
  └─ external: none

parser/errors.go
  ├─ imports: none
  └─ external: standard library

parser/options.go
  ├─ imports: none
  └─ external: standard library

tui/model.go
  ├─ imports: diff, tui/keys, tui/styles
  └─ external: github.com/charmbracelet/{bubbletea,bubbles,lipgloss}

tui/keys.go
  ├─ imports: none
  └─ external: github.com/charmbracelet/bubbletea

tui/styles.go
  ├─ imports: none
  └─ external: github.com/charmbracelet/lipgloss

diff/types.go
  ├─ imports: none
  └─ external: standard library
```

---

## Key Architectural Layers

```
┌─────────────────────────────────────────────────────┐
│              PRESENTATION LAYER                     │
│  ┌──────────────────────────────────────────────┐   │
│  │  tui/model.go    (UI State & Rendering)      │   │
│  └────────────────────────────────────────────  ┤   │
│  │  tui/keys.go     (Input Handling)            │   │
│  └─────────────────────────────────────────────┤   │
│  │  tui/styles.go   (Visual Styling)            │   │
│  └──────────────────────────────────────────────┘   │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│              BUSINESS LOGIC LAYER                   │
│  ┌──────────────────────────────────────────────┐   │
│  │  parser/parser.go   (Main Orchestrator)      │   │
│  └──────────────────────────────────────────────┤   │
│  │  parser/unified.go  (Diff Parsing)           │   │
│  └─────────────────────────────────────────────┤   │
│  │  parser/aligner.go  (Line Alignment)         │   │
│  └──────────────────────────────────────────────┘   │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│              DATA/UTILITY LAYER                     │
│  ┌──────────────────────────────────────────────┐   │
│  │  parser/git.go      (Git Integration)        │   │
│  └──────────────────────────────────────────────┤   │
│  │  parser/errors.go   (Error Types)            │   │
│  └─────────────────────────────────────────────┤   │
│  │  parser/options.go  (Configuration)          │   │
│  └──────────────────────────────────────────────┘   │
└────────────────────┬────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│              DOMAIN MODEL LAYER                     │
│  ┌──────────────────────────────────────────────┐   │
│  │  diff/types.go      (Core Data Structures)   │   │
│  │    - LineType       (Context, Add, Delete)   │   │
│  │    - Line           (Type + Content)         │   │
│  │    - FileDiff       (File Diff Metadata)     │   │
│  │    - Result         (Collection of Files)    │   │
│  └──────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────┐
│            EXTERNAL DEPENDENCIES                    │
│  - git (shell command)                              │
│  - bubbletea (TUI framework)                        │
│  - bubbles (TUI components)                         │
│  - lipgloss (Terminal styling)                      │
└─────────────────────────────────────────────────────┘
```

---

## Summary Statistics

| Metric | Count |
|--------|-------|
| Total Files | 11 |
| Core Packages | 3 (diff, parser, tui) |
| Structs Defined | 10 |
| Functions/Methods | 30+ |
| Error Types | 5 |
| Supported Enums | 2 (LineType, FocusedPanel) |
| External Dependencies | 3 (bubbletea, bubbles, lipgloss) |
| Test Files | 3 (parser_test, unified_test, aligner_test) |

