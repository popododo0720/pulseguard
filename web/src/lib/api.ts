const BASE_URL = '/api'

interface RequestOptions extends RequestInit {
  params?: Record<string, string>
}

class ApiClient {
  private baseUrl: string

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl
  }

  private buildUrl(path: string, params?: Record<string, string>): string {
    const url = new URL(`${this.baseUrl}${path}`, window.location.origin)
    if (params) {
      Object.entries(params).forEach(([key, value]) => {
        url.searchParams.set(key, value)
      })
    }
    return url.toString()
  }

  async get<T>(path: string, options?: RequestOptions): Promise<T> {
    const { params, ...fetchOptions } = options ?? {}
    const response = await fetch(this.buildUrl(path, params), {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      ...fetchOptions,
    })
    if (!response.ok) throw new Error(`GET ${path} failed: ${response.statusText}`)
    return response.json()
  }

  async post<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
    const { params, ...fetchOptions } = options ?? {}
    const response = await fetch(this.buildUrl(path, params), {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: body ? JSON.stringify(body) : undefined,
      ...fetchOptions,
    })
    if (!response.ok) throw new Error(`POST ${path} failed: ${response.statusText}`)
    const text = await response.text()
    return text ? JSON.parse(text) : (undefined as unknown as T)
  }

  async put<T>(path: string, body?: unknown, options?: RequestOptions): Promise<T> {
    const { params, ...fetchOptions } = options ?? {}
    const response = await fetch(this.buildUrl(path, params), {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: body ? JSON.stringify(body) : undefined,
      ...fetchOptions,
    })
    if (!response.ok) throw new Error(`PUT ${path} failed: ${response.statusText}`)
    return response.json()
  }

  async delete<T>(path: string, options?: RequestOptions): Promise<T> {
    const { params, ...fetchOptions } = options ?? {}
    const response = await fetch(this.buildUrl(path, params), {
      method: 'DELETE',
      headers: { 'Content-Type': 'application/json' },
      ...fetchOptions,
    })
    if (!response.ok) throw new Error(`DELETE ${path} failed: ${response.statusText}`)
    const text = await response.text()
    return text ? JSON.parse(text) : (undefined as unknown as T)
  }
}

export const api = new ApiClient(BASE_URL)
