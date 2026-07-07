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
      <div className="min-h-screen bg-[#0a0a0a]">
        <Navbar />
        <div className="flex items-center justify-center pt-32 text-gray-500">Loading...</div>
      </div>
    )
  }

  if (!challenge) {
    return (
      <div className="min-h-screen bg-[#0a0a0a]">
        <Navbar />
        <div className="flex items-center justify-center pt-32 text-gray-500">Challenge not found</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-[#0a0a0a]">
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
                <span className="text-sm text-[#00ff41] font-bold">{challenge.score} pts</span>
                <span className="text-sm text-gray-500 bg-[#1a1a1a] px-2 py-1 rounded">{challenge.category}</span>
              </div>

              <h1 className="text-2xl font-bold text-white mb-4">{challenge.title}</h1>

              <div className="prose prose-invert max-w-none text-gray-400 whitespace-pre-wrap">
                {challenge.description}
              </div>

              {challenge.tags && challenge.tags.length > 0 && (
                <div className="flex gap-2 mt-6">
                  {challenge.tags.map((tag) => (
                    <span key={tag.id} className="text-xs text-gray-500 bg-[#1a1a1a] px-2 py-1 rounded">
                      #{tag.tagName}
                    </span>
                  ))}
                </div>
              )}
            </div>

            {/* AI Hint Section */}
            <div className="glass rounded-xl p-6">
              <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <Lightbulb className="h-5 w-5 text-yellow-500" />
                AI Assistant
              </h2>
              <div className="space-y-3">
                <input
                  type="text"
                  value={hintQuestion}
                  onChange={(e) => setHintQuestion(e.target.value)}
                  placeholder="Ask for a hint... (e.g., 'How do I start?')"
                  className="w-full rounded-lg border border-[#262626] bg-[#111] px-4 py-2.5 text-sm text-white placeholder-gray-600 focus:border-[#00ff41] focus:outline-none"
                />
                <Button onClick={getHint} disabled={hintLoading} variant="outline" size="sm">
                  {hintLoading ? <Loader2 className="h-4 w-4 animate-spin mr-1" /> : null}
                  Get Hint
                </Button>
                {hint && (
                  <div className="rounded-lg bg-[#1a1a1a] p-4 text-sm text-gray-300">
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
              <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <Terminal className="h-5 w-5 text-[#00ff41]" />
                Environment
              </h2>

              {isAuthenticated ? (
                <div className="space-y-3">
                  {!containerStatus ? (
                    <Button onClick={startContainer} disabled={containerLoading}
                      className="w-full bg-[#00ff41] text-black hover:bg-[#00ff41]/90">
                      {containerLoading ? (
                        <><Loader2 className="h-4 w-4 animate-spin mr-2" /> Starting...</>
                      ) : (
                        <><Play className="h-4 w-4 mr-2" /> Start Challenge</>
                      )}
                    </Button>
                  ) : (
                    <div className="space-y-3">
                      <div className="flex items-center gap-2 text-sm">
                        <span className="h-2 w-2 rounded-full bg-green-500 animate-pulse" />
                        <span className="text-green-500">Running</span>
                        <Clock className="h-3 w-3 text-gray-500 ml-2" />
                        <span className="text-gray-500">{containerStatus.runningTime}</span>
                      </div>

                      {containerStatus.hostPort > 0 && (
                        <a
                          href={`http://lab.cyberlab.com:${containerStatus.hostPort}`}
                          target="_blank"
                          rel="noopener noreferrer"
                        >
                          <Button variant="outline" className="w-full border-[#00ff41]/30 text-[#00ff41]">
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
                <p className="text-sm text-gray-500">
                  <a href="/login" className="text-[#00ff41] hover:underline">Sign in</a> to start this challenge
                </p>
              )}
            </div>

            {/* Flag Submission */}
            <div className="glass rounded-xl p-6">
              <h2 className="text-lg font-semibold text-white mb-4 flex items-center gap-2">
                <Flag className="h-5 w-5 text-red-500" />
                Submit Flag
              </h2>

              <div className="space-y-3">
                <input
                  type="text"
                  value={flag}
                  onChange={(e) => setFlag(e.target.value)}
                  placeholder="flag{...}"
                  className="font-mono w-full rounded-lg border border-[#262626] bg-[#111] px-4 py-2.5 text-sm text-white placeholder-gray-600 focus:border-[#00ff41] focus:outline-none"
                />
                <Button onClick={submitFlag} disabled={flagLoading || !isAuthenticated}
                  className="w-full" variant="outline">
                  {flagLoading ? <Loader2 className="h-4 w-4 animate-spin mr-2" /> : null}
                  Submit Flag
                </Button>

                {flagResult && (
                  <div className={`rounded-lg p-3 text-sm text-center ${
                    flagResult.correct
                      ? 'bg-green-500/10 text-green-500 border border-green-500/20'
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
              <h3 className="text-sm font-semibold text-gray-400 mb-3">Challenge Info</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-gray-500">Category</span>
                  <span className="text-gray-300">{challenge.category}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-500">Difficulty</span>
                  <span className="text-gray-300">{challenge.difficulty}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-500">Points</span>
                  <span className="text-[#00ff41]">{challenge.score}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-500">Timeout</span>
                  <span className="text-gray-300">{challenge.timeoutMinutes} min</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
