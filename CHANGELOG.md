# Changelog

All notable changes to AutoDev will be documented in this file.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

---

## [0.6.2] - 2026-10-27

### Security

- **Go toolchain → go1.26.10**: clears 14 stdlib advisories flagged by osv-scanner across all 8 modules (114 alerts) — `crypto/x509` distrust of compromised intermediaries (CVE-2026-111943), HTTP/1 desync, HTTP/2 HPACK/Trailer/initial-window DoS, `crypto/tls` ECH rejection, `os` junction escapes, and `html/template` context bugs. GoReleaser, CI, CodeQL, and the local `go.work` now all build on 1.26.10.
- **Next.js 16.2.11 → 16.4.3**: fixes 5 critical RCEs (GHSA-2464-fhr4-w9xw image-optimization SVG, GHSA-vcvr-r3jv-pc5j `next/og`, GHSA-673q-mjm6-r5xr `use cache`, GHSA-2xp9-vwfh-vxw4 AVIF, GHSA-p293-qw3h-jr36 Windows-hosted) plus SSRF, App Router DoS, and SSG/ISR cache-poisoning advisories; `eslint-config-next` aligned to ^16.4.3.
- **Transitive npm advisories**: `pnpm audit` now reports **0 vulnerabilities** — bounded overrides for browserslist ≥4.28.7, zod ≥4.6.6, source-map-js ≥1.2.2, js-yaml ≥5.4.1, sharp ≥0.35.5, baseline-browser-mapping ≥2.11.0, and per-family `brace-expansion` pins (1.1.21 / 5.0.12) that keep `minimatch@3` consumers off 5.x.

### Changed

- v0.6.2 supersedes the GitHub-only v0.6.1 and is the first release published to npm since 0.6.0 — all shipped binaries are built with the patched toolchain.

---

## [0.6.1] - 2026-10-27

### Added

- **Repo Health Score**: New `autodev score` command grades your repository 0–100 (A–F) across weighted checks — secrets and supply-chain findings ×3, git/env/lockfile hygiene ×2, linter availability ×1. Supports `--badge` (shields.io README badge), `--json` (machine-readable report), and `--min` (CI quality gate).
- **Disciplined AI Chat**: `autodev chat` now runs with a bounded thinking budget (`--thinking-budget`), a scope guard that keeps answers inside AutoDev's domain, anti-hallucination checks, Google-search grounding, and a two-phase `<<REFINED>>` prompt refinement pass.
- **`.autodevs/` Agent Protocol**: Structured multi-agent workflow (Planner → Context Builder → Builder → Reviewer → Verifier) with memory retrieval and a Stage 2b Refiner for self-reprompting.

### Changed

- **Minimal Galaxy Theme**: Full TUI overhaul on a shared design system — blue accent (`#0070f3`), single-width icons (no emoji), normal-weight borders, no italics, and new `Panel` / `StatusRow` / `ProgressBar` / heading helpers used across every command view.
- Regenerated AI rule files (`.autodev-skills.md`, `.clinerules`, `.cursorrules`, `.github/copilot-instructions.md`) to document the `score` command.

---

## [0.2.0] - 2026-05-31

### Added

- **Ecosystem Expansion**: Added detection and automated setup for Docker and Kubernetes (Helm, kubectl).
- **Mobile/SDK Toolchains**: Native Flutter SDK installer, Android SDK, and emulator setup.
- **JVM & Build Systems**: OpenJDK installer (Java 17/21), Kotlin build tools, and Apache Maven.
- **AI Rule Files & 99.8% Context Saving**: Automatic generation of standard AI assistant rule files (`.autodev-skills.md`, `.cursorrules`, `.clinerules`, and `.github/copilot-instructions.md`) when running any AutoDev command. Reduces context payload by 99.8% (~350 tokens instead of 200,000+).
- **Interactive Visual Roadmap TUI**: A beautiful BubbleTea visual terminal interface for navigating skill roadmaps.
- **Model Context Protocol (MCP) Integration**: Native MCP server (`autodev mcp start` and auto-configuration via `autodev mcp setup`) exposing `autodev_scan`, `autodev_doctor`, `autodev_install`, and `autodev_audit` tools directly to Cursor, Claude Desktop, and other LLM agents.
- **Security Dependency Auditing**: Deep repository auditing against the OSV vulnerability database (`autodev audit`).
- **Reproducible Environments**: Exporting environment config states into reproducible `.autodev.lock.json` manifests.
- **Local Web Console**: Interactive browser console cockpit launched via `autodev ui`.
- **Benchmark Command**: CLI execution and token optimization metrics dashboard (`autodev benchmark`).

---

## [0.1.0] - 2026-05-29

### Added

- Initial release 🎉
- Core CLI with `autodev scan` and `autodev setup`
- Support for Node.js, Python, Go, Rust, Java, Docker detection
- HTML, JSON, and Markdown report generation
- Multi-platform binaries (Linux x86_64/arm64, macOS, Windows)

[Unreleased]: https://github.com/HEETMEHTA18/autodev/compare/v0.6.2...HEAD
[0.6.2]: https://github.com/HEETMEHTA18/autodev/compare/v0.6.1...v0.6.2
[0.6.1]: https://github.com/HEETMEHTA18/autodev/compare/v0.6.0...v0.6.1
[0.2.0]: https://github.com/HEETMEHTA18/autodev/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/HEETMEHTA18/autodev/releases/tag/v0.1.0
