"use client";
import React, { useEffect, useRef, useState } from "react";
import { usePathname } from "next/navigation";
import { useTheme } from "next-themes";

interface AsciiBackgroundProps {
  className?: string;
  speed?: number;
  density?: number;
  mouseSensi?: number; // 1 to 5, scales mouse radius and distortion reaction
}

interface FloatingText {
  text: string;
  c: number;
  r: number;
  life: number;
  maxLife: number;
}

interface TextSegment {
  char: string;
  opacity: number;
}

export default function AsciiBackground({
  className = "",
  speed = 1,
  density = 1,
  mouseSensi = 1.5,
}: AsciiBackgroundProps) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const pathname = usePathname();
  const { theme } = useTheme();
  const [mounted, setMounted] = useState(false);

  // Only run this animation on the landing page
  const isLandingPage = pathname === "/";

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    setMounted(true);
  }, []);

  useEffect(() => {
    if (!mounted || !isLandingPage) return;

    const canvas = canvasRef.current;
    if (!canvas) return;

    const ctx = canvas.getContext("2d", { alpha: true });
    if (!ctx) return;

    // Respect users who prefer reduced motion: render a single static frame
    const prefersReducedMotion =
      typeof window !== "undefined" &&
      window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    // Cap device pixel ratio to keep rasterization cheap on HiDPI screens
    const MAX_DPR = 1.5;

    let animationId: number;
    let width = 0;
    let height = 0;

    // Grid details
    const fontSize = 10;
    const charWidth = 10;
    const charHeight = 15;

    // ASCII Characters list from sparse/dim to dense/bright
    const asciiChars = [" ", ".", ":", "-", "=", "+", "*", "%", "#", "@"];
    
    // Telemetry items that float around
    const telemetryItems = [
      "DEPTH: -2450m",
      "PRESSURE: 24.5 MPa",
      "KELP_YIELD: 98.4%",
      "BALLAST: ACTIVE",
      "TIDAL_COEFF: 0.74",
      "SCANNER_SWEEP_HZ: 45.2",
      "TEMP: 104.5 C",
      "NODE_A4: STABLE",
      "TOMATO_GEN_0: OK",
      "go run scanner/main.go",
      "autodev scan .",
      "autodev doctor --fix",
      "autodev setup --yes",
      "SONAR: PING_SENT",
    ];

    const floatingTexts: FloatingText[] = [];

    // Mouse coordinates and velocity tracking
    let mouseX = -9999;
    let mouseY = -9999;
    let targetMouseX = -9999;
    let targetMouseY = -9999;
    let mouseSpeed = 0;
    let lastMouseX = -9999;
    let lastMouseY = -9999;

    const handleResize = () => {
      const rect = canvas.getBoundingClientRect();
      const dpr = Math.min(window.devicePixelRatio || 1, MAX_DPR);
      width = rect.width;
      height = rect.height;
      canvas.width = Math.round(rect.width * dpr);
      canvas.height = Math.round(rect.height * dpr);
      ctx.scale(dpr, dpr);
    };

    window.addEventListener("resize", handleResize);
    handleResize();

    const handleMouseMove = (e: MouseEvent) => {
      const rect = canvas.getBoundingClientRect();
      targetMouseX = e.clientX - rect.left;
      targetMouseY = e.clientY - rect.top;

      if (lastMouseX !== -9999) {
        const dx = targetMouseX - lastMouseX;
        const dy = targetMouseY - lastMouseY;
        const instantSpeed = Math.sqrt(dx * dx + dy * dy);
        // Smooth out the speed value
        mouseSpeed = mouseSpeed * 0.85 + instantSpeed * 0.15;
      }
      lastMouseX = targetMouseX;
      lastMouseY = targetMouseY;
    };

    const handleMouseLeave = () => {
      targetMouseX = -9999;
      targetMouseY = -9999;
      lastMouseX = -9999;
      lastMouseY = -9999;
      mouseSpeed = 0;
    };

    window.addEventListener("mousemove", handleMouseMove);
    canvas.addEventListener("mouseleave", handleMouseLeave);

    // Pause expensive rendering when the tab is hidden or the canvas is far
    // outside the viewport. This is the single biggest win for scroll smoothness.
    let isHidden = false;
    const handleVisibility = () => {
      isHidden = document.hidden;
    };
    document.addEventListener("visibilitychange", handleVisibility);

    let observer: IntersectionObserver | null = null;
    if (typeof IntersectionObserver !== "undefined") {
      observer = new IntersectionObserver(
        (entries) => {
          isHidden = !entries[0].isIntersecting;
        },
        { rootMargin: "400px 0px" },
      );
      observer.observe(canvas);
    }

    let frame = 0;
    let lastTime = performance.now();

    const drawFrame = (time: number) => {
      const dt = Math.min(4, (time - lastTime) / 16.666); // Cap dt to avoid huge jumps on tab wake
      lastTime = time;

      frame += speed * dt;
      ctx.clearRect(0, 0, width, height);

      // Decelerate mouse velocity feedback using dt
      mouseSpeed *= Math.pow(0.96, dt);

      // Frame-rate independent LERP for mouse coordinates
      if (targetMouseX === -9999) {
        mouseX = -9999;
        mouseY = -9999;
      } else {
        if (mouseX === -9999) {
          mouseX = targetMouseX;
          mouseY = targetMouseY;
        } else {
          const lerpFactor = 1 - Math.pow(1 - 0.15, dt);
          mouseX += (targetMouseX - mouseX) * lerpFactor;
          mouseY += (targetMouseY - mouseY) * lerpFactor;
        }
      }

      ctx.font = `${fontSize}px var(--font-mono), monospace`;

      // Determine color palette based on active theme
      const isDark = theme !== "light";
      const baseColor = isDark ? "200, 241, 53" : "148, 178, 32"; // Lime (dark) / Darker lime (light)

      const cols = Math.ceil(width / charWidth);
      const rows = Math.ceil(height / charHeight);

      // Manage floating telemetry text using dt for life ticks
      if (Math.random() < 0.008 * dt && floatingTexts.length < 5) {
        const text = telemetryItems[Math.floor(Math.random() * telemetryItems.length)];
        const length = text.length;
        const c = Math.max(1, Math.floor(Math.random() * (cols - length - 2)));
        const r = Math.max(1, Math.floor(Math.random() * (rows - 2)));
        floatingTexts.push({
          text,
          c,
          r,
          life: 0,
          maxLife: 200 + Math.random() * 150,
        });
      }

      // Filter out dead floating texts
      for (let i = floatingTexts.length - 1; i >= 0; i--) {
        floatingTexts[i].life += speed * dt;
        if (floatingTexts[i].life >= floatingTexts[i].maxLife) {
          floatingTexts.splice(i, 1);
        }
      }

      // Create a map of positions occupied by floating texts
      const textOverlayMap = new Map<string, TextSegment>();
      floatingTexts.forEach((ft) => {
        const textOpacity = Math.sin((ft.life / ft.maxLife) * Math.PI); // fade in and out
        for (let idx = 0; idx < ft.text.length; idx++) {
          const key = `${ft.c + idx},${ft.r}`;
          textOverlayMap.set(key, { char: ft.text[idx], opacity: textOpacity });
        }
      });

      // Horizontal sweep line position (circular scanner/sweep style)
      const sweepCol = (frame * 0.25) % (cols + 50) - 25;

      for (let r = 0; r < rows; r++) {
        const y = r * charHeight + fontSize;
        for (let c = 0; c < cols; c++) {
          const x = c * charWidth;

          // Check if this position is occupied by floating telemetry text
          const key = `${c},${r}`;
          const textItem = textOverlayMap.get(key);

          if (textItem) {
            ctx.fillStyle = `rgba(${baseColor}, ${textItem.opacity * 0.75 * density})`;
            ctx.fillText(textItem.char, x, y);
            continue;
          }

          // Technical grid structure: Draw a plus sign (+) at grid intersections
          if (c % 16 === 0 && r % 4 === 0) {
            ctx.fillStyle = isDark ? `rgba(${baseColor}, 0.08)` : `rgba(${baseColor}, 0.12)`;
            ctx.fillText("+", x, y);
            continue;
          }

          // Sonar scan line effect
          const distToSweep = Math.abs(c - sweepCol);
          let sweepInfluence = 0;
          if (distToSweep < 15) {
            sweepInfluence = 1 - distToSweep / 15;
          }

          // Organic wave effect using composite sine/cosine waves
          const wave1 = Math.sin(c * 0.08 + frame * 0.02) * Math.cos(r * 0.08 + frame * 0.015);
          const wave2 = Math.sin(c * 0.03 - frame * 0.01) * Math.cos(r * 0.05 + frame * 0.015);
          const wave = (wave1 + wave2) * 0.5;

          // Normalize wave value to [0, 1]
          let normWave = (wave + 1) * 0.5;

          // Distance-based mouse distortion with speed-dependent dynamic radius
          // Compare squared distances first to avoid a Math.sqrt for most cells
          const baseRadius = 120 * mouseSensi;
          const dynamicRadius = baseRadius + Math.min(120, mouseSpeed * 2.5 * mouseSensi);
          let mouseInfluence = 0;
          if (mouseX !== -9999) {
            const dx = x - mouseX;
            const dy = y - mouseY;
            const radiusSq = dynamicRadius * dynamicRadius;
            const distSq = dx * dx + dy * dy;
            if (distSq < radiusSq) {
              mouseInfluence = 1 - Math.sqrt(distSq) / dynamicRadius;
              // Magnify wave disruption when mouse moves fast
              const distortionFactor = 0.95 + mouseSpeed * 0.01;
              normWave = normWave * (1 - mouseInfluence) + mouseInfluence * distortionFactor;
            }
          }

          // Incorporate sweep line scanning intensity
          normWave = normWave * (1 - sweepInfluence * 0.4) + sweepInfluence * 0.4 * 0.85;

          // Choose character
          const charIndex = Math.min(
            asciiChars.length - 1,
            Math.floor(normWave * asciiChars.length),
          );
          const char = asciiChars[charIndex];

          if (char !== " ") {
            // Brighter baseline opacity
            let opacity = normWave * 0.12 * density;
            if (mouseInfluence > 0) {
              opacity = (normWave * 0.12 + mouseInfluence * 0.28) * density;
            }
            if (sweepInfluence > 0) {
              opacity = Math.max(opacity, sweepInfluence * 0.22 * density);
            }

            // Restrict bounds (increased maximum opacity limit for brighter highlighting)
            opacity = Math.max(0.01, Math.min(0.65, opacity));

            ctx.fillStyle = `rgba(${baseColor}, ${opacity})`;

            // Add subtle fluid animation offset
            const offsetX = Math.sin(r * 0.4 + frame * 0.03) * 1.5;
            const offsetY = Math.cos(c * 0.4 + frame * 0.03) * 1.5;

            ctx.fillText(char, x + offsetX, y + offsetY);
          }
        }
      }
    };

    const render = (time: number) => {
      animationId = requestAnimationFrame(render);
      if (isHidden) return;
      drawFrame(time);
    };

    if (prefersReducedMotion) {
      // Static frame — no animation loop for users who prefer reduced motion
      drawFrame(performance.now());
      animationId = 0;
    } else {
      animationId = requestAnimationFrame(render);
    }

    return () => {
      document.removeEventListener("visibilitychange", handleVisibility);
      if (observer) observer.disconnect();
      window.removeEventListener("resize", handleResize);
      window.removeEventListener("mousemove", handleMouseMove);
      if (canvas) {
        canvas.removeEventListener("mouseleave", handleMouseLeave);
      }
      cancelAnimationFrame(animationId);
    };
  }, [mounted, isLandingPage, theme, speed, density, mouseSensi]);

  if (!isLandingPage) return null;

  const isDark = theme !== "light";

  return (
    <canvas
      ref={canvasRef}
      className={`absolute inset-0 w-full h-full pointer-events-none select-none z-0 ${className}`}
      style={{ mixBlendMode: isDark ? "screen" : "multiply" }}
    />
  );
}
