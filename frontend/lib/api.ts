import axios from 'axios'

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

export const api = axios.create({
  baseURL: `${API_URL}/api`,
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Auth
export const auth = {
  register: (email: string, password: string) =>
    api.post('/auth/register', { email, password }),
  
  login: (email: string, password: string) =>
    api.post('/auth/login', { email, password }),
}

// Projects
export const projects = {
  list: () => api.get('/projects'),
  
  create: (data: { name: string; repo_url: string; build_command?: string; output_dir?: string }) =>
    api.post('/projects', data),
  
  get: (id: string) => api.get(`/projects/${id}`),
  
  update: (id: string, data: Partial<{ name: string; repo_url: string; build_command: string; output_dir: string }>) =>
    api.put(`/projects/${id}`, data),
  
  delete: (id: string) => api.delete(`/projects/${id}`),
}

// Deployments
export const deployments = {
  trigger: (project_id: string, commit_hash?: string) =>
    api.post('/deploy', { project_id, commit_hash }),
  
  getStatus: (id: string) => api.get(`/deploy/${id}`),
  
  connectLogs: (id: string) => {
    const WS_URL = process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080'
    return new WebSocket(`${WS_URL}/api/deploy/${id}/logs`)
  },
}

