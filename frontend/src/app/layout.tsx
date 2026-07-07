import type { Metadata, Viewport } from "next"
import { ThemeProvider } from "@/providers/theme-provider"
import { LanguageProvider } from "@/providers/language-provider"
import "./globals.css"

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
  themeColor: [
    { media: "(prefers-color-scheme: dark)", color: "#00ff41" },
    { media: "(prefers-color-scheme: light)", color: "#059669" },
  ],
}

export const metadata: Metadata = {
  title: "CyberLab — Practice Like a Real Hacker",
  description: "An online cybersecurity training platform providing isolated, dynamic challenge environments for security enthusiasts.",
  keywords: ["cybersecurity", "CTF", "penetration testing", "hacking", "lab", "security training"],
  authors: [{ name: "CyberLab" }],
  manifest: "/CyberLab/manifest.json",
  appleWebApp: {
    capable: true,
    title: "CyberLab",
    statusBarStyle: "black-translucent",
  },
  openGraph: {
    title: "CyberLab — Practice Like a Real Hacker",
    description: "Master cybersecurity through hands-on practice. Each challenge spawns an isolated, real-world environment.",
    type: "website",
    siteName: "CyberLab",
  },
}

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en" suppressHydrationWarning className="dark">
      <head>
        <link rel="apple-touch-icon" href="/CyberLab/icon-192.png" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="mobile-web-app-capable" content="yes" />
        <script dangerouslySetInnerHTML={{
          __html: `
            (function() {
              try {
                var theme = localStorage.getItem('cyberlab-theme') || 'dark';
                var resolved = theme === 'system'
                  ? (window.matchMedia('(prefers-color-scheme:light)').matches ? 'light' : 'dark')
                  : theme;
                document.documentElement.classList.add(resolved);
              } catch(e) {}
            })();
          `
        }} />
      </head>
      <body className="antialiased min-h-screen">
        <ThemeProvider>
          <LanguageProvider>
            {children}
          </LanguageProvider>
        </ThemeProvider>
      </body>
    </html>
  )
}
