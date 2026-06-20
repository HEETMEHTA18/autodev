"use client";
import { useState } from "react";
import { motion } from "framer-motion";
import { Copy, Check, Sparkles } from "lucide-react";
import { trackInstall } from "../utils/analytics";

const methods = [
  {
    label: "NPM",
    desc: "Install globally via npm",
    cmd: "npm install -g @heetmehta18/autodev",
    icon: "📦",
  },
  {
    label: "NPX",
    desc: "Run without local install",
    cmd: "npx @heetmehta18/autodev setup",
    icon: "🚀",
  },
  {
    label: "Shell",
    desc: "Linux & macOS — recommended",
    cmd: "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash",
    icon: "🐚",
  },
  {
    label: "Homebrew",
    desc: "macOS & Linux",
    cmd: "brew install HEETMEHTA18/tap/autodev",
    icon: "🍺",
  },
  {
    label: "Scoop",
    desc: "Windows — add bucket first",
    cmd: "scoop bucket add autodev https://github.com/HEETMEHTA18/scoop-bucket && scoop install autodev",
    icon: "🪣",
  },
  {
    label: "Docker",
    desc: "No local install required",
    cmd: "docker run --rm -v $(pwd):/workspace ghcr.io/heetmehta18/autodev setup",
    icon: "🐳",
  },
];

export default function InstallMethods() {
  const [copied, setCopied] = useState<string | null>(null);

  const copy = (cmd: string, label: string) => {
    navigator.clipboard.writeText(cmd);
    setCopied(cmd);
    trackInstall(label);
    setTimeout(() => setCopied(null), 1800);
  };

  return (
    <section id="install" className="py-24 px-6 max-w-7xl mx-auto">
      <motion.div
        className="mb-16 text-center"
        initial={{ opacity: 0, y: 30 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
      >
        <span className="inline-flex items-center gap-1.5 text-xs text-[#FFD700] font-bold uppercase tracking-widest">
          <Sparkles className="w-3 h-3" /> Install anywhere
        </span>
        <h2 className="text-5xl font-black text-white mt-2 mb-4">
          GET STARTED IN <span className="text-gradient-yellow">SECONDS</span>
        </h2>
        <p className="text-neutral-300">
          Pick your preferred installation method.
        </p>
      </motion.div>

      <motion.div
        className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4"
        initial="hidden"
        whileInView="visible"
        viewport={{ once: true, margin: "-50px" }}
        variants={{
          hidden: {},
          visible: { transition: { staggerChildren: 0.08 } },
        }}
      >
        {methods.map((m) => (
          <motion.div
            key={m.label}
            variants={{
              hidden: { opacity: 0, y: 30 },
              visible: { opacity: 1, y: 0, transition: { duration: 0.4 } },
            }}
            className="nb-card p-5 cursor-pointer group glow-yellow-hover"
            onClick={() => copy(m.cmd, m.label)}
          >
            <div className="flex items-center justify-between mb-1">
              <div className="flex items-center gap-2">
                <span className="text-xl">{m.icon}</span>
                <span className="font-bold text-white text-sm">{m.label}</span>
              </div>
              <div className="flex items-center gap-1 text-xs">
                {copied === m.cmd ? (
                  <>
                    <Check className="w-3.5 h-3.5 text-[#00FF87]" />
                    <span className="text-[10px] text-[#00FF87] font-mono">
                      Copied!
                    </span>
                  </>
                ) : (
                  <>
                    <Copy className="w-3.5 h-3.5 text-neutral-400 group-hover:text-[#FFD700] transition-colors" />
                    <span className="text-[10px] text-neutral-400 group-hover:text-[#FFD700] font-mono transition-colors">
                      Copy
                    </span>
                  </>
                )}
              </div>
            </div>
            <p className="text-[10px] text-neutral-400 font-sans mb-3 ml-7">
              {m.desc}
            </p>
            <div
              className="font-mono text-xs text-[#00FF87] bg-[#0D0D0D] border border-[#222] px-3 py-2 truncate"
              title={m.cmd}
            >
              {m.cmd}
            </div>
          </motion.div>
        ))}
      </motion.div>

      {/* OS badges */}
      <motion.div
        className="mt-12 flex flex-wrap gap-3 justify-center"
        initial={{ opacity: 0, y: 20 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5, delay: 0.3 }}
      >
        {["🐧 Linux", "🍎 macOS", "🪟 Windows", "🐳 Docker", "☁️ Cloud"].map(
          (os, i) => (
            <motion.span
              key={os}
              initial={{ opacity: 0, scale: 0.8 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ delay: 0.4 + i * 0.1 }}
              className="border-2 border-[#2A2A2A] text-neutral-400 text-sm font-semibold px-4 py-2 hover:border-[#FFD700] hover:text-[#FFD700] transition-colors"
            >
              {os}
            </motion.span>
          ),
        )}
      </motion.div>
    </section>
  );
}
