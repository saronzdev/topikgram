import { useState } from 'preact/hooks'
import { Link, useLocation } from 'wouter'
import { useAuth } from '../services/auth'
import styles from '../styles/Auth.module.css'

interface Props {
  mode: 'login' | 'register'
}

const REGISTER_FIELDS = ['username', 'email', 'name', 'birthday'] as const

const FIELD_LABELS: Record<string, string> = {
  username: 'Usuario', email: 'Email', name: 'Nombre', birthday: 'Fecha de nacimiento'
}

const FIELD_TYPES: Record<string, string> = {
  email: 'email', birthday: 'date'
}

export function Auth({ mode }: Props) {
  const [_, navigate] = useLocation()
  const { login, register, isLoading } = useAuth()
  const isLogin = mode === 'login'
  const [error, setError] = useState<string | null>(null)

  const [form, setForm] = useState({
    username: '', name: '', birthday: '', email: '', password: '', identifier: ''
  })

  async function onSubmit(e: Event) {
    e.preventDefault()
    setError(null)
    try {
      isLogin
        ? await login({ identifier: form.identifier, password: form.password })
        : await register({
            username: form.username, name: form.name, birthday: form.birthday,
            email: form.email, password: form.password
          })
      navigate('/')
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Error desconocido'
      if (msg.includes('already exists') || msg.includes('ya existe')) {
        setError('El usuario o email ya está registrado')
      } else if (msg.includes('invalid') || msg.includes('incorrect')) {
        setError('Credenciales incorrectas')
      } else {
        setError(msg)
      }
    }
  }

  return (
    <main class={styles.wrap}>
      <form class={styles.card} onSubmit={onSubmit}>
        <h1 class={styles.title}>Topikgram</h1>
        <p class={styles.subtitle}>{isLogin ? 'Inicia sesión en tu cuenta' : 'Crea una cuenta nueva'}</p>

        {isLogin ? (
          <label class={styles.field}>
            <span class={styles.label}>Usuario o email</span>
            <input class={styles.input} type="text" value={form.identifier}
              onInput={e => setForm(f => ({ ...f, identifier: e.currentTarget.value }))} required />
          </label>
        ) : (
          REGISTER_FIELDS.map(field => (
            <label key={field} class={styles.field}>
              <span class={styles.label}>{FIELD_LABELS[field]}</span>
              <input class={styles.input} type={FIELD_TYPES[field] || 'text'}
                value={form[field]} onInput={e => setForm(f => ({ ...f, [field]: e.currentTarget.value }))} required />
            </label>
          ))
        )}

        <label class={styles.field}>
          <span class={styles.label}>Contraseña</span>
          <input class={styles.input} type="password" value={form.password}
            onInput={e => setForm(f => ({ ...f, password: e.currentTarget.value }))} required />
        </label>

        {error && <p class={styles.error}>{error}</p>}

        <button class={styles.submit} type="submit" disabled={isLoading}>
          {isLoading ? 'Cargando…' : isLogin ? 'Iniciar sesión' : 'Registrarse'}
        </button>

        <p class={styles.switch}>
          {isLogin ? '¿No tienes cuenta? ' : '¿Ya tienes cuenta? '}
          <Link href={isLogin ? '/register' : '/login'} class={styles.link}>
            {isLogin ? 'Regístrate' : 'Inicia sesión'}
          </Link>
        </p>
      </form>
    </main>
  )
}
