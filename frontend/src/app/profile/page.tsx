"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { Navbar } from "@/components/layout/navbar"
import { Button } from "@/components/ui/button"
import { StatCard } from "@/components/cards/stat-card"
import { ChallengeCard } from "@/components/cards/challenge-card"
import { GrowthChart } from "@/components/charts/growth-chart"
import { useAuthStore } from "@/store/auth-store"
import { challengeApi } from "@/services/challenge"
import type { Challenge } from "@/types"
import { User, Trophy, Target, TrendingUp, Shield, Settings, LogOut } from "lucide-react"

export default function ProfilePage() {
  const { user, isAuthenticated, logout } = useAuthStore()
  const [solved, setSolved] = useState<Challenge[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isAuthenticated) return
    challengeApi.getSolvedChallenges()
      .then((data: any) => setSolved(data || []))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [isAuthenticated])

  if (!isAuthenticated || !user) {
    return (
      <div className="min-h-screen bg-background">
        <Navbar />
        <div className="flex items-center justify-center pt-32 text-muted-foreground">
          <Link href="/login" className="text-primary hover:underline">Sign in to view your profile</Link>
        </div>
      </div>
    )
  }

  const mockGrowth = [
    { date: "Week 1", score: 100 },
    { date: "Week 2", score: 250 },
    { date: "Week 3", score: 400 },
    { date: "Week 4", score: user.score || 0 },
  ]

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 pt-24 pb-20">
        {/* Header */}
        <div className="glass rounded-2xl p-8 mb-8">
          <div className="flex flex-col sm:flex-row items-start sm:items-center gap-6">
            <div className="h-20 w-20 rounded-full bg-primary/10 flex items-center justify-center text-2xl font-bold text-primary">
              {user.username.charAt(0).toUpperCase()}
            </div>
            <div className="flex-1">
              <h1 className="text-2xl font-bold text-foreground">{user.username}</h1>
              <p className="text-muted-foreground">{user.email}</p>
              {user.bio && <p className="text-sm text-muted-foreground mt-1">{user.bio}</p>}
              <div className="flex gap-2 mt-3">
                <span className={`text-xs px-2 py-1 rounded-full ${user.role === "admin" ? "bg-primary/10 text-primary" : "bg-muted text-muted-foreground"}`}>
                  {user.role}
                </span>
                {user.role === "admin" && (
                  <Link href="/admin">
                    <Button variant="outline" size="sm">
                      <Shield className="h-3 w-3 mr-1" />
                      Admin
                    </Button>
                  </Link>
                )}
              </div>
            </div>
            <div className="flex gap-2">
              <Link href="/settings">
                <Button variant="ghost" size="sm"><Settings className="h-4 w-4" /></Button>
              </Link>
              <Button variant="ghost" size="sm" onClick={logout}>
                <LogOut className="h-4 w-4 mr-1" /> Logout
              </Button>
            </div>
          </div>
        </div>

        {/* Stats */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
          <StatCard icon={Trophy} label="Score" value={user.score} color="text-yellow-500" />
          <StatCard icon={Target} label="Solved" value={user.solvedCount} color="text-primary" />
          <StatCard icon={TrendingUp} label="Rank" value="--" color="text-blue-500" />
          <StatCard icon={Shield} label="Role" value={user.role} color="text-purple-500" />
        </div>

        {/* Growth Chart */}
        <div className="glass rounded-xl p-6 mb-8">
          <h2 className="text-lg font-semibold text-foreground mb-4">Score Growth</h2>
          <GrowthChart
            data={mockGrowth}
            color="var(--primary)"
            loading={loading}
          />
        </div>

        {/* Solved Challenges */}
        <div>
          <h2 className="text-xl font-bold text-foreground mb-4">
            Completed Challenges ({solved.length})
          </h2>
          {loading ? (
            <div className="text-muted-foreground text-center py-12">Loading...</div>
          ) : solved.length === 0 ? (
            <div className="glass rounded-xl p-12 text-center">
              <Trophy className="h-12 w-12 mx-auto mb-3 text-muted-foreground" />
              <p className="text-muted-foreground">No challenges solved yet</p>
              <Link href="/challenges">
                <Button variant="outline" size="sm" className="mt-4">Browse Challenges</Button>
              </Link>
            </div>
          ) : (
            <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
              {solved.map((ch) => (
                <ChallengeCard key={ch.id} challenge={ch} solved />
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
