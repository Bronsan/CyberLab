import type { Metadata } from "next"
import { ThemeProvider } from "@/providers/theme-provider"
import { LanguageProvider } from "@/providers/language-provider"
import "./globals.css"

export const metadata: Metadata = {
  title: "CyberLab - Practice Like a Real Hacker",
  description: "An online cybersecurity training platform providing isolated, dynamic challenge environments for security enthusiasts.",
  keywords: ["cybersecurity", "CTF", "penetration testing", "hacking", "lab", "security training"],
  authors: [{ name: "CyberLab" }],
  openGraph: {
    title: "CyberLab - Practice Like a Real Hacker",
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
        {/* Prevent FOUC */}
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
