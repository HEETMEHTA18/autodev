"use client";
import { motion, useInView } from "framer-motion";
import { useRef } from "react";

interface SectionDividerProps {
  variant?: "wave" | "angle" | "curve";
  flip?: boolean;
}

export default function SectionDivider({ variant = "wave", flip = false }: SectionDividerProps) {
  const ref = useRef<HTMLDivElement>(null);
  const isInView = useInView(ref, { once: true });

  const paths = {
    wave: "M0,40 C40,10 60,70 120,40 C180,10 200,70 240,40 C300,10 320,70 360,40 L360,80 L0,80 Z",
    angle: "M0,60 L180,20 L360,60 L360,80 L0,80 Z",
    curve: "M0,60 Q180,10 360,60 L360,80 L0,80 Z",
  };

  return (
    <div ref={ref} className="section-divider">
      <motion.div
        initial={{ opacity: 0, scaleY: 0 }}
        animate={isInView ? { opacity: 1, scaleY: 1 } : {}}
        transition={{ duration: 0.5, ease: "easeOut" }}
        style={{ transformOrigin: flip ? "top" : "bottom" }}
      >
        <svg
          viewBox="0 0 360 80"
          fill="none"
          preserveAspectRatio="none"
          className={flip ? "rotate-180" : ""}
        >
          <path d={paths[variant]} fill="#0D0D0D" />
        </svg>
      </motion.div>
    </div>
  );
}
