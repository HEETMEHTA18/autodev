import type { MetadataRoute } from "next";

export const dynamic = "force-static";

export default function robots(): MetadataRoute.Robots {
  return {
    rules: {
      userAgent: "*",
      allow: "/",
      disallow: ["/cdn-cgi/"],
    },
    sitemap: "https://autodevs.dev/sitemap.xml",
  };
}
