import Link from "next/link"
import type { Challenge } from "@/types"
import { DIFFICULTY_COLORS } from "@/types"
import { Trophy, Clock } from "lucide-react"

interface ChallengeCardProps {
  challenge: Challenge
  solved?: boolean
}

export function ChallengeCard({ challenge, solved }: ChallengeCardProps) {
  return (
    <Link
      href={`/challenges/${challenge.id}`}
      className="glass rounded-xl p-5 hover:border-primary/20 transition-all group block"
    >
      <div className="flex items-start justify-between mb-3">
        <span className={`text-xs font-medium px-2 py-1 rounded-full ${
          DIFFICULTY_COLORS[challenge.difficulty] || "text-muted-foreground bg-muted"
        }`}>
          {challenge.difficulty}
        </span>
        <div className="flex items-center gap-3">
          {solved && <Trophy className="h-4 w-4 text-yellow-500" />}
          <span className="text-sm font-bold text-primary">{challenge.score} pts</span>
        </div>
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
        {challenge.tags?.map((tag) => (
          <span key={tag.id} className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded">
            #{tag.tagName}
          </span>
        ))}
        <span className="text-xs text-muted-foreground bg-muted px-2 py-1 rounded ml-auto flex items-center gap-1">
          <Clock className="h-3 w-3" />
          {challenge.timeoutMinutes}min
        </span>
      </div>
    </Link>
  )
}
