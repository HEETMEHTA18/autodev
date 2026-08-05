"use client";
import { useEffect, useState } from "react";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";
import { Copy, Check, Sparkles, Shield, Cpu, Bot } from "lucide-react";
import { trackTemplateView, trackTemplateCopy } from "../../utils/analytics";

export default function AiAgentTemplatePage() {
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    trackTemplateView("ai-chatbot");
  }, []);

  const handleCopy = () => {
    navigator.clipboard.writeText("autodev create ai-chatbot my-ai-agent");
    setCopied(true);
    trackTemplateCopy("ai-chatbot");
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <main className="min-h-screen bg-slate-955 text-slate-100 flex flex-col justify-between">
      <Navbar />

      <div className="flex-1 pt-32 pb-24 px-6 max-w-5xl mx-auto w-full">
        {/* Schema Markup */}
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              "@context": "https://schema.org",
              "@type": "SoftwareSourceCode",
              name: "AutoDev Gemini AI Chatbot Boilerplate",
              description:
                "Full-stack AI Chatbot starter environment preloaded with Express backend integrating Google Gemini 2.5 Flash API and React client dashboard.",
              programmingLanguage: "JavaScript, TypeScript",
              runtimePlatform: "Node.js, Gemini API",
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
            <Sparkles className="w-3 h-3" /> AI Chatbot Starter Kit
          </span>
        </div>

        {/* Heading */}
        <div className="text-center md:text-left mb-12">
          <h1 className="text-4xl md:text-6xl font-black tracking-tight text-white mb-6">
            Google Gemini API
            <br />
            <span className="text-[#C8F135]">AI Chatbot Boilerplate</span>
          </h1>
          <p className="text-lg md:text-xl text-slate-400 max-w-3xl leading-relaxed">
            Build your AI Agent products immediately using this pre-set
            monorepo. It features a React chatbot interface connected to an
            Express.js backend utilizing Google&apos;s official Gemini 2.5 Flash
            SDK.
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
              <span>autodev create ai-chatbot my-ai-agent</span>
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
              <Bot className="w-8 h-8 text-amber-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Gemini 2.5 Flash SDK
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Backend Express route fully integrated with Google&apos;s
                official `@google/genai` library, ready to accept prompt inputs.
              </p>
            </div>
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Shield className="w-8 h-8 text-emerald-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Futuristic Chat Dashboard
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                React chatbot UI with messaging bubbles, animations, prompt
                input forms, and dynamic message loaders.
              </p>
            </div>
            <div className="bg-[#111] border-2 border-[#2A2A2A] p-6">
              <Cpu className="w-8 h-8 text-cyan-400 mb-4" />
              <h3 className="text-lg font-bold text-white mb-2">
                Multi-stage Dockerfile
              </h3>
              <p className="text-xs text-slate-400 leading-relaxed">
                Includes Docker files to easily bundle both backend Node API and
                frontend static React client for distribution.
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
                4.5 seconds
              </p>
              <p className="text-[10px] text-[#555]">
                99% saved vs manual configuration
              </p>
            </div>
            <div>
              <p className="text-xs text-[#666] uppercase">API Overhead</p>
              <p className="text-2xl font-black text-emerald-400">1 Prompt</p>
              <p className="text-[10px] text-[#555]">Instead of 45 prompts</p>
            </div>
            <div>
              <p className="text-xs text-[#666] uppercase">Tokens Consumed</p>
              <p className="text-2xl font-black text-emerald-400">15,500</p>
              <p className="text-[10px] text-[#555]">
                87% savings ($0.18 cost)
              </p>
            </div>
          </div>
        </div>
      </div>

      <Footer />
    </main>
  );
}
