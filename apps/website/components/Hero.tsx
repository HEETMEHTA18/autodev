"use client";
import { useState, useEffect } from "react";
import { motion } from "framer-motion";
import { Copy, Check } from "lucide-react";
import { trackInstall } from "../utils/analytics";
import AsciiBackground from "./AsciiBackground";

const container = {
  hidden: {},
  show: { transition: { staggerChildren: 0.12 } },
};
const item = {
  hidden: { opacity: 0, y: 30 },
  show: { opacity: 1, y: 0, transition: { duration: 0.5 } },
};

const words = [
  "DEVELOPERS.",
  "ENGINEERS.",
  "BUILDERS.",
  "HACKERS.",
  "CREATORS.",
];
function TypingText() {
  const [text, setText] = useState("");
  const [wordIndex, setWordIndex] = useState(0);
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    const currentWord = words[wordIndex];

    const timer = setTimeout(
      () => {
        if (!isDeleting) {
          if (text !== currentWord) {
            setText(currentWord.substring(0, text.length + 1));
          } else {
            setIsDeleting(true);
          }
        } else {
          if (text !== "") {
            setText(currentWord.substring(0, text.length - 1));
          } else {
            setIsDeleting(false);
            setWordIndex((prev) => (prev + 1) % words.length);
          }
        }
      },
      isDeleting ? 80 : text === currentWord ? 2000 : text === "" ? 300 : 120,
    );

    return () => clearTimeout(timer);
  }, [text, isDeleting, wordIndex]);

  return (
    <span className="text-[#FFD700] inline-flex items-center min-h-[1.1em]">
      FOR {text}
      <span className="inline-block w-[4px] md:w-[8px] h-[0.8em] bg-[#FFD700] ml-2 align-middle animate-pulse" />
    </span>
  );
}

export default function Hero() {
  const [copiedQuickInstall, setCopiedQuickInstall] = useState(false);
  const [activeTab, setActiveTab] = useState<"npx" | "curl">("npx");

  const handleCopyQuickInstall = () => {
    const cmd =
      activeTab === "npx"
        ? "npx @heetmehta18/autodev"
        : "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash";
    navigator.clipboard.writeText(cmd);
    setCopiedQuickInstall(true);
    trackInstall(activeTab);
    setTimeout(() => setCopiedQuickInstall(false), 1800);
  };

  return (
    <section className="relative overflow-hidden w-full">
      <AsciiBackground className="opacity-55" mouseSensi={2.0} />
      <div className="max-w-7xl mx-auto pt-36 pb-24 px-6 relative z-10">
        <motion.div variants={container} initial="hidden" animate="show">
        {/* Badges */}
        <motion.div
          variants={item}
          className="mb-8 flex flex-wrap items-center gap-4"
        >
          <span className="inline-block border-2 border-[#FFD700] text-[#FFD700] text-xs font-bold px-3 py-1 uppercase tracking-widest">
            v0.3.2 — Open Source
          </span>
          <a
            href="https://www.producthunt.com/products/autodevs?embed=true&utm_source=badge-featured&utm_medium=badge&utm_campaign=badge-autodevs"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-block hover:opacity-90 transition-opacity"
          >
            <img
              alt="Autodevs - AI-powered development setup in minutes | Product Hunt"
              width="250"
              height="54"
              src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=1162368&theme=neutral&t=1780484994611"
              className="h-[36px] w-auto"
            />
          </a>
        </motion.div>

        {/* Headline */}
        <motion.h1
          variants={item}
          className="text-[clamp(3.5rem,10vw,8rem)] font-black leading-[0.9] tracking-tighter text-white mb-4"
        >
          THE APP STORE
          <br />
          <TypingText />
        </motion.h1>

        {/* Sub-headline */}
        <motion.p
          variants={item}
          className="text-xl text-neutral-300 max-w-2xl mb-4 font-medium"
        >
          Clone. Scan. Install. Build.
        </motion.p>
        <motion.p
          variants={item}
          className="text-neutral-400 max-w-xl mb-12 leading-relaxed"
        >
          Install any language, framework, database, or DevOps tool with a
          single command. Smart dependency resolution. Cross-platform. Fully
          open-source.
        </motion.p>

        {/* CTAs */}
        <motion.div variants={item} className="flex flex-wrap gap-4 mb-16">
          <a href="#install" className="nb-btn px-8 py-4 text-lg inline-block">
            ⚡ GET STARTED
          </a>
          <a
            href="https://github.com/HEETMEHTA18/autodev"
            target="_blank"
            rel="noreferrer"
            className="nb-btn-outline px-8 py-4 text-lg inline-block"
          >
            View on GitHub →
          </a>
        </motion.div>

        {/* Quick install */}
        <motion.div variants={item} className="w-full max-w-xl">
          <p className="text-xs text-neutral-400 mb-2 uppercase tracking-widest font-semibold">
            Quick install
          </p>
          <div className="terminal w-full rounded-none relative">
            <div className="terminal-bar flex justify-between items-center pr-3 w-full">
              <div className="flex items-center gap-1.5">
                <span className="terminal-dot bg-[#FF5F56]" />
                <span className="terminal-dot bg-[#FFBD2E]" />
                <span className="terminal-dot bg-[#27C93F]" />
                <div className="flex gap-2 ml-4">
                  <button
                    onClick={() => setActiveTab("npx")}
                    className={`text-xs px-2 py-0.5 font-mono rounded cursor-pointer transition-all border ${
                      activeTab === "npx"
                        ? "bg-[#FFD700] text-black font-bold border-[#FFD700]"
                        : "text-neutral-400 border-transparent hover:text-white"
                    }`}
                  >
                    npx
                  </button>
                  <button
                    onClick={() => setActiveTab("curl")}
                    className={`text-xs px-2 py-0.5 font-mono rounded cursor-pointer transition-all border ${
                      activeTab === "curl"
                        ? "bg-[#FFD700] text-black font-bold border-[#FFD700]"
                        : "text-neutral-400 border-transparent hover:text-white"
                    }`}
                  >
                    curl
                  </button>
                </div>
              </div>
              <button
                onClick={handleCopyQuickInstall}
                className="text-neutral-400 hover:text-[#FFD700] transition-colors p-1 flex items-center gap-1 rounded bg-[#1e1e1e] border border-[#2a2a2a] cursor-pointer"
                title="Copy install command"
              >
                {copiedQuickInstall ? (
                  <>
                    <Check className="w-3.5 h-3.5 text-[#00FF87]" />
                    <span className="text-[10px] text-[#00FF87] font-mono pr-0.5">
                      Copied!
                    </span>
                  </>
                ) : (
                  <>
                    <Copy className="w-3.5 h-3.5" />
                    <span className="text-[10px] text-neutral-400 font-mono pr-0.5">
                      Copy
                    </span>
                  </>
                )}
              </button>
            </div>
            <div className="px-6 py-4 font-mono text-sm text-[#00FF87] overflow-x-auto whitespace-nowrap bg-black">
              <span className="text-neutral-500">$ </span>
              {activeTab === "npx"
                ? "npx @heetmehta18/autodev"
                : "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash"}
            </div>
          </div>
        </motion.div>

        {/* Stats row */}
        <motion.div variants={item} className="flex flex-wrap gap-8 mt-16">
          {[
            { value: "40+", label: "Packages" },
            { value: "9", label: "Dev Profiles" },
            { value: "3", label: "Platforms" },
            { value: "100%", label: "Open Source" },
          ].map(({ value, label }) => (
            <div key={label} className="nb-card px-6 py-4 min-w-[120px]">
              <div className="text-3xl font-black text-[#FFD700]">{value}</div>
              <div className="text-xs text-neutral-400 mt-1 font-semibold uppercase tracking-wider">
                {label}
              </div>
            </div>
          ))}
        </motion.div>
      </motion.div>
      </div>
    </section>
  );
}
