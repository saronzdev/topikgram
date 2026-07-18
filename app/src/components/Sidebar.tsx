import { Link } from 'wouter'
import { getInitials } from '../utils/utils'
import { icons } from '../utils/svgs'
import { Icon } from './Icon'
import { useAuth } from '../services/auth'
import styles from '../styles/Sidebar.module.css'

type IconName = keyof typeof icons

const navItems: { icon: IconName; label: string; href: string; active?: boolean }[] = [
  { icon: 'home', label: 'Inicio', href: '/', active: true },
  { icon: 'search', label: 'Explorar', href: '/' },
  { icon: 'bookmark', label: 'Guardados', href: '/' }
]

interface Props {
  theme: string
  onToggleTheme: () => void
}

export function Sidebar({ theme, onToggleTheme }: Props) {
  const { user, isAuth } = useAuth()
  const initials = user ? getInitials(user.name) : '?'

  return (
    <aside class={styles.sidebar}>
      <Link to={isAuth && user ? '/u/' + user.username : '/login'} class={styles.profile}>
        <div class={styles.avatar}>
          <span>{initials}</span>
        </div>
        <span class={styles.name}>{user ? user.name : 'Invitado'}</span>
      </Link>

      <nav class={styles.nav}>
        {navItems.map((item) => (
          <Link key={item.label} href={item.href} class={`${styles.navItem} ${item.active ? styles.active : ''}`}>
            <Icon path={icons[item.icon]} />
            <span class={styles.navLabel}>{item.label}</span>
          </Link>
        ))}
      </nav>

      {!isAuth && (
        <Link to="/login" class={styles.login}>
          Iniciar Sesión
        </Link>
      )}

      <button class={styles.themeToggle} onClick={onToggleTheme}>
        <Icon path={icons.sun} />
        <span class={styles.navLabel}>Tema</span>
        <span class={styles.knob} data-dark={theme === 'dark'} />
      </button>
    </aside>
  )
}
