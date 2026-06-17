"use client";
import AsciiBackground from "./AsciiBackground";

export default function Footer() {
  return (
    <footer className="relative overflow-hidden border-t-2 border-[#2A2A2A] py-16 px-6">
      <AsciiBackground className="opacity-30" />
      <div className="relative z-10 max-w-7xl mx-auto">
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-10 mb-16">
          {/* Brand */}
          <div className="md:col-span-2">
            <div className="text-2xl font-black text-[#FFD700] mb-3">
              ⚡ AUTODEV
            </div>
            <p className="text-neutral-400 text-sm leading-relaxed max-w-xs">
              The App Store for Developers. Open-source, cross-platform, and
              built for thousands of contributors.
            </p>
            <div className="flex flex-wrap items-center gap-3 mt-5">
              <a
                href="https://github.com/HEETMEHTA18/autodev"
                target="_blank"
                rel="noreferrer"
                className="nb-btn px-4 py-2 text-xs"
              >
                GitHub →
              </a>
              <a
                href="https://www.producthunt.com/products/autodevs?embed=true&utm_source=badge-featured&utm_medium=badge&utm_campaign=badge-autodevs"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-block hover:opacity-90 transition-opacity"
              >
                <img
                  alt="Autodevs - AI-powered development setup in minutes | Product Hunt"
                  width="250"
                  height="54"
                  src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=1162368&theme=neutral&t=1780484994611"
                  className="h-[34px] w-auto"
                />
              </a>
            </div>
          </div>

          {/* Product */}
          <div>
            <h3 className="font-black text-white text-sm uppercase tracking-wider mb-4">
              Product
            </h3>
            <ul className="space-y-2 text-sm text-neutral-400">
              {[
                { name: "Features", href: "/#features" },
                { name: "Profiles", href: "/#profiles" },
                { name: "GitHub Scanner", href: "/#github-scanner" },
                { name: "Skills Roadmap", href: "/#skills" },
                {
                  name: "Changelog",
                  onClick: () =>
                    window.dispatchEvent(
                      new Event("autodev_open_update_modal"),
                    ),
                },
              ].map((l) => (
                <li key={l.name}>
                  {l.onClick ? (
                    <button
                      onClick={l.onClick}
                      className="hover:text-[#FFD700] transition-colors text-left bg-transparent border-0 p-0 cursor-pointer"
                    >
                      {l.name}
                    </button>
                  ) : (
                    <a
                      href={l.href}
                      className="hover:text-[#FFD700] transition-colors"
                    >
                      {l.name}
                    </a>
                  )}
                </li>
              ))}
            </ul>
          </div>

          {/* Templates */}
          <div>
            <h3 className="font-black text-white text-sm uppercase tracking-wider mb-4">
              Templates
            </h3>
            <ul className="space-y-2 text-sm text-neutral-400">
              {[
                { name: "React Template", href: "/react-template" },
                { name: "Next.js Template", href: "/nextjs-template" },
                { name: "MERN Stack Template", href: "/mern-template" },
                { name: "AI Agent Template", href: "/ai-agent-template" },
                { name: "Flutter Template", href: "/flutter-template" },
              ].map((l) => (
                <li key={l.name}>
                  <a
                    href={l.href}
                    className="hover:text-[#FFD700] transition-colors"
                  >
                    {l.name}
                  </a>
                </li>
              ))}
            </ul>
          </div>

          {/* Open Source */}
          <div>
            <h3 className="font-black text-white text-sm uppercase tracking-wider mb-4">
              Open Source
            </h3>
            <ul className="space-y-2 text-sm text-neutral-400">
              {[
                {
                  name: "Contributing",
                  href: "https://github.com/HEETMEHTA18/autodev/blob/main/CONTRIBUTING.md",
                },
                {
                  name: "Code of Conduct",
                  href: "https://github.com/HEETMEHTA18/autodev/blob/main/CODE_OF_CONDUCT.md",
                },
                {
                  name: "Security",
                  href: "https://github.com/HEETMEHTA18/autodev/blob/main/SECURITY.md",
                },
                {
                  name: "Roadmap",
                  href: "https://github.com/HEETMEHTA18/autodev/blob/main/ROADMAP.md",
                },
                {
                  name: "MIT License",
                  href: "https://github.com/HEETMEHTA18/autodev/blob/main/LICENSE",
                },
              ].map((l) => (
                <li key={l.name}>
                  <a
                    href={l.href}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="hover:text-[#FFD700] transition-colors"
                  >
                    {l.name}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <div className="border-t border-[#1A1A1A] pt-8 flex flex-col sm:flex-row justify-between items-center gap-4">
          <p className="text-xs text-neutral-400">
            © 2026 AutoDev Contributors — MIT License
          </p>
          <p className="text-xs text-neutral-400 font-mono">
            Clone. Scan. Install. Build.
          </p>
        </div>
      </div>
    </footer>
  );
}
