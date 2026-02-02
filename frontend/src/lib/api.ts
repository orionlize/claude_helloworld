import axios from 'axios'
import type { User, Project, Collection, Endpoint, Environment, ApiResponse } from '@/types'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Auth API
export const authApi = {
  register: (data: { username: string; email: string; password: string }) =>
    api.post<ApiResponse<{ token: string; user: User }>>('/auth/register', data),

  login: (data: { email: string; password: string }) =>
    api.post<ApiResponse<{ token: string; user: User }>>('/auth/login', data),

  refreshToken: () =>
    api.post<ApiResponse<{ token: string; user: User }>>('/auth/refresh'),

  getCurrentUser: () =>
    api.get<ApiResponse<User>>('/users/me'),
}

// Projects API
export const projectsApi = {
  list: () =>
    api.get<ApiResponse<Project[]>>('/projects'),

  create: (data: { name: string; description: string }) =>
    api.post<ApiResponse<Project>>('/projects', data),

  get: (id: string) =>
    api.get<ApiResponse<Project>>(`/projects/${id}`),

  update: (id: string, data: { name: string; description: string }) =>
    api.put<ApiResponse<Project>>(`/projects/${id}`, data),

  delete: (id: string) =>
    api.delete<ApiResponse>(`/projects/${id}`),
}

// Collections API
export const collectionsApi = {
  list: (projectId: string) =>
    api.get<ApiResponse<Collection[]>>(`/projects/${projectId}/collections`),

  create: (projectId: string, data: { name: string; description: string; parent_id?: string }) =>
    api.post<ApiResponse<Collection>>(`/projects/${projectId}/collections`, data),

  get: (id: string) =>
    api.get<ApiResponse<Collection>>(`/collections/${id}`),

  update: (id: string, data: { name: string; description: string; parent_id?: string }) =>
    api.put<ApiResponse<Collection>>(`/collections/${id}`, data),

  delete: (id: string) =>
    api.delete<ApiResponse>(`/collections/${id}`),
}

// Endpoints API
export const endpointsApi = {
  list: (collectionId: string) =>
    api.get<ApiResponse<Endpoint[]>>(`/collections/${collectionId}/endpoints`),

  create: (collectionId: string, data: {
    name: string
    method: string
    url: string
    headers?: Record<string, string>
    body?: string
    description?: string
  }) =>
    api.post<ApiResponse<Endpoint>>(`/collections/${collectionId}/endpoints`, data),

  get: (id: string) =>
    api.get<ApiResponse<Endpoint>>(`/endpoints/${id}`),

  update: (id: string, data: {
    name: string
    method: string
    url: string
    headers?: Record<string, string>
    body?: string
    description?: string
  }) =>
    api.put<ApiResponse<Endpoint>>(`/endpoints/${id}`, data),

  delete: (id: string) =>
    api.delete<ApiResponse>(`/endpoints/${id}`),
}

// Environments API
export const environmentsApi = {
  list: (projectId: string) =>
    api.get<ApiResponse<Environment[]>>(`/projects/${projectId}/environments`),

  create: (projectId: string, data: { name: string; variables: Record<string, string>; is_default?: boolean }) =>
    api.post<ApiResponse<Environment>>(`/projects/${projectId}/environments`, data),

  get: (id: string) =>
    api.get<ApiResponse<Environment>>(`/environments/${id}`),

  update: (id: string, data: { name: string; variables: Record<string, string>; is_default?: boolean }) =>
    api.put<ApiResponse<Environment>>(`/environments/${id}`, data),

  delete: (id: string) =>
    api.delete<ApiResponse>(`/environments/${id}`),
}

// Test Request API
export const testApi = {
  send: (data: {
    method: string
    url: string
    headers?: Record<string, string>
    body?: string
  }) =>
    api.post<ApiResponse<{
      status_code: number
      status: string
      headers: Record<string, string>
      body: string
      duration: number
    }>>('/test/request', data),
}

export default api
