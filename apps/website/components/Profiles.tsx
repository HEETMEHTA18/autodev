"use client";
import { useState } from "react";
import { motion } from "framer-motion";

const profiles = [
  {
    id: "web-dev",
    icon: "🌐",
    name: "Web Developer",
    desc: "Full-stack web development",
    color: "#4A90E2",
    packages: [
      "git",
      "nodejs",
      "pnpm",
      "react",
      "nextjs",
      "docker",
      "postgresql",
      "redis",
    ],
  },
  {
    id: "ml-engineer",
    icon: "🤖",
    name: "ML Engineer",
    desc: "Machine learning & data science",
    color: "#FF6B6B",
    packages: [
      "git",
      "python",
      "jupyter",
      "pytorch",
      "tensorflow",
      "langchain",
      "docker",
    ],
  },
  {
    id: "flutter-dev",
    icon: "💙",
    name: "Flutter Developer",
    desc: "Cross-platform mobile apps",
    color: "#00B4D8",
    packages: [
      "git",
      "java",
      "dart",
      "flutter",
      "android-sdk",
      "android-studio",
    ],
  },
  {
    id: "devops-engineer",
    icon: "⚙️",
    name: "DevOps Engineer",
    desc: "Infrastructure & CI/CD",
    color: "#00FF87",
    packages: [
      "git",
      "docker",
      "docker-compose",
      "kubernetes",
      "terraform",
      "nginx",
    ],
  },
  {
    id: "backend-dev",
    icon: "🔧",
    name: "Backend Developer",
    desc: "Server-side APIs & services",
    color: "#FFD700",
    packages: [
      "git",
      "nodejs",
      "python",
      "go",
      "docker",
      "postgresql",
      "redis",
    ],
  },
  {
    id: "fullstack-ai",
    icon: "🧠",
    name: "Full-Stack AI Dev",
    desc: "AI-powered full-stack apps",
    color: "#C77DFF",
    packages: [
      "git",
      "nodejs",
      "python",
      "nextjs",
      "fastapi",
      "pytorch",
      "langchain",
      "docker",
      "postgresql",
    ],
  },
];

export default function Profiles() {
  const [active, setActive] = useState<string | null>(null);
  const current = profiles.find((p) => p.id === active);

  return (
    <section
      id="profiles"
      className="py-24 px-6 bg-[#0D0D0D] border-y-2 border-[#2A2A2A]"
    >
      <div className="max-w-7xl mx-auto">
        <motion.div
          className="mb-16"
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
        >
          <span className="inline-flex items-center gap-1.5 text-xs text-[#FFD700] font-bold uppercase tracking-widest">
            <span className="w-1.5 h-1.5 bg-[#FFD700] rounded-full" /> One command
          </span>
          <h2 className="text-5xl font-black text-white mt-2 mb-4">
            DEVELOPER <span className="text-gradient-yellow">PROFILES</span>
          </h2>
          <p className="text-neutral-400 max-w-xl">
            Pick your role. AutoDev installs everything — with smart dependency
            resolution.
          </p>
        </motion.div>

        <motion.div
          className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5"
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: "-50px" }}
          variants={{
            hidden: {},
            visible: { transition: { staggerChildren: 0.08 } },
          }}
        >
          {profiles.map((prof) => (
            <motion.button
              key={prof.id}
              variants={{
                hidden: { opacity: 0, y: 30 },
                visible: { opacity: 1, y: 0, transition: { duration: 0.4 } },
              }}
              onClick={() => setActive(active === prof.id ? null : prof.id)}
              whileTap={{ scale: 0.97 }}
              className={`text-left p-6 border-2 transition-all duration-200 nb-card glow-yellow-hover
                ${active === prof.id ? "border-[#FFD700] shadow-[4px_4px_0_#FFD700]" : ""}
              `}
            >
              <div className="text-4xl mb-4">{prof.icon}</div>
              <h3 className="font-black text-lg text-white mb-1">
                {prof.name}
              </h3>
              <p className="text-sm text-neutral-400 mb-4">{prof.desc}</p>
              <div className="flex flex-wrap gap-1">
                {prof.packages.slice(0, 5).map((pkg) => (
                  <span
                    key={pkg}
                    className="text-[10px] px-2 py-0.5 border border-[#333] text-neutral-400 font-mono"
                  >
                    {pkg}
                  </span>
                ))}
                {prof.packages.length > 5 && (
                  <span className="text-[10px] px-2 py-0.5 border border-[#333] text-[#FFD700] font-mono">
                    +{prof.packages.length - 5} more
                  </span>
                )}
              </div>
            </motion.button>
          ))}
        </motion.div>

        {/* Command preview */}
        {current && (
          <motion.div
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            className="mt-8 terminal"
          >
            <div className="terminal-bar">
              <span className="terminal-dot bg-[#FF5F56]" />
              <span className="terminal-dot bg-[#FFBD2E]" />
              <span className="terminal-dot bg-[#27C93F]" />
              <span className="text-xs text-neutral-400 ml-3">
                {current.name} Profile
              </span>
            </div>
            <div className="px-6 py-5 font-mono text-sm space-y-1 bg-black">
              <div className="text-[#00FF87]">
                $ autodev profile {current.id}
              </div>
              <div className="text-neutral-400 mt-2">
                {current.icon} {current.name}
              </div>
              <div className="text-neutral-400"> {current.desc}</div>
              <div className="mt-3 text-[#4A90E2]">
                Packages ({current.packages.length}):
              </div>
              {current.packages.map((p) => (
                <div key={p} className="text-[#00FF87] pl-3">
                  {" "}
                  ✓ {p}
                </div>
              ))}
            </div>
          </motion.div>
        )}
      </div>
    </section>
  );
}
