"use client";
import { useEffect, useState } from "react";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";
import { Copy, Check, Sparkles, Shield, Cpu, Layers } from "lucide-react";
import { trackTemplateView, trackTemplateCopy } from "../../utils/analytics";

export default function FlutterTemplatePage() {
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    trackTemplateView("flutter");
  }, []);

  const handleCopy = () => {
    navigator.clipboard.writeText("autodev create flutter my-flutter-app");
    setCopied(true);
    trackTemplateCopy("flutter");
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
              name: "AutoDev Flutter Clean Architecture Template",
              description:
                "Production-ready Flutter clean structure template preloaded with multi-stage Docker builds and CI workflows.",
              programmingLanguage: "Dart",
              runtimePlatform: "Flutter",
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
          <span className="inline-flex items-center gap-1.5 px-3 py-1 border border-cyan-500/30 bg-cyan-500/10 text-cyan-400 text-xs font-bold uppercase tracking-wider rounded-full">
            <Sparkles className="w-3 h-3" /> Mobile & Web Boilerplate
          </span>
        </div>

        {/* Heading */}
        <div className="text-center md:text-left mb-12">
          <h1 className="text-4xl md:text-6xl font-black tracking-tight text-white mb-6">
            Flutter Clean
            <br />
            <span className="text-[#C8F135]">Architecture Template</span>
          </h1>
          <p className="text-lg md:text-xl text-slate-400 max-w-3xl leading-relaxed">
            Initialize cross-platform Flutter mobile, desktop, and web
            applications instantly, structured according to clean architecture
            directory conventions.
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
              <span>autodev create flutter my-flutter-app</span>
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

        {/* What's Pre-Configured */}
        <div className="mb-20">
          <h2 className="text-2xl font-bold text-white mb-8 border-b border-slate-900 pb-3">
            What&apos;s Pre-Configured
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Layers className="w-8 h-8 text-amber-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Clean Directory Structure
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Organized structure separating widgets, screens, and platform
                entrypoints for quick scalability.
              </p>
            </div>
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Shield className="w-8 h-8 text-emerald-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Web Nginx Serving
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Includes Dockerfile configurations to compile the codebase for
                the web and serve it using Nginx.
              </p>
            </div>
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Cpu className="w-8 h-8 text-cyan-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Automated CI Actions
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Pre-configured GitHub actions pipeline to check compilation
                status for Flutter web application targets.
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
                3.5 seconds
              </p>
              <p className="text-[10px] text-[#555]">
                99% saved vs manual configuration
              </p>
            </div>
            <div>
              <p className="text-xs text-[#666] uppercase">API Overhead</p>
              <p className="text-2xl font-black text-emerald-400">1 Prompt</p>
              <p className="text-[10px] text-[#555]">Instead of 24 prompts</p>
            </div>
            <div>
              <p className="text-xs text-[#666] uppercase">Tokens Consumed</p>
              <p className="text-2xl font-black text-emerald-400">8,500</p>
              <p className="text-[10px] text-[#555]">
                84% savings ($0.08 cost)
              </p>
            </div>
          </div>
        </div>
      </div>

      <Footer />
    </main>
  );
}
