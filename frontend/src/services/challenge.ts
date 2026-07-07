import api, { apiCall } from './api'
import type {
  Challenge,
  ChallengeListResponse,
  StartContainerResponse,
  ContainerStatus,
  SubmitFlagResponse,
  Submission,
  RankingEntry,
  Announcement,
  AIHintRequest,
  AIHintResponse,
  AuditResponse,
} from '@/types'

export const challengeApi = {
  // Challenges
  getList: (params?: {
    page?: number
    pageSize?: number
    category?: string
    difficulty?: string
    search?: string
  }) =>
    apiCall<ChallengeListResponse>(() => api.get('/challenges', { params })),

  getDetail: (id: number) =>
    apiCall<Challenge>(() => api.get(`/challenges/${id}`)),

  search: (q: string) =>
    apiCall<Challenge[]>(() => api.get('/challenges/search', { params: { q } })),

  getCategories: () =>
    apiCall<string[]>(() => api.get('/challenges/categories')),

  // Container
  startContainer: (challengeId: number) =>
    apiCall<StartContainerResponse>(() => api.post('/container/start', { challengeId })),

  getContainerStatus: (instanceId: number) =>
    apiCall<ContainerStatus>(() => api.get(`/container/status/${instanceId}`)),

  stopContainer: (instanceId: number) =>
    apiCall<void>(() => api.post('/container/stop', { instanceId })),

  getMyInstances: () =>
    apiCall<unknown[]>(() => api.get('/container/instances')),

  // Submission
  submitFlag: (challengeId: number, flag: string) =>
    apiCall<SubmitFlagResponse>(() => api.post('/submit', { challengeId, flag })),

  getSubmissionHistory: (params?: { page?: number; pageSize?: number }) =>
    apiCall<{ list: Submission[]; total: number }>(() =>
      api.get('/submit/history', { params })
    ),

  getSolvedChallenges: () =>
    apiCall<unknown[]>(() => api.get('/submit/solved')),

  // Ranking
  getGlobalRanking: () =>
    apiCall<{ list: RankingEntry[] }>(() => api.get('/ranking/global')),

  getWeeklyRanking: () =>
    apiCall<{ list: RankingEntry[] }>(() => api.get('/ranking/week')),

  getMonthlyRanking: () =>
    apiCall<{ list: RankingEntry[] }>(() => api.get('/ranking/month')),

  // Announcements
  getAnnouncements: (params?: { page?: number; pageSize?: number }) =>
    apiCall<{ list: Announcement[]; total: number }>(() =>
      api.get('/announcement', { params })
    ),

  getAnnouncement: (id: number) =>
    apiCall<Announcement>(() => api.get(`/announcement/${id}`)),

  // AI
  getAIHint: (data: AIHintRequest) =>
    apiCall<AIHintResponse>(() => api.post('/ai/hint', data)),

  auditCode: (file: File) =>
    apiCall<AuditResponse>(() => {
      const formData = new FormData()
      formData.append('file', file)
      return api.post('/ai/audit', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      })
    }),
}
