"use client";
import { motion } from "framer-motion";
import { Star, MessageSquare } from "lucide-react";

const testimonials = [
  {
    quote: "I can't imagine setting up a new repository without running a scan now. It finds the missing test coverages, security holes, and tells me exactly what to implement. AutoDevs is like Cursor + GitHub Analysis + AI Project Planning in one developer tool.",
    author: "Alex Rivera",
    role: "Open Source Maintainer",
    avatar: "💻",
    rating: 5,
  },
  {
    quote: "AutoDevs scans my repository, finds all the deprecated configurations and package issues, and structures them into actionable tasks. Slashed our onboarding time for junior devs by half.",
    author: "Sarah Chen",
    role: "Engineering Lead @ TechFlow",
    avatar: "🚀",
    rating: 5,
  },
  {
    quote: "As a freelancer managing multiple client codebases, autodev is an absolute lifesaver. Perfect for keeping track of all bugs and security issues, estimating the fix times, and getting them solved.",
    author: "Marc Dupont",
    role: "Freelance Fullstack Developer",
    avatar: "🦀",
    rating: 5,
  },
];

export default function Testimonials() {
  return (
    <section id="testimonials" className="py-24 px-6 bg-[#0A0A0A] border-b-2 border-[#2A2A2A]">
      <div className="max-w-7xl mx-auto">
        
        {/* Header */}
        <div className="text-center mb-16">
          <span className="inline-flex items-center gap-1.5 border border-[#FFD700] text-[#FFD700] text-xs font-bold px-3 py-1 uppercase tracking-widest mb-4">
            <MessageSquare className="w-3 h-3" /> Testimonials
          </span>
          <h2 className="text-4xl md:text-6xl font-black text-white mb-4 tracking-tighter">
            DEVELOPER <span className="text-gradient-yellow">LOVE</span>
          </h2>
          <p className="text-neutral-400 text-lg max-w-2xl mx-auto font-medium">
            Hear from developers, startup founders, and students shipping faster code with AutoDevs.
          </p>
        </div>

        {/* Grid layout */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {testimonials.map((t, idx) => (
            <motion.div
              key={idx}
              initial={{ opacity: 0, y: 20 }}
              whileInView={{ opacity: 1, y: 0 }}
              viewport={{ once: true }}
              transition={{ duration: 0.4, delay: idx * 0.15 }}
              className="nb-card p-6 flex flex-col justify-between hover:border-[#FFD700] hover:shadow-[6px_6px_0px_#FFD700] bg-[#111]/45 transition-all"
            >
              {/* Rating stars */}
              <div>
                <div className="flex gap-1 mb-4 text-[#FFD700]">
                  {[...Array(t.rating)].map((_, i) => (
                    <Star key={i} className="w-4 h-4 fill-current" />
                  ))}
                </div>

                <p className="text-sm md:text-base text-neutral-300 font-medium italic leading-relaxed mb-6">
                  &ldquo;{t.quote}&rdquo;
                </p>
              </div>

              {/* Author Info */}
              <div className="flex items-center gap-3 pt-4 border-t border-[#2A2A2A]">
                <div className="w-10 h-10 rounded-none border-2 border-[#2A2A2A] bg-black flex items-center justify-center text-lg">
                  {t.avatar}
                </div>
                <div>
                  <h4 className="text-sm font-black text-white uppercase tracking-tight">
                    {t.author}
                  </h4>
                  <p className="text-xs text-neutral-500 font-mono">
                    {t.role}
                  </p>
                </div>
              </div>
            </motion.div>
          ))}
        </div>

      </div>
    </section>
  );
}
