"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { Navbar } from "@/components/layout/navbar"
import { Button } from "@/components/ui/button"
import { StatCard } from "@/components/cards/stat-card"
import { useAuthStore } from "@/store/auth-store"
import { challengeApi } from "@/services/challenge"
import type { AdminDashboard, ContainerInstance } from "@/types"
import {
  Shield, Users, Server, BookOpen,
  Megaphone, Plus, Trash2, RefreshCw,
  LayoutDashboard, Container, Terminal
} from "lucide-react"

type TabKey = "overview" | "users" | "challenges" | "containers" | "announcements"

const tabs: { key: TabKey; label: string; icon: any }[] = [
  { key: "overview", label: "Overview", icon: LayoutDashboard },
  { key: "users", label: "Users", icon: Users },
  { key: "challenges", label: "Challenges", icon: Server },
  { key: "containers", label: "Containers", icon: Container },
  { key: "announcements", label: "Notices", icon: Megaphone },
]

export default function AdminPage() {
  const { user, isAuthenticated } = useAuthStore()
  const [activeTab, setActiveTab] = useState<TabKey>("overview")
  const [dashboard, setDashboard] = useState<AdminDashboard | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!isAuthenticated || user?.role !== "admin") return

    // Fetch via our service - for now use a simple GET
    fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/v1/admin/dashboard`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("cyberlab-auth")
          ? JSON.parse(localStorage.getItem("cyberlab-auth") || "{}")?.state?.token
          : ""}`,
      },
    })
      .then((r) => r.json())
      .then((res) => {
        if (res.code === 200) setDashboard(res.data)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [isAuthenticated, user])

  if (!isAuthenticated || user?.role !== "admin") {
    return (
      <div className="min-h-screen bg-background">
        <Navbar />
        <div className="flex flex-col items-center justify-center pt-32 gap-4">
          <Shield className="h-16 w-16 text-muted-foreground" />
          <h1 className="text-2xl font-bold text-foreground">Access Denied</h1>
          <p className="text-muted-foreground">You need admin privileges to access this page.</p>
          <Link href="/"><Button variant="outline">Go Home</Button></Link>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 pt-24 pb-20">
        <div className="flex items-center gap-3 mb-8">
          <Shield className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold text-foreground">Admin Panel</h1>
        </div>

        {/* Tabs */}
        <div className="flex gap-2 mb-8 overflow-x-auto pb-2">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              onClick={() => setActiveTab(tab.key)}
              className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium whitespace-nowrap transition-colors ${
                activeTab === tab.key
                  ? "bg-primary text-primary-foreground"
                  : "bg-muted text-muted-foreground hover:text-foreground"
              }`}
            >
              <tab.icon className="h-4 w-4" />
              {tab.label}
            </button>
          ))}
        </div>

        {/* Tab Content */}
        {activeTab === "overview" && <OverviewTab dashboard={dashboard} loading={loading} />}
        {activeTab === "users" && <UsersTab />}
        {activeTab === "challenges" && <ChallengesTab />}
        {activeTab === "containers" && <ContainersTab />}
        {activeTab === "announcements" && <AnnouncementsTab />}
      </div>
    </div>
  )
}

function OverviewTab({ dashboard, loading }: { dashboard: AdminDashboard | null; loading: boolean }) {
  return (
    <div>
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <StatCard icon={Users} label="Total Users" value={dashboard?.totalUsers ?? 0} color="text-blue-500" loading={loading} />
        <StatCard icon={Server} label="Total Challenges" value={dashboard?.totalChallenges ?? 0} color="text-primary" loading={loading} />
        <StatCard icon={Container} label="Running Containers" value={dashboard?.runningContainers ?? 0} color="text-yellow-500" loading={loading} />
        <StatCard icon={Terminal} label="System Status" value="Online" color="text-green-500" />
      </div>

      <div className="glass rounded-xl p-6">
        <h2 className="text-lg font-semibold text-foreground mb-4">Quick Actions</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
          {[
            { label: "New Challenge", icon: Plus, action: () => {} },
            { label: "New Notice", icon: Megaphone, action: () => {} },
            { label: "Refresh", icon: RefreshCw, action: () => window.location.reload() },
            { label: "View Logs", icon: BookOpen, action: () => {} },
          ].map((action) => (
            <button
              key={action.label}
              onClick={action.action}
              className="flex flex-col items-center gap-2 p-4 rounded-xl bg-muted hover:bg-muted/70 transition-colors"
            >
              <action.icon className="h-6 w-6 text-primary" />
              <span className="text-sm text-muted-foreground">{action.label}</span>
            </button>
          ))}
        </div>
      </div>
    </div>
  )
}

function UsersTab() {
  return (
    <div className="glass rounded-xl p-8 text-center">
      <Users className="h-12 w-12 mx-auto mb-3 text-muted-foreground" />
      <h3 className="text-lg font-semibold text-foreground mb-2">User Management</h3>
      <p className="text-muted-foreground text-sm">
        Connect to the backend API to manage user accounts, view activity, and modify permissions.
      </p>
      <p className="text-xs text-muted-foreground mt-2 font-mono">
        GET /api/v1/admin/users · GET /api/v1/admin/logs
      </p>
    </div>
  )
}

function ChallengesTab() {
  return (
    <div className="glass rounded-xl p-8 text-center">
      <Server className="h-12 w-12 mx-auto mb-3 text-muted-foreground" />
      <h3 className="text-lg font-semibold text-foreground mb-2">Challenge Management</h3>
      <p className="text-muted-foreground text-sm">
        Create, edit, and manage challenges. Each challenge requires a Docker image and a valid flag.
      </p>
      <p className="text-xs text-muted-foreground mt-2 font-mono">
        POST/PUT/DELETE /api/v1/admin/challenge
      </p>
      <Link href="/admin/challenge/new">
        <Button size="sm" className="mt-4"><Plus className="h-4 w-4 mr-1" /> New Challenge</Button>
      </Link>
    </div>
  )
}

function ContainersTab() {
  return (
    <div className="glass rounded-xl p-8 text-center">
      <Container className="h-12 w-12 mx-auto mb-3 text-muted-foreground" />
      <h3 className="text-lg font-semibold text-foreground mb-2">Container Management</h3>
      <p className="text-muted-foreground text-sm">
        View all running challenge instances, monitor resource usage, and forcefully destroy containers.
      </p>
      <p className="text-xs text-muted-foreground mt-2 font-mono">
        GET /api/v1/admin/container · DELETE /api/v1/admin/container/:id
      </p>
    </div>
  )
}

function AnnouncementsTab() {
  return (
    <div className="glass rounded-xl p-8 text-center">
      <Megaphone className="h-12 w-12 mx-auto mb-3 text-muted-foreground" />
      <h3 className="text-lg font-semibold text-foreground mb-2">Announcements</h3>
      <p className="text-muted-foreground text-sm">
        Post and manage system announcements visible to all users on the platform.
      </p>
      <p className="text-xs text-muted-foreground mt-2 font-mono">
        POST /api/v1/admin/announcement · DELETE /api/v1/admin/announcement/:id
      </p>
    </div>
  )
}
