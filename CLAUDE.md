# laravel-docgen — CLAUDE.md

A Go tool that statically analyzes Laravel PHP projects and generates PlantUML-compatible JSON diagrams to visualize architecture and system flow.

---

## Project Mission

Parse Laravel PHP projects and generate structured JSON stored in `./docs` that represents:
- **Sequence diagrams** → `./docs/sequence/`
- **Object diagrams** → `./docs/objects/`
- **Use case diagrams** → `./docs/usecase/`

---

## Project Structure

```
.
├── CLAUDE.md
├── go.mod
├── go.sum
├── main.go
├── cmd/
│   └── root.go                  # Cobra CLI entry point
├── internal/
│   ├── parser/
│   │   ├── parser.go            # Filesystem walker, orchestrates sub-parsers
│   │   ├── routes.go            # routes/web.php, routes/api.php
│   │   ├── controllers.go       # app/Http/Controllers/**
│   │   ├── models.go            # app/Models/** (Eloquent relationships)
│   │   ├── middleware.go        # app/Http/Middleware/**
│   │   ├── services.go          # app/Services/**
│   │   ├── jobs.go              # app/Jobs/**
│   │   ├── events.go            # app/Events/**, app/Listeners/**
│   │   └── policies.go          # app/Policies/**
│   ├── analyzer/
│   │   ├── analyzer.go          # Builds IR from parsed data, validates relationships
│   │   ├── sequence.go          # Derives sequence diagram IR
│   │   ├── objects.go           # Derives object diagram IR
│   │   └── usecase.go           # Derives use case diagram IR
│   ├── generator/
│   │   ├── generator.go         # Orchestrates JSON file writing
│   │   ├── sequence.go          # Serializes sequence IR → JSON
│   │   ├── objects.go           # Serializes object IR → JSON
│   │   └── usecase.go           # Serializes use case IR → JSON
│   ├── model/
│   │   ├── ir.go                # Intermediate representation structs
│   │   └── plantuml.go          # PlantUML JSON schema structs
│   └── cache/
│       └── cache.go             # File-level parse cache for incremental runs
├── pkg/
│   └── phpast/
│       └── phpast.go            # PHP AST parsing wrapper (see Parsing Strategy)
├── tests/
│   └── fixtures/
│       └── laravel/             # Minimal Laravel fixture project for integration tests
│           ├── routes/
│           ├── app/
│           └── ...
└── docs/                        # Generated output
    ├── sequence/
    ├── objects/
    └── usecase/
```

---

## Architecture & Data Flow

```
Laravel Project (PHP files)
        │
        ▼
  internal/parser           ← AST-based PHP extraction per Laravel component type
        │
        ▼
  internal/analyzer         ← Build IR, resolve namespaces, validate relationships
        │
        ├──► analyzer/sequence.go   → SequenceDiagram IR
        ├──► analyzer/objects.go    → ObjectDiagram IR
        └──► analyzer/usecase.go    → UseCaseDiagram IR
                │
                ▼
  internal/generator        ← Marshal IR → deterministic JSON → write to ./docs/**
```

Parser and diagram logic must remain **strictly decoupled** — the analyzer IR is the only contract between them.

---

## Laravel Components to Parse

### Files & Directories
| Component      | Path                              |
|----------------|-----------------------------------|
| Web routes     | `routes/web.php`                  |
| API routes     | `routes/api.php`                  |
| Controllers    | `app/Http/Controllers/**/*.php`   |
| Middleware     | `app/Http/Middleware/**/*.php`     |
| Models         | `app/Models/**/*.php`             |
| Services       | `app/Services/**/*.php`           |
| Jobs           | `app/Jobs/**/*.php`               |
| Events         | `app/Events/**/*.php`             |
| Listeners      | `app/Listeners/**/*.php`          |
| Policies       | `app/Policies/**/*.php`           |

### Relationships to Detect
- Route → Middleware → Controller → Method
- Controller → Service → Model
- Model Eloquent: `hasOne`, `hasMany`, `belongsTo`, `belongsToMany`, `morphTo`, etc.
- Event dispatch (`event(new SomeEvent())`) → Listener
- Job dispatch (`dispatch(new SomeJob())`)
- Policy authorization (`$this->authorize()`)

---

## Parsing Strategy

### Use AST Parsing — Not Regex-Only
- Wrap a PHP AST parser via subprocess: invoke a small PHP CLI helper script (`php parse-helper.php <file>`) that uses `nikic/php-parser` and outputs a JSON AST to stdout. The Go layer reads and interprets the JSON AST. This avoids CGO while enabling real AST parsing.
- Regex is acceptable only for **simple top-level extraction** (e.g., finding `Route::get(...)` lines) when AST overhead is clearly not justified
- Never use regex as the sole mechanism for resolving namespaces, method calls, or dependency injection
- Gracefully handle parse errors — collect warnings, do not abort the full run

### Namespace & Import Resolution
- Track `use` statements and `namespace` declarations
- Resolve short class names to fully qualified names
- Detect constructor injection and `app()->make()` calls

### Laravel Magic Methods
- Recognize Eloquent relationship method signatures by name convention (`hasMany`, `belongsTo`, etc.) since they return dynamic types
- Recognize `__construct` DI patterns

---

## Intermediate Representation (IR)

All parsed data normalizes into an IR before any diagram is generated. Key types in `internal/model/ir.go`:

```go
type Project struct {
    Routes      []Route
    Controllers []Controller
    Models      []Model
    Services    []Service
    Middleware  []Middleware
    Jobs        []Job
    Events      []Event
    Listeners   []Listener
    Policies    []Policy
}

type Route struct {
    Method     string   // GET, POST, etc.
    Path       string
    Middleware []string
    Controller string
    Action     string
}

type Controller struct {
    Name      string
    Namespace string
    Actions   []Action
}

type Action struct {
    Name         string
    Dependencies []string // services, models referenced
}

type Model struct {
    Name          string
    Namespace     string
    Fields        []Field
    Relationships []Relationship
}

type Relationship struct {
    Type    string // hasMany, belongsTo, etc.
    Related string // target model name
    Name    string // method name
}
```

---

## JSON Schema (PlantUML Output)

All output is **deterministic** — keys are sorted, arrays are sorted by stable identifiers. This ensures clean git diffs.

### Sequence Diagram (`docs/sequence/<n>.json`)
```json
{
  "type": "sequence",
  "title": "POST /login",
  "participants": [
    { "alias": "Client", "label": "Client" },
    { "alias": "AuthMiddleware", "label": "AuthMiddleware" },
    { "alias": "AuthController", "label": "AuthController" },
    { "alias": "UserService", "label": "UserService" },
    { "alias": "User", "label": "User" }
  ],
  "messages": [
    { "from": "Client", "to": "AuthMiddleware", "label": "POST /login", "type": "sync" },
    { "from": "AuthMiddleware", "to": "AuthController", "label": "login()", "type": "sync" },
    { "from": "AuthController", "to": "UserService", "label": "authenticate()", "type": "sync" },
    { "from": "UserService", "to": "User", "label": "findByEmail()", "type": "sync" },
    { "from": "AuthController", "to": "Client", "label": "redirect()", "type": "return" }
  ]
}
```

### Object Diagram (`docs/objects/<n>.json`)
```json
{
  "type": "object",
  "title": "User Model",
  "objects": [
    {
      "name": "User",
      "fields": [
        { "name": "id", "type": "int" },
        { "name": "email", "type": "string" },
        { "name": "name", "type": "string" }
      ]
    }
  ],
  "relationships": [
    { "from": "User", "to": "Post", "label": "hasMany" },
    { "from": "User", "to": "Role", "label": "belongsToMany" }
  ]
}
```

### Use Case Diagram (`docs/usecase/<n>.json`)
```json
{
  "type": "usecase",
  "title": "Authentication",
  "actors": [
    { "name": "Guest" },
    { "name": "AuthenticatedUser" }
  ],
  "usecases": [
    { "id": "UC1", "label": "Login" },
    { "id": "UC2", "label": "Register" },
    { "id": "UC3", "label": "Logout" }
  ],
  "relationships": [
    { "actor": "Guest", "usecase": "UC1", "type": "association" },
    { "actor": "Guest", "usecase": "UC2", "type": "association" },
    { "actor": "AuthenticatedUser", "usecase": "UC3", "type": "association" }
  ]
}
```

---

## CLI Design

Binary name: `laravel-docgen`

```bash
# Analyze and generate all diagram types
laravel-docgen analyze ./my-laravel-app

# Generate specific diagram type only
laravel-docgen generate sequence --input ./my-laravel-app
laravel-docgen generate objects  --input ./my-laravel-app
laravel-docgen generate usecase  --input ./my-laravel-app
laravel-docgen generate all      --input ./my-laravel-app

# Flags
--input        string   Path to Laravel project root (required)
--output       string   Output directory (default: ./docs)
--verbose               Enable verbose/debug logging
--incremental           Skip files unchanged since last run (uses cache)
--validate              Validate generated JSON against schema after writing
```

Error messages must be actionable — include the file path and reason for any parse failure.

---

## Development Commands

```bash
# Run
go run main.go analyze ./path/to/laravel-app

# Build
go build -o laravel-docgen ./...

# Test all
go test ./...

# Test with coverage
go test -cover ./...

# Test specific package
go test ./internal/parser/...
go test ./internal/analyzer/...
go test ./internal/generator/...

# Regenerate golden files after intentional output changes
go test ./tests/... -update-golden

# Lint
golangci-lint run

# Format
gofmt -w .
```

---

## Coding Standards

- Write idiomatic Go — clarity over cleverness
- Use interfaces to decouple parsing from analysis from generation
- Document all exported functions and types
- Meaningful, wrapped error handling: `fmt.Errorf("parsing controller %s: %w", name, err)`
- No hardcoded Laravel paths — all paths configurable or derived from conventions
- Avoid overengineering — solve today's problem, design for tomorrow's extension

---

## Testing Strategy

### Unit Tests
- One `_test.go` per parser file
- Use `.php` fixture files in `internal/parser/testdata/`
- Test each component type: routes, controllers, models, middleware, events, etc.

### Golden File Tests
- A minimal but representative Laravel fixture project lives in `tests/fixtures/laravel/`
- Each test run generates JSON output and compares against committed golden files in `tests/golden/`
- Run `go test ./tests/... -update-golden` to regenerate after intentional changes
- Golden file comparison is the primary regression guard for output format stability

### Integration Tests
- Run the full pipeline on the fixture Laravel project
- Assert presence and shape of all output JSON files
- Assert determinism: run twice, output must be byte-identical

---

## Performance

- **Parallel file scanning** — use `golang.org/x/sync/errgroup` for concurrent directory walks
- **Parse caching** — cache parsed IR per file keyed by path + modtime in `internal/cache/`; skip re-parsing unchanged files on `--incremental` runs
- Large projects (500+ PHP files) must complete in under 30 seconds
- Do not load entire file trees into memory — stream and process

---

## Dependencies

```
github.com/spf13/cobra              # CLI framework
github.com/spf13/viper              # Config management
go.uber.org/zap                     # Structured logging
golang.org/x/sync                   # errgroup for parallel scanning
```

---

## Output File Naming

Files are named after their subject: lowercase, hyphens for spaces.

Examples: `auth-controller.json`, `user-model.json`, `post-login.json`

---

## Future Roadmap

Design decisions should not block these planned capabilities:

- Real-time architecture visualization (watch mode)
- IDE plugin integration
- CI/CD documentation pipeline (fail on architecture drift)
- Architecture drift detection (diff two snapshots)
- OpenAPI spec generation from routes
- GraphQL flow diagrams
- Multi-framework support (Symfony, CodeIgniter)
