"use client"

import { useEffect, useState } from "react"
import Link from "next/link"
import { Navbar } from "@/components/layout/navbar"
import { Button } from "@/components/ui/button"
import { challengeApi } from "@/services/challenge"
import type { Challenge } from "@/types"
import { DIFFICULTY_COLORS } from "@/types"
import { Search } from "lucide-react"

export default function ChallengesPage() {
  const [challenges, setChallenges] = useState<Challenge[]>([])
  const [total, setTotal] = useState(0)
  const [page, setPage] = useState(1)
  const [category, setCategory] = useState("")
  const [difficulty, setDifficulty] = useState("")
  const [search, setSearch] = useState("")
  const [categories, setCategories] = useState<string[]>([])
  const [loading, setLoading] = useState(true)

  const pageSize = 20

  useEffect(() => {
    challengeApi.getCategories().then(setCategories).catch(() => {})
  }, [])

  useEffect(() => {
    setLoading(true)
    challengeApi.getList({ page, pageSize, category, difficulty, search })
      .then((res) => {
        setChallenges(res.list)
        setTotal(res.total)
      })
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [page, category, difficulty, search])

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 pt-24 pb-20">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-foreground">Challenges</h1>
          <p className="text-muted-foreground mt-1">{total} challenges available</p>
        </div>

        <div className="flex flex-wrap gap-3 mb-8">
          <div className="relative flex-1 min-w-[200px]">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <input
              type="text"
              placeholder="Search challenges..."
              value={search}
              onChange={(e) => { setSearch(e.target.value); setPage(1) }}
              className="w-full rounded-lg border border-border bg-card pl-10 pr-4 py-2.5 text-sm text-foreground placeholder-muted-foreground/50 focus:border-primary focus:outline-none transition-colors"
            />
          </div>

          <select
            value={category}
            onChange={(e) => { setCategory(e.target.value); setPage(1) }}
            className="rounded-lg border border-border bg-card px-4 py-2.5 text-sm text-foreground focus:border-primary focus:outline-none transition-colors"
          >
            <option value="">All Categories</option>
            {categories.map((c) => (
              <option key={c} value={c}>{c}</option>
            ))}
          </select>

          <select
            value={difficulty}
            onChange={(e) => { setDifficulty(e.target.value); setPage(1) }}
            className="rounded-lg border border-border bg-card px-4 py-2.5 text-sm text-foreground focus:border-primary focus:outline-none transition-colors"
          >
            <option value="">All Difficulties</option>
            <option value="Easy">Easy</option>
            <option value="Medium">Medium</option>
            <option value="Hard">Hard</option>
            <option value="Insane">Insane</option>
          </select>
        </div>

        {loading ? (
          <div className="text-center py-20 text-muted-foreground">Loading challenges...</div>
        ) : challenges.length === 0 ? (
          <div className="text-center py-20 text-muted-foreground">No challenges found</div>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {challenges.map((challenge) => (
              <Link
                key={challenge.id}
                href={`/challenges/${challenge.id}`}
                className="glass rounded-xl p-5 hover:border-primary/20 transition-all group"
              >
                <div className="flex items-start justify-between mb-3">
                  <span className={`text-xs font-medium px-2 py-1 rounded-full ${DIFFICULTY_COLORS[challenge.difficulty] || 'text-muted-foreground bg-muted'}`}>
                    {challenge.difficulty}
                  </span>
                  <span className="text-sm font-bold text-primary">{challenge.score} pts</span>
                </div>

                <h3 className="text-lg font-semibold text-foreground mb-2 group-hover:text-primary transition-colors">
                  {challenge.title}
                </h3>

                <p className="text-sm text-muted-foreground line-clamp-2 mb-4">
                  {challenge.description}
                </p>

                <div className="flex items-center gap-2 flex-wrap">
                  <span className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded">
                    {challenge.category}
                  </span>
                  {challenge.tags?.slice(0, 3).map((tag) => (
                    <span key={tag.id} className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded">
                      {tag.tagName}
                    </span>
                  ))}
                </div>
              </Link>
            ))}
          </div>
        )}

        {total > pageSize && (
          <div className="flex justify-center gap-2 mt-8">
            <Button variant="outline" size="sm" disabled={page <= 1}
              onClick={() => setPage(page - 1)}>Previous</Button>
            <span className="flex items-center px-4 text-sm text-muted-foreground">
              Page {page} of {Math.ceil(total / pageSize)}
            </span>
            <Button variant="outline" size="sm" disabled={page >= Math.ceil(total / pageSize)}
              onClick={() => setPage(page + 1)}>Next</Button>
          </div>
        )}
      </div>
    </div>
  )
}
