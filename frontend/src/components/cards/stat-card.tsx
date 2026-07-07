import { type LucideIcon } from "lucide-react"

interface StatCardProps {
  icon: LucideIcon
  label: string
  value: number | string
  color?: string
  loading?: boolean
}

export function StatCard({ icon: Icon, label, value, color = "text-primary", loading }: StatCardProps) {
  return (
    <div className="glass rounded-xl p-5 hover:border-primary/20 transition-all">
      <div className="flex items-center gap-3 mb-2">
        <div className={`p-2 rounded-lg bg-muted ${color}`}>
          <Icon className="h-5 w-5" />
        </div>
        <span className="text-sm text-muted-foreground">{label}</span>
      </div>
      <div className="text-3xl font-bold text-foreground">
        {loading ? "..." : typeof value === "number" ? value.toLocaleString() : value}
      </div>
    </div>
  )
}
