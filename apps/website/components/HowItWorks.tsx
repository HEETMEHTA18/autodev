"use client";
import { motion } from "framer-motion";
import { Download, Search, Brain, Zap, ArrowRight, ArrowDown } from "lucide-react";
import Image from "next/image";

const steps = [
  {
    num: "01",
    title: "Install AutoDevs",
    desc: "Install AutoDevs or run it directly using npx on your machine.",
    icon: Download,
    color: "#C8F135",
  },
  {
    num: "02",
    title: "Scan Repository",
    desc: "Point it to any codebase to scan dependencies, structure, and configurations.",
    icon: Search,
    color: "#00FF87",
  },
  {
    num: "03",
    title: "Get AI Analysis",
    desc: "Receive structured security audits, performance suggestions, and fixes.",
    icon: Brain,
    color: "#2563eb",
  },
  {
    num: "04",
    title: "Fix Issues Faster",
    desc: "Generate developer plans and apply changes with guided assistance.",
    icon: Zap,
    color: "#dc2626",
  },
];

export default function HowItWorks() {
  return (
    <section id="how-it-works" className="py-24 px-6 bg-[#0A0A0A] border-b-2 border-[#2A2A2A]">
      <div className="max-w-7xl mx-auto">
        {/* Section Header */}
        <div className="text-center mb-16">
          <span className="inline-flex items-center gap-1.5 border border-[#C8F135] text-[#C8F135] text-xs font-bold px-3 py-1 uppercase tracking-widest mb-4">
            ⚡ Workflow Demo
          </span>
          <h2 className="text-4xl md:text-6xl font-black text-white mb-4 tracking-tighter">
            HOW <span className="text-gradient-yellow">AUTODEVS</span> WORKS
          </h2>
          <p className="text-neutral-400 text-lg max-w-2xl mx-auto font-medium">
            Turn any codebase into actionable insights and structured developer plans in 4 simple steps.
          </p>
        </div>

        {/* Content Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
          {/* Steps Column */}
          <div className="lg:col-span-5 space-y-6">
            {steps.map((step, idx) => {
              const Icon = step.icon;
              return (
                <div key={idx} className="relative">
                  {/* Step Card */}
                  <motion.div
                    initial={{ opacity: 0, x: -30 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    viewport={{ once: true }}
                    transition={{ duration: 0.4, delay: idx * 0.1 }}
                    className="nb-card p-6 flex gap-4 items-start relative z-10 hover:border-white hover:shadow-[4px_4px_0px_#FFFFFF] transition-all"
                  >
                    <div
                      className="p-3 border-2 border-current shrink-0"
                      style={{ color: step.color }}
                    >
                      <Icon className="w-6 h-6" />
                    </div>
                    <div>
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-xs font-mono font-bold text-neutral-500 uppercase">
                          Step {step.num}
                        </span>
                        <span className="w-1.5 h-1.5 rounded-full bg-neutral-600" />
                        <h3 className="text-lg font-black text-white uppercase tracking-tight">
                          {step.title}
                        </h3>
                      </div>
                      <p className="text-sm text-neutral-400 font-medium leading-relaxed">
                        {step.desc}
                      </p>
                    </div>
                  </motion.div>

                  {/* Flow Arrow */}
                  {idx < steps.length - 1 && (
                    <div className="flex justify-center lg:justify-start lg:pl-10 py-2 text-neutral-600">
                      <ArrowDown className="w-5 h-5 animate-bounce" />
                    </div>
                  )}
                </div>
              );
            })}
          </div>

          {/* Video/GIF Demo Column */}
          <div className="lg:col-span-7">
            <motion.div
              initial={{ opacity: 0, scale: 0.95 }}
              whileInView={{ opacity: 1, scale: 1 }}
              viewport={{ once: true }}
              transition={{ duration: 0.5 }}
              className="terminal overflow-hidden"
            >
              {/* Terminal Window Header */}
              <div className="terminal-bar flex justify-between items-center px-4 py-3">
                <div className="flex items-center gap-2">
                  <span className="terminal-dot bg-[#FF5F56]" />
                  <span className="terminal-dot bg-[#FFBD2E]" />
                  <span className="terminal-dot bg-[#27C93F]" />
                  <span className="text-xs text-neutral-400 ml-3 font-mono font-semibold">
                    autodev scan — 30s Demo
                  </span>
                </div>
                <div className="flex items-center gap-1.5 border border-[#2A2A2A] bg-black/40 px-2 py-0.5 rounded text-[10px] text-neutral-400 font-mono font-bold">
                  <span className="w-1.5 h-1.5 rounded-full bg-[#00FF87] animate-ping" />
                  LIVE DEMO
                </div>
              </div>

              {/* Terminal Body with Demo GIF */}
              <div className="relative bg-black aspect-video flex items-center justify-center border-t border-[#1A1A1A]">
                <Image
                  src="/autodev-demo.gif"
                  alt="AutoDevs command-line interface scan demonstration"
                  fill
                  sizes="(max-width: 768px) 100vw, 800px"
                  className="object-cover opacity-90 hover:opacity-100 transition-opacity"
                  unoptimized
                />
                
                {/* Overlay Banner */}
                <div className="absolute bottom-4 left-4 right-4 bg-black/85 backdrop-blur-sm border-2 border-[#C8F135] p-3 text-center shadow-[4px_4px_0px_#000]">
                  <p className="text-xs md:text-sm font-mono font-bold text-[#C8F135]">
                    $ npx @heetmehta18/autodev scan
                  </p>
                </div>
              </div>
            </motion.div>
            
            <p className="text-center text-xs text-neutral-500 mt-4 font-mono">
              * AutoDevs is like Cursor + GitHub Analysis + AI Project Planning in one developer tool.
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
