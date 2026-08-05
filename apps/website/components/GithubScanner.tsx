"use client";
import { useState, FormEvent } from "react";
import { motion } from "framer-motion";
import { Sparkles } from "lucide-react";

interface ScanData {
  username: string;
  totalRepos: number;
  languages: { [key: string]: number };
  detected: string[];
  recommended: string[];
  skills: string[];
}

const defaultData: ScanData = {
  username: "HEETMEHTA18",
  totalRepos: 12,
  languages: { TypeScript: 6, Go: 3, Python: 2, Dart: 1 },
  detected: ["React", "TypeScript", "Node.js", "Docker", "Firebase"],
  recommended: ["Node.js 22", "Go 1.22", "Docker", "Flutter SDK"],
  skills: ["Docker", "Kubernetes", "CI/CD", "Go"],
};

interface GithubRepo {
  name: string;
  fork: boolean;
  language: string | null;
  description: string | null;
  topics: string[];
}

function checkScannerCooldown(): { waitSecs: number | null; now: number } {
  if (typeof window === "undefined") {
    return { waitSecs: null, now: 0 };
  }
  const lastScan = localStorage.getItem("autodev_last_scan_time");
  const now = Date.now();
  if (lastScan && now - parseInt(lastScan, 10) < 10000) {
    // 10 second cooldown
    const waitSecs = Math.ceil((10000 - (now - parseInt(lastScan, 10))) / 1000);
    return { waitSecs, now };
  }
  return { waitSecs: null, now };
}

export default function GithubScanner() {
  const [username, setUsername] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [scanResult, setScanResult] = useState<ScanData | null>(defaultData);
  const [terminalLines, setTerminalLines] = useState<string[]>([
    `🐙 GitHub Scanner — @${defaultData.username}`,
    `  Fetching public repositories...`,
    `  Found 12 repositories (excluding forks).`,
    `  Detecting languages and frameworks...`,
    `  Recommended Environment compiled.`,
    `  Install all required tools? [Y/n] y`,
    `✓ Environment ready!`,
  ]);

  const handleScanSubmit = async (e: FormEvent) => {
    e.preventDefault();
    const targetUser = username.trim();
    if (!targetUser) return;

    setLoading(true);
    setError(null);
    setScanResult(null);
    setTerminalLines([
      `🐙 GitHub Scanner — @${targetUser}`,
      `  Initializing security check...`,
    ]);

    try {
      // 1. Sanitize GitHub username (GitHub username rules: alphanumeric and single hyphens, 1-39 chars)
      const githubUsernameRegex =
        /^[a-zA-Z0-9](?:[a-zA-Z0-9]|-(?=[a-zA-Z0-9])){0,38}$/;
      if (!githubUsernameRegex.test(targetUser)) {
        throw new Error("Invalid GitHub username format.");
      }

      // 2. Cooldown check (prevent bot spamming / API abuse)
      const { waitSecs, now } = checkScannerCooldown();
      if (waitSecs !== null) {
        throw new Error(
          `Rate limit: Please wait ${waitSecs}s before scanning again.`,
        );
      }
      if (typeof window !== "undefined") {
        localStorage.setItem("autodev_last_scan_time", now.toString());
      }

      await delay(200);
      setTerminalLines((prev) => [...prev, `  Connecting to GitHub API...`]);
      await delay(300);

      // 3. Fetch repositories (1st page, up to 100 repositories to optimize rate limits and speed)
      const res = await fetch(
        `https://api.github.com/users/${encodeURIComponent(targetUser)}/repos?per_page=100&type=public`,
      );

      if (!res.ok) {
        if (res.status === 404) {
          throw new Error("GitHub user not found");
        }
        if (res.status === 403) {
          throw new Error("GitHub rate limit exceeded. Try again later.");
        }
        throw new Error(`GitHub API error: ${res.statusText}`);
      }

      const repos = await res.json();

      // Simulate terminal progress steps for better immersion
      await delay(400);
      setTerminalLines((prev) => [
        ...prev,
        `  Fetching public repositories...`,
      ]);
      await delay(400);

      const allRepos = repos as GithubRepo[];
      if (allRepos.length === 0) {
        throw new Error("No public repositories found for this user.");
      }

      // 2. Count languages
      const languages: { [key: string]: number } = {};
      let totalRepos = 0;

      allRepos.forEach((repo) => {
        if (repo.fork) return; // skip forks
        totalRepos++;
        if (repo.language) {
          languages[repo.language] = (languages[repo.language] || 0) + 1;
        }
      });

      setTerminalLines((prev) => [
        ...prev,
        `  Found ${totalRepos} repositories (excluding forks).`,
      ]);
      await delay(450);
      setTerminalLines((prev) => [
        ...prev,
        `  Detecting languages and frameworks...`,
      ]);

      // Detect frameworks/runtimes
      const detectedSet = new Set<string>();
      const recommendedSet = new Set<string>();

      // Mapping rules
      const envMap: { [key: string]: string } = {
        JavaScript: "Node.js 22",
        TypeScript: "Node.js 22 + TypeScript",
        Python: "Python 3.12 + pip",
        Go: "Go 1.22",
        Rust: "Rust (rustup)",
        Java: "OpenJDK 21 + Maven/Gradle",
        Kotlin: "Kotlin + OpenJDK 21",
        PHP: "PHP 8.3 + Composer",
        Ruby: "Ruby 3.3 + Bundler",
        Dart: "Flutter SDK + Dart",
        Swift: "Xcode + Swift Toolchain",
        "C#": ".NET 8 SDK",
        "C++": "GCC / Clang + CMake",
        Shell: "Bash / Zsh",
        HTML: "Web Browser",
        CSS: "Web Browser",
      };

      Object.keys(languages).forEach((lang) => {
        detectedSet.add(lang);
        if (envMap[lang]) {
          recommendedSet.add(envMap[lang]);
        }
      });

      // Check topics and descriptions for other tech/frameworks
      const frameworksMap: { [key: string]: string[] } = {
        react: ["React", "Node.js 22"],
        nextjs: ["Next.js", "Node.js 22"],
        vue: ["Vue", "Node.js 22"],
        angular: ["Angular", "Node.js 22"],
        svelte: ["Svelte", "Node.js 22"],
        django: ["Django", "Python 3.12 + pip"],
        fastapi: ["FastAPI", "Python 3.12 + pip"],
        laravel: ["Laravel", "PHP 8.3 + Composer"],
        "spring-boot": ["Spring Boot", "OpenJDK 21 + Maven/Gradle"],
        springboot: ["Spring Boot", "OpenJDK 21 + Maven/Gradle"],
        flutter: ["Flutter", "Flutter SDK + Dart"],
        docker: ["Docker", "Docker Desktop"],
        kubernetes: ["Kubernetes", "Docker Desktop"],
        k8s: ["Kubernetes", "Docker Desktop"],
      };

      allRepos.forEach((repo) => {
        const desc = (repo.description || "").toLowerCase();
        const name = repo.name.toLowerCase();
        const topics = (repo.topics || []).map((t: string) => t.toLowerCase());

        Object.keys(frameworksMap).forEach((key) => {
          if (
            topics.includes(key) ||
            name.includes(key) ||
            desc.includes(key)
          ) {
            const [frameworkName, runtimeName] = frameworksMap[key];
            detectedSet.add(frameworkName);
            recommendedSet.add(runtimeName);
          }
        });
      });

      // Build skill gaps
      const gaps: { [key: string]: string[] } = {
        JavaScript: ["TypeScript", "Docker", "CI/CD"],
        TypeScript: ["Docker", "Kubernetes", "Go"],
        Python: ["Docker", "FastAPI", "Celery"],
        Go: ["Kubernetes", "gRPC", "Terraform"],
        Java: ["Spring Boot", "Docker", "Kubernetes"],
        Rust: ["WebAssembly", "async-std", "tokio"],
      };

      const skillsSet = new Set<string>();
      Object.keys(languages).forEach((lang) => {
        if (gaps[lang]) {
          gaps[lang].forEach((s) => {
            if (!detectedSet.has(s)) {
              skillsSet.add(s);
            }
          });
        }
      });

      if (skillsSet.size === 0) {
        skillsSet.add("Docker");
        skillsSet.add("Kubernetes");
        skillsSet.add("CI/CD");
      }

      const detected = Array.from(detectedSet).slice(0, 8);
      const recommended = Array.from(recommendedSet).slice(0, 6);
      const skills = Array.from(skillsSet).slice(0, 5);

      const result: ScanData = {
        username: targetUser,
        totalRepos,
        languages,
        detected,
        recommended,
        skills,
      };

      await delay(450);
      setTerminalLines((prev) => [
        ...prev,
        `  Recommended Environment compiled.`,
      ]);
      await delay(400);
      setTerminalLines((prev) => [
        ...prev,
        `  Install all required tools? [Y/n] y`,
      ]);
      await delay(300);
      setTerminalLines((prev) => [...prev, `✓ Environment ready!`]);

      setScanResult(result);
    } catch (err: unknown) {
      console.error(err);
      const errMsg =
        err instanceof Error ? err.message : "Something went wrong";
      setError(errMsg);
      setTerminalLines((prev) => [
        ...prev,
        `❌ Error: ${errMsg}`,
        `Please try again later.`,
      ]);
    } finally {
      setLoading(false);
    }
  };

  const delay = (ms: number) =>
    new Promise((resolve) => setTimeout(resolve, ms));

  // Determine which data to render
  const data = scanResult || defaultData;
  const showStats = !loading && data;

  return (
    <section id="github-scanner" className="py-24 px-6 max-w-7xl mx-auto">
      <motion.div
        className="mb-16"
        initial={{ opacity: 0, y: 30 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
      >
        <span className="inline-flex items-center gap-1.5 text-xs text-[#C8F135] font-bold uppercase tracking-widest font-mono">
          <Sparkles className="w-3 h-3" /> GitHub Scanner
        </span>
        <h2 className="text-5xl font-black text-white mt-2 mb-4 font-mono">
          SCAN ANY <span className="text-gradient-yellow">GITHUB</span> USER
        </h2>
        <p className="text-neutral-400 max-w-xl font-mono text-sm">
          AutoDev reads all public repos, detects languages, and generates a
          setup plan.
        </p>
      </motion.div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-start">
        {/* Input */}
        <div>
          <form onSubmit={handleScanSubmit} className="flex gap-0 mb-4">
            <div className="border-2 border-r-0 border-[#2A2A2A] px-4 py-3 text-neutral-400 font-mono text-sm bg-[#111] whitespace-nowrap select-none">
              autodev github
            </div>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="USERNAME"
              disabled={loading}
              className="flex-1 border-2 border-[#C8F135] bg-black text-white font-mono text-sm px-4 py-3 outline-none focus:shadow-[4px_4px_0_#C8F135] transition-shadow placeholder:text-neutral-500 min-w-0"
            />
            <button
              type="submit"
              disabled={loading}
              className="nb-btn px-6 py-3 text-sm disabled:opacity-50 whitespace-nowrap select-none font-mono"
            >
              {loading ? "SCANNING..." : "SCAN →"}
            </button>
          </form>
          <p className="text-xs text-neutral-400 font-mono">
            Or run:{" "}
            <span className="text-[#C8F135]">
              autodev github {username || "USERNAME"}
            </span>
          </p>

          {/* Language bars */}
          <div className="mt-8 border-2 border-[#2A2A2A] p-5">
            <h3 className="font-bold text-white mb-4 text-sm uppercase tracking-wider font-mono">
              Languages Detected
            </h3>
            {showStats && Object.keys(data.languages).length > 0 ? (
              Object.entries(data.languages)
                .sort((a, b) => b[1] - a[1])
                .slice(0, 6)
                .map(([lang, count]) => {
                  const pct = Math.round((count / data.totalRepos) * 100);
                  return (
                    <div key={lang} className="mb-3">
                      <div className="flex justify-between text-xs mb-1 font-mono">
                        <span className="text-neutral-400">{lang}</span>
                        <span className="text-neutral-400">{count} repos</span>
                      </div>
                      <div className="h-2 bg-[#1A1A1A] border border-[#2A2A2A]">
                        <motion.div
                          key={`${data.username}-${lang}`} // forces re-render/animation on new search
                          initial={{ width: 0 }}
                          animate={{ width: `${Math.min(pct, 100)}%` }}
                          transition={{ duration: 0.6, delay: 0.1 }}
                          className="h-full bg-[#C8F135]"
                        />
                      </div>
                    </div>
                  );
                })
            ) : loading ? (
              <div className="py-8 text-center text-neutral-400 font-mono text-xs animate-pulse">
                Analyzing repositories...
              </div>
            ) : (
              <div className="py-8 text-center text-neutral-400 font-mono text-xs">
                No languages detected.
              </div>
            )}
          </div>
        </div>

        {/* Output terminal */}
        <div className="terminal border-2 border-[#C8F135] bg-black">
          <div className="terminal-bar border-b-2 border-[#2A2A2A] px-4 py-2 flex items-center bg-[#111]">
            <span className="terminal-dot bg-[#FF5F56] w-3 h-3 rounded-full mr-1.5 inline-block" />
            <span className="terminal-dot bg-[#FFBD2E] w-3 h-3 rounded-full mr-1.5 inline-block" />
            <span className="terminal-dot bg-[#27C93F] w-3 h-3 rounded-full mr-1.5 inline-block" />
            <span className="text-xs text-neutral-400 ml-3 font-mono select-none">
              autodev github {loading ? username : data.username}
            </span>
          </div>
          <div className="px-6 py-5 font-mono text-sm space-y-1 min-h-[360px] flex flex-col justify-between">
            <div className="space-y-1">
              {terminalLines.map((line, idx) => {
                const isError = line.startsWith("❌");
                const isSuccess = line.startsWith("✓");
                let colorClass = "text-neutral-400";
                if (isError) colorClass = "text-[#FF5F56]";
                else if (isSuccess) colorClass = "text-[#00FF87]";
                else if (idx === 0) colorClass = "text-[#00FF87] font-bold";

                return (
                  <div key={idx} className={colorClass}>
                    {line}
                  </div>
                );
              })}
            </div>

            {/* Complete output breakdown (shown when scan completes successfully) */}
            {showStats && !error && (
              <div className="mt-4 pt-4 border-t border-[#1F1F1F] space-y-3">
                <div>
                  <div className="text-[#2563eb] text-xs uppercase tracking-wider mb-1 font-bold">
                    Technologies Detected:
                  </div>
                  <div className="flex flex-wrap gap-1.5 pl-1">
                    {data.detected.map((t) => (
                      <span
                        key={t}
                        className="text-xs bg-[#111] border border-[#2A2A2A] px-2 py-0.5 text-[#00FF87]"
                      >
                        • {t}
                      </span>
                    ))}
                  </div>
                </div>

                <div>
                  <div className="text-[#2563eb] text-xs uppercase tracking-wider mb-1 font-bold">
                    Recommended Tools:
                  </div>
                  <div className="flex flex-wrap gap-1.5 pl-1">
                    {data.recommended.map((r) => (
                      <span
                        key={r}
                        className="text-xs bg-[#111] border border-[#2A2A2A] px-2 py-0.5 text-[#C8F135]"
                      >
                        → {r}
                      </span>
                    ))}
                  </div>
                </div>

                <div>
                  <div className="text-[#2563eb] text-xs uppercase tracking-wider mb-1 font-bold">
                    Suggested Skills:
                  </div>
                  <div className="flex flex-wrap gap-1.5 pl-1">
                    {data.skills.map((s) => (
                      <span
                        key={s}
                        className="text-xs bg-[#111] border border-[#2A2A2A] px-2 py-0.5 text-neutral-400"
                      >
                        ↗ {s}
                      </span>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
