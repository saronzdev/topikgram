import { useState, useEffect } from 'preact/hooks'
import { Sidebar } from '../components/Sidebar'
import { Home } from '../components/Home'
import styles from '../styles/Layout.module.css'

export function Layout() {
  const [theme, setTheme] = useState(() => localStorage.getItem('theme') || 'light')

  useEffect(() => {
    document.documentElement.dataset.theme = theme
    localStorage.setItem('theme', theme)
  }, [theme])

  return (
    <div className={styles.layout}>
      <Sidebar theme={theme} onToggleTheme={() => setTheme((t) => (t === 'dark' ? 'light' : 'dark'))} />
      <Home />
    </div>
  )
}
