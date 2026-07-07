// API Response types
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data?: T
}

export interface ErrorResponse {
  code: number
  message: string
}

// User types
export interface User {
  id: number
  username: string
  email: string
  avatar?: string
  bio?: string
  role: 'user' | 'admin'
  score: number
  solvedCount: number
  createdAt?: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

// Challenge types
export interface Challenge {
  id: number
  title: string
  description: string
  category: string
  difficulty: 'Easy' | 'Medium' | 'Hard' | 'Insane'
  score: number
  dockerImage?: string
  timeoutMinutes: number
  isActive: boolean
  createdAt: string
  tags?: ChallengeTag[]
}

export interface ChallengeTag {
  id: number
  challengeId: number
  tagName: string
}

export interface ChallengeListResponse {
  list: Challenge[]
  total: number
}

// Container types
export interface ContainerInstance {
  id: number
  userId: number
  challengeId: number
  containerId?: string
  hostPort: number
  status: 'CREATED' | 'RUNNING' | 'STOPPED' | 'DESTROYED' | 'ERROR'
  startTime: string
  expireTime: string
  challenge?: Challenge
}

export interface StartContainerResponse {
  instanceId: number
  url: string
}

export interface ContainerStatus {
  status: string
  runningTime: string
  hostPort: number
}

// Submission types
export interface Submission {
  id: number
  userId: number
  challengeId: number
  submittedFlag: string
  isCorrect: boolean
  submittedAt: string
  challenge?: Challenge
}

export interface SubmitFlagResponse {
  correct: boolean
  score: number
}

// Ranking types
export interface RankingEntry {
  rank: number
  username: string
  score: number
  solvedCount: number
}

// Announcement types
export interface Announcement {
  id: number
  title: string
  content: string
  createdAt: string
}

// AI types
export interface AIHintRequest {
  challengeId: number
  question: string
}

export interface AIHintResponse {
  answer: string
}

export interface Vulnerability {
  type: string
  severity: string
  file: string
  line: number
  description: string
  suggestion: string
}

export interface AuditResponse {
  riskCount: number
  vulnerabilities: Vulnerability[]
}

// Admin types
export interface AdminDashboard {
  totalUsers: number
  totalChallenges: number
  runningContainers: number
}

// Difficulty color mapping
export const DIFFICULTY_COLORS: Record<string, string> = {
  Easy: 'text-green-500 bg-green-500/10',
  Medium: 'text-yellow-500 bg-yellow-500/10',
  Hard: 'text-red-500 bg-red-500/10',
  Insane: 'text-purple-500 bg-purple-500/10',
}
