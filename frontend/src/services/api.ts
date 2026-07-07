import axios from 'axios'
import type { ApiResponse } from '@/types'

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

const api = axios.create({
  baseURL: `${BASE_URL}/api/v1`,
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor - add auth token
api.interceptors.request.use((config) => {
  if (typeof window !== 'undefined') {
    const stored = localStorage.getItem('cyberlab-auth')
    if (stored) {
      try {
        const { state } = JSON.parse(stored)
        if (state.token) {
          config.headers.Authorization = `Bearer ${state.token}`
        }
      } catch {
        // ignore parse errors
      }
    }
  }
  return config
})

// Response interceptor - handle errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Token expired or invalid
      if (typeof window !== 'undefined') {
        localStorage.removeItem('cyberlab-auth')
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

// Helper to unwrap API response
export async function apiCall<T>(
  call: () => Promise<{ data: ApiResponse<T> }>
): Promise<T> {
  const response = await call()
  const data = response.data

  if (data.code !== 200) {
    throw new Error(data.message || 'Unknown error')
  }

  return data.data as T
}

export default api
