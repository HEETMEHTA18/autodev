"use client";
import { useState } from "react";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";
import { Copy, Check } from "lucide-react";
import Image from "next/image";

interface DocSection {
  id: string;
  title: string;
  category: string;
  content: React.ReactNode;
}

const Callout = ({
  type,
  title,
  children,
}: {
  type: "info" | "warning" | "tip";
  title: string;
  children: React.ReactNode;
}) => {
  const styles = {
    info: {
      border: "border-l-4 border-l-[#2563eb] border-[#2A2A2A]",
      bg: "bg-[#2563eb]/5",
      icon: "ℹ️",
      color: "text-[#2563eb]",
    },
    warning: {
      border: "border-l-4 border-l-[#dc2626] border-[#2A2A2A]",
      bg: "bg-[#dc2626]/5",
      icon: "⚠️",
      color: "text-[#dc2626]",
    },
    tip: {
      border: "border-l-4 border-l-[#00FF87] border-[#2A2A2A]",
      bg: "bg-[#00FF87]/5",
      icon: "⚡",
      color: "text-[#00FF87]",
    },
  };

  const currentStyle = styles[type];

  return (
    <div
      className={`p-4 my-6 rounded border ${currentStyle.border} ${currentStyle.bg}`}
    >
      <div className="flex items-center gap-2 mb-2 font-bold font-mono text-sm uppercase tracking-wider">
        <span className="text-base">{currentStyle.icon}</span>
        <span className={currentStyle.color}>{title}</span>
      </div>
      <div className="text-neutral-400 text-sm font-sans">{children}</div>
    </div>
  );
};

export default function DocsClient() {
  const [activeSection, setActiveSection] = useState("introduction");
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const handleCopy = (text: string, id: string) => {
    navigator.clipboard.writeText(text);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  const renderCodeBlock = (id: string, code: string) => (
    <div
      key={id}
      className="docs-code-block relative group border-2 border-[#2A2A2A] bg-[#050505] p-4 font-mono text-sm my-4 rounded-md overflow-x-auto pr-24"
    >
      <button
        onClick={() => handleCopy(code, id)}
        className="docs-code-copy absolute top-2.5 right-2.5 text-xs bg-[#111] hover:bg-[#C8F135] hover:text-black hover:border-[#C8F135] border border-[#2A2A2A] px-2.5 py-1.5 rounded transition-all text-neutral-400 font-mono flex items-center gap-1.5 cursor-pointer"
      >
        {copiedId === id ? (
          <>
            <Check className="w-3.5 h-3.5 text-[#00FF87]" />
            <span className="text-[#00FF87]">Copied!</span>
          </>
        ) : (
          <>
            <Copy className="w-3.5 h-3.5" />
            <span>Copy</span>
          </>
        )}
      </button>
      <pre className="text-neutral-200 bg-transparent border-0 p-0 m-0">
        {code}
      </pre>
    </div>
  );

  const sections: DocSection[] = [
    {
      id: "introduction",
      category: "Getting Started",
      title: "Introduction",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            Introduction
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            AutoDev is an open-source, cross-platform developer environment
            bootstrapper. It acts as an
            <strong> &quot;App Store for Developers,&quot;</strong> simplifying
            complex toolchain setups through intelligent, profile-based
            automation.
          </p>
          <p className="text-neutral-300 leading-relaxed mb-4">
            Setting up local environments often involves manually downloading
            compilers, installing dependencies, configuring paths, and matching
            versions. AutoDev handles this end-to-end: it scans your project
            files, detects the programming languages, frameworks, package
            managers, and databases used, and installs any missing pieces with a
            single command.
          </p>
          <Callout type="tip" title="Why use AutoDev?">
            Bring team onboarding down from hours to minutes. Run{" "}
            <code className="text-[#00FF87] font-mono font-bold bg-[#111] px-1.5 py-0.5 rounded">
              autodev setup
            </code>{" "}
            in any newly cloned repository to configure the exact toolchains
            required.
          </Callout>
        </>
      ),
    },
    {
      id: "quick-start",
      category: "Getting Started",
      title: "Quick Start",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            Quick Start
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            Get up and running with AutoDev instantly. Run the bootstrap script
            for your platform, run it directly via Node.js (npx), or run it
            inside Docker without installation.
          </p>
          <h2 className="text-2xl font-bold text-[#C8F135] mt-8 mb-4">
            One-line Shell Install
          </h2>
          <p className="text-neutral-300 mb-2">
            For Linux and macOS, install using curl:
          </p>
          {renderCodeBlock(
            "qs-curl",
            "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash",
          )}

          <h2 className="text-2xl font-bold text-[#C8F135] mt-8 mb-4">
            Run via NPX
          </h2>
          <p className="text-neutral-300 mb-2">
            Run the interactive bootstrapper setup instantly using Node.js:
          </p>
          {renderCodeBlock("qs-npx", "npx @heetmehta18/autodev setup")}

          <h2 className="text-2xl font-bold text-[#C8F135] mt-8 mb-4">
            Interactive Setup
          </h2>
          <p className="text-neutral-300 mb-2">
            Once installed via shell/package manager, enter your project&apos;s
            directory and run:
          </p>
          {renderCodeBlock("qs-setup", "autodev setup")}

          <p className="text-neutral-300 leading-relaxed mt-4">
            AutoDev will analyze your codebase, display a list of recommended
            tools and libraries, and prompt you to install them.
          </p>
        </>
      ),
    },
    {
      id: "installation",
      category: "Getting Started",
      title: "Installation Methods",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            Installation Methods
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-6">
            AutoDev is available via NPM, Shell scripts, Homebrew, Scoop, and
            Docker. Pick the method that fits your workflow.
          </p>

          <Callout type="warning" title="Important Npm Naming Notice">
            The unscoped{" "}
            <code className="font-mono text-neutral-300">autodev</code> package
            on npm is owned by an unrelated author. Do <strong>not</strong>{" "}
            install <code className="font-mono text-neutral-300">autodev</code>.
            Instead, always use the scoped package:{" "}
            <strong>
              <code className="font-mono text-neutral-200">
                @heetmehta18/autodev
              </code>
            </strong>
            .
          </Callout>

          <div className="space-y-6">
            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                NPM (Global Install)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Install the official CLI globally via npm:
              </p>
              {renderCodeBlock(
                "inst-npm",
                "npm install -g @heetmehta18/autodev",
              )}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                NPX (Run instantly)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Run the interactive setup on the fly without permanent
                installation:
              </p>
              {renderCodeBlock("inst-npx", "npx @heetmehta18/autodev setup")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                Shell Script (Linux & macOS) — Recommended
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Downloads and installs the pre-compiled binary directly from
                GitHub:
              </p>
              {renderCodeBlock(
                "inst-curl",
                "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash",
              )}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                Homebrew (macOS / Linux)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Install globally via Brew tap:
              </p>
              {renderCodeBlock(
                "inst-brew",
                "brew install HEETMEHTA18/tap/autodev",
              )}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                Scoop (Windows)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Add the AutoDev bucket, then install:
              </p>
              {renderCodeBlock(
                "inst-scoop-bucket",
                "scoop bucket add autodev https://github.com/HEETMEHTA18/scoop-bucket",
              )}
              {renderCodeBlock("inst-scoop", "scoop install autodev")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                Docker Container
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Run inside a container mapping your repository directory:
              </p>
              {renderCodeBlock(
                "inst-docker",
                "docker run --rm -v $(pwd):/workspace ghcr.io/heetmehta18/autodev setup",
              )}
            </div>
          </div>
        </>
      ),
    },

    {
      id: "cmd-doctor",
      category: "Commands",
      title: "autodev doctor",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev doctor
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              doctor
            </code>{" "}
            command inspects your system&apos;s specifications, active package
            managers, and checks the installation status of managed runtimes.
          </p>
          {renderCodeBlock("cmd-doc-run", "autodev doctor")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It performs a diagnostic check of languages, framework execution
            environments, databases, container configurations, DevOps
            infrastructure, and mobile development SDKs.
          </p>
          <Callout type="info" title="Diagnostic Output">
            It lists:
            <ul className="list-disc pl-5 mt-2 space-y-1 text-neutral-400 font-sans">
              <li>System OS, Kernel, CPU Core Count, Memory Capacity</li>
              <li>
                Status of compilers & runtimes (e.g. Go, Python, Node, Java,
                Docker, Rust)
              </li>
              <li>
                Missing components, coupled with recommendations on how to
                install them
              </li>
            </ul>
          </Callout>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-doctor.png"
            alt="autodev doctor CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-scan",
      category: "Commands",
      title: "autodev scan",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev scan
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              scan
            </code>{" "}
            command analyzes your current working directory for configuration
            markers, lockfiles, and structures.
          </p>
          {renderCodeBlock("cmd-scan-run", "autodev scan")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            Instead of simply looking at file extensions, AutoDev queries
            package registries (like{" "}
            <code className="font-mono text-neutral-400">package.json</code>,{" "}
            <code className="font-mono text-neutral-400">go.mod</code>, or{" "}
            <code className="font-mono text-neutral-400">requirements.txt</code>
            ) to build a precise map of frameworks, databases, and secondary
            development tools.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-scan.png"
            alt="autodev scan CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-audit",
      category: "Commands",
      title: "autodev audit",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev audit
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              audit
            </code>{" "}
            command scans your project&apos;s lockfiles and dependency lists and
            queries the OSV (Open Source Vulnerabilities) database.
          </p>
          {renderCodeBlock("cmd-audit-run", "autodev audit")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It checks Python, Node.js, and Go dependencies for known
            supply-chain vulnerabilities, compromised packages, and security
            advisories, listing the CVE and severity level for any threats.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-audit.png"
            alt="autodev audit CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-setup",
      category: "Commands",
      title: "autodev setup",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev setup
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              setup
            </code>{" "}
            command scans the project and aligns your local development
            environment.
          </p>
          {renderCodeBlock("cmd-setup-run", "autodev setup")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It performs download caching, automated binary extraction,
            environment variable linking (PATH configuration), and triggers
            package-manager installs (like running{" "}
            <code className="font-mono text-neutral-400">npm install</code>,{" "}
            <code className="font-mono text-neutral-400">go mod download</code>,
            or setting up virtual environments).
          </p>
          <Callout type="warning" title="Sudo Permissions">
            Depending on your installation path (e.g. installing global binaries
            to{" "}
            <code className="font-mono text-neutral-300">/usr/local/bin</code>),
            AutoDev might request elevated execution permissions (sudo) to
            register links.
          </Callout>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-setup.png"
            alt="autodev setup CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-profile",
      category: "Commands",
      title: "autodev profile",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev profile
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            Bootstrap a pre-defined set of developer tools based on your job
            profile or team role. This is useful for configuring new machines.
          </p>
          {renderCodeBlock("cmd-prof-web", "autodev profile web-dev")}
          <p className="text-neutral-300 mb-4">Supported profiles include:</p>
          <ul className="list-disc pl-5 mb-4 space-y-2 text-neutral-300 font-sans">
            <li>
              <strong>web-dev</strong>: Node.js, pnpm, Docker, VS Code
              integrations
            </li>
            <li>
              <strong>data-science</strong>: Python, Jupyter Notebooks, pandas,
              NumPy, Docker
            </li>
            <li>
              <strong>devops</strong>: Terraform, kubectl, helm, Docker, AWS/GCP
              CLIs
            </li>
            <li>
              <strong>mobile-dev</strong>: Flutter SDK, Android Studio tools,
              CocoaPods
            </li>
          </ul>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-profile.png"
            alt="autodev profile CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-containerize",
      category: "Commands",
      title: "autodev containerize",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev containerize
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            Scaffold container environments for remote DevContainer or Cloud IDE
            development.
          </p>
          {renderCodeBlock("cmd-containerize-run", "autodev containerize")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            AutoDev scans the workspace to identify tech stacks and
            automatically:
          </p>
          <ul className="list-disc pl-5 mb-4 space-y-2 text-neutral-300 font-sans">
            <li>
              Generates a{" "}
              <code className="font-mono text-neutral-400">
                .devcontainer.json
              </code>{" "}
              configured with specific container features and IDE extensions
              matching your project.
            </li>
            <li>
              Creates{" "}
              <code className="font-mono text-neutral-400">
                .vscode/extensions.json
              </code>{" "}
              to recommend necessary extensions.
            </li>
          </ul>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-containerize.png"
            alt="autodev containerize CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-migrate",
      category: "Commands",
      title: "autodev migrate",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev migrate
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            Upgrade legacy configuration files to the standard YAML profile
            schema.
          </p>
          {renderCodeBlock("cmd-migrate-run", "autodev migrate")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            Converts deprecated{" "}
            <code className="font-mono text-neutral-400">.autodev.json</code>{" "}
            and global JSON config configurations into YAML, creating a{" "}
            <code className="font-mono text-neutral-400">.bak</code> file as
            backup.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-migrate.png"
            alt="autodev migrate CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-create",
      category: "Commands",
      title: "autodev create",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev create
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              create
            </code>{" "}
            command generates a new pre-configured boilerplate project matching
            standard developer profiles with all build configurations, linters,
            layout conventions, and git hooks already in place.
          </p>
          {renderCodeBlock(
            "cmd-create-run",
            "autodev create [template] [project-name]",
          )}
          <p className="text-neutral-300 mb-4">Supported templates include:</p>
          <ul className="list-disc pl-5 mb-4 space-y-2 text-neutral-300 font-sans">
            <li>
              <strong>react-ts</strong> (or <strong>react</strong>): Vite +
              React + TypeScript + Tailwind CSS
            </li>
            <li>
              <strong>nextjs</strong> (or <strong>next</strong>): Next.js App
              Router + TypeScript + Tailwind + Docker
            </li>
            <li>
              <strong>ai-chatbot</strong> (or <strong>ai-agent</strong>): React
              chatbot UI + Express.js backend + Google Gemini 2.5 Flash SDK
            </li>
            <li>
              <strong>mern-stack</strong> (or <strong>mern</strong>): Mongo +
              Express + React + Node.js multi-container setup via Docker Compose
            </li>
            <li>
              <strong>flutter-app</strong> (or <strong>flutter</strong>):
              Cross-platform Flutter App structured with Clean Architecture
            </li>
          </ul>
          <Callout type="tip" title="Setup Time Comparison">
            Generating a template with{" "}
            <code className="text-[#00FF87] font-mono font-bold bg-[#111] px-1.5 py-0.5 rounded">
              autodev create
            </code>{" "}
            is up to <strong>99% faster</strong> than prompts-based AI
            configuration, reducing typical API costs by{" "}
            <strong>80% or more</strong>.
          </Callout>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-create.png"
            alt="autodev create CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-clone",
      category: "Commands",
      title: "autodev clone",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev clone
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              clone
            </code>{" "}
            command clones a Git repository, runs a deep stack technology scan,
            and automatically installs all missing dependencies, runtimes, and
            compiler toolchains.
          </p>
          {renderCodeBlock(
            "cmd-clone-run",
            "autodev clone <repository-url> [target-directory]",
          )}
          <p className="text-neutral-300 leading-relaxed mb-4">
            This streamlines developer onboarding to a single command. By
            pairing cloning with automated environment detection and
            installation, a newly cloned repository is immediately ready to run
            and build.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-clone.png"
            alt="autodev clone CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-install",
      category: "Commands",
      title: "autodev install",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev install
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              install
            </code>{" "}
            command installs a specific developer tool, compiler runtime, or SDK
            directly to your system path.
          </p>
          {renderCodeBlock("cmd-install-run", "autodev install <tool>")}
          <p className="text-neutral-300 mb-4">
            Supported tools and runtimes include:
          </p>
          <div className="grid grid-cols-2 gap-2 text-sm text-neutral-400 font-mono mb-4 max-w-md">
            <div>• nodejs</div>
            <div>• go</div>
            <div>• python</div>
            <div>• rust</div>
            <div>• docker</div>
            <div>• bun</div>
            <div>• pnpm</div>
            <div>• java</div>
            <div>• terraform</div>
            <div>• kubectl</div>
            <div>• php</div>
            <div>• ruby</div>
          </div>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-install.png"
            alt="autodev install CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "agent-orchestrator",
      category: "AI Agents",
      title: "Agent Orchestrator (v0.5.0)",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            Agent Orchestrator
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            AutoDev 0.5.0 turns the CLI into a{" "}
            <strong className="text-white">
              unified orchestrator for AI coding agents
            </strong>
            . Install, launch, and manage OpenCode, Claude Code, Codex, Gemini
            CLI, and Aider from a single terminal.
          </p>
          {renderCodeBlock("agent-stack", "autodev tools install opencode claude codex gemini aider")}
          <Callout type="tip" title="One-command install">
            <code className="text-[#00FF87] font-mono font-bold bg-[#111] px-1.5 py-0.5 rounded">
              autodev setup
            </code>{" "}
            detects your machine and the{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              ai-agents
            </code>{" "}
            profile installs the full recommended agent stack.
          </Callout>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Manage agents
          </h2>
          {renderCodeBlock("cmd-tools", "autodev tools list\nautodev tools install <agent>\nautodev tools remove <agent>")}
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Launch an agent
          </h2>
          <p className="text-neutral-300 leading-relaxed mb-4">
            Agents run inside a managed pseudo-terminal — arrow keys, colors and
            window resizing work exactly like a real shell.
          </p>
          {renderCodeBlock("cmd-run", "autodev run opencode\nautodev run codex exec \"fix the auth bug\"")}
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Sessions
          </h2>
          {renderCodeBlock("cmd-session", "autodev session\nautodev session new opencode\nautodev session list\nautodev session stop opencode")}
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Task routing
          </h2>
          <p className="text-neutral-300 leading-relaxed mb-4">
            AutoDev classifies your request and launches the best-fit agent:
          </p>
          {renderCodeBlock("cmd-agent", "autodev agent \"fix the authentication bug\"        # → OpenCode (coding)\nautodev agent \"research the postgres release\"      # → Gemini (research)\nautodev agent \"review this pull request\"          # → Codex (review)\nautodev agent \"harden this codebase\"              # → AutoDevs security")}
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            The full workflow
          </h2>
          {renderCodeBlock("agent-flow", "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash\nautodev setup              # scan machine + install stack\nautodev tools list         # verify agents\nautodev agent \"ship my feature\"  # route task to the best agent")}
        </>
      ),
    },
    {
      id: "cmd-skills",
      category: "Commands",
      title: "autodev skills",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev skills
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              skills
            </code>{" "}
            command analyzes your developer profiles and commits in local git
            history to map out and output a personalized skills progression and
            learning path directly in the terminal (TUI).
          </p>
          {renderCodeBlock("cmd-skills-run", "autodev skills")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It integrates with skills.sh profiles and exports a telemetry
            summary `.autodev-skills.md` to feed relevant skills and framework
            contexts directly to AI Coding Assistants, improving accuracy.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-skills.png"
            alt="autodev skills CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-mcp",
      category: "Commands",
      title: "autodev mcp",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev mcp
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              mcp
            </code>{" "}
            command launches the AutoDev Model Context Protocol (MCP) server.
          </p>
          {renderCodeBlock("cmd-mcp-run", "autodev mcp")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            This server exposes environment diagnostics, security auditing, and
            tool installation commands as MCP tools that AI assistants (like
            Claude, Cursor, Windsurf, Roo-Cline) can interact with directly and
            securely.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-mcp.png"
            alt="autodev mcp CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-benchmark",
      category: "Commands",
      title: "autodev benchmark",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev benchmark
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              benchmark
            </code>{" "}
            command measures and prints AI efficiency comparisons comparing
            standard prompting vs AutoDev&apos;s telemetry integration.
          </p>
          {renderCodeBlock("cmd-benchmark-run", "autodev benchmark")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It demonstrates the token usage, response latency, and cost savings
            achieved when using AutoDev metadata rules files to keep context
            windows small and compact.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-benchmark.png"
            alt="autodev benchmark CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-report",
      category: "Commands",
      title: "autodev report",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev report
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              report
            </code>{" "}
            command generates a comprehensive configuration report of your
            developer environment in HTML, PDF, or JSON format.
          </p>
          {renderCodeBlock("cmd-report-run", "autodev report")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            This is extremely useful for compliance audits, team alignment, or
            archiving developer environment requirements alongside your
            codebase.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-report.png"
            alt="autodev report CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-github",
      category: "Commands",
      title: "autodev github",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev github
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              github
            </code>{" "}
            command scans all public repositories belonging to a specified
            GitHub user or organization to map out the collective technology
            footprint.
          </p>
          {renderCodeBlock("cmd-github-run", "autodev github <username>")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It provides valuable aggregate statistics on primary languages,
            framework prevalence, and tooling usage across their entire public
            code catalog.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-github.png"
            alt="autodev github CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-exec",
      category: "Commands",
      title: "autodev exec",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev exec
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              exec
            </code>{" "}
            command runs a specific command directly inside a configured AutoDev
            virtual sandboxed or path-resolved environment.
          </p>
          {renderCodeBlock("cmd-exec-run", "autodev exec <command>")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It dynamically maps local path runtimes, environment variables, and
            configurations on the fly to ensure commands execute with the
            correct tool versions without polluting your host shell environment.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-exec.png"
            alt="autodev exec CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "cmd-prompts",
      category: "Commands",
      title: "autodev prompts",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            autodev prompts
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The{" "}
            <code className="text-[#C8F135] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">
              prompts
            </code>{" "}
            command manages the prompt capture engine. It lets you view, track, and replay prompts typed into your local AI coding sessions.
          </p>
          {renderCodeBlock("cmd-prompts-run", "autodev prompts")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            By default, running this command opens the master <code className="font-mono text-neutral-400">prompts.md</code> file in your terminal pager (less), displaying a complete history of all captured AI coding prompts.
          </p>
          <Callout type="tip" title="Today's Session Summary">
            Run <code className="text-[#00FF87] font-mono font-bold bg-[#111] px-1.5 py-0.5 rounded">autodev prompts --today</code> to print a timeline of today&apos;s active captured prompts directly in the console.
          </Callout>

          <h2 className="text-2xl font-bold text-[#C8F135] mt-8 mb-4">
            Available Subcommands
          </h2>
          
          <div className="space-y-6">
            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                autodev prompts chat (or autodev chat)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Starts an interactive AI coding session that automatically records all developer prompts, AI responses, created files, and commands executed. Integrates directly with real Gemini models when <code className="font-mono text-neutral-300">GEMINI_API_KEY</code> is set.
              </p>
              {renderCodeBlock("cmd-chat-run", "autodev prompts chat")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                autodev prompts capture (or autodev capture)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Intercepts and captures prompts entered into other terminal-based AI assistant CLI tools (e.g. Copilot, Claude, Gemini CLIs) and automatically appends them to your local project logs.
              </p>
              {renderCodeBlock("cmd-capture-run", "autodev prompts capture <command> [args...]")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                autodev prompts daemon (or autodev daemon)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Starts a background monitoring daemon that automatically records command-line prompts and activity from active editor/AI assistant sessions without needing to manually run wrapper commands.
              </p>
              {renderCodeBlock("cmd-daemon-run", "autodev prompts daemon")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                autodev prompts replay (or autodev replay)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Lists all recently captured prompts from the latest session, allowing you to select and re-execute a prompt (along with any suggested file creations or shell commands) against your current codebase.
              </p>
              {renderCodeBlock("cmd-replay-run", "autodev prompts replay")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                autodev prompts export-prompts (or autodev export-prompts)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Consolidates and exports all captured prompt logs into a single Markdown or JSON file.
              </p>
              {renderCodeBlock("cmd-export-run-cli", "autodev prompts export-prompts -o backup.md -f markdown")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">
                autodev prompts sync (or autodev sync)
              </h3>
              <p className="text-neutral-400 text-sm mb-2">
                Synchronizes any offline-queued prompt events (such as prompts written while offline) to the remote DevMentor API.
              </p>
              {renderCodeBlock("cmd-sync-run", "autodev prompts sync")}
            </div>
          </div>

          <h2 className="text-xl font-bold text-white mt-8 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-prompts.png"
            alt="autodev prompts CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
    {
      id: "lockfiles",
      category: "Advanced",
      title: "Reproducible Lockfiles",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            Reproducible Lockfiles
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            To lock environment tooling configurations across your engineering
            team, you can generate a lockfile config.
          </p>
          <p className="text-neutral-300 mb-2">
            Export your current environment config using:
          </p>
          {renderCodeBlock("adv-export", "autodev export")}
          <p className="text-neutral-300 leading-relaxed mt-4">
            This creates a{" "}
            <code className="font-mono text-[#C8F135] bg-[#111] px-1 py-0.5 rounded">
              .autodev.yaml
            </code>{" "}
            file at your project root. When other developers run{" "}
            <code className="font-mono text-neutral-300">autodev setup</code>,
            the CLI reads this file to guarantee identical compiler and tool
            versions across everyone&apos;s workstation.
          </p>
        </>
      ),
    },
    {
      id: "docker-usage",
      category: "Advanced",
      title: "Docker Integration",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">
            Running in Docker
          </h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            If you want to keep your host machine clean, you can run AutoDev
            scanner inside a Docker sandbox.
          </p>
          {renderCodeBlock(
            "adv-dock",
            "docker run --rm -v $(pwd):/workspace ghcr.io/heetmehta18/autodev scan",
          )}
          <p className="text-neutral-300 leading-relaxed mt-4">
            This will mount your local repository into the container workspace,
            run the dependency scanners, and print the resulting plan without
            altering files on your host OS.
          </p>
          <h2 className="text-xl font-bold text-white mt-6 mb-2">
            Terminal Output
          </h2>
          <Image
            src="/screenshot-docker-integration.png"
            alt="autodev docker run CLI Output"
            className="my-4 border-2 border-[#2A2A2A] rounded-md max-w-full shadow-lg"
           width={800} height={450}/>
        </>
      ),
    },
  ];

  // Group sections by category
  const categories = ["Getting Started", "Commands", "Advanced"];

  const currentSection =
    sections.find((s) => s.id === activeSection) || sections[0];

  return (
    <main className="min-h-screen bg-black flex flex-col">
      <Navbar />

      {/* Main Documentation Wrapper */}
      <div className="flex-1 flex pt-16 max-w-7xl mx-auto w-full px-6">
        {/* Left Sidebar (Vercel Style) */}
        <aside className="w-64 shrink-0 border-r border-[#1F1F1F] pr-8 pt-10 hidden md:block sticky top-16 h-[calc(100vh-4rem)] overflow-y-auto">
          <div className="space-y-8">
            {categories.map((cat) => (
              <div key={cat}>
                <h4 className="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-3 font-mono">
                  {cat}
                </h4>
                <ul className="space-y-1.5">
                  {sections
                    .filter((s) => s.category === cat)
                    .map((s) => {
                      const isActive = s.id === activeSection;
                      return (
                        <li key={s.id}>
                          <button
                            onClick={() => setActiveSection(s.id)}
                            className={`w-full text-left px-3 py-1.5 rounded text-sm transition-all font-mono font-medium ${
                              isActive
                                ? "bg-[#111] text-[#C8F135] border-l-2 border-[#C8F135] pl-4"
                                : "text-neutral-400 hover:text-white hover:bg-neutral-900/50"
                            }`}
                          >
                            {s.title}
                          </button>
                        </li>
                      );
                    })}
                </ul>
              </div>
            ))}
          </div>
        </aside>

        {/* Right Content Area */}
        <article className="flex-1 min-w-0 pl-0 md:pl-12 pt-10 pb-20">
          <div className="max-w-3xl">
            {/* Mobile Navigation Selector */}
            <div className="md:hidden mb-8 border-2 border-[#2A2A2A] bg-[#111] p-3 rounded flex flex-col gap-2">
              <label className="text-xs font-bold text-[#C8F135] uppercase font-mono tracking-wider">
                Select Page:
              </label>
              <select
                value={activeSection}
                onChange={(e) => setActiveSection(e.target.value)}
                className="bg-black text-white border border-[#2A2A2A] p-2 font-mono text-sm outline-none rounded"
              >
                {sections.map((s) => (
                  <option key={s.id} value={s.id}>
                    {s.category} → {s.title}
                  </option>
                ))}
              </select>
            </div>

            {/* Render Active Content */}
            <div className="prose prose-invert prose-yellow max-w-none">
              {currentSection.content}
            </div>

            {/* Navigation Footer Buttons */}
            <div className="mt-16 pt-8 border-t border-[#1F1F1F] flex justify-between gap-4">
              {(() => {
                const curIdx = sections.findIndex(
                  (s) => s.id === activeSection,
                );
                const prevSec = curIdx > 0 ? sections[curIdx - 1] : null;
                const nextSec =
                  curIdx < sections.length - 1 ? sections[curIdx + 1] : null;

                return (
                  <>
                    {prevSec ? (
                      <button
                        onClick={() => setActiveSection(prevSec.id)}
                        className="nb-btn-outline px-4 py-2.5 text-xs text-left flex flex-col font-mono"
                      >
                        <span className="text-neutral-500 font-sans uppercase font-bold text-[10px] tracking-widest">
                          Previous
                        </span>
                        <span>← {prevSec.title}</span>
                      </button>
                    ) : (
                      <div />
                    )}

                    {nextSec ? (
                      <button
                        onClick={() => setActiveSection(nextSec.id)}
                        className="nb-btn px-4 py-2.5 text-xs text-right flex flex-col font-mono"
                      >
                        <span className="text-black/50 font-sans uppercase font-bold text-[10px] tracking-widest">
                          Next
                        </span>
                        <span>{nextSec.title} →</span>
                      </button>
                    ) : (
                      <div />
                    )}
                  </>
                );
              })()}
            </div>
          </div>
        </article>
      </div>

      <Footer />
    </main>
  );
}
