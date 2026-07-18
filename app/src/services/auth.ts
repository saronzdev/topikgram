import { useCallback } from 'preact/hooks'
import useSWR, { mutate } from 'swr'
import { API_BASE, apiFetch } from './fetch'

export interface RegisterData {
  username: string
  name: string
  birthday: string
  email: string
  password: string
}

const USER_KEY = `${API_BASE}/auth/me`

export function useAuth() {
  const {
    data: user,
    error,
    isLoading
  } = useSWR<User>(USER_KEY, (url: string) => apiFetch(url), {
    revalidateOnFocus: true,
    dedupingInterval: 5000,
    onErrorRetry: (err) => {
      if (err && (err as Error & { status: number }).status === 401) return
    }
  })

  const isAuth = !!user && !(error)

  const login = useCallback(async (credentials: { identifier: string; password: string }) => {
    const res = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(credentials)
    })
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: 'Login failed' }))
      throw new Error(err.error || 'Login failed')
    }
    await mutate(USER_KEY)
  }, [])

  const register = useCallback(async (data: RegisterData) => {
    const res = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify(data)
    })
    if (!res.ok) {
      const body = await res.json().catch(() => ({ error: 'Register failed' }))
      throw new Error(body.error || 'Register failed')
    }
    await mutate(USER_KEY)
  }, [])

  const logout = useCallback(() => {
    mutate(USER_KEY, null, false)
  }, [])

  return { user, isLoading, error, login, register, logout, isAuth }
}
