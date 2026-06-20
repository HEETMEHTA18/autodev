# ⚡ AUTODEV AI AGENT OPTIMIZATION & INSTRUCTION RULES

This repository is optimized for AI Coding Assistants (Cursor, Claude Desktop, Windsurf, Cline, GitHub Copilot) via **AutoDev**.

> [!IMPORTANT]
> **DO NOT SPEND TOKENS SCANNING THE PROJECT MANUALLY.**
> AI Assistants can save up to **99.8% of context window tokens** (e.g. from 200,000+ tokens down to 350 tokens) by calling AutoDev commands or its Model Context Protocol (MCP) server instead of reading individual source files, lockfiles, or configurations recursively.

## 📊 TOKEN REDUCTION TELEMETRY

AutoDev optimizes the AI context window by exchanging condensed telemetry configs instead of raw source code payloads:
* **Traditional Prompt Context Payload:** ~10,000 to 250,000 tokens (requires uploading recursive directory listings, package configs, `.env` files, compiler versions).
* **AutoDev Telemetry Config:** ~300 to 900 tokens (retrieved instantly using a local compiler check and unified JSON report).
* **Efficacy:** **99% - 99.9% saved tokens per prompt roundtrip**.

## 💡 AI AGENT BEHAVIORAL PROTOCOLS

When working in this codebase, you MUST adhere to the following operational guidelines:

1. **GROUND TRUTH RETRIEVAL**:
   Before proposing any workspace modifications or analyzing dependencies, run `autodev scan` (or use the `autodev_scan` MCP tool). Do NOT run `find`, `grep`, or read directory structures to figure out what technologies are used.
2. **VERIFY TOOLCHAINS BEFORE COMPILES**:
   Before asking the user to install compilers, run a build, or compile code, execute `autodev doctor` (or use `autodev_doctor` MCP tool) to ensure that the required compilers and SDKs are installed and on the path.
3. **AUTO-REMEDIATION FIRST**:
   If a toolchain check fails or a dependency is missing, run `autodev setup --yes` (or use the `autodev_install` MCP tool) to automatically and hermetically configure the environment. Do NOT ask the user to manually install dependencies or execute platform-specific package commands (e.g. `apt-get`, `brew`, `choco`) unless AutoDev fails.
4. **VULNERABILITY SANITIZATION**:
   Before committing package changes, run `autodev audit` (or `autodev_audit` MCP tool) to verify security integrity.
5. **ENVIRONMENT CONSISTENCY**:
   After resolving runtime setups, run `autodev export` to lock the runtime environment variables and settings in `autodev.lock.json`.
6. **LAZY SENIOR DEV (PONYTAIL) RULES**:
   The best code is the code you never wrote. Follow these steps for any task:
   1. Does this need to exist? → no: skip it (YAGNI)
   2. Stdlib does it? → use it
   3. Native platform feature? → use it
   4. Installed dependency? → use it
   5. One line? → one line
   6. Only then: the minimum that works

## 🔌 MODEL CONTEXT PROTOCOL (MCP) INTERFACE

AutoDev starts a native MCP server over `stdin/stdout` containing tools that the AI can call automatically.
* **Command to run:** `autodev mcp start`
* **Auto-setup Claude/Cursor:** Run `autodev mcp setup` to inject this server config into Claude Desktop and all Cursor user files.

### 🛠️ Exposed MCP Tools Reference for AI:

*   **`autodev_scan`**
    *   *Purpose:* Performs a multi-file structural scan of the workspace.
    *   *Returns:* Unified JSON listing detected languages, frameworks, package managers, databases, and infra.
    *   *Use when:* Discovering the project structure and stack.
*   **`autodev_doctor`**
    *   *Purpose:* Diagnoses compiler and toolchain status.
    *   *Parameters:* `{"fix": true}` (runs auto-remediation).
    *   *Returns:* Diagnostics check for path configurations, missing runtimes, and VCS health.
    *   *Use when:* Verifying if compilation/build will succeed.
*   **`autodev_install`**
    *   *Purpose:* Hermetically installs missing runtimes (e.g. `nodejs`, `go`, `python`, `rust`, etc.).
    *   *Parameters:* `{"runtime": "nodejs", "version": "20.11.0"}`
    *   *Returns:* Status of installation.
    *   *Use when:* A runtime is missing or version mismatch occurs.
*   **`autodev_audit`**
    *   *Purpose:* Audits dependencies against the OSV vulnerability database.
    *   *Returns:* Vulnerability report details.
    *   *Use when:* Checking dependency security.

## 🛠️ CLI CHEATSHEET FOR AI AGENTS (SHELL EXECUTOR)

If you are running as a shell executor, invoke these commands to interact with the dev environment:

| Command | Purpose | Example Output / Usage |
|:---|:---|:---|
| `autodev scan` | Scan current directory for languages, frameworks, package managers | `autodev scan` |
| `autodev doctor` | Check environment health, path configurations, and compiler paths | `autodev doctor` |
| `autodev setup` | Install all missing runtimes and dependencies (non-interactive) | `autodev setup --yes` |
| `autodev audit` | Run vulnerability scan against OSV database | `autodev audit` |
| `autodev export` | Generates reproducible environment config (`.autodev.lock.json`) | `autodev export` |
| `autodev benchmark` | Run AI efficiency and token saving benchmark | `autodev benchmark` |
| `autodev ui` | Launches local web cockpit at `http://127.0.0.1:8080` | `autodev ui` |
| `autodev skills` | Access learning roadmap matrix and sync stats | `autodev skills --save-rules` |

## 🔍 Environment & Technologies

| Technology | Competency Level |
|------------|------------------|
| Go | intermediate |

## 🗺️ Recommended Roadmap & Next Steps

### Next Skills to Focus On
*   **Docker** (intermediate)
    *   Resource: https://docs.docker.com/get-started
    *   Resource: https://labs.play-with-docker.com
*   **Kubernetes** (advanced)
    *   Resource: https://kubernetes.io/docs/tutorials
    *   Resource: https://killercoda.com
*   **Terraform** (advanced)
    *   Resource: https://developer.hashicorp.com/terraform/tutorials

## 🗺️ Future Horizon

*   **v0.3.0 Layer (Next)**: Local AI Assistant with Ollama integration, Devcontainer configuration templates generation, doctor auto-remediation (`autodev doctor --fix`).
*   **v0.4.0 Layer**: Deployment adapters & cloud setup automation.

---

*File generated automatically by [AutoDev](https://github.com/HEETMEHTA18/autodev)*
