"use client"

import Link from "next/link"
import { useState } from "react"
import { useAuthStore } from "@/store/auth-store"
import { Button } from "@/components/ui/button"
import { Terminal, Menu, X, User, Trophy, BookOpen, LogOut } from "lucide-react"

export function Navbar() {
  const { isAuthenticated, user, logout } = useAuthStore()
  const [mobileOpen, setMobileOpen] = useState(false)

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 border-b border-[#1a1a1a] bg-[#0a0a0a]/80 backdrop-blur-xl">
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2">
            <Terminal className="h-6 w-6 text-[#00ff41]" />
            <span className="text-lg font-bold text-white">CyberLab</span>
          </Link>

          {/* Desktop nav */}
          <div className="hidden md:flex items-center gap-6">
            <Link
              href="/challenges"
              className="text-sm text-gray-400 hover:text-white transition-colors"
            >
              Challenges
            </Link>
            <Link
              href="/ranking"
              className="text-sm text-gray-400 hover:text-white transition-colors"
            >
              <div className="flex items-center gap-1">
                <Trophy className="h-4 w-4" />
                Ranking
              </div>
            </Link>

            {isAuthenticated ? (
              <div className="flex items-center gap-4">
                <Link
                  href="/profile"
                  className="flex items-center gap-2 text-sm text-gray-400 hover:text-white transition-colors"
                >
                  <User className="h-4 w-4" />
                  <span className="text-[#00ff41]">{user?.score || 0} pts</span>
                </Link>
                {user?.role === 'admin' && (
                  <Link href="/admin">
                    <Button variant="outline" size="sm">Admin</Button>
                  </Link>
                )}
                <Button variant="ghost" size="sm" onClick={logout}>
                  <LogOut className="h-4 w-4" />
                </Button>
              </div>
            ) : (
              <div className="flex items-center gap-3">
                <Link href="/login">
                  <Button variant="ghost" size="sm">Login</Button>
                </Link>
                <Link href="/register">
                  <Button size="sm">Register</Button>
                </Link>
              </div>
            )}
          </div>

          {/* Mobile menu button */}
          <button
            className="md:hidden text-gray-400"
            onClick={() => setMobileOpen(!mobileOpen)}
          >
            {mobileOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
          </button>
        </div>
      </div>

      {/* Mobile nav */}
      {mobileOpen && (
        <div className="md:hidden border-t border-[#1a1a1a] bg-[#0a0a0a]">
          <div className="space-y-1 px-4 py-4">
            <Link
              href="/challenges"
              className="flex items-center gap-2 py-2 text-sm text-gray-400"
              onClick={() => setMobileOpen(false)}
            >
              <BookOpen className="h-4 w-4" />
              Challenges
            </Link>
            <Link
              href="/ranking"
              className="flex items-center gap-2 py-2 text-sm text-gray-400"
              onClick={() => setMobileOpen(false)}
            >
              <Trophy className="h-4 w-4" />
              Ranking
            </Link>
            {isAuthenticated ? (
              <>
                <Link
                  href="/profile"
                  className="flex items-center gap-2 py-2 text-sm text-gray-400"
                  onClick={() => setMobileOpen(false)}
                >
                  <User className="h-4 w-4" />
                  Profile ({user?.score} pts)
                </Link>
                <button
                  onClick={() => { logout(); setMobileOpen(false) }}
                  className="flex items-center gap-2 py-2 text-sm text-red-400"
                >
                  <LogOut className="h-4 w-4" />
                  Logout
                </button>
              </>
            ) : (
              <div className="flex gap-3 pt-2">
                <Link href="/login" onClick={() => setMobileOpen(false)}>
                  <Button variant="ghost" size="sm">Login</Button>
                </Link>
                <Link href="/register" onClick={() => setMobileOpen(false)}>
                  <Button size="sm">Register</Button>
                </Link>
              </div>
            )}
          </div>
        </div>
      )}
    </nav>
  )
}
