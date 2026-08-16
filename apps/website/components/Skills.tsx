"use client";
import { useState, useRef } from "react";
import { motion } from "framer-motion";
import { Copy, Check } from "lucide-react";

interface SkillDetection {
  name: string;
  icon: string;
  level: "Detected" | "Proficient" | "Expert";
  confidence: number;
  files: number;
  commits: number;
}

interface RoadmapStep {
  from: string;
  to: string;
  reason: string;
  priority: "High" | "Medium" | "Low";
}

// Simulated scan output lines
const scanLines = [
  { text: "$ autodev skills --deep --ai", color: "#00FF87", delay: 0 },
  { text: "", color: "#888", delay: 200 },
  { text: "  ⚡ AutoDev Skills Engine v0.6.0", color: "#C8F135", delay: 400 },
  {
    text: "  Powered by skills.sh — https://skills.sh",
    color: "#888",
    delay: 600,
  },
  { text: "", color: "#888", delay: 800 },
  {
    text: "  [scan] Indexing git log (2,847 commits across 14 repos)...",
    color: "#2563eb",
    delay: 1000,
  },
  {
    text: "  [scan] Parsing package.json, go.mod, Cargo.toml, requirements.txt...",
    color: "#2563eb",
    delay: 1400,
  },
  {
    text: "  [scan] Analyzing Dockerfile, docker-compose.yml, .github/workflows...",
    color: "#2563eb",
    delay: 1800,
  },
  {
    text: "  [scan] Resolving framework patterns: Next.js, FastAPI, Spring Boot...",
    color: "#2563eb",
    delay: 2200,
  },
  {
    text: "  [scan] Computing skill confidence from code frequency + recency...",
    color: "#2563eb",
    delay: 2600,
  },
  { text: "", color: "#888", delay: 3000 },
  {
    text: "  ✓ Skill matrix generated: 12 technologies across 6 categories",
    color: "#00FF87",
    delay: 3200,
  },
  {
    text: "  ✓ Personalized roadmap: 4 upgrade paths identified",
    color: "#00FF87",
    delay: 3400,
  },
  { text: "  ✓ Synced to skills.sh profile", color: "#00FF87", delay: 3600 },
  { text: "", color: "#888", delay: 3800 },
  {
    text: "  CURRENT SKILLS (from repo analysis):",
    color: "#C8F135",
    delay: 4000,
  },
  {
    text: "  ┌─────────────────────┬───────────┬────────────┬─────────┐",
    color: "#333",
    delay: 4200,
  },
  {
    text: "  │ Technology          │ Level     │ Confidence │ Files   │",
    color: "#888",
    delay: 4300,
  },
  {
    text: "  ├─────────────────────┼───────────┼────────────┼─────────┤",
    color: "#333",
    delay: 4400,
  },
  {
    text: "  │ 🐍 Python           │ Expert    │ 94%        │ 247     │",
    color: "#00FF87",
    delay: 4500,
  },
  {
    text: "  │ ⚛️  React/TypeScript │ Expert    │ 91%        │ 189     │",
    color: "#00FF87",
    delay: 4600,
  },
  {
    text: "  │ 🟢 Node.js          │ Proficient│ 82%        │ 134     │",
    color: "#2563eb",
    delay: 4700,
  },
  {
    text: "  │ 🐳 Docker           │ Proficient│ 76%        │ 28      │",
    color: "#2563eb",
    delay: 4800,
  },
  {
    text: "  │ 🐘 PostgreSQL       │ Detected  │ 58%        │ 12      │",
    color: "#C8F135",
    delay: 4900,
  },
  {
    text: "  │ 🔵 Go               │ Detected  │ 45%        │ 8       │",
    color: "#C8F135",
    delay: 5000,
  },
  {
    text: "  └─────────────────────┴───────────┴────────────┴─────────┘",
    color: "#333",
    delay: 5100,
  },
  { text: "", color: "#888", delay: 5200 },
  { text: "  RECOMMENDED NEXT STEPS:", color: "#C8F135", delay: 5400 },
  {
    text: "  → Learn Kubernetes (complements your Docker proficiency)",
    color: "#F0F0F0",
    delay: 5600,
  },
  {
    text: "  → Deepen PostgreSQL (upgrade from Detected → Proficient)",
    color: "#F0F0F0",
    delay: 5800,
  },
  {
    text: "  → Explore Terraform (natural DevOps progression)",
    color: "#F0F0F0",
    delay: 6000,
  },
  {
    text: "  → Invest in Go (build on existing foundations)",
    color: "#F0F0F0",
    delay: 6200,
  },
  { text: "", color: "#888", delay: 6400 },
  {
    text: "  [AI-POWERED INSIGHTS — Perplexity]",
    color: "#C8F135",
    delay: 6600,
  },
  {
    text: "  💡 Your React+TypeScript expertise positions you for senior frontend roles",
    color: "#F0F0F0",
    delay: 6800,
  },
  {
    text: "  💡 Docker → Kubernetes is the highest-ROI upgrade for your stack",
    color: "#F0F0F0",
    delay: 7000,
  },
  {
    text: "  💡 Full-stack Go+React profiles are in top 5% demand on job markets",
    color: "#F0F0F0",
    delay: 7200,
  },
  {
    text: "  💡 Consider adding CI/CD — it unlocks DevOps career trajectories",
    color: "#F0F0F0",
    delay: 7400,
  },
  { text: "", color: "#888", delay: 7600 },
  {
    text: "  Run 'autodev skills --export md' to save as Markdown.",
    color: "#888",
    delay: 7800,
  },
];

const detectedSkills: SkillDetection[] = [
  {
    name: "Python",
    icon: "🐍",
    level: "Expert",
    confidence: 94,
    files: 247,
    commits: 1420,
  },
  {
    name: "React / TypeScript",
    icon: "⚛️",
    level: "Expert",
    confidence: 91,
    files: 189,
    commits: 1180,
  },
  {
    name: "Node.js",
    icon: "🟢",
    level: "Proficient",
    confidence: 82,
    files: 134,
    commits: 640,
  },
  {
    name: "Docker",
    icon: "🐳",
    level: "Proficient",
    confidence: 76,
    files: 28,
    commits: 210,
  },
  {
    name: "PostgreSQL",
    icon: "🐘",
    level: "Detected",
    confidence: 58,
    files: 12,
    commits: 85,
  },
  {
    name: "Go",
    icon: "🔵",
    level: "Detected",
    confidence: 45,
    files: 8,
    commits: 42,
  },
];

const roadmapSteps: RoadmapStep[] = [
  {
    from: "Docker",
    to: "Kubernetes",
    reason: "Natural container orchestration progression",
    priority: "High",
  },
  {
    from: "PostgreSQL",
    to: "Advanced SQL",
    reason: "Upgrade from Detected → Proficient",
    priority: "High",
  },
  {
    from: "Node.js",
    to: "Terraform",
    reason: "DevOps pipeline automation",
    priority: "Medium",
  },
  {
    from: "Go",
    to: "Go Microservices",
    reason: "Build on existing Go foundations",
    priority: "Low",
  },
];

const cliCommands = [
  { cmd: "autodev skills", desc: "Scan repos and generate your skill matrix" },
  {
    cmd: "autodev skills --deep",
    desc: "Deep analysis with git history + commit frequency",
  },
  {
    cmd: "autodev skills --deep --ai",
    desc: "AI-powered insights via Perplexity API",
  },
  {
    cmd: "autodev skills --export md",
    desc: "Export roadmap as Markdown / JSON / HTML",
  },
  { cmd: "autodev skills --sync", desc: "Sync profile to skills.sh" },
];

export default function Skills() {
  const [visibleLines, setVisibleLines] = useState(0);
  const [isPlaying, setIsPlaying] = useState(false);
  const [hasPlayed, setHasPlayed] = useState(false);
  const [activeTab, setActiveTab] = useState<"terminal" | "matrix" | "roadmap">(
    "terminal",
  );
  const terminalRef = useRef<HTMLDivElement>(null);
  const [copiedCmd, setCopiedCmd] = useState<string | null>(null);
  const [copiedRoadmapExport, setCopiedRoadmapExport] = useState(false);

  const startDemo = () => {
    if (isPlaying) return;
    setIsPlaying(true);
    setVisibleLines(0);
    setActiveTab("terminal");

    scanLines.forEach((line, i) => {
      setTimeout(() => {
        setVisibleLines(i + 1);
        // Auto-scroll terminal
        if (terminalRef.current) {
          terminalRef.current.scrollTop = terminalRef.current.scrollHeight;
        }
        if (i === scanLines.length - 1) {
          setIsPlaying(false);
          setHasPlayed(true);
        }
      }, line.delay);
    });
  };

  const getLevelColor = (level: string) => {
    switch (level) {
      case "Expert":
        return "#00FF87";
      case "Proficient":
        return "#2563eb";
      case "Detected":
        return "#C8F135";
      default:
        return "#A3A3A3";
    }
  };

  const getPriorityBadge = (p: string) => {
    switch (p) {
      case "High":
        return "border-[#dc2626] text-[#dc2626]";
      case "Medium":
        return "border-[#C8F135] text-[#C8F135]";
      case "Low":
        return "border-[#2563eb] text-[#2563eb]";
      default:
        return "border-neutral-400 text-neutral-400";
    }
  };

  return (
    <section
      id="skills"
      className="py-24 px-6 bg-[#0D0D0D] border-y-2 border-[#2A2A2A] scroll-mt-24"
    >
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <motion.div
          className="mb-16"
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
        >
          <span className="inline-flex items-center gap-1.5 text-xs text-[#C8F135] font-bold uppercase tracking-widest">
            <span className="w-1.5 h-1.5 bg-[#C8F135] rounded-full animate-pulse" /> Skills.sh Integration
          </span>
          <h2 className="text-5xl font-black text-white mt-2 mb-4">
            YOUR LEARNING <span className="text-gradient-yellow">ROADMAP</span>
          </h2>
          <p className="text-neutral-400 max-w-2xl">
            AutoDev scans your repositories — git history, package manifests,
            Dockerfiles, CI configs — and builds a{" "}
            <strong className="text-white">
              confidence-scored skill matrix
            </strong>{" "}
            with personalized upgrade paths. Powered by{" "}
            <a
              href="https://skills.sh"
              className="text-[#C8F135] underline hover:text-white transition-colors"
            >
              skills.sh
            </a>
            .
          </p>
        </motion.div>

        <div className="grid grid-cols-1 lg:grid-cols-[1fr_340px] gap-8 items-start">
          {/* Main Panel — Tabbed Interface */}
          <div className="border-2 border-[#2A2A2A] bg-[#0A0A0A] overflow-hidden">
            {/* Tab Bar */}
            <div className="border-b-2 border-[#2A2A2A] bg-[#111] flex items-center justify-between">
              <div className="flex">
                {[
                  {
                    id: "terminal" as const,
                    label: "Live Terminal Demo",
                    icon: "▶",
                  },
                  { id: "matrix" as const, label: "Skill Matrix", icon: "◆" },
                  { id: "roadmap" as const, label: "Upgrade Paths", icon: "→" },
                ].map((tab) => (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id)}
                    className={`px-5 py-3 text-xs font-bold uppercase tracking-wider border-b-2 transition-colors
                      ${
                        activeTab === tab.id
                          ? "text-[#C8F135] border-[#C8F135] bg-[#C8F13508]"
                          : "text-neutral-400 border-transparent hover:text-white"
                      }`}
                  >
                    <span className="mr-1.5">{tab.icon}</span>
                    {tab.label}
                  </button>
                ))}
              </div>

              {activeTab === "terminal" && (
                <button
                  onClick={startDemo}
                  disabled={isPlaying}
                  className={`mr-4 px-4 py-1.5 text-xs font-bold border-2 uppercase tracking-wider transition-all
                    ${
                      isPlaying
                        ? "border-[#888] text-neutral-400 cursor-not-allowed"
                        : "border-[#C8F135] bg-[#C8F135] text-black hover:bg-transparent hover:text-[#C8F135]"
                    }`}
                >
                  {isPlaying ? "Running..." : hasPlayed ? "Replay" : "Run Demo"}
                </button>
              )}
            </div>

            {/* Tab Content */}
            <div className="min-h-[480px]">
              {/* Terminal Tab */}
              {activeTab === "terminal" && (
                <div
                  ref={terminalRef}
                  className="p-6 font-mono text-sm leading-6 overflow-y-auto max-h-[520px] select-text selection:bg-[#C8F135] selection:text-black"
                >
                  {visibleLines === 0 && !isPlaying && (
                    <div className="flex flex-col items-center justify-center h-[400px] text-center gap-4">
                      <div className="text-6xl mb-2">⚡</div>
                      <p className="text-neutral-400 text-sm max-w-md">
                        Click{" "}
                        <strong className="text-[#C8F135]">Run Demo</strong> to
                        see AutoDev&apos;s Skills Engine scan a developer&apos;s
                        repositories and generate a live skill matrix.
                      </p>
                    </div>
                  )}
                  {scanLines.slice(0, visibleLines).map((line, i) => (
                    <motion.div
                      key={i}
                      initial={{ opacity: 0 }}
                      animate={{ opacity: 1 }}
                      transition={{ duration: 0.1 }}
                      style={{ color: line.color }}
                    >
                      {line.text || "\u00A0"}
                    </motion.div>
                  ))}
                  {isPlaying && (
                    <span className="inline-block w-2 h-4 bg-[#C8F135] animate-pulse mt-1" />
                  )}
                </div>
              )}

              {/* Skill Matrix Tab */}
              {activeTab === "matrix" && (
                <div className="p-6 space-y-4">
                  <div className="flex items-center justify-between mb-2">
                    <p className="text-xs text-neutral-400">
                      Confidence scores computed from code frequency, recency,
                      and repository diversity.
                    </p>
                    <div className="flex gap-3 text-[10px] font-mono">
                      <span className="text-[#00FF87]">● Expert</span>
                      <span className="text-[#2563eb]">● Proficient</span>
                      <span className="text-[#C8F135]">● Detected</span>
                    </div>
                  </div>

                  {detectedSkills.map((skill) => (
                    <div
                      key={skill.name}
                      className="border border-[#2A2A2A] bg-[#111] p-4 hover:border-[#444] transition-colors"
                    >
                      <div className="flex items-center justify-between mb-3">
                        <div className="flex items-center gap-3">
                          <span className="text-2xl w-10 h-10 bg-black border border-[#2A2A2A] flex items-center justify-center">
                            {skill.icon}
                          </span>
                          <div>
                            <h3 className="font-bold text-white text-sm">
                              {skill.name}
                            </h3>
                            <span
                              className="text-[10px] font-mono font-bold uppercase tracking-wider"
                              style={{ color: getLevelColor(skill.level) }}
                            >
                              {skill.level}
                            </span>
                          </div>
                        </div>
                        <div className="text-right">
                          <div
                            className="text-lg font-black"
                            style={{ color: getLevelColor(skill.level) }}
                          >
                            {skill.confidence}%
                          </div>
                          <div className="text-[10px] text-neutral-400 font-mono">
                            {skill.files} files · {skill.commits} commits
                          </div>
                        </div>
                      </div>

                      {/* Confidence bar */}
                      <div className="w-full bg-black border border-[#2A2A2A] h-2">
                        <div
                          className="h-full transition-all duration-700"
                          style={{
                            width: `${skill.confidence}%`,
                            backgroundColor: getLevelColor(skill.level),
                          }}
                        />
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* Roadmap Tab */}
              {activeTab === "roadmap" && (
                <div className="p-6 space-y-1">
                  <p className="text-xs text-neutral-400 mb-6">
                    Upgrade paths ranked by impact. Recommendations based on
                    your existing skill adjacencies.
                  </p>

                  {roadmapSteps.map((step, i) => (
                    <motion.div
                      key={i}
                      initial={{ opacity: 0, x: -10 }}
                      animate={{ opacity: 1, x: 0 }}
                      transition={{ delay: i * 0.1 }}
                      className="border border-[#2A2A2A] bg-[#111] p-5 hover:border-[#444] transition-colors flex items-center justify-between gap-4"
                    >
                      <div className="flex items-center gap-5 flex-1">
                        <div className="flex items-center gap-2 min-w-0">
                          <span className="text-sm font-bold text-neutral-400 whitespace-nowrap">
                            {step.from}
                          </span>
                          <span className="text-[#C8F135] text-lg font-bold">
                            →
                          </span>
                          <span className="text-sm font-bold text-white whitespace-nowrap">
                            {step.to}
                          </span>
                        </div>
                        <span className="text-xs text-neutral-400 hidden md:block">
                          {step.reason}
                        </span>
                      </div>
                      <span
                        className={`text-[9px] font-mono px-2 py-0.5 border font-bold uppercase shrink-0 ${getPriorityBadge(step.priority)}`}
                      >
                        {step.priority}
                      </span>
                    </motion.div>
                  ))}

                  <div className="mt-6 border-t border-[#2A2A2A] pt-5">
                    <div className="terminal relative">
                      <div className="terminal-bar py-1 px-3 flex justify-between items-center pr-3">
                        <div className="flex gap-1">
                          <span className="w-1.5 h-1.5 rounded-full bg-[#FF5F56]" />
                          <span className="w-1.5 h-1.5 rounded-full bg-[#FFBD2E]" />
                          <span className="w-1.5 h-1.5 rounded-full bg-[#27C93F]" />
                        </div>
                        <button
                          onClick={() => {
                            navigator.clipboard.writeText(
                              "autodev skills --export md > ~/dev-roadmap.md",
                            );
                            setCopiedRoadmapExport(true);
                            setTimeout(
                              () => setCopiedRoadmapExport(false),
                              1800,
                            );
                          }}
                          className="text-neutral-400 hover:text-[#C8F135] transition-colors p-1 flex items-center gap-1 rounded bg-[#1e1e1e] border border-[#2a2a2a] cursor-pointer"
                          title="Copy command"
                        >
                          {copiedRoadmapExport ? (
                            <>
                              <Check className="w-3 h-3 text-[#00FF87]" />
                              <span className="text-[9px] text-[#00FF87] font-mono">
                                Copied!
                              </span>
                            </>
                          ) : (
                            <>
                              <Copy className="w-3 h-3" />
                              <span className="text-[9px] text-neutral-400 font-mono">
                                Copy
                              </span>
                            </>
                          )}
                        </button>
                      </div>
                      <div className="px-4 py-3 font-mono text-xs text-[#00FF87] bg-black">
                        $ autodev skills --export md &gt; ~/dev-roadmap.md
                        <br />
                        <span className="text-neutral-400">
                          {" "}
                          ✓ Roadmap exported. Share it on skills.sh or add to
                          your README.
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Right Sidebar — CLI Commands Reference */}
          <div className="space-y-6">
            <div className="border-2 border-[#2A2A2A] bg-[#111] p-5">
              <h3 className="font-bold text-white text-sm uppercase tracking-wider mb-4 flex items-center gap-2">
                <span className="text-[#C8F135]">$</span> CLI Commands
              </h3>
              <div className="space-y-3">
                {cliCommands.map((c) => (
                  <div
                    key={c.cmd}
                    className="group cursor-pointer"
                    onClick={() => {
                      navigator.clipboard.writeText(c.cmd);
                      setCopiedCmd(c.cmd);
                      setTimeout(() => setCopiedCmd(null), 1800);
                    }}
                  >
                    <div className="font-mono text-xs text-[#00FF87] bg-black border border-[#2A2A2A] px-3 py-2 group-hover:border-[#444] transition-colors flex justify-between items-center gap-2">
                      <span>{c.cmd}</span>
                      <span className="shrink-0 flex items-center gap-1">
                        {copiedCmd === c.cmd ? (
                          <Check className="w-3.5 h-3.5 text-[#00FF87]" />
                        ) : (
                          <Copy className="w-3.5 h-3.5 text-neutral-400 group-hover:text-[#C8F135] transition-colors" />
                        )}
                      </span>
                    </div>
                    <p className="text-[10px] text-neutral-400 mt-1 pl-1">
                      {c.desc}
                    </p>
                  </div>
                ))}
              </div>
            </div>

            <div className="border-2 border-[#2A2A2A] bg-[#111] p-5">
              <h3 className="font-bold text-white text-sm uppercase tracking-wider mb-3">
                How It Works
              </h3>
              <div className="space-y-3 text-xs text-neutral-400 leading-relaxed">
                {[
                  {
                    step: "1",
                    title: "Scan",
                    desc: "Indexes package manifests, Dockerfiles, CI configs across all local repos",
                  },
                  {
                    step: "2",
                    title: "Score",
                    desc: "Computes confidence from code frequency, git recency, and file diversity",
                  },
                  {
                    step: "3",
                    title: "Map",
                    desc: "Generates personalized upgrade paths based on skill adjacency graphs",
                  },
                  {
                    step: "4",
                    title: "Sync",
                    desc: "Pushes your skill profile to skills.sh for portfolio display",
                  },
                ].map((s) => (
                  <div key={s.step} className="flex gap-3">
                    <span className="w-5 h-5 bg-[#C8F135] text-black font-black text-[10px] flex items-center justify-center shrink-0">
                      {s.step}
                    </span>
                    <div>
                      <span className="font-bold text-white">{s.title}</span>
                      <span className="text-neutral-400"> — {s.desc}</span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <a
              href="https://skills.sh"
              target="_blank"
              rel="noopener noreferrer"
              className="block border-2 border-[#C8F135] bg-[#C8F13508] p-4 text-center hover:bg-[#C8F13515] transition-colors group"
            >
              <div className="text-[#C8F135] font-black text-sm mb-1">
                skills.sh
              </div>
              <p className="text-[10px] text-neutral-400 group-hover:text-[#aaa] transition-colors">
                Developer learning platform. Create your public skill profile →
              </p>
            </a>
          </div>
        </div>
      </div>
    </section>
  );
}
