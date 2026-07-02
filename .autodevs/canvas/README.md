# Developer Knowledge Canvas

**Version:** 1.0.0  
**Path:** `.autodevs/canvas/`  
**Primary file:** `canvas.json`

## What is the Knowledge Canvas?

The Developer Knowledge Canvas is a comprehensive, machine-readable representation of your entire codebase. It goes far beyond a simple dependency graph — it captures file relationships, architecture, API routes, database models, security findings, git history, and AI-generated context that helps any AI agent understand your project in seconds.

Unlike a static diagram, the canvas is:
- **Generated** from the actual codebase (not a hand-maintained doc)
- **Multi-lens** — same data viewed as Technology, Files, Functions, APIs, Database, Security, Git, Runtime, or Architecture
- **AI-optimized** — designed to minimize token consumption while maximizing context
- **User-editable** — supports sticky notes, colored groups, and manual connections
- **Persistent** — saved to `.autodevs/canvas/canvas.json` and read by every AI session

## Quick Start

```bash
# Generate the canvas (terminal output)
autodev canvas

# Save canvas for AI sessions
autodev canvas --save

# Print compact AI-readable summary
autodev canvas --summary

# Output raw JSON
autodev canvas --json
```

## How AI Uses the Canvas

When an AI agent (AutoDevs/Ralph Loop) starts a session, it reads the canvas in this order:

1. `.autodevs/canvas/canvas.json` — full project knowledge
2. `ai_context` section — one-liner, key insights, patterns, gotchas, reading order
3. `files` — specific file cards for files relevant to the current task
4. `dependencies` — import graph for understanding impact of changes
5. `security` — vulnerabilities that might affect the task
6. `git` — recent changes and hot files related to the task

This replaces scanning the entire codebase and reduces token consumption by ~60-80%.

## Canvas Structure

```
autodevs/canvas/
├── canvas.json           # The generated canvas
├── canvas.schema.json    # JSON Schema for validation
├── notes.json            # User sticky notes (optional)
├── groups.json           # User-defined groups (optional)
├── connections.json      # User-drawn connections (optional)
└── README.md             # This file
```

## Sections

### `project` — Overview
Name, languages, frameworks, total files/LOC, monorepo status.

### `files` — File Cards
Every source file with:
- Language, LOC, imports, exports, functions
- File type (entry, component, util, config, route, test, etc.)
- Dependencies (what it imports) and dependents (what imports it)
- Security findings, TODOs, dead code status

### `architecture` — Layers & Entry Points
Detected architectural layers, entry points, directory tree, data flow.

### `dependencies` — Import Graph
File-level dependency edges, circular dependencies, dead files, impact map.

### `api` — Routes
Detected API routes from common frameworks (Express, Gin, Echo, Fastify, Django, Flask, etc.).

### `database` — Models
Detected database models from Go structs, Django models, Sequelize schemas, etc.

### `security` — Findings
Static security analysis: hardcoded secrets, code injection, XSS, SQL injection, etc.

### `git` — History Summary
Commit count, contributors, hot files, churn map.

### `ai_context` — AI-Optimized Summary
One-liner, key insights, patterns, conventions, gotchas, recommended reading order, critical files, FAQs, glossary.

### `lenses` — Multi-Lens Views
9 different views of the same project:
- Technology, File, Function, API, Database, Security, Git, Runtime, Architecture

### `notes`, `groups`, `connections` — User Canvas Layer
Editable sticky notes, colored groups, and manual connections for project planning.

## Regeneration

The canvas should be regenerated:
- Before starting a new Ralph Loop task
- After significant file changes
- When adding/removing dependencies
- On `autodev canvas --save`

## Schema Validation

Validate the canvas against the schema:
```bash
# Using jsonschema CLI
jsonschema -i .autodevs/canvas/canvas.json .autodevs/canvas/canvas.schema.json

# Or using ajv-cli
ajv validate -s .autodevs/canvas/canvas.schema.json -d .autodevs/canvas/canvas.json
```
