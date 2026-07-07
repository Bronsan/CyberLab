'use client'

import { useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import { useTranslation } from '@/providers/language-provider'
import { useTheme } from '@/providers/theme-provider'
import { LanguageSwitcher } from '@/components/ui/language-switcher'
import { ThemeToggle } from '@/components/ui/theme-toggle'
import {
  Terminal,
  Shield,
  Zap,
  Server,
  Code2,
  Trophy,
  Sparkles,
  ChevronRight,
  Star,
  GitBranch,
  MessageCircle,
  Globe,
  BookOpen,
  Layers,
  Lock,
  Cpu,
  ArrowRight,
  Menu,
  X,
  ChevronDown,
  Brain,
  Container,
  LayoutDashboard,
} from 'lucide-react'

// ─── Animated Counter ───
function AnimatedCounter({ end, duration = 2000 }: { end: number; duration?: number }) {
  const [count, setCount] = useState(0)
  const ref = useRef<HTMLSpanElement>(null)

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          let start = 0
          const step = Math.ceil(end / (duration / 16))
          const timer = setInterval(() => {
            start += step
            if (start >= end) {
              setCount(end)
              clearInterval(timer)
            } else {
              setCount(start)
            }
          }, 16)
          observer.disconnect()
        }
      },
      { threshold: 0.3 }
    )
    if (ref.current) observer.observe(ref.current)
    return () => observer.disconnect()
  }, [end, duration])

  return <span ref={ref}>{count.toLocaleString()}</span>
}

// ─── Scroll Reveal ───
function ScrollReveal({ children, className = '' }: { children: React.ReactNode; className?: string }) {
  const ref = useRef<HTMLDivElement>(null)
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => { if (entry.isIntersecting) { setVisible(true); observer.disconnect() } },
      { threshold: 0.1 }
    )
    if (ref.current) observer.observe(ref.current)
    return () => observer.disconnect()
  }, [])

  return (
    <div ref={ref} className={`transition-all duration-700 ${visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'} ${className}`}>
      {children}
    </div>
  )
}

// ─── Feature Card ───
function FeatureCard({ icon: Icon, title, description, index }: { icon: any; title: string; description: string; index: number }) {
  return (
    <ScrollReveal>
      <div className="group relative p-6 sm:p-8 rounded-2xl border border-border bg-card hover:border-primary/20 transition-all duration-300 hover:shadow-lg hover:shadow-primary/5">
        <div className={`absolute inset-0 rounded-2xl bg-gradient-to-br from-primary/[0.03] to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500`} />
        <div className="relative">
          <div className="inline-flex p-3 rounded-xl bg-primary/10 text-primary mb-4 group-hover:scale-110 transition-transform duration-300">
            <Icon className="h-6 w-6" />
          </div>
          <h3 className="text-lg font-semibold mb-2">{title}</h3>
          <p className="text-sm text-muted-foreground leading-relaxed">{description}</p>
        </div>
      </div>
    </ScrollReveal>
  )
}

// ─── Tech Badge ───
function TechBadge({ name, icon: Icon }: { name: string; icon: any }) {
  return (
    <div className="inline-flex items-center gap-2 px-4 py-2 rounded-xl border border-border bg-card text-sm font-medium hover:border-primary/20 hover:bg-primary/5 transition-all duration-300">
      <Icon className="h-4 w-4 text-primary" />
      <span>{name}</span>
    </div>
  )
}

// ─── Navbar ───
function LandingNavbar() {
  const { t } = useTranslation()
  const { resolvedTheme } = useTheme()
  const [mobileOpen, setMobileOpen] = useState(false)
  const [scrolled, setScrolled] = useState(false)

  useEffect(() => {
    const handler = () => setScrolled(window.scrollY > 20)
    window.addEventListener('scroll', handler)
    return () => window.removeEventListener('scroll', handler)
  }, [])

  const scrollTo = (id: string) => {
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth' })
    setMobileOpen(false)
  }

  const navItems = [
    { label: t.nav.features, id: 'features' },
    { label: t.nav.techStack, id: 'tech' },
    { label: t.nav.challenges, href: '/challenges' },
    { label: t.nav.ranking, href: '/ranking' },
  ]

  return (
    <nav className={`fixed top-0 left-0 right-0 z-50 transition-all duration-300 ${
      scrolled ? 'bg-background/80 backdrop-blur-xl border-b border-border' : 'bg-transparent'
    }`}>
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2 group">
            <div className="relative">
              <div className="absolute inset-0 bg-primary/20 blur-xl rounded-full group-hover:blur-2xl transition-all" />
              <Terminal className="relative h-6 w-6 text-primary" />
            </div>
            <span className="text-lg font-bold">{t.common.brand}</span>
          </Link>

          {/* Desktop Nav */}
          <div className="hidden md:flex items-center gap-1">
            {navItems.map((item) =>
              item.href ? (
                <Link
                  key={item.label}
                  href={item.href}
                  className="px-3 py-2 text-sm text-muted-foreground hover:text-foreground rounded-lg hover:bg-muted transition-colors"
                >
                  {item.label}
                </Link>
              ) : (
                <button
                  key={item.label}
                  onClick={() => item.id && scrollTo(item.id)}
                  className="px-3 py-2 text-sm text-muted-foreground hover:text-foreground rounded-lg hover:bg-muted transition-colors"
                >
                  {item.label}
                </button>
              )
            )}
          </div>

          {/* Right side */}
          <div className="flex items-center gap-2">
            <div className="hidden md:flex items-center gap-1">
              <LanguageSwitcher />
              <ThemeToggle />
            </div>
            <Link href="/login" className="hidden md:inline-flex px-4 py-2 text-sm font-medium rounded-lg text-muted-foreground hover:text-foreground hover:bg-muted transition-colors">
              {t.common.signIn}
            </Link>
            <Link
              href="/register"
              className="hidden md:inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors"
            >
              {t.common.getStarted}
              <ArrowRight className="h-3.5 w-3.5" />
            </Link>

            {/* Mobile menu */}
            <button
              className="md:hidden p-2 rounded-lg hover:bg-muted transition-colors"
              onClick={() => setMobileOpen(!mobileOpen)}
            >
              {mobileOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
            </button>
          </div>
        </div>
      </div>

      {/* Mobile panel */}
      {mobileOpen && (
        <div className="md:hidden border-b border-border bg-background/95 backdrop-blur-xl">
          <div className="px-4 py-4 space-y-1">
            {navItems.map((item) =>
              item.href ? (
                <Link key={item.label} href={item.href} onClick={() => setMobileOpen(false)}
                  className="block px-3 py-2.5 text-sm rounded-lg hover:bg-muted transition-colors">
                  {item.label}
                </Link>
              ) : (
                <button key={item.label} onClick={() => item.id && scrollTo(item.id)}
                  className="block w-full text-left px-3 py-2.5 text-sm rounded-lg hover:bg-muted transition-colors">
                  {item.label}
                </button>
              )
            )}
            <div className="flex items-center gap-2 pt-3 px-3">
              <LanguageSwitcher />
              <ThemeToggle />
            </div>
            <div className="pt-2 px-3 flex gap-2">
              <Link href="/login" onClick={() => setMobileOpen(false)}
                className="flex-1 text-center px-4 py-2 text-sm font-medium rounded-lg border border-border hover:bg-muted transition-colors">
                {t.common.signIn}
              </Link>
              <Link href="/register" onClick={() => setMobileOpen(false)}
                className="flex-1 text-center px-4 py-2 text-sm font-medium rounded-lg bg-primary text-primary-foreground hover:bg-primary/90 transition-colors">
                {t.common.getStarted}
              </Link>
            </div>
          </div>
        </div>
      )}
    </nav>
  )
}

// ─── Footer ───
function LandingFooter() {
  const { t } = useTranslation()
  const { resolvedTheme } = useTheme()

  return (
    <footer className="border-t border-border">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-16">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
          {/* Brand */}
          <div className="col-span-2 md:col-span-1">
            <div className="flex items-center gap-2 mb-4">
              <Terminal className="h-5 w-5 text-primary" />
              <span className="font-bold text-lg">{t.common.brand}</span>
            </div>
            <p className="text-sm text-muted-foreground leading-relaxed mb-4">
              {t.footer.description}
            </p>
            <div className="flex gap-3">
              <a href="https://github.com" target="_blank" rel="noopener noreferrer"
                className="p-2 rounded-lg hover:bg-muted text-muted-foreground hover:text-foreground transition-colors">
                <GitBranch className="h-5 w-5" />
              </a>
              <a href="https://twitter.com" target="_blank" rel="noopener noreferrer"
                className="p-2 rounded-lg hover:bg-muted text-muted-foreground hover:text-foreground transition-colors">
                <MessageCircle className="h-5 w-5" />
              </a>
            </div>
          </div>

          {/* Links */}
          <div>
            <h4 className="text-sm font-semibold mb-4">{t.footer.links}</h4>
            <ul className="space-y-2.5">
              <li><Link href="/challenges" className="text-sm text-muted-foreground hover:text-foreground transition-colors">{t.nav.features}</Link></li>
              <li><Link href="/challenges" className="text-sm text-muted-foreground hover:text-foreground transition-colors">{t.nav.challenges}</Link></li>
              <li><Link href="/ranking" className="text-sm text-muted-foreground hover:text-foreground transition-colors">{t.nav.ranking}</Link></li>
            </ul>
          </div>

          {/* Community */}
          <div>
            <h4 className="text-sm font-semibold mb-4">{t.footer.community}</h4>
            <ul className="space-y-2.5">
              <li><a href="https://github.com" target="_blank" rel="noopener noreferrer" className="text-sm text-muted-foreground hover:text-foreground transition-colors">GitHub</a></li>
              <li><a href="https://twitter.com" target="_blank" rel="noopener noreferrer" className="text-sm text-muted-foreground hover:text-foreground transition-colors">Twitter</a></li>
            </ul>
          </div>

          {/* Legal */}
          <div>
            <h4 className="text-sm font-semibold mb-4">{t.footer.legal}</h4>
            <ul className="space-y-2.5">
              <li><a href="#" className="text-sm text-muted-foreground hover:text-foreground transition-colors">{t.footer.privacyPolicy}</a></li>
              <li><a href="#" className="text-sm text-muted-foreground hover:text-foreground transition-colors">{t.footer.termsOfService}</a></li>
            </ul>
          </div>
        </div>

        <div className="mt-12 pt-8 border-t border-border flex flex-col sm:flex-row items-center justify-between gap-4">
          <p className="text-sm text-muted-foreground">
            &copy; {new Date().getFullYear()} {t.common.brand}. {t.footer.copyright}
          </p>
          <p className="text-sm text-muted-foreground flex items-center gap-1">
            {t.footer.madeWith} <span className="text-red-500">❤</span>
          </p>
        </div>
      </div>
    </footer>
  )
}

// ─── Main Page ───
export default function HomePage() {
  const { t, locale } = useTranslation()

  const stats = [
    { label: t.stats.challenges, value: 42, icon: Code2 },
    { label: t.stats.users, value: 1280, icon: Star },
    { label: t.stats.online, value: 47, icon: Zap },
    { label: t.stats.containers, value: 89, icon: Server },
  ]

  const features = t.features.items.map((item, i) => ({
    ...item,
    icon: [Server, Shield, Brain, Trophy, Cpu, GitBranch][i],
  }))

  const techStack = [
    { name: 'Go', icon: Terminal },
    { name: 'Gin', icon: Zap },
    { name: 'GORM', icon: Layers },
    { name: 'Next.js', icon: LayoutDashboard },
    { name: 'TypeScript', icon: Code2 },
    { name: 'Docker', icon: Container },
    { name: 'Redis', icon: Cpu },
    { name: 'MySQL', icon: Server },
    { name: 'WebSocket', icon: Globe },
    { name: 'Tailwind CSS', icon: Sparkles },
  ]

  return (
    <div className="min-h-screen bg-background">
      <LandingNavbar />

      {/* ───── HERO ───── */}
      <section className="relative min-h-[90vh] flex items-center overflow-hidden">
        {/* Background effects */}
        <div className="absolute inset-0">
          <div className="absolute top-1/4 left-1/4 w-96 h-96 bg-primary/5 rounded-full blur-[120px] animate-pulse-glow" />
          <div className="absolute bottom-1/4 right-1/4 w-80 h-80 bg-primary/5 rounded-full blur-[100px] animate-pulse-glow" style={{ animationDelay: '1.5s' }} />
        </div>

        <div className="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-32 md:py-40">
          <div className="max-w-4xl mx-auto text-center">
            {/* Badge */}
            <ScrollReveal>
              <div className="inline-flex items-center gap-2 px-4 py-1.5 rounded-full border border-primary/20 bg-primary/5 text-primary text-sm mb-8">
                <Terminal className="h-4 w-4" />
                <span>{t.common.tagline}</span>
              </div>
            </ScrollReveal>

            {/* Title */}
            <ScrollReveal>
              <h1 className="text-5xl sm:text-6xl md:text-7xl lg:text-8xl font-bold tracking-tight leading-tight mb-6">
                {t.hero.title}{' '}
                <span className="gradient-text">{t.hero.titleAccent}</span>
              </h1>
            </ScrollReveal>

            {/* Subtitle */}
            <ScrollReveal>
              <p className="text-lg md:text-xl text-muted-foreground max-w-2xl mx-auto mb-10 leading-relaxed">
                {t.hero.subtitle}
              </p>
            </ScrollReveal>

            {/* CTA Buttons */}
            <ScrollReveal>
              <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
                <Link
                  href="/register"
                  className="inline-flex items-center gap-2 px-8 py-3.5 rounded-xl bg-primary text-primary-foreground font-semibold hover:bg-primary/90 transition-all duration-300 shadow-lg shadow-primary/20 hover:shadow-xl hover:shadow-primary/30 active:scale-[0.98]"
                >
                  {t.hero.cta}
                  <ArrowRight className="h-4 w-4" />
                </Link>
                <a
                  href="#features"
                  onClick={(e) => { e.preventDefault(); document.getElementById('features')?.scrollIntoView({ behavior: 'smooth' }) }}
                  className="inline-flex items-center gap-2 px-8 py-3.5 rounded-xl border border-border font-semibold hover:bg-muted transition-all duration-300 active:scale-[0.98]"
                >
                  {t.hero.ctaSecondary}
                  <ChevronDown className="h-4 w-4" />
                </a>
              </div>
            </ScrollReveal>
          </div>
        </div>
      </section>

      {/* ───── STATS ───── */}
      <section className="border-y border-border">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-12">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            {stats.map((stat) => (
              <div key={stat.label} className="text-center">
                <stat.icon className="h-6 w-6 mx-auto mb-2 text-primary" />
                <div className="text-3xl sm:text-4xl font-bold mb-1">
                  <AnimatedCounter end={stat.value} />
                </div>
                <div className="text-sm text-muted-foreground">{stat.label}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* ───── FEATURES ───── */}
      <section id="features" className="py-24 sm:py-32">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <ScrollReveal>
            <div className="text-center max-w-2xl mx-auto mb-16">
              <h2 className="text-3xl sm:text-4xl font-bold mb-4">{t.features.title}</h2>
              <p className="text-muted-foreground text-lg">{t.features.subtitle}</p>
            </div>
          </ScrollReveal>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((feature, i) => (
              <FeatureCard
                key={feature.title}
                icon={feature.icon}
                title={feature.title}
                description={feature.description}
                index={i}
              />
            ))}
          </div>
        </div>
      </section>

      {/* ───── TECH STACK ───── */}
      <section id="tech" className="py-24 sm:py-32 bg-muted/30">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <ScrollReveal>
            <div className="text-center max-w-2xl mx-auto mb-16">
              <h2 className="text-3xl sm:text-4xl font-bold mb-4">{t.techStack.title}</h2>
              <p className="text-muted-foreground text-lg">{t.techStack.subtitle}</p>
            </div>
          </ScrollReveal>

          <ScrollReveal>
            <div className="flex flex-wrap justify-center gap-3 max-w-3xl mx-auto">
              {techStack.map((tech) => (
                <TechBadge key={tech.name} name={tech.name} icon={tech.icon} />
              ))}
            </div>
          </ScrollReveal>

          {/* Architecture diagram using CSS */}
          <ScrollReveal>
            <div className="mt-16 glass rounded-2xl p-8 sm:p-12">
              <div className="flex flex-col items-center gap-6">
                {/* User */}
                <div className="px-6 py-3 rounded-xl border-2 border-primary/30 bg-primary/5 text-sm font-semibold">
                  🌐 Browser / Client
                </div>
                <ChevronDown className="h-4 w-4 text-muted-foreground" />
                {/* Nginx */}
                <div className="px-6 py-3 rounded-xl border-2 border-border bg-card text-sm font-medium">
                  Nginx Reverse Proxy
                </div>
                <div className="flex items-center gap-4">
                  <ChevronDown className="h-4 w-4 text-muted-foreground" />
                  <ChevronDown className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex flex-wrap justify-center gap-4">
                  <div className="px-5 py-3 rounded-xl border-2 border-primary/20 bg-primary/5 text-sm font-medium">Next.js Frontend</div>
                  <div className="px-5 py-3 rounded-xl border-2 border-primary/20 bg-primary/5 text-sm font-medium">Go Gin API</div>
                </div>
                <div className="flex items-center gap-4">
                  <ChevronDown className="h-4 w-4 text-muted-foreground" />
                  <ChevronDown className="h-4 w-4 text-muted-foreground" />
                  <ChevronDown className="h-4 w-4 text-muted-foreground" />
                </div>
                <div className="flex flex-wrap justify-center gap-4">
                  <div className="px-4 py-2 rounded-lg border border-border bg-card text-xs text-muted-foreground">MySQL</div>
                  <div className="px-4 py-2 rounded-lg border border-border bg-card text-xs text-muted-foreground">Redis</div>
                  <div className="px-4 py-2 rounded-lg border border-border bg-card text-xs text-muted-foreground">Docker</div>
                  <div className="px-4 py-2 rounded-lg border border-border bg-card text-xs text-muted-foreground">WebSocket</div>
                  <div className="px-4 py-2 rounded-lg border border-border bg-card text-xs text-muted-foreground">LLM AI</div>
                </div>
              </div>
            </div>
          </ScrollReveal>
        </div>
      </section>

      {/* ───── CTA ───── */}
      <section className="py-24 sm:py-32">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <ScrollReveal>
            <div className="glass rounded-3xl p-12 sm:p-16 text-center relative overflow-hidden">
              <div className="absolute inset-0 bg-gradient-to-br from-primary/[0.03] to-transparent" />
              <div className="relative">
                <h2 className="text-3xl sm:text-4xl font-bold mb-4">{t.cta.title}</h2>
                <p className="text-muted-foreground text-lg mb-8 max-w-xl mx-auto">{t.cta.subtitle}</p>
                <Link
                  href="/register"
                  className="inline-flex items-center gap-2 px-8 py-3.5 rounded-xl bg-primary text-primary-foreground font-semibold hover:bg-primary/90 transition-all duration-300 shadow-lg shadow-primary/20 hover:shadow-xl active:scale-[0.98]"
                >
                  {t.cta.button}
                  <ArrowRight className="h-4 w-4" />
                </Link>
              </div>
            </div>
          </ScrollReveal>
        </div>
      </section>

      <LandingFooter />
    </div>
  )
}
