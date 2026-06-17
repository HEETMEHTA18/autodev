"use client";
import React, { useEffect, useRef, useState } from "react";
import { usePathname } from "next/navigation";
import { useTheme } from "next-themes";

interface AsciiBackgroundProps {
  className?: string;
  speed?: number;
  density?: number;
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

    const ctx = canvas.getContext("2d");
    if (!ctx) return;

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

    // Mouse coordinates tracking
    let mouseX = -9999;
    let mouseY = -9999;
    let targetMouseX = -9999;
    let targetMouseY = -9999;

    const handleResize = () => {
      const rect = canvas.getBoundingClientRect();
      width = canvas.width = rect.width;
      height = canvas.height = rect.height;
    };

    window.addEventListener("resize", handleResize);
    handleResize();

    const handleMouseMove = (e: MouseEvent) => {
      const rect = canvas.getBoundingClientRect();
      targetMouseX = e.clientX - rect.left;
      targetMouseY = e.clientY - rect.top;
    };

    const handleMouseLeave = () => {
      targetMouseX = -9999;
      targetMouseY = -9999;
    };

    window.addEventListener("mousemove", handleMouseMove);
    canvas.addEventListener("mouseleave", handleMouseLeave);

    let frame = 0;

    const render = () => {
      frame += speed;
      ctx.clearRect(0, 0, width, height);

      // Smooth lag interpolation for mouse coordinates
      if (targetMouseX === -9999) {
        mouseX = -9999;
        mouseY = -9999;
      } else {
        if (mouseX === -9999) {
          mouseX = targetMouseX;
          mouseY = targetMouseY;
        } else {
          mouseX += (targetMouseX - mouseX) * 0.15;
          mouseY += (targetMouseY - mouseY) * 0.15;
        }
      }

      ctx.font = `${fontSize}px var(--font-mono), monospace`;
      
      // Determine color palette based on active theme
      const isDark = theme === "dark" || !theme;
      const baseColor = isDark ? "255, 215, 0" : "180, 140, 10"; // Gold/Yellow (autodev brand color)
      
      const cols = Math.ceil(width / charWidth);
      const rows = Math.ceil(height / charHeight);

      // Manage floating telemetry text
      if (frame % 150 === 0 && floatingTexts.length < 5) {
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
        floatingTexts[i].life += speed;
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
            ctx.fillStyle = `rgba(${baseColor}, ${textItem.opacity * 0.45 * density})`;
            ctx.fillText(textItem.char, x, y);
            continue;
          }

          // Technical grid structure: Draw a plus sign (+) at grid intersections
          if (c % 16 === 0 && r % 4 === 0) {
            ctx.fillStyle = `rgba(${baseColor}, 0.04)`;
            ctx.fillText("+", x, y);
            continue;
          }

          // Sonar scan line effect
          const distToSweep = Math.abs(c - sweepCol);
          let sweepInfluence = 0;
          if (distToSweep < 15) {
            sweepInfluence = (1 - distToSweep / 15);
          }

          // Organic wave effect using composite sine/cosine waves
          const wave1 = Math.sin(c * 0.08 + frame * 0.02) * Math.cos(r * 0.08 + frame * 0.015);
          const wave2 = Math.sin(c * 0.03 - frame * 0.01) * Math.cos(r * 0.05 + frame * 0.015);
          const wave = (wave1 + wave2) * 0.5;

          // Normalize wave value to [0, 1]
          let normWave = (wave + 1) * 0.5;

          // Distance-based mouse distortion
          let distance = 9999;
          if (mouseX !== -9999) {
            const dx = x - mouseX;
            const dy = y - mouseY;
            distance = Math.sqrt(dx * dx + dy * dy);
          }

          const mouseRadius = 120;
          let mouseInfluence = 0;
          if (distance < mouseRadius) {
            mouseInfluence = (1 - distance / mouseRadius);
            normWave = normWave * (1 - mouseInfluence) + mouseInfluence * 0.95;
          }

          // Incorporate sweep line scanning intensity
          normWave = normWave * (1 - sweepInfluence * 0.4) + sweepInfluence * 0.4 * 0.85;

          // Choose character
          const charIndex = Math.min(
            asciiChars.length - 1,
            Math.floor(normWave * asciiChars.length)
          );
          const char = asciiChars[charIndex];

          if (char !== " ") {
            // Calculate opacity based on density, wave, mouse and sweep presence
            let opacity = normWave * 0.06 * density;
            if (mouseInfluence > 0) {
              opacity = (normWave * 0.06 + mouseInfluence * 0.18) * density;
            }
            if (sweepInfluence > 0) {
              opacity = Math.max(opacity, sweepInfluence * 0.12 * density);
            }
            
            // Restrict bounds
            opacity = Math.max(0.005, Math.min(0.35, opacity));

            ctx.fillStyle = `rgba(${baseColor}, ${opacity})`;
            
            // Add subtle fluid animation offset
            const offsetX = Math.sin(r * 0.4 + frame * 0.03) * 1.5;
            const offsetY = Math.cos(c * 0.4 + frame * 0.03) * 1.5;

            ctx.fillText(char, x + offsetX, y + offsetY);
          }
        }
      }

      animationId = requestAnimationFrame(render);
    };

    render();

    return () => {
      window.removeEventListener("resize", handleResize);
      window.removeEventListener("mousemove", handleMouseMove);
      if (canvas) {
        canvas.removeEventListener("mouseleave", handleMouseLeave);
      }
      cancelAnimationFrame(animationId);
    };
  }, [mounted, isLandingPage, theme, speed, density]);

  if (!isLandingPage) return null;

  return (
    <canvas
      ref={canvasRef}
      className={`absolute inset-0 w-full h-full pointer-events-none select-none z-0 ${className}`}
      style={{ mixBlendMode: "screen" }}
    />
  );
}
