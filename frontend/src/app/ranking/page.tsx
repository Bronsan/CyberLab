"use client"

import { useEffect, useState } from "react"
import { Navbar } from "@/components/layout/navbar"
import { challengeApi } from "@/services/challenge"
import type { RankingEntry } from "@/types"
import { Trophy, Medal } from "lucide-react"

const rankMedals = ["text-yellow-500", "text-gray-400", "text-amber-700"]

export default function RankingPage() {
  const [globalRank, setGlobalRank] = useState<RankingEntry[]>([])
  const [activeTab, setActiveTab] = useState<"global" | "weekly" | "monthly">("global")
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    setLoading(true)
    const fetchRanking = activeTab === "global"
      ? challengeApi.getGlobalRanking()
      : activeTab === "weekly"
        ? challengeApi.getWeeklyRanking()
        : challengeApi.getMonthlyRanking()

    fetchRanking
      .then((res) => setGlobalRank(res.list))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [activeTab])

  return (
    <div className="min-h-screen bg-[#0a0a0a]">
      <Navbar />
      <div className="mx-auto max-w-4xl px-4 pt-24 pb-20">
        <div className="flex items-center gap-3 mb-8">
          <Trophy className="h-8 w-8 text-yellow-500" />
          <h1 className="text-3xl font-bold text-white">Leaderboard</h1>
        </div>

        {/* Tabs */}
        <div className="flex gap-1 mb-8 bg-[#111] rounded-lg p-1 w-fit">
          {(["global", "weekly", "monthly"] as const).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                activeTab === tab
                  ? "bg-[#00ff41] text-black"
                  : "text-gray-400 hover:text-white"
              }`}
            >
              {tab.charAt(0).toUpperCase() + tab.slice(1)}
            </button>
          ))}
        </div>

        {/* Ranking List */}
        {loading ? (
          <div className="text-center py-20 text-gray-500">Loading...</div>
        ) : globalRank.length === 0 ? (
          <div className="text-center py-20 text-gray-500">No rankings yet</div>
        ) : (
          <div className="glass rounded-xl overflow-hidden">
            <div className="grid grid-cols-12 gap-4 px-6 py-4 border-b border-[#1a1a1a] text-sm text-gray-500">
              <div className="col-span-1">#</div>
              <div className="col-span-6">User</div>
              <div className="col-span-2 text-right">Score</div>
              <div className="col-span-3 text-right">Solved</div>
            </div>

            {globalRank.map((entry) => (
              <div
                key={entry.rank}
                className="grid grid-cols-12 gap-4 px-6 py-4 border-b border-[#1a1a1a] last:border-0 hover:bg-[#1a1a1a]/50 transition-colors"
              >
                <div className="col-span-1 flex items-center">
                  {entry.rank <= 3 ? (
                    <Medal className={`h-5 w-5 ${rankMedals[entry.rank - 1]}`} />
                  ) : (
                    <span className="text-sm text-gray-500">{entry.rank}</span>
                  )}
                </div>
                <div className="col-span-6 flex items-center gap-2">
                  <div className="h-8 w-8 rounded-full bg-[#1a1a1a] flex items-center justify-center text-sm text-gray-400">
                    {entry.username.charAt(0).toUpperCase()}
                  </div>
                  <span className="text-sm font-medium text-white">{entry.username}</span>
                </div>
                <div className="col-span-2 flex items-center justify-end">
                  <span className="text-sm font-bold text-[#00ff41]">{entry.score}</span>
                </div>
                <div className="col-span-3 flex items-center justify-end">
                  <span className="text-sm text-gray-400">{entry.solvedCount}</span>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
