# AutoDev — Comprehensive Security & Quality Audit

**Date:** 2026-05-31
**Date:** 2026-06-01

---

## Executive Summary

This audit reviewed the repository code, frontend (Next.js), backend (Go CLI packages), Dockerfile, CI release workflow, dependency manifests, and installation/catalog scripts. The project contains reasonable structure and security awareness, but several high-severity supply-chain and command-execution risks exist that would prevent safe production use without remediation.

Summary scores (subjective, evidence-based):
- Security Score: 5/10
- Code Quality Score: 7/10
- Architecture Score: 7/10
- Documentation Score: 8/10
- Cross-Platform Score: 6/10
- Production Readiness Score: 5/10
- Overall Project Grade: B-

---

## Critical Findings (Immediate action required)

1) Remote script execution / supply-chain risk
- Evidence:
  - `packages/installer/installer.go` contains many install commands that pipe remote scripts into a shell, e.g.:
    - `curl -fsSL https://deb.nodesource.com/setup_22.x | sudo -E bash -` (installer runtimes map)
    - `curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y`
    - `curl -fsSL https://get.docker.com | sh`
    - `curl -fsSL https://bun.sh/install | bash`
    (See `packages/installer/installer.go`, runtimes map)
  - `packages/catalog/catalog.yaml` contains install `script` entries that invoke remote scripts (e.g. `curl ... | sudo -E bash -` and `curl --proto ... | sh`).
- Severity: Critical
- Impact: Remote code execution on user machines; supply-chain compromise allows arbitrary code execution with elevated privileges (`sudo`) during installation.
- Remediation:
  - Remove or avoid piping remote content directly to shell. Prefer package manager installs (signed packages) or download + checksum verification before executing.
  - Where remote installers are unavoidable, validate signatures/checksums and present them to users; run installers unprivileged when possible; require explicit user consent and show human-readable provenance.
  - Add a strict allowlist for remote hosts and fail safe if a script or URL changes.

2) Shell command execution with `sh -c` and dynamic inputs (command injection risk)
- Evidence:
  - `Installer.Install` runs each `cmdStr` with `exec.Command("sh", "-c", cmdStr)` (packages/installer/installer.go).
  - `catalog.Package.IsInstalled()` runs `exec.CommandContext(ctx, "sh", "-c", p.Verify)` where `p.Verify` is loaded from `catalog.yaml` (`packages/catalog/catalog.go`), allowing arbitrary verify commands to execute.
  - `catalog` and `installer` use constructs like `cmd := exec.Command(parts[0], parts[1:]...)` in some places, but accept string-based commands in others.
- Severity: High
- Impact: If any of the catalog YAML entries or runtime commands are modified by an attacker (or a malicious third-party dependency), arbitrary shell execution can occur on end-user systems.
- Remediation:
  - Avoid `sh -c` for untrusted inputs. Use `exec.Command` with explicit args when possible.
  - Treat YAML-sourced commands as untrusted: validate/whitelist allowed commands and arguments, and restrict the set of executable verbs.
  - Provide dry-run and explicit confirmation flows; do not auto-run commands requiring `sudo` without explicit user action.

---

## High Severity Findings

1) Use of `sudo` in install scripts (privilege escalation risk and brittle UX)
- Evidence: `sudo apt-get install -y ...`, `sudo tar -C /usr/local -xzf -`, `sudo usermod -aG docker $USER` in `packages/installer/installer.go` and `packages/catalog/catalog.yaml`.
- Impact: Scripts assume the user has `sudo` and will run commands as root; this is both a security risk and causes failures on non-standard systems.
- Remediation: Detect privilege level and require explicit consent; prefer unprivileged installation methods or clearly document required privileges.

2) Catalog-driven Verify/Script execution (trust boundary violation)
- Evidence: `packages/catalog/catalog.go` uses embedded `catalog.yaml` `verify` and `script` fields to run platform commands (`exec.CommandContext(ctx, "sh", "-c", p.Verify)`).
- Impact: The catalog is a central trust boundary; if modified upstream, it can run arbitrary commands on user's hosts.
- Remediation: Harden the catalog: sign the YAML, pin catalog versions, validate script contents against allowed patterns, and avoid executing arbitrary shell scripts directly.

3) OS detection assumptions and deprecated tools
- Evidence: `packages/core/osinfo/detect.go` uses Windows `wmic` calls (`wmic ComputerSystem get TotalPhysicalMemory /Value`) and uses `sw_vers` and `sysctl` on macOS. (`packages/core/osinfo/detect.go`)
- Impact: `wmic` is deprecated on some Windows versions; detection may fail or produce inaccurate results on modern Windows. Many Linux distros lack `lsb_release` or `sudo` in minimal containers.
- Remediation: Use cross-platform libraries where possible or robust fallbacks. Document behavior when a command is missing.

---

## Medium Severity Findings

1) CI/Secrets — use of GitHub Secrets appears correct but surface review is required
- Evidence: `.github/workflows/release.yml` uses `GITHUB_TOKEN` and `HOMEBREW_TAP_GITHUB_TOKEN` from GitHub Actions secrets. (`.github/workflows/release.yml` lines ~34-36)
- Impact: No hardcoded secrets found in repository files. Ensure secrets in CI are scoped and rotated.
- Remediation: Keep secrets in repository secrets, avoid printing them in logs, and restrict permissions. Add audit to ensure no accidental secret writes to artifacts.

2) Docker runtime includes several packages (image size and surface area)
- Evidence: `Dockerfile` installs `bash`, `curl`, `wget`, `git`, `ca-certificates` in runtime image (`Dockerfile`).
- Impact: Larger image and larger attacker surface. Installing `git` and `curl` in runtime could be unnecessary.
- Remediation: Consider a minimal runtime (scratch/distroless) if distribution does not require these tools; or create separate dev/prod images.

3) Network timeouts and resilience
- Evidence: `packages/scanner/security.go` uses `client := &http.Client{Timeout: 3 * time.Second}` for OSV queries and `AuditRepository` uses an 8s global context. (`packages/scanner/security.go`). `packages/github/client.go` uses `30s` timeout.
- Impact: Very short timeouts can cause false negatives on slow networks; error handling usually returns errors but could be surfaced to users.
- Remediation: Increase timeouts for vulnerable network calls or implement retry/backoff.

---

## Low Severity Findings

1) `dangerouslySetInnerHTML` usage in `apps/website/app/layout.tsx` — used with static JSON
- Evidence: `layout.tsx` injects JSON-LD via `dangerouslySetInnerHTML` with `JSON.stringify({...})`. (`apps/website/app/layout.tsx`)
- Impact: Low — content is static and not user-supplied. Approve but document why safe.

2) Some `exec.Command` uses for convenience in `packages/cli/*` (e.g., `clone` command using `git clone`) — acceptable but should check inputs.

---

## Dependency Risks & Recommendations

Findings:
- JavaScript: `apps/website/package.json` lists `next@16.2.6`, `react@19.2.4`. Root project uses `pnpm` workspace and `turbo`.
- Go: Multiple modules present under `packages/*` with `go 1.22.2` across modules; Go packages include `github.com/spf13/viper`, `cobra`, etc.
- Lockfiles present: `pnpm-lock.yaml` (very large) — use `pnpm audit` or supply-chain scanners.

Recommendations (evidence-based actions):
- Run `pnpm audit` / `npm audit` and review `pnpm-lock.yaml` for known vulnerabilities.
- Run `govulncheck` and `gosec` against Go code.
- Use a dependency scanning pipeline (OSV, Dependabot, GitHub code scanning) to create PRs for fixes.
- Pin direct dependencies where appropriate and avoid `latest` in production-critical paths.

---

## Code Quality Review (selected notes)

- Code is generally well structured: packages map logically (`core`, `cli`, `scanner`, `installer`, `catalog`, `github`, `skills`).
- Some code smells:
  - Large `catalog.yaml` with embedded commands increases surface area and makes code review of runtime commands harder.
  - Repeated patterns for `exec.Command("sh","-c", ...)` — centralize execution helper with validation.
  - Hard-coded install commands across two places (`packages/installer/installer.go` and `packages/catalog/catalog.yaml`) — duplication.

Recommendations:
- Consolidate installation logic into a single module with safe execution helper.
- Reduce duplication and apply input validation and whitelisting.
- Add static analysis in CI (gofmt/go vet/gosec for Go; eslint and type checking for Next.js).

Maintainability & Readability: 7/10

---

## Architecture Review

- Strengths: Clear separation of concerns (CLI core vs installer vs catalog vs website). Package layout is idiomatic for Go modules.
- Risks:
  - Catalog YAML currently is a single source-of-truth but also a dangerous execution surface (scripts + verify commands). Treat it as untrusted input.
  - Installer executes platform-specific shell snippets with `sudo` and piped scripts — not safe for unattended production use.

Suggested patterns:
- Use the Builder pattern or Strategy pattern for platform install strategies (encapsulate platform behavior behind interfaces).
- Use Signed Manifests: sign the catalog and verify signatures before execution.

---

## Cross-Platform Issues

Observed issues and fixes:
- `sudo` and `apt-get` assumptions (Linux): minimal distributions, containers, or non-debian distros may lack these tools.
  - Fix: Detect package manager earlier and provide a fallback or explicit instructions.
- `wmic` used for Windows detection may be absent; prefer querying via Windows APIs or PowerShell (or gracefully fallback).
  - File: `packages/core/osinfo/detect.go`.
- Shell assumptions on Windows: many `sh -c` executions assume a POSIX shell. For Windows, use PowerShell or `cmd` equivalents, or require WSL.
  - Fix: Use platform-conditional execution and clear user messaging.

---

## DevOps & CI Review

- `.github/workflows/release.yml` looks standard (GoReleaser & Docker build/push using `GITHUB_TOKEN`). (`.github/workflows/release.yml`)
- No `.github/workflows` entries found that expose secrets directly in logs (observed usage uses `${{ secrets.GITHUB_TOKEN }}`).

Recommendations:
- Add CI job for static analysis (gosec, govet, pnpm audit, pnpm test, lint) and supply-chain scanning.
- Consider using reproducible builds and signed artifacts for release.
- Minimize runtime image by splitting dev vs prod Dockerfile or adopt distroless.

---

## Installation Experience (first-time user)

Observations:
- README, ROADMAP, and docs are present and helpful. Default `autodev` config uses `viper` and `AUTODEV_` env prefix (`packages/core/config/config.go`).
- Installer runs scripts that require root/sudo and network; this may fail in restricted environments (CI, enterprise endpoints).

Recommendations:
- Improve instructions and preflight checks (detect missing prerequisites and offer scripted remediation or explicit `autodev setup --yes`).
- Confirm network requirements and clearly document required environment variables (e.g., `AUTODEV_GITHUB_TOKEN` if desired to avoid rate limits). Evidence: `packages/github/client.go` mentions using `AUTODEV_GITHUB_TOKEN` to increase rate limits.

---

## Testing & Coverage

- There are some tests (e.g., `packages/catalog/catalog_test.go`, `packages/scanner/scanner_test.go`), but coverage appears incomplete for critical code paths (installer, catalog execution paths, CI pipelines). Run `go test ./...` to measure coverage.

Recommendations:
- Add unit/integration tests for install path logic but avoid actually invoking package installs in tests — use mocks.
- Add security-focused tests: ensure untrusted catalog contents are rejected.

---

## Prioritized Fix Roadmap (actionable)

1. (Critical) Eliminate `curl | sh` usage or gate it strongly with checksum/signature verification and explicit user confirmation.
2. (Critical) Replace `sh -c` executions for untrusted inputs; implement an execution helper with argument whitelisting.
3. (High) Harden `catalog.yaml`: sign and verify the catalog, or move shell scripts to versioned, reviewable asset files.
4. (High) Add CI jobs for `gosec`, `govulncheck`, `pnpm audit`, `eslint`, and unit tests coverage gate.
5. (Medium) Reduce runtime Docker image surface (move dev tools out of runtime image).
6. (Medium) Improve cross-platform detection fallbacks and document required tools and privileges.
7. (Low) Add better error messages and retry/backoff for external network calls.

---

## Final Scorecard

- Security Score: 5/10 (supply-chain & shell execution risks)
- Code Quality Score: 7/10
- Architecture Score: 7/10
- Documentation Score: 8/10
- Cross-Platform Score: 6/10
- Production Readiness Score: 5/10
- Overall Project Grade: B-

---

## Conclusion & Next Steps

This repository shows solid engineering practices and good documentation. The principal blocker for production is the installer/catalog's execution of remote scripts and shell commands. Addressing the supply-chain execution model and hardening the execution path will significantly improve security posture and production readiness.

If you'd like, I can:
- Produce focused patches that replace `sh -c` with safe `exec.Command` calls and add allowlists.
- Add a CI job to run `gosec` and `pnpm audit` and surface results.
- Implement catalog signing and verification helper.


---

**Appendix — Quick Evidence Index**
- `packages/installer/installer.go` — runtimes map and `Install()` (exec with `sh -c`), e.g. `curl ... | sh` (search for `curl` occurrences).
- `packages/catalog/catalog.yaml` — many `script` entries that run remote installers and `verify` fields executed via `sh -c`.
- `packages/catalog/catalog.go` — `IsInstalled()` executes `p.Verify` with `sh -c`.
- `packages/core/osinfo/detect.go` — OS detection using `sw_vers`, `wmic`, `sysctl`.
- `.github/workflows/release.yml` — uses `${{ secrets.GITHUB_TOKEN }}` appropriately.
- `apps/website/app/layout.tsx` — `dangerouslySetInnerHTML` usage with static JSON-LD.


Generated by automated repo scan and manual review.
