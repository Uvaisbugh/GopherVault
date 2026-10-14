# 🐹 GopherVault

> A fast, local-first command-line knowledge manager built in Go.

Developers repeatedly search for old terminal commands, code snippets, config files, and bug fixes they've already solved. **GopherVault** provides a local command-line vault to store, categorize, tag, and instantly search that knowledge offline.

---

## 🚀 Features

- `gopher add` — Save a note, terminal command, or code snippet (flags or interactive).
- `gopher search <query>` — Full-text search across titles, snippet contents, and tags.
- `gopher list` — List entries filtered by `--tag` and/or `--category` in a clean table.
- `gopher show <id>` — Inspect full note, tags, timestamps, and formatted snippet.
- `gopher delete <id>` — Remove outdated notes (alias: `gopher rm`).
- `gopher export <file>` — Export vault entries to JSON for backups or sharing.
- `gopher import <file>` — Restore and merge entries from backup files.
- `gopher stats` — Display vault metrics, category breakdown, and popular tags.
- `--db <path>` — Global flag to specify a custom storage file for team or test vaults.

---

## 🛠️ Architecture

GopherVault follows clean architectural separation:

```
User Terminal Command (e.g. gopher add --title "Docker prune")
                        │
                        ▼
                [ cmd/ (Cobra) ]
           Parses CLI flags & arguments
                        │
                        ▼
            [ internal/service/ (VaultService) ]
       Validates input, aggregates stats, business rules
                        │
                        ▼
          [ internal/storage/ (Repository) ]
          Go interface separating data from logic
                        │
                        ▼
          [ internal/storage/ (JSONStore) ]
       Reads/writes thread-safe JSON to disk
```

### Why this design?
By decoupling `VaultService` via the `storage.Repository` interface, the storage engine can be swapped with SQLite, PostgreSQL, or an in-memory test store without modifying a single line of business or CLI logic!

---

## 📚 Go Concepts Masterclass (Learn as You Build)

| Go Concept | Where It's Used | What It Teaches |
|---|---|---|
| **Packages & Modules** | `go.mod`, `main.go`, `internal/` | How Go organizes namespaces and privacy. Packages in `internal/` cannot be imported by external modules. |
| **Structs & Struct Tags** | `internal/models/entry.go` | How to model objects. Backtick annotations (e.g. `` `json:"id"` ``) dictate JSON serialization mapping. |
| **Methods & Receivers** | `func (e *Entry) HasTag(...)` | Functions attached to types. Pointer receiver `(e *Entry)` avoids copying large structs in memory. |
| **Interfaces** | `internal/storage/repository.go` | Go's implicit interfaces. No `implements` keyword — if a type has the required methods, it satisfies the interface. |
| **Error Handling** | Throughout (`errors.New`, `fmt.Errorf("%w")`) | Explicit error propagation using Go's `error` interface and `%w` error wrapping. |
| **Slices & Slicing** | `json_store.go`, `vault.go` | Dynamic arrays, `append()`, sub-slicing `slice[:i]`, and unpacking `...`. |
| **Maps & Sorting** | `service/vault.go` (`GetStats`) | Key-value hash tables and custom comparator sorting via `sort.Slice`. |
| **JSON & File I/O** | `json_store.go`, `vault.go` | `json.MarshalIndent`, `json.Unmarshal`, `os.ReadFile`, `os.WriteFile`, and directory management. |
| **Concurrency Safety** | `internal/storage/json_store.go` | `sync.RWMutex` with `mu.Lock()` and `mu.RLock()` for thread-safe concurrent reads/writes. |
| **Unit Testing** | `*_test.go` | Standard library `testing` package, table-driven tests, and `t.TempDir()`. |
| **Cobra CLI** | `cmd/` | Industry-standard CLI command trees, persistent flags, and subcommands. |

---

## 💻 Installation & Usage

### 1. Build the Binary

```bash
# Build for your current OS
go build -o gopher.exe main.go

# Or cross-compile for Linux:
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o bin/gopher-linux-amd64 main.go
```

### 2. Examples

```bash
# Add an entry with flags
./gopher add --title "Docker cleanup" --tag devops --body "docker system prune -a" --category devops

# Add a Git note
./gopher add --title "Git undo last commit" --tag "git,vcs" --body "git reset --soft HEAD~1" --category git

# Search notes
./gopher search "docker"

# List notes
./gopher list
./gopher list --tag devops
./gopher list --category git

# Show note details
./gopher show <id>

# View statistics
./gopher stats

# Back up and restore
./gopher export my-backup.json
./gopher import my-backup.json

# Delete a note
./gopher delete <id>
```

---

## 🧪 Running Automated Tests

Run the test suite with coverage:

```bash
go test -v ./...
```

---

## 📁 Project Directory Structure

```
GopherVault/
├── main.go                       # Application entrypoint
├── go.mod                        # Go module definition
├── go.sum                        # Cryptographic dependency checksums
├── .gitignore                    # Git ignore file
├── README.md                     # Documentation & learning guide
├── backup-example.json           # Sample exported backup
│
├── cmd/                          # CLI commands (Cobra)
│   ├── root.go                   # Root command & storage initialization
│   ├── add.go                    # gopher add
│   ├── search.go                 # gopher search
│   ├── list.go                   # gopher list
│   ├── show.go                   # gopher show
│   ├── delete.go                 # gopher delete / rm
│   ├── export.go                 # gopher export
│   ├── import.go                 # gopher import
│   └── stats.go                  # gopher stats
│
└── internal/                     # Private application packages
    ├── models/                   # Data structures
    │   ├── entry.go              # Entry struct & methods
    │   └── entry_test.go         # Model unit tests
    │
    ├── storage/                  # Storage abstractions & implementations
    │   ├── repository.go         # Repository interface
    │   └── json_store.go         # JSON file store implementation
    │
    └── service/                  # Business logic layer
        ├── vault.go              # VaultService implementation
        └── vault_test.go         # Service unit tests
```

---

## 📄 License

MIT License. Built for learning Go!
