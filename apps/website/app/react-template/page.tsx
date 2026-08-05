"use client";
import { useEffect, useState } from "react";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";
import { Terminal, Copy, Check, Sparkles, Shield, Cpu } from "lucide-react";
import { trackTemplateView, trackTemplateCopy } from "../../utils/analytics";

export default function ReactTemplatePage() {
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    trackTemplateView("react");
  }, []);

  const handleCopy = () => {
    navigator.clipboard.writeText("autodev create react-ts my-react-app");
    setCopied(true);
    trackTemplateCopy("react");
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <main className="min-h-screen bg-slate-950 text-slate-100 flex flex-col justify-between">
      <Navbar />

      <div className="flex-1 pt-32 pb-24 px-6 max-w-5xl mx-auto w-full">
        {/* Schema Markup for Google SEO */}
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              "@context": "https://schema.org",
              "@type": "SoftwareSourceCode",
              name: "AutoDev React + TypeScript + Tailwind Template",
              description:
                "Production-ready React starter boilerplate configured with Vite, TypeScript, Tailwind CSS, ESLint, Prettier, and Docker.",
              programmingLanguage: "TypeScript",
              runtimePlatform: "Node.js",
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
          <span className="inline-flex items-center gap-1.5 px-3 py-1 border border-amber-500/30 bg-amber-500/10 text-amber-400 text-xs font-bold uppercase tracking-wider rounded-full">
            <Sparkles className="w-3 h-3" /> Best React Starter Template
          </span>
        </div>

        {/* Heading */}
        <div className="text-center md:text-left mb-12">
          <h1 className="text-4xl md:text-6xl font-black tracking-tight text-white mb-6">
            React + TypeScript + Tailwind CSS
            <br />
            <span className="text-[#C8F135]">Starter Boilerplate</span>
          </h1>
          <p className="text-lg md:text-xl text-slate-400 max-w-3xl leading-relaxed">
            Get a production-grade React environment set up in 3 seconds.
            Pre-configured with modular folders, code styles, lint rules, Docker
            orchestration, and CI/CD pipelines.
          </p>
        </div>

        {/* Copy Command Terminal */}
        <div className="mb-16 max-w-2xl">
          <p className="text-xs text-[#888] mb-2 uppercase tracking-widest font-semibold font-mono">
            Run in your terminal
          </p>
          <div className="terminal w-full p-4 bg-[#0d0d0d] border-2 border-[#2A2A2A] relative flex items-center justify-between font-mono text-sm text-[#00FF87]">
            <div className="flex items-center gap-2 overflow-x-auto pr-8">
              <span className="text-amber-500 font-bold">$</span>
              <span>autodev create react-ts my-react-app</span>
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

        {/* Features list */}
        <div className="mb-20">
          <h2 className="text-2xl font-bold text-white mb-8 border-b border-slate-900 pb-3">
            What&apos;s Pre-Configured
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Cpu className="w-8 h-8 text-amber-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Vite + React + TS
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Superfast builds with TypeScript strict checking, modular layout
                guidelines, and hot module reloading.
              </p>
            </div>
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Shield className="w-8 h-8 text-emerald-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Tailwind & PostCSS
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Fully pre-configured Tailwind classes, modular theme styling,
                and autoprefixer compilation configs.
              </p>
            </div>
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Terminal className="w-8 h-8 text-cyan-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Code Quality Audit
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Pre-configured ESLint and Prettier rules to enforce team
                standards and prevent regression bugs.
              </p>
            </div>
          </div>
        </div>

        {/* Benchmark / Funnel Card */}
        <div className="bg-[#111] border-2 border-[#2A2A2A] p-8">
          <h3 className="text-xl font-bold text-[#C8F135] mb-4">
            🚀 AutoDev Efficiency Stats
          </h3>
          <div className="grid grid-cols-1 sm:grid-cols-3 gap-6 text-left">
            <div>
              <p className="text-xs text-[#666] uppercase">Setup Time</p>
              <p className="text-2xl font-black text-emerald-400">
                3.2 seconds
              </p>
              <p className="text-[10px] text-[#555]">99% saved vs manual</p>
            </div>
            <div>
              <p className="text-xs text-[#666] uppercase">API Overhead</p>
              <p className="text-2xl font-black text-emerald-400">1 Prompt</p>
              <p className="text-[10px] text-[#555]">
                Instead of 18 prompt steps
              </p>
            </div>
            <div>
              <p className="text-xs text-[#666] uppercase">Tokens Consumed</p>
              <p className="text-2xl font-black text-emerald-400">9,000</p>
              <p className="text-[10px] text-[#555]">
                78% savings ($0.10 cost)
              </p>
            </div>
          </div>
        </div>
      </div>

      <Footer />
    </main>
  );
}
