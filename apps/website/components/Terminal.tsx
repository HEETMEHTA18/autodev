"use client";
import { useState, useEffect, useRef } from "react";
import { motion, useInView } from "framer-motion";
import { Sparkles } from "lucide-react";

const lines = [
  { delay: 0, text: "$ autodev setup", color: "#00FF87" },
  { delay: 600, text: "", color: "#888" },
  {
    delay: 800,
    text: "  ⚡ AUTODEV · Clone. Scan. Install. Build.",
    color: "#C8F135",
  },
  { delay: 1000, text: "", color: "#888" },
  { delay: 1100, text: "  What do you want to install?", color: "#F0F0F0" },
  { delay: 1300, text: "", color: "#888" },
  { delay: 1400, text: "  1. Languages", color: "#2563eb" },
  { delay: 1500, text: "  2. Frameworks", color: "#2563eb" },
  { delay: 1600, text: "  3. Databases", color: "#2563eb" },
  { delay: 1700, text: "  4. DevOps", color: "#2563eb" },
  { delay: 1800, text: "  5. Mobile Development", color: "#2563eb" },
  { delay: 1900, text: "  6. AI / ML", color: "#2563eb" },
  { delay: 2000, text: "  7. Install by Profile", color: "#C8F135" },
  { delay: 2100, text: "", color: "#888" },
  {
    delay: 2300,
    text: "> [User selects: Languages → Python, Go, Node.js]",
    color: "#A3A3A3",
  },
  { delay: 2800, text: "", color: "#888" },
  { delay: 2900, text: "  Installing 🐍 Python...", color: "#C8F135" },
  { delay: 3400, text: "  ✓ Python 3.12 installed", color: "#00FF87" },
  { delay: 3600, text: "  Installing 🔵 Go...", color: "#C8F135" },
  { delay: 4100, text: "  ✓ Go 1.22 installed", color: "#00FF87" },
  { delay: 4300, text: "  Installing 🟢 Node.js...", color: "#C8F135" },
  { delay: 4800, text: "  ✓ Node.js LTS installed", color: "#00FF87" },
  { delay: 5000, text: "", color: "#888" },
  {
    delay: 5100,
    text: "  ✓ All done! Run 'autodev doctor' to verify your setup.",
    color: "#00FF87",
  },
];

export default function Terminal() {
  const [visible, setVisible] = useState(0);
  const [runKey, setRunKey] = useState(0);

  const containerRef = useRef<HTMLDivElement>(null);
  const outputRef = useRef<HTMLDivElement>(null);
  const isInView = useInView(containerRef, { once: true, amount: 0.1 });

  useEffect(() => {
    if (!isInView) return;

    const timers = lines.map((l, i) =>
      setTimeout(() => setVisible(i + 1), l.delay),
    );
    return () => timers.forEach(clearTimeout);
  }, [isInView, runKey]);

  // Auto-scroll to bottom of terminal output as text prints
  useEffect(() => {
    if (outputRef.current) {
      outputRef.current.scrollTop = outputRef.current.scrollHeight;
    }
  }, [visible]);

  const isRunning = visible > 0 && visible < lines.length;

  return (
    <section
      ref={containerRef}
      className="py-20 px-6 bg-[#0D0D0D] border-y-2 border-[#2A2A2A]"
    >
      <div className="max-w-5xl mx-auto">
        <motion.div
          className="text-center mb-12"
          initial={{ opacity: 0, y: 30 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
        >
          <span className="inline-flex items-center gap-1.5 border border-[#C8F135] text-[#C8F135] text-xs font-bold px-3 py-1 uppercase tracking-widest mb-4">
            <Sparkles className="w-3 h-3" /> Live Demo
          </span>
          <h2 className="text-4xl md:text-5xl font-black text-white mb-3">
            SEE IT IN <span className="text-gradient-yellow">ACTION</span>
          </h2>
          <p className="text-neutral-400 max-w-xl mx-auto">
            Run <code className="text-[#C8F135] font-mono">autodev setup</code>{" "}
            in any project and the interactive installer opens.
          </p>
        </motion.div>

        <div className="terminal">
          {/* Title bar */}
          <div className="terminal-bar flex justify-between items-center">
            <div className="flex items-center gap-2">
              <span className="terminal-dot bg-[#FF5F56]" />
              <span className="terminal-dot bg-[#FFBD2E]" />
              <span className="terminal-dot bg-[#27C93F]" />
              <span className="text-xs text-neutral-400 ml-3 font-mono">
                autodev — bash
              </span>
            </div>
            <button
              onClick={() => {
                setVisible(0);
                setRunKey((prev) => prev + 1);
              }}
              disabled={isRunning}
              className="text-xs text-[#C8F135] disabled:text-[#555] disabled:border-[#222] disabled:cursor-not-allowed hover:text-[#00FF87] hover:border-[#00FF87] font-mono bg-[#1E1E1E] px-2.5 py-1 rounded border border-[#333] transition-colors cursor-pointer"
            >
              {isRunning ? "Running..." : "Run Command"}
            </button>
          </div>

          {/* Output */}
          <div
            ref={outputRef}
            className="px-6 py-5 space-y-1 h-[400px] overflow-y-auto scroll-smooth"
          >
            {lines.slice(0, visible).map((line, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.15 }}
                className="font-mono text-sm leading-6"
                style={{ color: line.color }}
              >
                {line.text || "\u00A0"}
              </motion.div>
            ))}
            {/* Blinking cursor */}
            {visible < lines.length && (
              <span className="inline-block w-2 h-4 bg-[#C8F135] animate-pulse" />
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
