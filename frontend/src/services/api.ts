const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? ''
const TOKEN_KEY = 'birdnest_token'

interface ApiErrorBody {
  msg?: string
  message?: string
  error?: string
}

interface LoginResponse {
  token: string
  msg: string
}

interface MessageResponse {
  msg: string
}

interface PublishResponse extends MessageResponse {
  author: string
  tags: Array<{ tagName: string }>
}

export interface Passage {
  ID: number
  content: string
  tags?: Array<{ tagName: string }>
  CreatedAt?: string
  UpdatedAt?: string
}

interface PassageListResponse extends MessageResponse {
  data: Passage[]
}

interface PassageResponse extends MessageResponse {
  data: Passage
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })

  const body = (await response.json().catch(() => ({}))) as T & ApiErrorBody

  if (!response.ok) {
    throw new Error(body.msg || body.message || body.error || '请求失败，请稍后重试')
  }

  return body
}

export const auth = {
  token: () => localStorage.getItem(TOKEN_KEY),
  save: (token: string) => localStorage.setItem(TOKEN_KEY, token),
  clear: () => localStorage.removeItem(TOKEN_KEY),
}

export const api = {
  login(email: string, password: string) {
    return request<LoginResponse>('/api/v1/user/login', {
      method: 'POST',
      body: JSON.stringify({ email, password }),
    })
  },

  sendCode(name: string, email: string) {
    return request<MessageResponse>('/api/v1/code', {
      method: 'POST',
      body: JSON.stringify({ name, email }),
    })
  },

  register(name: string, email: string, password: string, code: string) {
    return request<MessageResponse>('/api/v1/user', {
      method: 'POST',
      body: JSON.stringify({ name, email, password, code }),
    })
  },

  publish(content: string, tags: string[]) {
    return request<PublishResponse>('/api/v1/passage', {
      method: 'POST',
      headers: {
        token: auth.token() ?? '',
      },
      body: JSON.stringify({
        content,
        tags: tags.map((tagName) => ({ tagName })),
      }),
    })
  },

  getAllPassages() {
    return request<PassageListResponse>('/api/v1/passage/all', {
      headers: {
        token: auth.token() ?? '',
      },
    })
  },

  getPassage(id: number) {
    return request<PassageResponse>(`/api/v1/passage/${id}`, {
      headers: {
        token: auth.token() ?? '',
      },
    })
  },

  getPassagesByTag(tag: string) {
    return request<PassageListResponse>(`/api/v1/passage/tag/${encodeURIComponent(tag)}`, {
      headers: {
        token: auth.token() ?? '',
      },
    })
  },

  updatePassage(id: number, content: string, tags: string[]) {
    return request<PassageResponse>(`/api/v1/passage/${id}`, {
      method: 'PUT',
      headers: {
        token: auth.token() ?? '',
      },
      body: JSON.stringify({
        content,
        tags: tags.map((tagName) => ({ tagName })),
      }),
    })
  },

  deletePassage(id: number) {
    return request<MessageResponse>(`/api/v1/passage/${id}`, {
      method: 'DELETE',
      headers: {
        token: auth.token() ?? '',
      },
    })
  },
}
