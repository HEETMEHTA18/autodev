"use client";
import { useEffect, useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  X,
  Copy,
  Check,
  Sparkles,
  Terminal as TermIcon,
  Download,
} from "lucide-react";

const updateMethods = [
  {
    label: "NPM (Global)",
    cmd: "npm install -g @heetmehta18/autodev@latest",
    desc: "Upgrade to the latest NPM wrapper",
  },
  {
    label: "Homebrew",
    cmd: "brew upgrade autodev",
    desc: "Update via macOS/Linux brew tap",
  },
  {
    label: "Scoop",
    cmd: "scoop update autodev",
    desc: "Upgrade on Windows via Scoop bucket",
  },
];

export default function UpdatePopup() {
  const [isOpen, setIsOpen] = useState(false);
  const [activeTab, setActiveTab] = useState(0);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    // Detect Lighthouse / headless / automated browser testing environments
    const isLighthouse =
      typeof window !== "undefined" &&
      (navigator.webdriver ||
        /lighthouse|chrome-lighthouse|headless/i.test(navigator.userAgent));
    if (isLighthouse) return;

    // Check if user has already dismissed this update notification
    const isDismissed = localStorage.getItem("autodev_update_dismissed_v0.5.2");
    if (!isDismissed) {
      // Delay showing the popup slightly for better UX
      const timer = setTimeout(() => {
        setIsOpen(true);
      }, 1500);
      return () => clearTimeout(timer);
    }
  }, []);

  useEffect(() => {
    const handleOpen = () => {
      setIsOpen(true);
    };
    window.addEventListener("autodev_open_update_modal", handleOpen);
    return () =>
      window.removeEventListener("autodev_open_update_modal", handleOpen);
  }, []);

  const handleDismiss = () => {
    localStorage.setItem("autodev_update_dismissed_v0.5.2", "true");
    setIsOpen(false);
  };

  const copyCommand = (cmd: string) => {
    navigator.clipboard.writeText(cmd);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0, y: 100, scale: 0.95 }}
          animate={{ opacity: 1, y: 0, scale: 1 }}
          exit={{ opacity: 0, y: 100, scale: 0.95 }}
          transition={{ type: "spring", stiffness: 300, damping: 30 }}
          className="fixed bottom-6 right-6 z-[100] w-full max-w-[420px] p-1.5"
        >
          {/* Card Border wrapper (Neo-brutalist) */}
          <div className="bg-black border-4 border-[#C8F135] text-white shadow-[8px_8px_0px_0px_#111] p-5 relative">
            {/* Close Button */}
            <button
              onClick={handleDismiss}
              className="absolute top-3 right-3 text-neutral-400 hover:text-[#C8F135] transition-colors p-1 border-2 border-transparent hover:border-[#2A2A2A] bg-neutral-900"
              aria-label="Dismiss notification"
            >
              <X className="w-4 h-4" />
            </button>

            {/* Header */}
            <div className="flex items-center gap-2 mb-3">
              <span className="bg-[#C8F135] text-black text-[10px] font-black px-2 py-0.5 tracking-wider uppercase animate-pulse">
                New Release
              </span>
              <span className="text-neutral-400 text-xs font-mono font-bold">
                v0.5.2 is live!
              </span>
            </div>

            <h3 className="text-xl font-black tracking-tight text-white mb-2 flex items-center gap-1.5">
              <Sparkles className="w-5 h-5 text-[#C8F135] shrink-0" />
              UPGRADE TO AUTODEV v0.5.2
            </h3>

            <p className="text-neutral-300 text-xs leading-relaxed mb-4">
              AutoDevs is now a unified agent orchestrator: install, launch, and
              manage multiple AI CLIs (OpenCode, Claude Code, Codex, Gemini
              CLI, Aider) with <code className="text-[#C8F135]">autodev tools</code>,{" "}
              <code className="text-[#C8F135]">autodev run</code>,{" "}
              <code className="text-[#C8F135]">autodev session</code>, and{" "}
              <code className="text-[#C8F135]">autodev agent</code> task routing.
            </p>

            {/* Selector Tabs */}
            <div className="border-b border-[#2A2A2A] flex mb-3 text-[11px] font-bold font-mono">
              {updateMethods.map((m, idx) => (
                <button
                  key={m.label}
                  onClick={() => {
                    setActiveTab(idx);
                    setCopied(false);
                  }}
                  className={`pb-1.5 px-2.5 border-b-2 transition-colors cursor-pointer ${
                    activeTab === idx
                      ? "border-[#C8F135] text-[#C8F135]"
                      : "border-transparent text-neutral-400 hover:text-neutral-200"
                  }`}
                >
                  {m.label.split(" ")[0]}
                </button>
              ))}
            </div>

            {/* Command Display */}
            <div className="bg-[#0A0A0A] border-2 border-[#2A2A2A] p-3 mb-4">
              <div className="flex items-center justify-between gap-2 mb-1.5">
                <span className="text-[10px] text-neutral-400 font-mono font-bold">
                  {updateMethods[activeTab].desc}
                </span>
                <button
                  onClick={() => copyCommand(updateMethods[activeTab].cmd)}
                  className="flex items-center gap-1 text-[10px] text-neutral-400 hover:text-[#00FF87] font-mono font-bold transition-colors"
                >
                  {copied ? (
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
              </div>

              <div className="flex items-center gap-2 font-mono text-xs text-[#00FF87] truncate select-all">
                <TermIcon className="w-3.5 h-3.5 text-neutral-400 shrink-0" />
                <span>{updateMethods[activeTab].cmd}</span>
              </div>
            </div>

            {/* Action buttons */}
            <div className="flex gap-2">
              <a
                href="https://github.com/HEETMEHTA18/autodev/releases/tag/v0.5.2"
                target="_blank"
                rel="noreferrer"
                className="nb-btn-small flex items-center justify-center gap-1.5 flex-1 py-1.5 bg-[#C8F135] text-black text-xs font-black tracking-wider uppercase"
              >
                <Download className="w-3.5 h-3.5" />
                Release Notes
              </a>
              <button
                onClick={handleDismiss}
                className="flex-1 py-1.5 text-center text-xs font-mono font-bold border-2 border-[#2A2A2A] bg-[#111] hover:border-neutral-500 text-neutral-400 hover:text-white transition-colors cursor-pointer"
              >
                Maybe Later
              </button>
            </div>
          </div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
