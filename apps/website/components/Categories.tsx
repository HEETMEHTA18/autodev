"use client";
import { useState, useMemo } from "react";
import { motion, AnimatePresence } from "framer-motion";
import {
  Search,
  Check,
  X,
  Terminal as TerminalIcon,
  Info,
  Sparkles,
  Layers,
  CheckCircle,
  Copy,
} from "lucide-react";

// Package definition
interface DevPackage {
  id: string;
  name: string;
  icon: string;
  desc: string;
  longDesc: string;
  version: string;
  size: string;
  category: string;
  developer: string;
  license: string;
  dependencies: string[];
  snippet: string;
}

// Categories definitions
const categories = [
  { id: "all", icon: "🌐", label: "All Items" },
  { id: "languages", icon: "🔤", label: "Languages" },
  { id: "databases", icon: "🗄️", label: "Databases" },
  { id: "devops", icon: "⚙️", label: "DevOps / Containers" },
  { id: "tools", icon: "⚡", label: "Editors & CLI Tools" },
  { id: "aiml", icon: "🧠", label: "AI / ML" },
];

const packages: DevPackage[] = [
  // Languages
  {
    id: "python",
    name: "Python",
    icon: "🐍",
    category: "languages",
    desc: "Dynamic, object-oriented programming language.",
    longDesc:
      "Python is an interpreted, high-level, general-purpose programming language. Its design philosophy emphasizes code readability with its use of significant indentation. Its language constructs as well as its object-oriented approach aim to help programmers write clear, logical code for small and large-scale projects.",
    version: "3.12.3",
    size: "28.4 MB",
    developer: "Python Software Foundation",
    license: "PSF License",
    dependencies: ["openssl", "sqlite", "zlib"],
    snippet: "python --version && python -c 'print(\"Hello from AutoDev!\")'",
  },
  {
    id: "nodejs",
    name: "Node.js",
    icon: "🟢",
    category: "languages",
    desc: "JavaScript runtime built on Chrome's V8 engine.",
    longDesc:
      "Node.js is an open-source, cross-platform JavaScript runtime environment that executes JavaScript code outside a web browser. It allows developers to use JavaScript to write command line tools and for server-side scripting—running scripts server-side to produce dynamic web page content before the page is sent to the user's web browser.",
    version: "20.12.2 (LTS)",
    size: "32.1 MB",
    developer: "OpenJS Foundation",
    license: "MIT",
    dependencies: ["libuv", "openssl", "v8"],
    snippet: "node -v && node -e 'console.log(\"Node is up and running!\")'",
  },
  {
    id: "go",
    name: "Go",
    icon: "🔵",
    category: "languages",
    desc: "Open source programming language developed by Google.",
    longDesc:
      "Go (also known as Golang) is a statically typed, compiled programming language designed at Google. Go is syntactically similar to C, but with memory safety, garbage collection, structural typing, and CSP-style concurrency.",
    version: "1.22.2",
    size: "64.8 MB",
    developer: "Google & Contributors",
    license: "BSD 3-Clause",
    dependencies: ["libc", "libpthread"],
    snippet: "go version && go run main.go",
  },
  {
    id: "rust",
    name: "Rust",
    icon: "🦀",
    category: "languages",
    desc: "Language empowering everyone to build reliable and efficient software.",
    longDesc:
      "Rust is a multi-paradigm, general-purpose programming language designed for performance and safety, especially safe concurrency. Rust is syntactically similar to C++, but can guarantee memory safety by using a borrow checker to validate references.",
    version: "1.77.2",
    size: "142.5 MB",
    developer: "Rust Foundation",
    license: "MIT / Apache 2.0",
    dependencies: ["gcc", "llvm", "make"],
    snippet: "rustc --version && cargo --version",
  },
  // Databases
  {
    id: "postgresql",
    name: "PostgreSQL",
    icon: "🐘",
    category: "databases",
    desc: "Powerful, open source object-relational database system.",
    longDesc:
      "PostgreSQL is a free and open-source relational database management system emphasizing extensibility and SQL compliance. It is known for its proven architecture, reliability, data integrity, robust feature set, extensibility, and the dedication of the open-source community behind it.",
    version: "16.2",
    size: "45.0 MB",
    developer: "PostgreSQL Global Development Group",
    license: "PostgreSQL License",
    dependencies: ["openssl", "readline", "zlib"],
    snippet: "psql --version && pg_ctl status",
  },
  {
    id: "redis",
    name: "Redis",
    icon: "🔴",
    category: "databases",
    desc: "In-memory data structure store used as a database, cache, and message broker.",
    longDesc:
      "Redis is an in-memory data structure store, used as a distributed, in-memory key-value database, cache and message broker, with optional durability. Redis supports different kinds of abstract data structures such as strings, lists, maps, sets, sorted sets, HyperLogLogs, bitmaps, streams, and spatial indexes.",
    version: "7.2.4",
    size: "6.2 MB",
    developer: "Redis Ltd. & Community",
    license: "RSALv2 / SSPLv1",
    dependencies: ["libc", "openssl"],
    snippet: "redis-cli ping",
  },
  {
    id: "mongodb",
    name: "MongoDB",
    icon: "🍃",
    category: "databases",
    desc: "Document-based, distributed database designed for modern apps.",
    longDesc:
      "MongoDB is a source-available cross-platform document-oriented database program. Classified as a NoSQL database program, MongoDB uses JSON-like documents with optional schemas. MongoDB is developed by MongoDB Inc. and licensed under the Server Side Public License (SSPL).",
    version: "7.0.8",
    size: "95.2 MB",
    developer: "MongoDB Inc.",
    license: "SSPL",
    dependencies: ["openssl", "curl", "snappy"],
    snippet: "mongod --version && mongosh",
  },
  // DevOps
  {
    id: "docker",
    name: "Docker CLI",
    icon: "🐳",
    category: "devops",
    desc: "Pack, ship and run any application as a lightweight container.",
    longDesc:
      "Docker is a set of platform as a service products that use OS-level virtualization to deliver software in packages called containers. Containers are isolated from one another and bundle their own software, libraries and configuration files; they can communicate with each other through well-defined channels.",
    version: "26.0.1",
    size: "38.6 MB",
    developer: "Docker Inc.",
    license: "Apache 2.0",
    dependencies: ["containerd", "runc", "iptables"],
    snippet: "docker --version && docker ps",
  },
  {
    id: "kubectl",
    name: "kubectl",
    icon: "☸️",
    category: "devops",
    desc: "Command line tool for controlling Kubernetes clusters.",
    longDesc:
      "Kubectl is the command line utility for communicating with a Kubernetes cluster's control plane, using the Kubernetes API. It allows you to run commands against Kubernetes clusters to deploy applications, inspect and manage cluster resources, and view logs.",
    version: "1.29.3",
    size: "18.4 MB",
    developer: "CNCF",
    license: "Apache 2.0",
    dependencies: ["glibc"],
    snippet: "kubectl version --client",
  },
  {
    id: "terraform",
    name: "Terraform",
    icon: "🏗️",
    category: "devops",
    desc: "Infrastructure as Code tool to build, change, and version resources.",
    longDesc:
      "HashiCorp Terraform is an open-source infrastructure as code software tool created by HashiCorp. Users define and provide data center infrastructure using a declarative configuration language known as HashiCorp Configuration Language, or optionally JSON.",
    version: "1.8.0",
    size: "42.0 MB",
    developer: "HashiCorp",
    license: "BSL 1.1",
    dependencies: ["libc"],
    snippet: "terraform --version && terraform init",
  },
  // Editors & CLI Tools
  {
    id: "neovim",
    name: "Neovim",
    icon: "⚡",
    category: "tools",
    desc: "Vim-fork focused on extensibility and usability.",
    longDesc:
      "Neovim is a refactor, and a logical successor, of Vim. It is designed for users who want the good parts of Vim, and more: built-in LSP support, an asynchronous plugin architecture, powerful Lua configuration capabilities, and modern UI integration.",
    version: "0.10.0",
    size: "14.2 MB",
    developer: "Neovim Community",
    license: "Apache 2.0 / Vim License",
    dependencies: ["luajit", "libuv", "msgpack"],
    snippet: "nvim --version",
  },
  {
    id: "git",
    name: "Git",
    icon: "📝",
    category: "tools",
    desc: "Fast, scalable, distributed revision control system.",
    longDesc:
      "Git is a distributed version control system: tracking changes in any set of files, usually used for coordinating work among programmers collaboratively developing source code during software development. Its goals include speed, data integrity, and support for distributed, non-linear workflows.",
    version: "2.44.0",
    size: "18.5 MB",
    developer: "Software Freedom Conservancy",
    license: "GPL v2",
    dependencies: ["curl", "expat", "gettext"],
    snippet: "git --version && git status",
  },
  {
    id: "tmux",
    name: "tmux",
    icon: "🔲",
    category: "tools",
    desc: "Terminal multiplexer enabling multiple windows in one terminal.",
    longDesc:
      "tmux is a terminal multiplexer: it enables a number of terminals to be created, accessed, and controlled from a single screen. tmux may be detached from a screen and continue running in the background, then later reattached.",
    version: "3.4",
    size: "1.8 MB",
    developer: "Nicholas Marriott",
    license: "ISC License",
    dependencies: ["libevent", "ncurses"],
    snippet: "tmux -V && tmux new-session -d",
  },
  // AI / ML
  {
    id: "ollama",
    name: "Ollama",
    icon: "🦙",
    category: "aiml",
    desc: "Run large language models locally on your system.",
    longDesc:
      "Ollama is a lightweight, extensible framework for building and running language models on your local machine. It manages model weights, configuration, and runtimes, and provides a clean API and CLI interface to interact with local LLMs like Llama 3, Mistral, and Phi.",
    version: "0.1.33",
    size: "180.2 MB",
    developer: "Ollama Inc.",
    license: "MIT",
    dependencies: ["nvidia-cuda-toolkit", "rocm"],
    snippet: "ollama --version && ollama run llama3",
  },
  {
    id: "jupyter",
    name: "JupyterLab",
    icon: "📓",
    category: "aiml",
    desc: "Web-based interactive development environment for notebooks.",
    longDesc:
      "JupyterLab is the next-generation user interface for Project Jupyter, offering all the familiar building blocks of the classic Jupyter Notebook (notebook, terminal, text editor, file browser, rich outputs, etc.) in a flexible and powerful user interface.",
    version: "4.1.5",
    size: "12.4 MB",
    developer: "Project Jupyter",
    license: "BSD 3-Clause",
    dependencies: ["python", "nodejs", "npm"],
    snippet: "jupyter lab --version",
  },
];

export default function Categories() {
  const [activeCategory, setActiveCategory] = useState("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedApp, setSelectedApp] = useState<DevPackage | null>(null);
  const [copiedCompiler, setCopiedCompiler] = useState(false);
  const [copiedModalSnippet, setCopiedModalSnippet] = useState(false);

  // App states: "uninstalled", "installing", "installed"
  const [appStates, setAppStates] = useState<
    Record<string, "uninstalled" | "installing" | "installed">
  >({
    python: "uninstalled",
    nodejs: "installed", // Node.js starts pre-installed as the demo
    go: "uninstalled",
    rust: "uninstalled",
    postgresql: "uninstalled",
    redis: "uninstalled",
    mongodb: "uninstalled",
    docker: "uninstalled",
    kubectl: "uninstalled",
    terraform: "uninstalled",
    neovim: "uninstalled",
    git: "installed", // Git starts pre-installed as standard developer setup
    tmux: "uninstalled",
    ollama: "uninstalled",
    jupyter: "uninstalled",
  });

  // Terminal logging state
  const [terminalLogs, setTerminalLogs] = useState<string[]>([
    "AutoDev Store Service v0.1.0 Initialized.",
    "System architectures scanned: Linux x86_64, Darwin arm64, Windows amd64",
    "Pre-installed components detected: Node.js (v20.12.2), Git (v2.44.0)",
    "Type an app search or click [GET] to execute mock commands...",
  ]);
  const [isTerminalRunning, setIsTerminalRunning] = useState(false);

  // Filter packages based on active category and search query
  const filteredPackages = useMemo(() => {
    return packages.filter((pkg) => {
      const matchesCategory =
        activeCategory === "all" || pkg.category === activeCategory;
      const matchesSearch =
        pkg.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        pkg.desc.toLowerCase().includes(searchQuery.toLowerCase());
      return matchesCategory && matchesSearch;
    });
  }, [activeCategory, searchQuery]);

  // Install trigger (terminal simulation)
  const installPackage = (pkg: DevPackage) => {
    if (appStates[pkg.id] !== "uninstalled") return;

    setAppStates((prev) => ({ ...prev, [pkg.id]: "installing" }));
    setIsTerminalRunning(true);

    const logLines = [
      `$ autodev install ${pkg.id}`,
      `[autodev] Resolving metadata for "${pkg.name}"...`,
      `[autodev] Dependency check: ${pkg.dependencies.length > 0 ? pkg.dependencies.join(", ") : "None required"}`,
      `[autodev] Downloading packages... [${pkg.size}]`,
      `[autodev] Extracting tarball / binary archives...`,
      `[autodev] Binding path: /usr/local/bin/${pkg.id}`,
      `[autodev] Testing: executing "${pkg.snippet.split(" && ")[0]}"`,
      `[autodev] ✓ Successfully installed ${pkg.name} v${pkg.version}!`,
    ];

    // Log lines one by one to simulate terminal output
    logLines.forEach((line, index) => {
      setTimeout(() => {
        setTerminalLogs((prev) => [...prev, line]);
        if (index === logLines.length - 1) {
          setAppStates((prev) => ({ ...prev, [pkg.id]: "installed" }));
          setIsTerminalRunning(false);
        }
      }, index * 400);
    });
  };

  // Uninstall trigger
  const uninstallPackage = (id: string, name: string) => {
    setAppStates((prev) => ({ ...prev, [id]: "uninstalled" }));
    setTerminalLogs((prev) => [
      ...prev,
      `$ autodev uninstall ${id}`,
      `[autodev] Removing environment paths for ${name}...`,
      `[autodev] Cleanup complete. Removed ${name}.`,
    ]);
  };

  const installedCount = useMemo(() => {
    return Object.values(appStates).filter((state) => state === "installed")
      .length;
  }, [appStates]);

  const selectedPackagesForCliCommand = useMemo(() => {
    return Object.entries(appStates)
      .filter(([, state]) => state === "installed")
      .map(([id]) => id)
      .join(" ");
  }, [appStates]);

  return (
    <section
      id="features"
      className="py-20 px-6 max-w-7xl mx-auto scroll-mt-24"
    >
      {/* Brutalist Section Header */}
      <motion.div
        className="mb-12 text-center lg:text-left flex flex-col lg:flex-row items-center lg:items-end justify-between gap-6"
        initial={{ opacity: 0, y: 30 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.5 }}
      >
        <div>
          <div className="inline-flex items-center gap-2 border border-[#FFD700] text-[#FFD700] text-xs font-bold px-3 py-1 uppercase tracking-widest bg-[#FFD70010] mb-3">
            <Sparkles className="w-3.5 h-3.5" /> Interactive Store Dashboard
          </div>
          <h2 className="text-[2.8rem] md:text-5xl font-black leading-none tracking-tight uppercase">
            <span className="text-gradient-white">THE APP STORE</span> <br className="hidden md:block" />
            <span className="text-gradient-yellow">FOR DEVELOPERS.</span>
          </h2>
        </div>
        <p className="text-neutral-300 max-w-md text-sm md:text-base leading-relaxed lg:text-right font-medium">
          Select developer tools, databases, and language runtimes. Install them
          cleanly with a single command on your terminal. Try it live below.
        </p>
      </motion.div>

      {/* Main Store Frame */}
      <motion.div
        className="border-4 border-[#2A2A2A] bg-[#0A0A0A] nb-shadow flex flex-col overflow-hidden rounded-none"
        initial={{ opacity: 0, y: 40 }}
        whileInView={{ opacity: 1, y: 0 }}
        viewport={{ once: true }}
        transition={{ duration: 0.6, ease: "easeOut" }}
      >
        {/* Header Bar */}
        <div className="border-b-4 border-[#2A2A2A] bg-[#111111] p-4 flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-3 w-full md:w-auto">
            <div className="flex gap-1.5 mr-2">
              <span className="w-3 h-3 rounded-full bg-[#FF5F56]" />
              <span className="w-3 h-3 rounded-full bg-[#FFBD2E]" />
              <span className="w-3 h-3 rounded-full bg-[#27C93F]" />
            </div>
            <div className="text-xs font-mono bg-black border border-[#2A2A2A] text-neutral-400 px-3 py-1 font-bold">
              STOREFRONT VERIFIED
            </div>
            <div className="text-xs text-neutral-400 font-mono hidden sm:block">
              Installed:{" "}
              <span className="text-[#00FF87] font-bold">
                {installedCount} / {packages.length}
              </span>
            </div>
          </div>

          {/* Search bar */}
          <div className="relative w-full md:w-80">
            <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-neutral-400" />
            <input
              type="text"
              placeholder="Search CLI compilers, DBs, tools..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full bg-black border-2 border-[#2A2A2A] py-1.5 pl-10 pr-4 text-sm font-mono text-white placeholder-neutral-600 focus:outline-none focus:border-[#FFD700] transition-colors"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery("")}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-neutral-400 hover:text-white"
              >
                <X className="w-4 h-4" />
              </button>
            )}
          </div>
        </div>

        {/* Mid-Row: Main Columns (Sidebar Navigation + App Grid) */}
        <div className="grid grid-cols-1 lg:grid-cols-[240px_1fr] min-h-[500px]">
          {/* Sidebar */}
          <div className="border-b-4 lg:border-b-0 lg:border-r-4 border-[#2A2A2A] bg-[#0E0E0E] p-4 flex flex-row lg:flex-col overflow-x-auto lg:overflow-x-visible gap-2 lg:space-y-2 scrollbar-none whitespace-nowrap lg:whitespace-normal">
            <p className="hidden lg:block text-[10px] font-bold text-neutral-400 uppercase tracking-wider pl-2 mb-2 font-mono">
              Categories
            </p>
            {categories.map((cat) => (
              <button
                key={cat.id}
                onClick={() => setActiveCategory(cat.id)}
                className={`flex-shrink-0 px-3 py-2.5 font-bold text-xs flex items-center justify-between border-2 transition-all gap-4
                  ${
                    activeCategory === cat.id
                      ? "border-[#FFD700] text-[#FFD700] bg-[#FFD70010]"
                      : "border-transparent text-neutral-400 hover:border-[#2A2A2A] hover:text-white"
                  }`}
              >
                <div className="flex items-center gap-2.5">
                  <span className="text-base">{cat.icon}</span>
                  <span>{cat.label}</span>
                </div>
                <span className="font-mono text-[10px] bg-black border border-[#2A2A2A] text-neutral-400 px-1.5 py-0.5 font-bold rounded ml-2">
                  {cat.id === "all"
                    ? packages.length
                    : packages.filter((p) => p.category === cat.id).length}
                </span>
              </button>
            ))}

            {/* Smart Collections Divider / Text */}
            <div className="hidden lg:block pt-4 border-t border-[#2A2A2A] mt-4" />
            <p className="hidden lg:block text-[10px] font-bold text-neutral-400 uppercase tracking-wider pl-2 mb-2 font-mono">
              Smart Collections
            </p>

            <button
              onClick={() => {
                setActiveCategory("installed");
                setSearchQuery("");
              }}
              className={`flex-shrink-0 px-3 py-2.5 font-bold text-xs flex items-center justify-between border-2 transition-all gap-4
                ${
                  activeCategory === "installed"
                    ? "border-[#00FF87] text-[#00FF87] bg-[#00FF8710]"
                    : "border-transparent text-neutral-400 hover:border-[#2A2A2A] hover:text-white"
                }`}
            >
              <div className="flex items-center gap-2.5">
                <span>📦</span>
                <span>Installed</span>
              </div>
              <span className="font-mono text-[10px] bg-black border border-[#2A2A2A] text-neutral-400 px-1.5 py-0.5 font-bold rounded ml-2">
                {installedCount}
              </span>
            </button>
          </div>

          {/* Main Grid */}
          <div className="p-6 bg-[#080808] overflow-y-auto">
            {/* Banner (Visible on Discover "all" state and empty search) */}
            {activeCategory === "all" && !searchQuery && (
              <div className="mb-6 border-2 border-[#FFD700] bg-[#FFD70010] p-5 flex flex-col md:flex-row items-center justify-between gap-6 relative overflow-hidden">
                <div className="absolute right-4 top-4 text-6xl opacity-10 pointer-events-none select-none">
                  ⚡
                </div>
                <div className="space-y-2 max-w-xl">
                  <div className="inline-block text-[10px] font-black bg-[#FFD700] text-black px-2 py-0.5 uppercase tracking-wide">
                    Featured Environment Tool
                  </div>
                  <h3 className="text-2xl font-black text-white">
                    NEOVIM v0.10.0
                  </h3>
                  <p className="text-xs text-neutral-300 leading-relaxed">
                    A hyper-extensible, fast, Lua-configurable terminal text
                    editor. Autodev configures it complete with pre-bundled
                    compilers and LSP servers out of the box.
                  </p>
                </div>
                <button
                  onClick={() => {
                    const nvimPkg = packages.find((p) => p.id === "neovim");
                    if (nvimPkg) setSelectedApp(nvimPkg);
                  }}
                  className="px-4 py-2 border-2 border-white bg-white text-black text-xs font-black uppercase tracking-wide hover:bg-transparent hover:text-white transition-colors whitespace-nowrap self-stretch md:self-center text-center"
                >
                  View Details
                </button>
              </div>
            )}

            {/* App Cards List */}
            {filteredPackages.length === 0 &&
              activeCategory === "installed" && (
                <div className="py-12 text-center border-2 border-dashed border-[#2A2A2A] text-neutral-400 font-mono text-sm">
                  No apps installed yet in this session. <br />
                  Click on the [GET] button next to any package to simulate an
                  installation!
                </div>
              )}

            {filteredPackages.length === 0 &&
              activeCategory !== "installed" && (
                <div className="py-12 text-center border-2 border-dashed border-[#2A2A2A] text-neutral-400 font-mono text-sm">
                  No developer tools match your search criteria.
                </div>
              )}

            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
              {(activeCategory === "installed"
                ? packages.filter((p) => appStates[p.id] === "installed")
                : filteredPackages
              ).map((pkg) => {
                const state = appStates[pkg.id];
                return (
                  <motion.div
                    key={pkg.id}
                    initial={{ opacity: 0, y: 20 }}
                    whileInView={{ opacity: 1, y: 0 }}
                    viewport={{ once: true, margin: "-30px" }}
                    transition={{ duration: 0.35 }}
                    className={`border-2 p-4 flex flex-col justify-between transition-all duration-200 hover:border-[#FFD700] bg-[#111] hover:bg-[#151515] relative group glow-yellow-hover
                      ${state === "installed" ? "border-[#FFD70018]" : "border-[#2A2A2A]"}
                    `}
                  >
                    <div>
                      {/* Top Header */}
                      <div className="flex items-start justify-between gap-3 mb-2">
                        <div className="flex items-center gap-3">
                          <span className="text-3xl bg-black border border-[#2A2A2A] w-12 h-12 flex items-center justify-center font-bold">
                            {pkg.icon}
                          </span>
                          <div>
                            <h3 className="font-black text-sm text-white group-hover:text-[#FFD700] transition-colors">
                              {pkg.name}
                            </h3>
                            <span className="text-[10px] text-neutral-400 uppercase font-bold tracking-wider font-mono">
                              {pkg.category}
                            </span>
                          </div>
                        </div>

                        {/* CTA button (GET/INSTALL/OPEN) */}
                        <div>
                          {state === "uninstalled" && (
                            <button
                              onClick={() => installPackage(pkg)}
                              disabled={isTerminalRunning}
                              className="px-3.5 py-1 text-xs font-black border-2 border-[#FFD700] bg-[#FFD700] text-black hover:bg-transparent hover:text-[#FFD700] transition-colors disabled:opacity-50"
                            >
                              GET
                            </button>
                          )}
                          {state === "installing" && (
                            <button
                              disabled
                              className="px-2 py-1 text-[10px] font-black border-2 border-[#888] text-neutral-400 animate-pulse"
                            >
                              LOADING...
                            </button>
                          )}
                          {state === "installed" && (
                            <div className="flex gap-1.5 items-center">
                              <span className="text-[10px] font-bold text-[#00FF87] flex items-center gap-1">
                                <CheckCircle className="w-3.5 h-3.5" />
                              </span>
                              <button
                                onClick={() => setSelectedApp(pkg)}
                                className="px-2.5 py-1 text-xs font-black border-2 border-[#2A2A2A] hover:border-white text-white transition-colors"
                              >
                                OPEN
                              </button>
                            </div>
                          )}
                        </div>
                      </div>

                      <p className="text-xs text-neutral-400 leading-normal line-clamp-2 mt-2 min-h-[2rem]">
                        {pkg.desc}
                      </p>
                    </div>

                    {/* Bottom row info */}
                    <div className="mt-4 pt-3 border-t border-[#2A2A2A] flex items-center justify-between text-[10px] font-mono text-neutral-400 font-semibold">
                      <span>Size: {pkg.size}</span>
                      <button
                        onClick={() => setSelectedApp(pkg)}
                        className="text-neutral-400 hover:text-white flex items-center gap-1"
                      >
                        <Info className="w-3 h-3" /> More
                      </button>
                    </div>
                  </motion.div>
                );
              })}
            </div>
          </div>
        </div>

        {/* Bottom Panel: Interactive Terminal Simulator & Command Compiler */}
        <div className="border-t-4 border-[#2A2A2A] bg-[#0E0E0E] p-5 grid grid-cols-1 lg:grid-cols-[1fr_400px] gap-6">
          {/* CLI Builder */}
          <div className="space-y-4">
            <div className="flex items-center gap-2">
              <Layers className="w-4 h-4 text-[#FFD700]" />
              <h3 className="text-xs font-black text-white uppercase tracking-wider">
                Installer CLI Command Builder
              </h3>
            </div>
            <p className="text-xs text-neutral-400 leading-relaxed max-w-xl">
              Install the selected environment packages on your local machine
              instantly. Copy the dynamically compiled command below and run it
              in your terminal.
            </p>

            <div className="terminal w-full">
              <div className="terminal-bar py-1.5 px-3 flex justify-between items-center pr-3">
                <div className="flex items-center gap-1.5">
                  <div className="flex gap-1">
                    <span className="w-2 h-2 rounded-full bg-[#FF5F56]" />
                    <span className="w-2 h-2 rounded-full bg-[#FFBD2E]" />
                    <span className="w-2 h-2 rounded-full bg-[#27C93F]" />
                  </div>
                  <span className="text-[10px] text-neutral-400 ml-2 font-mono">
                    autodev cli compiler
                  </span>
                </div>
                <button
                  onClick={() => {
                    const cmd = `curl -fsSL https://autodevs.dev/install.sh | bash -s -- --install ${selectedPackagesForCliCommand || "nodejs git"}`;
                    navigator.clipboard.writeText(cmd);
                    setCopiedCompiler(true);
                    setTimeout(() => setCopiedCompiler(false), 1800);
                  }}
                  className="text-neutral-400 hover:text-[#FFD700] transition-colors p-1.5 flex items-center gap-1 rounded bg-[#1e1e1e] border border-[#2a2a2a] cursor-pointer"
                  title="Copy installation command"
                >
                  {copiedCompiler ? (
                    <>
                      <Check className="w-3.5 h-3.5 text-[#00FF87]" />
                      <span className="text-[10px] text-[#00FF87] font-mono pr-0.5">
                        Copied!
                      </span>
                    </>
                  ) : (
                    <>
                      <Copy className="w-3.5 h-3.5" />
                      <span className="text-[10px] text-neutral-400 font-mono pr-0.5">
                        Copy
                      </span>
                    </>
                  )}
                </button>
              </div>
              <div className="px-4 py-3 font-mono text-xs md:text-sm text-[#00FF87] flex flex-col md:flex-row md:items-center justify-between gap-4 overflow-x-auto bg-black">
                <code className="whitespace-nowrap">
                  <span className="text-neutral-500">$</span> curl -fsSL
                  https://autodevs.dev/install.sh | bash -s -- --install{" "}
                  <span className="text-[#FFD700] font-bold">
                    {selectedPackagesForCliCommand || "nodejs git"}
                  </span>
                </code>
              </div>
            </div>
          </div>

          {/* Local Installation Terminal Logs */}
          <div className="border-2 border-[#2A2A2A] bg-black flex flex-col h-[180px] lg:h-auto overflow-hidden">
            <div className="bg-[#111] px-3 py-1.5 border-b border-[#2A2A2A] flex items-center justify-between text-[10px] font-mono text-neutral-400 font-bold">
              <span className="flex items-center gap-1.5 text-[#00FF87]">
                <TerminalIcon className="w-3.5 h-3.5" /> STDOUT MONITOR
              </span>
              <span>{isTerminalRunning ? "RUNNING" : "READY"}</span>
            </div>
            <div className="p-3.5 font-mono text-[11px] leading-5 text-neutral-400 overflow-y-auto flex-1 space-y-0.5 select-text selection:bg-[#FFD700] selection:text-black">
              {terminalLogs.map((log, i) => (
                <div
                  key={i}
                  className={
                    log.startsWith("$")
                      ? "text-[#00FF87]"
                      : log.includes("✓")
                        ? "text-[#FFD700]"
                        : ""
                  }
                >
                  {log}
                </div>
              ))}
              {isTerminalRunning && (
                <span className="inline-block w-1.5 h-3 bg-[#00FF87] animate-pulse" />
              )}
            </div>
          </div>
        </div>
      </motion.div>

      {/* Package Detail Modal */}
      <AnimatePresence>
        {selectedApp && (
          <div className="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/80 backdrop-blur-sm">
            <motion.div
              initial={{ scale: 0.95, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0.95, opacity: 0 }}
              className="w-full max-w-2xl border-4 border-[#FFD700] bg-[#0A0A0A] p-6 md:p-8 relative overflow-y-auto max-h-[90vh]"
            >
              {/* Close button */}
              <button
                onClick={() => setSelectedApp(null)}
                className="absolute right-4 top-4 p-1.5 border-2 border-[#2A2A2A] bg-black text-neutral-400 hover:text-white hover:border-[#FFD700] transition-all"
              >
                <X className="w-5 h-5" />
              </button>

              {/* Modal Body */}
              <div className="space-y-6">
                {/* Header */}
                <div className="flex items-start gap-4">
                  <span className="text-5xl bg-[#111] border-2 border-[#2A2A2A] w-20 h-20 flex items-center justify-center">
                    {selectedApp.icon}
                  </span>
                  <div>
                    <div className="flex items-center gap-2.5">
                      <h3 className="text-2xl font-black text-white">
                        {selectedApp.name}
                      </h3>
                      <span className="px-2 py-0.5 border border-[#444] bg-[#111] text-[10px] font-mono text-neutral-400">
                        v{selectedApp.version}
                      </span>
                    </div>
                    <p className="text-xs text-neutral-400 uppercase tracking-wider font-mono mt-1">
                      {selectedApp.category} · {selectedApp.license}
                    </p>
                    <div className="flex gap-2 mt-3">
                      {appStates[selectedApp.id] === "installed" ? (
                        <span className="inline-flex items-center gap-1.5 text-xs text-[#00FF87] font-bold">
                          <CheckCircle className="w-4 h-4" /> Installed locally
                        </span>
                      ) : appStates[selectedApp.id] === "installing" ? (
                        <span className="text-xs text-neutral-400 font-bold animate-pulse">
                          Installing via CLI...
                        </span>
                      ) : (
                        <button
                          onClick={() => {
                            installPackage(selectedApp);
                          }}
                          className="px-4 py-1.5 text-xs font-black border-2 border-[#FFD700] bg-[#FFD700] text-black hover:bg-transparent hover:text-[#FFD700] transition-colors"
                        >
                          INSTALL NOW
                        </button>
                      )}

                      {appStates[selectedApp.id] === "installed" && (
                        <button
                          onClick={() => {
                            uninstallPackage(selectedApp.id, selectedApp.name);
                          }}
                          className="px-3 py-1.5 text-xs font-bold border-2 border-[#FF4444] hover:bg-[#FF4444] hover:text-white text-[#FF4444] transition-all"
                        >
                          Uninstall
                        </button>
                      )}
                    </div>
                  </div>
                </div>

                {/* Description */}
                <div className="space-y-2">
                  <h4 className="text-xs font-bold uppercase tracking-wider text-[#FFD700] font-mono">
                    Overview
                  </h4>
                  <p className="text-sm text-neutral-300 leading-relaxed">
                    {selectedApp.longDesc}
                  </p>
                </div>

                {/* Metadata details */}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 border-y border-[#2A2A2A] py-4 my-2">
                  <div>
                    <span className="block text-[10px] font-mono text-neutral-400 uppercase">
                      Developer
                    </span>
                    <span className="text-xs font-bold text-white block mt-0.5 truncate">
                      {selectedApp.developer}
                    </span>
                  </div>
                  <div>
                    <span className="block text-[10px] font-mono text-neutral-400 uppercase">
                      File Size
                    </span>
                    <span className="text-xs font-bold text-white block mt-0.5">
                      {selectedApp.size}
                    </span>
                  </div>
                  <div>
                    <span className="block text-[10px] font-mono text-neutral-400 uppercase">
                      License
                    </span>
                    <span className="text-xs font-bold text-white block mt-0.5">
                      {selectedApp.license}
                    </span>
                  </div>
                  <div>
                    <span className="block text-[10px] font-mono text-neutral-400 uppercase">
                      Category
                    </span>
                    <span className="text-xs font-bold text-[#FFD700] block mt-0.5 capitalize">
                      {selectedApp.category}
                    </span>
                  </div>
                </div>

                {/* Dependencies */}
                {selectedApp.dependencies.length > 0 && (
                  <div className="space-y-2">
                    <h4 className="text-xs font-bold uppercase tracking-wider text-[#FFD700] font-mono">
                      Dependencies
                    </h4>
                    <div className="flex flex-wrap gap-1.5">
                      {selectedApp.dependencies.map((dep) => (
                        <span
                          key={dep}
                          className="px-2.5 py-0.5 border border-[#333] bg-[#111] text-xs font-mono text-neutral-300"
                        >
                          {dep}
                        </span>
                      ))}
                    </div>
                  </div>
                )}

                {/* Execution command snippet */}
                <div className="space-y-2">
                  <h4 className="text-xs font-bold uppercase tracking-wider text-[#FFD700] font-mono">
                    Run / Execute Example
                  </h4>
                  <div className="bg-[#111] border border-[#2A2A2A] p-3 text-xs font-mono text-neutral-400 flex items-center justify-between gap-4">
                    <code className="break-all">{selectedApp.snippet}</code>
                    <button
                      onClick={() => {
                        navigator.clipboard.writeText(selectedApp.snippet);
                        setCopiedModalSnippet(true);
                        setTimeout(() => setCopiedModalSnippet(false), 1800);
                      }}
                      className="text-xs text-[#FFD700] hover:text-[#00FF87] font-mono bg-[#1E1E1E] px-2.5 py-1.5 rounded border border-[#333] flex items-center gap-1.5 transition-colors cursor-pointer shrink-0"
                    >
                      {copiedModalSnippet ? (
                        <>
                          <Check className="w-3.5 h-3.5 text-[#00FF87]" />
                          <span>Copied!</span>
                        </>
                      ) : (
                        <>
                          <Copy className="w-3.5 h-3.5" />
                          <span>Copy</span>
                        </>
                      )}
                    </button>
                  </div>
                </div>
              </div>
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </section>
  );
}
