import type { Metadata } from "next";
import { Space_Grotesk, JetBrains_Mono } from "next/font/google";
import { ThemeProvider } from "next-themes";
import "./globals.css";

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  variable: "--font-space",
  display: "swap",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
  display: "swap",
});

export async function generateMetadata(): Promise<Metadata> {
  return {
    metadataBase: new URL("https://autodevs.dev"),
    title: "AutoDevs — Turn any GitHub repository into actionable development tasks",
    description:
      "AutoDevs analyzes your codebase, finds issues, generates improvements, and helps you ship faster. Turn any GitHub repository into actionable development tasks in seconds.",
    keywords: [
      "developer tools",
      "CLI",
      "package manager",
      "environment setup",
      "autodev",
      "autodevs",
    ],
    alternates: { canonical: "https://autodevs.dev/" },
    robots: {
      index: true,
      follow: true,
      googleBot: {
        index: true,
        follow: true,
        "max-video-preview": -1,
        "max-image-preview": "large",
        "max-snippet": -1,
      },
    },
    icons: {
      icon: "/favicon.ico",
      shortcut: "/favicon.ico",
      apple: "/apple-touch-icon.png",
    },
    openGraph: {
      title: "AutoDevs — Turn GitHub repos into development tasks",
      description: "AutoDevs analyzes your codebase, finds issues, generates improvements, and helps you ship faster.",
      url: "https://autodevs.dev/",
      siteName: "AutoDevs",
      type: "website",
    },
    twitter: {
      card: "summary_large_image",
      title: "AutoDevs — Turn GitHub repos into development tasks",
      description: "AutoDevs analyzes your codebase, finds issues, generates improvements, and helps you ship faster.",
    },
  };
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={`${spaceGrotesk.variable} ${jetbrainsMono.variable}`}
    >
      <head>
        <link rel="preconnect" href="https://api.github.com" crossOrigin="anonymous" />
        <link rel="preconnect" href="https://api.counterapi.dev" crossOrigin="anonymous" />
        <script
          async
          src="https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-1639394894393563"
          crossOrigin="anonymous"
        />
      </head>
      <body className="font-space antialiased">
        <script
          type="application/ld+json"
          dangerouslySetInnerHTML={{
            __html: JSON.stringify({
              "@context": "https://schema.org",
              "@type": "SoftwareApplication",
              name: "AutoDevs",
              applicationCategory: "DeveloperApplication",
              operatingSystem: "Linux, macOS, Windows",
              url: "https://autodevs.dev",
              offers: {
                "@type": "Offer",
                price: "0.00",
                priceCurrency: "USD",
              },
            }),
          }}
        />
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem={false}
          disableTransitionOnChange={false}
        >
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
