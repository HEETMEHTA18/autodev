"use client";
import { useState, useEffect, useRef } from "react";
import { motion, useInView } from "framer-motion";
import { Terminal, CheckSquare, Square, Clock, Sparkles, AlertTriangle } from "lucide-react";

export default function RealExample() {
  const [step, setStep] = useState(0);
  const [typedCommand, setTypedCommand] = useState("");
  const [checkedTasks, setCheckedTasks] = useState<boolean[]>([false, false, false, false, false]);
  
  const containerRef = useRef<HTMLDivElement>(null);
  const isInView = useInView(containerRef, { once: true, amount: 0.1 });

  // Typist effect for $ autodev scan my-react-app
  useEffect(() => {
    if (!isInView) return;
    
    const commandText = "autodev scan my-react-app";
    let index = 0;
    
    const typingTimer = setInterval(() => {
      if (index < commandText.length) {
        setTypedCommand((prev) => prev + commandText.charAt(index));
        index++;
      } else {
        clearInterval(typingTimer);
        // Start streaming output steps after command typing completes
        setTimeout(() => setStep(1), 500); // Scan start
        setTimeout(() => setStep(2), 1200); // Progress details
        setTimeout(() => setStep(3), 2000); // Findings summary
        setTimeout(() => setStep(4), 2800); // Final task summary
      }
    }, 60);

    return () => clearInterval(typingTimer);
  }, [isInView]);

  const toggleTask = (index: number) => {
    const updated = [...checkedTasks];
    updated[index] = !updated[index];
    setCheckedTasks(updated);
  };

  const sampleTasks = [
    { id: 1, text: "Add React Error Boundary wrapper to root layout", category: "Stability" },
    { id: 2, text: "Configure Webpack bundle splitting & dynamic imports", category: "Performance" },
    { id: 3, text: "Prune unused node_modules (lodash, moment) from packages", category: "Cleanup" },
    { id: 4, text: "Write unit tests for authentication hooks & providers", category: "Testing" },
    { id: 5, text: "Resolve critical CVE in active package-lock dependencies", category: "Security" },
  ];

  return (
    <section id="real-example" ref={containerRef} className="py-24 px-6 bg-[#0D0D0D] border-b-2 border-[#2A2A2A]">
      <div className="max-w-7xl mx-auto">
        
        {/* Header */}
        <div className="text-center mb-16">
          <span className="inline-flex items-center gap-1.5 border border-[#00FF87] text-[#00FF87] text-xs font-bold px-3 py-1 uppercase tracking-widest mb-4">
            <Sparkles className="w-3 h-3" /> Real Example
          </span>
          <h2 className="text-4xl md:text-6xl font-black text-white mb-4 tracking-tighter">
            SEE A <span className="text-gradient-yellow">REAL SCAN</span>
          </h2>
          <p className="text-neutral-400 text-lg max-w-2xl mx-auto font-medium">
            Watch AutoDevs analyze a standard React app and transform unstructured code issues into discrete developer tasks.
          </p>
        </div>

        {/* Layout Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-stretch">
          
          {/* Left Column: Command CLI Terminal */}
          <div className="lg:col-span-6 flex flex-col justify-between">
            <div className="flex items-center gap-2 mb-4">
              <Terminal className="w-5 h-5 text-[#FFD700]" />
              <span className="text-sm font-mono font-bold text-white uppercase tracking-wider">
                Interactive CLI Simulation
              </span>
            </div>

            <div className="terminal w-full rounded-none flex-grow flex flex-col bg-black">
              {/* Header bar */}
              <div className="terminal-bar flex justify-between items-center px-4 py-3">
                <div className="flex items-center gap-2">
                  <span className="terminal-dot bg-[#FF5F56]" />
                  <span className="terminal-dot bg-[#FFBD2E]" />
                  <span className="terminal-dot bg-[#27C93F]" />
                  <span className="text-xs text-neutral-400 ml-3 font-mono">
                    ~ - autodev
                  </span>
                </div>
                <div className="text-[10px] text-neutral-600 font-mono">
                  bash
                </div>
              </div>

              {/* Stream output */}
              <div className="p-6 font-mono text-sm leading-6 space-y-2 flex-grow overflow-y-auto min-h-[360px]">
                <div className="text-[#00FF87]">
                  <span className="text-neutral-500">$ </span>
                  {typedCommand}
                  {typedCommand.length < 25 && (
                    <span className="inline-block w-1.5 h-4 bg-[#FFD700] ml-0.5 animate-pulse" />
                  )}
                </div>

                {step >= 1 && (
                  <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="text-neutral-400">
                    🔍 Scan initiated in <span className="text-[#FFD700]">./my-react-app</span>
                  </motion.div>
                )}

                {step >= 2 && (
                  <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-1">
                    <div className="text-[#4A90E2]">● Indexing source files (128 source files found)...</div>
                    <div className="text-[#4A90E2]">● Resolving dependency graph configuration...</div>
                    <div className="text-[#4A90E2]">● Evaluating code smells, test-coverage, and lockfiles...</div>
                  </motion.div>
                )}

                {step >= 3 && (
                  <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="pt-4 space-y-1 border-t border-[#1A1A1A] mt-4">
                    <div className="text-[#FF4444] font-bold">Found 12 Issues:</div>
                    <div className="text-neutral-300 ml-3">1. Missing error boundaries (stability danger)</div>
                    <div className="text-neutral-300 ml-3">2. Large bundle size (slow page rendering)</div>
                    <div className="text-neutral-300 ml-3">3. Unused dependencies (unneeded overhead)</div>
                    <div className="text-neutral-300 ml-3">4. Missing tests (low regression coverage)</div>
                    <div className="text-neutral-300 ml-3">5. Security vulnerability (dep-lock CVE warning)</div>
                  </motion.div>
                )}

                {step >= 4 && (
                  <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="pt-4 border-t border-dashed border-[#2A2A2A] mt-4 space-y-1">
                    <div className="text-[#00FF87] font-bold">Generated 8 Improvement Tasks</div>
                    <div className="text-[#FFD700] font-bold">Estimated Fix Time: 2 hours</div>
                  </motion.div>
                )}
              </div>
            </div>
          </div>

          {/* Right Column: Generated Task Dashboard Output */}
          <div className="lg:col-span-6 flex flex-col justify-between">
            <div className="flex items-center gap-2 mb-4 justify-between">
              <div className="flex items-center gap-2">
                <CheckSquare className="w-5 h-5 text-[#00FF87]" />
                <span className="text-sm font-mono font-bold text-white uppercase tracking-wider">
                  Generated Task Board
                </span>
              </div>
              {step >= 4 && (
                <div className="flex items-center gap-1.5 px-3 py-1 bg-yellow/10 border border-yellow/20 rounded-full text-xs font-mono font-bold text-[#FFD700]">
                  <Clock className="w-3.5 h-3.5" />
                  2 hrs left
                </div>
              )}
            </div>

            {/* Task Checklist Dashboard */}
            <div className="nb-card p-6 bg-[#111] flex-grow flex flex-col justify-between">
              <div>
                <div className="flex justify-between items-center mb-6 pb-4 border-b border-[#2A2A2A]">
                  <div>
                    <h4 className="text-lg font-black text-white uppercase tracking-tight">
                      Generated Plan: my-react-app
                    </h4>
                    <p className="text-xs text-neutral-400 font-mono mt-0.5">
                      Ready to execute · 5/8 critical shown
                    </p>
                  </div>
                  <span className="text-xs px-2.5 py-1 border border-[#FFD700] text-[#FFD700] font-mono font-bold bg-[#FFD700]/5">
                    EST: 2H
                  </span>
                </div>

                {/* List items */}
                {step < 3 ? (
                  <div className="h-[240px] flex flex-col items-center justify-center text-center p-6 border-2 border-dashed border-[#2A2A2A]">
                    <Clock className="w-8 h-8 text-neutral-600 mb-2 animate-spin" />
                    <p className="text-sm font-mono text-neutral-500">
                      Waiting for CLI scan to complete...
                    </p>
                  </div>
                ) : (
                  <div className="space-y-4">
                    {sampleTasks.map((task, idx) => (
                      <motion.div
                        initial={{ opacity: 0, y: 10 }}
                        animate={{ opacity: 1, y: 0 }}
                        transition={{ delay: idx * 0.1 }}
                        key={task.id}
                        onClick={() => toggleTask(idx)}
                        className={`flex gap-3 items-start p-3 border-2 transition-all cursor-pointer select-none ${
                          checkedTasks[idx]
                            ? "bg-black/40 border-[#00FF87]/40 opacity-60"
                            : "bg-[#161616] border-[#2A2A2A] hover:border-white hover:shadow-[3px_3px_0px_#2A2A2A]"
                        }`}
                      >
                        <div className="mt-0.5 shrink-0">
                          {checkedTasks[idx] ? (
                            <CheckSquare className="w-4 h-4 text-[#00FF87]" />
                          ) : (
                            <Square className="w-4 h-4 text-neutral-500" />
                          )}
                        </div>
                        <div className="flex-grow">
                          <p className={`text-xs font-semibold ${
                            checkedTasks[idx] ? "line-through text-neutral-500" : "text-neutral-200"
                          }`}>
                            {task.text}
                          </p>
                          <div className="flex gap-2 items-center mt-1">
                            <span className="text-[10px] font-mono uppercase bg-[#222] text-neutral-400 px-1.5 py-0.5">
                              {task.category}
                            </span>
                          </div>
                        </div>
                      </motion.div>
                    ))}
                  </div>
                )}
              </div>

              {/* Bottom Quote Banner */}
              <div className="mt-6 pt-4 border-t border-[#2A2A2A] flex justify-between items-center text-xs font-mono">
                <span className="text-neutral-500">AutoDevs integrates with any CI framework.</span>
                <span className="text-[#00FF87] font-bold">100% Automated</span>
              </div>
            </div>
          </div>
        </div>

        {/* Simple sentence at the bottom */}
        <p className="text-center font-mono text-xs md:text-sm text-neutral-400 mt-12 bg-black/60 p-4 border border-[#2A2A2A] max-w-3xl mx-auto">
          💡 <span className="font-bold text-white">Compare:</span> AutoDevs is like Cursor + GitHub Analysis + AI Project Planning in one developer tool.
        </p>

      </div>
    </section>
  );
}
