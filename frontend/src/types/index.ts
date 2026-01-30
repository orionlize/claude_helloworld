export interface User {
  id: string
  username: string
  email: string
  created_at: string
  updated_at: string
}

export interface Project {
  id: string
  name: string
  description: string
  user_id: string
  created_at: string
  updated_at: string
}

export interface Collection {
  id: string
  project_id: string
  name: string
  description: string
  parent_id: string | null
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Endpoint {
  id: string
  collection_id: string
  name: string
  method: string
  url: string
  headers: Record<string, string>
  body: string | null
  description: string
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Environment {
  id: string
  project_id: string
  name: string
  variables: Record<string, string>
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface ApiResponse<T = any> {
  success: boolean
  message?: string
  data?: T
  error?: string
}
