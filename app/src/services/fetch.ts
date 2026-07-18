export const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:3000/api/v1'

export async function apiFetch(url: string, options: RequestInit = {}) {
  const res = await fetch(url, {
    ...options,
    credentials: 'include',
    headers: { ...(options.headers as Record<string, string>) }
  })
  if (!res.ok) {
    const body = await res.json().catch(() => ({ error: res.statusText }))
    const err = new Error(body.error || res.statusText) as Error & { status: number }
    err.status = res.status
    throw err
  }
  const text = await res.text()
  return text ? JSON.parse(text) : {}
}

export const fetcher = (url: string) => apiFetch(url)

export async function getUser(username: string): Promise<UserPublic | null> {
  try {
    return await apiFetch(`${API_BASE}/users?username=${username}`)
  } catch {
    return null
  }
}

export async function getUserByID(id: number): Promise<UserPublic | null> {
  try {
    return await apiFetch(`${API_BASE}/users/${id}`)
  } catch {
    return null
  }
}

export async function followUser(id: number) {
  await apiFetch(`${API_BASE}/users/${id}/follow`, { method: 'POST' })
}

export async function unfollowUser(id: number) {
  await apiFetch(`${API_BASE}/users/${id}/follow`, { method: 'DELETE' })
}

export async function getFollowers(id: number): Promise<UserPublic[]> {
  return await apiFetch(`${API_BASE}/users/${id}/followers`)
}

export async function getFollowing(id: number): Promise<UserPublic[]> {
  return await apiFetch(`${API_BASE}/users/${id}/following`)
}
