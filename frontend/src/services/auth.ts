import api, { apiCall } from './api'
import type {
  LoginRequest,
  RegisterRequest,
  LoginResponse,
  User,
} from '@/types'

export const authApi = {
  login: (data: LoginRequest) =>
    apiCall<LoginResponse>(() => api.post('/auth/login', data)),

  register: (data: RegisterRequest) =>
    apiCall<User>(() => api.post('/auth/register', data)),

  getProfile: () =>
    apiCall<User>(() => api.get('/auth/profile')),

  updateProfile: (data: { avatar?: string; bio?: string }) =>
    apiCall<void>(() => api.put('/auth/profile', data)),
}
