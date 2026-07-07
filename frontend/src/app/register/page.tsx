"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Navbar } from "@/components/layout/navbar"
import { Button } from "@/components/ui/button"
import { authApi } from "@/services/auth"
import { Terminal } from "lucide-react"

export default function RegisterPage() {
  const router = useRouter()
  const [form, setForm] = useState({ username: "", email: "", password: "", confirmPassword: "" })
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")

    if (form.password !== form.confirmPassword) {
      setError("Passwords do not match")
      return
    }
    if (form.password.length < 6) {
      setError("Password must be at least 6 characters")
      return
    }

    setLoading(true)
    try {
      await authApi.register({
        username: form.username,
        email: form.email,
        password: form.password,
      })
      router.push("/login?registered=true")
    } catch (err: any) {
      setError(err?.response?.data?.message || err.message || "Registration failed")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="flex items-center justify-center px-4 pt-24 pb-20">
        <div className="w-full max-w-md">
          <div className="glass rounded-2xl p-8">
            <div className="flex items-center justify-center gap-2 mb-8">
              <Terminal className="h-6 w-6 text-primary" />
              <span className="text-xl font-bold text-foreground">CyberLab</span>
            </div>

            <h1 className="text-2xl font-bold text-foreground text-center mb-2">
              Create Account
            </h1>
            <p className="text-muted-foreground text-center text-sm mb-8">
              Start your cybersecurity journey
            </p>

            {error && (
              <div className="mb-4 rounded-lg bg-destructive/10 border border-destructive/20 px-4 py-3 text-sm text-destructive">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-muted-foreground mb-1.5">Username</label>
                <input type="text" value={form.username} onChange={(e) => setForm({...form, username: e.target.value})}
                  placeholder="lokumi" required minLength={3} maxLength={50}
                  className="w-full rounded-lg border border-border bg-card px-4 py-2.5 text-sm text-foreground placeholder-muted-foreground/50 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors" />
              </div>
              <div>
                <label className="block text-sm font-medium text-muted-foreground mb-1.5">Email</label>
                <input type="email" value={form.email} onChange={(e) => setForm({...form, email: e.target.value})}
                  placeholder="you@example.com" required
                  className="w-full rounded-lg border border-border bg-card px-4 py-2.5 text-sm text-foreground placeholder-muted-foreground/50 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors" />
              </div>
              <div>
                <label className="block text-sm font-medium text-muted-foreground mb-1.5">Password</label>
                <input type="password" value={form.password} onChange={(e) => setForm({...form, password: e.target.value})}
                  placeholder="••••••••" required minLength={6}
                  className="w-full rounded-lg border border-border bg-card px-4 py-2.5 text-sm text-foreground placeholder-muted-foreground/50 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors" />
              </div>
              <div>
                <label className="block text-sm font-medium text-muted-foreground mb-1.5">Confirm Password</label>
                <input type="password" value={form.confirmPassword} onChange={(e) => setForm({...form, confirmPassword: e.target.value})}
                  placeholder="••••••••" required
                  className="w-full rounded-lg border border-border bg-card px-4 py-2.5 text-sm text-foreground placeholder-muted-foreground/50 focus:border-primary focus:outline-none focus:ring-1 focus:ring-primary transition-colors" />
              </div>
              <Button type="submit" disabled={loading}
                className="w-full bg-primary text-primary-foreground hover:bg-primary/90">
                {loading ? "Creating account..." : "Create Account"}
              </Button>
            </form>

            <p className="mt-6 text-center text-sm text-muted-foreground">
              Already have an account?{" "}
              <Link href="/login" className="text-primary hover:underline">Sign in</Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
