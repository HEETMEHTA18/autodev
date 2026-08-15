import { ArrowRight, Bot, Boxes, CheckCircle2, ScanSearch, Settings2, Wrench } from "lucide-react";

const actions = [
  { icon: Settings2, title: "Setup", text: "Install the languages, frameworks and tools your project needs." },
  { icon: ScanSearch, title: "Scan", text: "Understand dependencies, configuration and environment state." },
  { icon: Wrench, title: "Doctor", text: "Find problems and get actionable fixes before they block you." },
  { icon: Bot, title: "Agent", text: "Use an AI agent that can reason about your workspace and tools." },
];

export default function CommandCenterPreview() {
  return (
    <section className="max-w-6xl mx-auto px-6 py-24">
      <div className="grid lg:grid-cols-[1fr_1.15fr] gap-10 items-center">
        <div>
          <div className="text-xs font-semibold uppercase tracking-[0.2em] text-[#22D3EE] mb-4">One entry point</div>
          <h2 className="text-3xl md:text-5xl font-bold tracking-tight mb-5">Stop memorizing commands.</h2>
          <p className="text-slate-400 text-lg leading-relaxed mb-7">
            Run <span className="font-mono text-slate-200">autodev</span> and use the Command Center. Every major capability is discoverable from the same workflow, while the full CLI remains available for scripts and power users.
          </p>
          <div className="flex flex-wrap gap-3 text-sm text-slate-300">
            {['Detect', 'Plan', 'Execute', 'Verify'].map((step) => (
              <span key={step} className="inline-flex items-center gap-2 rounded-full border border-[#24304A] bg-[#0F1524] px-3 py-1.5">
                <CheckCircle2 className="w-4 h-4 text-[#34D399]" /> {step}
              </span>
            ))}
          </div>
        </div>

        <div className="rounded-2xl border border-[#24304A] bg-[#0F1524]/90 p-5 shadow-2xl">
          <div className="flex items-center gap-2 border-b border-[#24304A] pb-4 mb-4">
            <div className="w-2.5 h-2.5 rounded-full bg-[#FB7185]" />
            <div className="w-2.5 h-2.5 rounded-full bg-amber-400" />
            <div className="w-2.5 h-2.5 rounded-full bg-[#34D399]" />
            <span className="ml-2 text-xs text-slate-500 font-mono">autodev</span>
          </div>
          <div className="font-mono text-sm text-slate-300 mb-5">
            <span className="text-[#22D3EE]">$</span> autodev
          </div>
          <div className="grid sm:grid-cols-2 gap-3">
            {actions.map(({ icon: Icon, title, text }, index) => (
              <div key={title} className={`rounded-xl border p-4 transition ${index === 0 ? 'border-[#7C3AED]/60 bg-[#7C3AED]/10' : 'border-[#24304A] bg-[#0B1020]'}`}>
                <div className="flex items-center gap-3 mb-2">
                  <Icon className="w-4 h-4 text-[#8B5CF6]" />
                  <span className="font-semibold text-slate-100">{title}</span>
                  <ArrowRight className="w-3.5 h-3.5 ml-auto text-slate-600" />
                </div>
                <p className="text-xs leading-relaxed text-slate-500">{text}</p>
              </div>
            ))}
          </div>
          <div className="mt-4 rounded-xl border border-[#24304A] bg-[#0B1020] p-3 text-xs text-slate-500 font-mono">
            <span className="text-[#34D399]">✓</span> Environment ready · 0 blocking issues
          </div>
        </div>
      </div>
    </section>
  );
}
