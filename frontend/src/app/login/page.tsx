"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { Navbar } from "@/components/layout/navbar"
import { Button } from "@/components/ui/button"
import { useAuthStore } from "@/store/auth-store"
import { authApi } from "@/services/auth"
import { Terminal, Eye, EyeOff } from "lucide-react"

export default function LoginPage() {
  const router = useRouter()
  const { setAuth } = useAuthStore()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError("")
    setLoading(true)

    try {
      const result = await authApi.login({ email, password })
      setAuth(result.token, result.user)
      router.push("/challenges")
    } catch (err: any) {
      setError(err?.response?.data?.message || err.message || "Login failed")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-[#0a0a0a]">
      <Navbar />
      <div className="flex items-center justify-center px-4 pt-24 pb-20">
        <div className="w-full max-w-md">
          <div className="glass rounded-2xl p-8">
            <div className="flex items-center justify-center gap-2 mb-8">
              <Terminal className="h-6 w-6 text-[#00ff41]" />
              <span className="text-xl font-bold text-white">CyberLab</span>
            </div>

            <h1 className="text-2xl font-bold text-white text-center mb-2">
              Welcome Back
            </h1>
            <p className="text-gray-500 text-center text-sm mb-8">
              Sign in to continue your journey
            </p>

            {error && (
              <div className="mb-4 rounded-lg bg-red-500/10 border border-red-500/20 px-4 py-3 text-sm text-red-400">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1.5">
                  Email
                </label>
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  placeholder="you@example.com"
                  required
                  className="w-full rounded-lg border border-[#262626] bg-[#111] px-4 py-2.5 text-sm text-white placeholder-gray-600 focus:border-[#00ff41] focus:outline-none focus:ring-1 focus:ring-[#00ff41]"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-400 mb-1.5">
                  Password
                </label>
                <div className="relative">
                  <input
                    type={showPassword ? "text" : "password"}
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    placeholder="••••••••"
                    required
                    className="w-full rounded-lg border border-[#262626] bg-[#111] px-4 py-2.5 pr-10 text-sm text-white placeholder-gray-600 focus:border-[#00ff41] focus:outline-none focus:ring-1 focus:ring-[#00ff41]"
                  />
                  <button
                    type="button"
                    onClick={() => setShowPassword(!showPassword)}
                    className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500"
                  >
                    {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                  </button>
                </div>
              </div>

              <Button
                type="submit"
                disabled={loading}
                className="w-full bg-[#00ff41] text-black hover:bg-[#00ff41]/90"
              >
                {loading ? "Signing in..." : "Sign In"}
              </Button>
            </form>

            <p className="mt-6 text-center text-sm text-gray-500">
              Don&apos;t have an account?{" "}
              <Link href="/register" className="text-[#00ff41] hover:underline">
                Register
              </Link>
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}
