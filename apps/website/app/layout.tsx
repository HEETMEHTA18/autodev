import type { Metadata } from "next";
import { Space_Grotesk, JetBrains_Mono } from "next/font/google";
import { ThemeProvider } from "next-themes";
import "./globals.css";
import "./professional-theme.css";

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
    title: "AutoDevs — Developer environment intelligence",
    description:
      "AutoDevs understands your project and development environment, helps you set up tools, diagnose problems, and ship with an AI-assisted workflow.",
    keywords: [
      "developer tools",
      "developer environment",
      "CLI",
      "AI developer tools",
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
      title: "AutoDevs — Developer environment intelligence",
      description: "Understand, set up, diagnose and operate your development environment with one workflow.",
      url: "https://autodevs.dev/",
      siteName: "AutoDevs",
      type: "website",
    },
    twitter: {
      card: "summary_large_image",
      title: "AutoDevs — Developer environment intelligence",
      description: "Understand, set up, diagnose and operate your development environment with one workflow.",
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
