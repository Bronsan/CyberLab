"use client"

import { useEffect, useState } from "react"
import { useParams } from "next/navigation"
import { Navbar } from "@/components/layout/navbar"
import { Button } from "@/components/ui/button"
import { challengeApi } from "@/services/challenge"
import { useAuthStore } from "@/store/auth-store"
import type { Challenge, ContainerStatus } from "@/types"
import { DIFFICULTY_COLORS } from "@/types"
import { Terminal, Play, Square, Flag, Lightbulb, Clock, ExternalLink, Loader2 } from "lucide-react"

export default function ChallengeDetailClient() {
  const params = useParams()
  const { isAuthenticated } = useAuthStore()
  const [challenge, setChallenge] = useState<Challenge | null>(null)
  const [loading, setLoading] = useState(true)
  const [containerStatus, setContainerStatus] = useState<ContainerStatus | null>(null)
  const [instanceId, setInstanceId] = useState<number | null>(null)
  const [containerLoading, setContainerLoading] = useState(false)
  const [flag, setFlag] = useState("")
  const [flagResult, setFlagResult] = useState<{ correct: boolean; score: number } | null>(null)
  const [flagLoading, setFlagLoading] = useState(false)
  const [hint, setHint] = useState("")
  const [hintQuestion, setHintQuestion] = useState("")
  const [hintLoading, setHintLoading] = useState(false)

  useEffect(() => {
    const id = Number(params.id)
    if (!id) return

    challengeApi.getDetail(id)
      .then(setChallenge)
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [params.id])

  const startContainer = async () => {
    if (!challenge) return
    setContainerLoading(true)
    try {
      const result = await challengeApi.startContainer(challenge.id)
      setInstanceId(result.instanceId)
      setContainerStatus({ status: "RUNNING", runningTime: "0m", hostPort: 0 })
    } catch (err: any) {
      alert(err.message || "Failed to start container")
    } finally {
      setContainerLoading(false)
    }
  }

  const stopContainer = async () => {
    if (!instanceId) return
    setContainerLoading(true)
    try {
      await challengeApi.stopContainer(instanceId)
      setContainerStatus(null)
      setInstanceId(null)
    } catch (err: any) {
      alert(err.message || "Failed to stop container")
    } finally {
      setContainerLoading(false)
    }
  }

  const submitFlag = async () => {
    if (!challenge || !flag) return
    setFlagLoading(true)
    try {
      const result = await challengeApi.submitFlag(challenge.id, flag)
      setFlagResult(result)
    } catch (err: any) {
      alert(err.message || "Failed to submit flag")
    } finally {
      setFlagLoading(false)
    }
  }

  const getHint = async () => {
    if (!challenge) return
    setHintLoading(true)
    try {
      const result = await challengeApi.getAIHint({
        challengeId: challenge.id,
        question: hintQuestion || "Give me a hint to get started",
      })
      setHint(result.answer)
    } catch (err: any) {
      alert(err.message || "Failed to get hint")
    } finally {
      setHintLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-background">
        <Navbar />
        <div className="flex items-center justify-center pt-32 text-muted-foreground">Loading...</div>
      </div>
    )
  }

  if (!challenge) {
    return (
      <div className="min-h-screen bg-background">
        <Navbar />
        <div className="flex items-center justify-center pt-32 text-muted-foreground">Challenge not found</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="mx-auto max-w-7xl px-4 pt-24 pb-20">
        <div className="grid lg:grid-cols-3 gap-8">
          {/* Left: Challenge info */}
          <div className="lg:col-span-2 space-y-6">
            <div className="glass rounded-xl p-6">
              <div className="flex items-center gap-3 mb-4">
                <span className={`text-xs font-medium px-2.5 py-1 rounded-full ${DIFFICULTY_COLORS[challenge.difficulty]}`}>
                  {challenge.difficulty}
                </span>
                <span className="text-sm text-primary font-bold">{challenge.score} pts</span>
                <span className="text-sm text-muted-foreground bg-muted px-2 py-1 rounded">{challenge.category}</span>
              </div>

              <h1 className="text-2xl font-bold text-foreground mb-4">{challenge.title}</h1>

              <div className="prose prose-invert max-w-none text-muted-foreground whitespace-pre-wrap">
                {challenge.description}
              </div>

              {challenge.tags && challenge.tags.length > 0 && (
                <div className="flex gap-2 mt-6">
                  {challenge.tags.map((tag) => (
                    <span key={tag.id} className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded">
                      #{tag.tagName}
                    </span>
                  ))}
                </div>
              )}
            </div>

            {/* AI Hint Section */}
            <div className="glass rounded-xl p-6">
              <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
                <Lightbulb className="h-5 w-5 text-yellow-500" />
                AI Assistant
              </h2>
              <div className="space-y-3">
                <input
                  type="text"
                  value={hintQuestion}
                  onChange={(e) => setHintQuestion(e.target.value)}
                  placeholder="Ask for a hint... (e.g., 'How do I start?')"
                  className="w-full rounded-lg border border-border bg-card px-4 py-2.5 text-sm text-foreground placeholder-muted-foreground/50 focus:border-primary focus:outline-none transition-colors"
                />
                <Button onClick={getHint} disabled={hintLoading} variant="outline" size="sm">
                  {hintLoading ? <Loader2 className="h-4 w-4 animate-spin mr-1" /> : null}
                  Get Hint
                </Button>
                {hint && (
                  <div className="rounded-lg bg-muted p-4 text-sm text-muted-foreground">
                    {hint}
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Right: Actions */}
          <div className="space-y-6">
            {/* Container Controls */}
            <div className="glass rounded-xl p-6">
              <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
                <Terminal className="h-5 w-5 text-primary" />
                Environment
              </h2>

              {isAuthenticated ? (
                <div className="space-y-3">
                  {!containerStatus ? (
                    <Button onClick={startContainer} disabled={containerLoading}
                      className="w-full bg-primary text-primary-foreground hover:bg-primary/90">
                      {containerLoading ? (
                        <><Loader2 className="h-4 w-4 animate-spin mr-2" /> Starting...</>
                      ) : (
                        <><Play className="h-4 w-4 mr-2" /> Start Challenge</>
                      )}
                    </Button>
                  ) : (
                    <div className="space-y-3">
                      <div className="flex items-center gap-2 text-sm">
                        <span className="h-2 w-2 rounded-full bg-primary animate-pulse" />
                        <span className="text-primary">Running</span>
                        <Clock className="h-3 w-3 text-muted-foreground ml-2" />
                        <span className="text-muted-foreground">{containerStatus.runningTime}</span>
                      </div>

                      {containerStatus.hostPort > 0 && (
                        <a
                          href={`http://lab.cyberlab.com:${containerStatus.hostPort}`}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          <Button variant="outline" className="w-full border-primary/30 text-primary">
                            <ExternalLink className="h-4 w-4 mr-2" />
                            Open Environment
                          </Button>
                        </a>
                      )}

                      <Button onClick={stopContainer} disabled={containerLoading}
                        variant="destructive" className="w-full">
                        <Square className="h-4 w-4 mr-2" />
                        Stop Environment
                      </Button>
                    </div>
                  )}
                </div>
              ) : (
                <p className="text-sm text-muted-foreground">
                  <a href="/login" className="text-primary hover:underline">Sign in</a> to start this challenge
                </p>
              )}
            </div>

            {/* Flag Submission */}
            <div className="glass rounded-xl p-6">
              <h2 className="text-lg font-semibold text-foreground mb-4 flex items-center gap-2">
                <Flag className="h-5 w-5 text-red-500" />
                Submit Flag
              </h2>

              <div className="space-y-3">
                <input
                  type="text"
                  value={flag}
                  onChange={(e) => setFlag(e.target.value)}
                  placeholder="flag{...}"
                  className="font-mono w-full rounded-lg border border-border bg-card px-4 py-2.5 text-sm text-foreground placeholder-muted-foreground/50 focus:border-primary focus:outline-none transition-colors"
                />
                <Button onClick={submitFlag} disabled={flagLoading || !isAuthenticated}
                  className="w-full" variant="outline">
                  {flagLoading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                  Submit Flag
                </Button>

                {flagResult && (
                  <div className={`rounded-lg p-3 text-sm text-center ${
                    flagResult.correct
                      ? 'bg-green-500/10 text-primary border border-green-500/20'
                      : 'bg-red-500/10 text-red-500 border border-red-500/20'
                  }`}>
                    {flagResult.correct
                      ? `Correct! +${flagResult.score} points`
                      : 'Incorrect flag. Try again!'}
                  </div>
                )}
              </div>
            </div>

            {/* Challenge Info */}
            <div className="glass rounded-xl p-6">
              <h3 className="text-sm font-semibold text-muted-foreground mb-3">Challenge Info</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Category</span>
                  <span className="text-muted-foreground">{challenge.category}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Difficulty</span>
                  <span className="text-muted-foreground">{challenge.difficulty}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Points</span>
                  <span className="text-primary">{challenge.score}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Timeout</span>
                  <span className="text-muted-foreground">{challenge.timeoutMinutes} min</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
