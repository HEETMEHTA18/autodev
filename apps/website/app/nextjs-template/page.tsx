"use client";
import { useEffect, useState } from "react";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";
import { Copy, Check, Sparkles, Shield, Cpu, Zap } from "lucide-react";
import { trackTemplateView, trackTemplateCopy } from "../../utils/analytics";

export default function NextJsTemplatePage() {
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    trackTemplateView("nextjs");
  }, []);

  const handleCopy = () => {
    navigator.clipboard.writeText("autodev create nextjs my-next-app");
    setCopied(true);
    trackTemplateCopy("nextjs");
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <main className="min-h-screen bg-slate-950 text-slate-100 flex flex-col justify-between">
      <Navbar />

      <div className="flex-1 pt-32 pb-24 px-6 max-w-5xl mx-auto w-full">
        {/* Schema Markup */}
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              "@context": "https://schema.org",
              "@type": "SoftwareSourceCode",
              name: "AutoDev Next.js Production Starter Template",
              description:
                "Production-grade Next.js starter boilerplate optimized with standalone Docker, TypeScript, Tailwind CSS, ESLint, and GitHub Actions.",
              programmingLanguage: "TypeScript",
              runtimePlatform: "Next.js",
              targetProduct: {
                "@type": "SoftwareApplication",
                name: "AutoDev CLI",
                operatingSystem: "Linux, macOS, Windows",
              },
            }),
          }}
        />

        {/* Badge */}
        <div className="mb-8 text-center md:text-left">
          <span className="inline-flex items-center gap-1.5 px-3 py-1 border border-emerald-500/30 bg-emerald-500/10 text-emerald-400 text-xs font-bold uppercase tracking-wider rounded-full">
            <Sparkles className="w-3 h-3" /> Best Next.js Boilerplate
          </span>
        </div>

        {/* Heading */}
        <div className="text-center md:text-left mb-12">
          <h1 className="text-4xl md:text-6xl font-black tracking-tight text-white mb-6">
            Next.js Standalone
            <br />
            <span className="text-[#C8F135]">Production Boilerplate</span>
          </h1>
          <p className="text-lg md:text-xl text-slate-400 max-w-3xl leading-relaxed">
            Boost your product startup using Next.js App Router, fully prepared
            for containerized orchestration (Docker/K8s) and continuous
            integrations pipeline checks.
          </p>
        </div>

        {/* Copy Command Terminal */}
        <div className="mb-16 max-w-2xl">
          <p className="text-xs text-[#888] mb-2 uppercase tracking-widest font-semibold font-mono">
            Execute command in terminal
          </p>
          <div className="terminal w-full p-4 bg-[#0d0d0d] border-2 border-[#2A2A2A] relative flex items-center justify-between font-mono text-sm text-[#00FF87]">
            <div className="flex items-center gap-2 overflow-x-auto pr-8">
              <span className="text-amber-500 font-bold">$</span>
              <span>autodev create nextjs my-next-app</span>
            </div>
            <button
              onClick={handleCopy}
              className="text-[#555] hover:text-[#C8F135] transition-colors p-2 shrink-0 cursor-pointer border border-[#2a2a2a] bg-[#1a1a1a]"
              title="Copy command"
            >
              {copied ? (
                <Check className="w-4 h-4 text-emerald-400" />
              ) : (
                <Copy className="w-4 h-4" />
              )}
            </button>
          </div>
        </div>

        {/* Details section */}
        <div className="mb-20">
          <h2 className="text-2xl font-bold text-white mb-8 border-b border-slate-900 pb-3">
            What&apos;s Inside
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Zap className="w-8 h-8 text-emerald-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Next.js App Router
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Clean and scalable page layouts, layouts metadata optimization,
                and layout variables ready for server-side fetches.
              </p>
            </div>
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Shield className="w-8 h-8 text-amber-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Standalone Docker
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Multi-stage build Dockerfile generating standalone node outputs,
                lowering compiled container memory footprints down to ~80MB.
              </p>
            </div>
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Cpu className="w-8 h-8 text-cyan-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                CI/CD Pipeline
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Pre-configured GitHub Action running automatic code formatting
                checks, ESLint verification audits, and build compiles on pull
                actions.
              </p>
            </div>
          </div>
        </div>

        {/* Stats card */}
        <div className="bg-[#111] border-2 border-[#2A2A2A] p-8">
          <h3 className="text-xl font-bold text-[#C8F135] mb-4">
            🚀 AutoDev Efficiency Stats
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 text-left">
            <div>
              <p className="text-xs text-[#666] uppercase">Setup Time</p>
              <p className="text-2xl font-black text-emerald-400">
                4.1 seconds
              </p>
              <p className="text-[10px] text-[#555]">
                99% saved vs manual setup
              </p>
            </div>
            <div>
              <p className="text-xs text-[#666] uppercase">API Overhead</p>
              <p className="text-2xl font-black text-emerald-400">1 Prompt</p>
              <p className="text-[10px] text-[#555]">Instead of 32 prompts</p>
            </div>
            <div>
              <p className="text-xs text-[#666] uppercase">Tokens Consumed</p>
              <p className="text-2xl font-black text-emerald-400">12,000</p>
              <p className="text-[10px] text-[#555]">
                86% savings ($0.14 cost)
              </p>
            </div>
          </div>
        </div>
      </div>

      <Footer />
    </main>
  );
}
