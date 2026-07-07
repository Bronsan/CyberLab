// Server Component — provides generateStaticParams for static export
import ChallengeDetailClient from "./client"

export function generateStaticParams() {
  // Pre-generate some sample challenge pages for static export
  // In production, this should be populated from the API
  return [{ id: "1" }, { id: "2" }, { id: "3" }]
}

export default function ChallengeDetailPage() {
  return <ChallengeDetailClient />
}
