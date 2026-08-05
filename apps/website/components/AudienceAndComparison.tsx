"use client";
import { motion } from "framer-motion";
import { CheckCircle2, XCircle, AlertTriangle, Users, ArrowRightLeft } from "lucide-react";

export default function AudienceAndComparison() {
  return (
    <section id="audience-comparison" className="py-24 px-6 bg-[#0A0A0A] border-b-2 border-[#2A2A2A]">
      <div className="max-w-7xl mx-auto">
        
        {/* Header */}
        <div className="text-center mb-16">
          <span className="inline-flex items-center gap-1.5 border border-[#dc2626] text-[#dc2626] text-xs font-bold px-3 py-1 uppercase tracking-widest mb-4">
            ⚖️ Product Comparison
          </span>
          <h2 className="text-4xl md:text-6xl font-black text-white mb-4 tracking-tighter">
            THE <span className="text-gradient-yellow">TRANSFORMATION</span>
          </h2>
          <p className="text-neutral-400 text-lg max-w-2xl mx-auto font-medium">
            How AutoDevs shifts your development workflow from manual searching to autonomous shipping.
          </p>
        </div>

        {/* Double Column Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
          
          {/* Column 1: Before / After */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5 }}
            className="flex flex-col h-full"
          >
            <div className="flex items-center gap-2 mb-6">
              <ArrowRightLeft className="w-6 h-6 text-[#C8F135]" />
              <h3 className="text-2xl font-black text-white uppercase tracking-tight">
                Before & After AutoDevs
              </h3>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 flex-grow">
              {/* Before Card */}
              <div className="nb-card p-6 border-[#dc2626] hover:shadow-[4px_4px_0px_#dc2626] bg-[#111]/40">
                <div className="flex items-center gap-2 mb-4 pb-3 border-b border-[#2A2A2A]">
                  <XCircle className="w-5 h-5 text-[#dc2626]" />
                  <span className="text-sm font-mono font-bold text-[#dc2626] uppercase tracking-wider">
                    Before AutoDevs
                  </span>
                </div>
                <ul className="space-y-4">
                  {[
                    "Manual code reviews & tedious spot checks",
                    "Random TODOs scattered and forgotten in files",
                    "Constant context switching to inspect dependencies",
                    "Missing hidden bugs, config leaks, and vulnerabilities",
                  ].map((item, i) => (
                    <li key={i} className="flex gap-2.5 items-start">
                      <span className="text-[#dc2626] font-bold text-sm mt-0.5 select-none">❌</span>
                      <span className="text-sm font-medium text-neutral-400 leading-relaxed">
                        {item}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>

              {/* After Card */}
              <div className="nb-card p-6 border-[#00FF87] hover:shadow-[4px_4px_0px_#00FF87] bg-[#111]/40">
                <div className="flex items-center gap-2 mb-4 pb-3 border-b border-[#2A2A2A]">
                  <CheckCircle2 className="w-5 h-5 text-[#00FF87]" />
                  <span className="text-sm font-mono font-bold text-[#00FF87] uppercase tracking-wider">
                    After AutoDevs
                  </span>
                </div>
                <ul className="space-y-4">
                  {[
                    "Structured tasks prepared instantly on commands",
                    "AI recommendations generated from codebase scans",
                    "Faster, focused development pipelines",
                    "Better repositories, cleaner configs, and secure code",
                  ].map((item, i) => (
                    <li key={i} className="flex gap-2.5 items-start">
                      <span className="text-[#00FF87] font-bold text-sm mt-0.5 select-none">✅</span>
                      <span className="text-sm font-medium text-neutral-200 leading-relaxed font-semibold">
                        {item}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </motion.div>

          {/* Column 2: Who is this for? */}
          <motion.div
            initial={{ opacity: 0, y: 30 }}
            whileInView={{ opacity: 1, y: 0 }}
            viewport={{ once: true }}
            transition={{ duration: 0.5, delay: 0.2 }}
            className="flex flex-col h-full"
          >
            <div className="flex items-center gap-2 mb-6">
              <Users className="w-6 h-6 text-[#C8F135]" />
              <h3 className="text-2xl font-black text-white uppercase tracking-tight">
                Target Audience
              </h3>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 flex-grow">
              {/* Perfect For Card */}
              <div className="nb-card p-6 border-[#C8F135] hover:shadow-[4px_4px_0px_#C8F135] bg-[#111]/40">
                <div className="flex items-center gap-2 mb-4 pb-3 border-b border-[#2A2A2A]">
                  <CheckCircle2 className="w-5 h-5 text-[#C8F135]" />
                  <span className="text-sm font-mono font-bold text-[#C8F135] uppercase tracking-wider">
                    Perfect For
                  </span>
                </div>
                <ul className="space-y-4">
                  {[
                    "Students building complex projects",
                    "Open source contributors analyzing new repos",
                    "Startup developers moving at lightning speed",
                    "Freelancers handling multiple client codebases",
                    "Engineering teams maintaining codebase health",
                  ].map((item, i) => (
                    <li key={i} className="flex gap-2.5 items-start">
                      <span className="text-[#C8F135] font-bold text-sm mt-0.5 select-none">✅</span>
                      <span className="text-sm font-medium text-neutral-200 leading-relaxed">
                        {item}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>

              {/* Not For Card */}
              <div className="nb-card p-6 border-neutral-700 hover:shadow-[4px_4px_0px_#A3A3A3] bg-[#111]/10">
                <div className="flex items-center gap-2 mb-4 pb-3 border-b border-[#2A2A2A]">
                  <AlertTriangle className="w-5 h-5 text-neutral-500" />
                  <span className="text-sm font-mono font-bold text-neutral-500 uppercase tracking-wider">
                    Not For
                  </span>
                </div>
                <ul className="space-y-4">
                  {[
                    "Non-technical users who don't run terminals",
                    "No-code builders looking for drag-and-drop",
                    "Teams avoiding AI automation entirely",
                  ].map((item, i) => (
                    <li key={i} className="flex gap-2.5 items-start">
                      <span className="text-neutral-600 font-bold text-sm mt-0.5 select-none">❌</span>
                      <span className="text-sm font-medium text-neutral-500 leading-relaxed">
                        {item}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </motion.div>

        </div>

        {/* Highlight Quote Banner */}
        <div className="mt-16 bg-[#111] border-2 border-[#2A2A2A] p-6 text-center max-w-4xl mx-auto shadow-[6px_6px_0px_#C8F135] hover:border-[#C8F135] transition-colors">
          <p className="text-base md:text-lg font-mono font-bold text-[#C8F135]">
            &ldquo;AutoDevs is like Cursor + GitHub Analysis + AI Project Planning in one developer tool.&rdquo;
          </p>
        </div>

      </div>
    </section>
  );
}
