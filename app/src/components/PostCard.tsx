import { useState, useEffect, useRef } from 'preact/hooks'
import { getInitials, getAgeAndDaysSince } from '../utils/utils'
import { getTopicsLabels } from '../utils/topics'
import { Icon } from './Icon'
import { icons } from '../utils/svgs'
import { useAuth } from '../services/auth'
import { removePost, likePost, unlikePost, savePost, unsavePost, POSTS_KEY } from '../services/posts'
import { mutate } from 'swr'
import { followUser, unfollowUser } from '../services/fetch'
import { Link } from 'wouter'
import styles from '../styles/PostCard.module.css'

interface Props {
  post: Post
}

export function PostCard({ post }: Props) {
  const { user: authUser, isAuth } = useAuth()
  const id = isAuth ? authUser?.id : 0
  const [liked, setLiked] = useState(post.liked)
  const [saved, setSaved] = useState(post.saved)
  const [likesCount, setLikesCount] = useState(post.likes)
  const [following, setFollowing] = useState(false)
  const [menuOpen, setMenuOpen] = useState(false)
  const menuRef = useRef<HTMLDivElement>(null)
  const isOwn = id === post.user_id

  useEffect(() => {
    if (!menuOpen) return
    function onClick(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) setMenuOpen(false)
    }
    document.addEventListener('mousedown', onClick)
    return () => document.removeEventListener('mousedown', onClick)
  }, [menuOpen])

  async function toggleFollow() {
    if (!isAuth) return
    const prev = following
    setFollowing((v) => !v)
    try {
      prev ? await unfollowUser(post.user_id) : await followUser(post.user_id)
    } catch {
      setFollowing(prev)
    }
  }

  async function toggleLike() {
    if (!isAuth) return
    const prevLiked = liked
    const prevCount = likesCount
    setLiked((v) => !v)
    setLikesCount((n) => (prevLiked ? n - 1 : n + 1))
    try {
      prevLiked ? await unlikePost(post.id) : await likePost(post.id)
    } catch {
      setLiked(prevLiked)
      setLikesCount(prevCount)
    }
  }

  async function toggleSaved() {
    if (!isAuth) return
    const prevSaved = saved
    setSaved((v) => !v)
    try {
      prevSaved ? await unsavePost(post.id) : await savePost(post.id)
    } catch {
      setSaved(prevSaved)
    }
  }

  async function handleDelete() {
    setMenuOpen(false)
    mutate(
      POSTS_KEY,
      (prev?: PostListResponse) => (prev ? { ...prev, posts: prev.posts.filter((p: Post) => p.id !== post.id) } : prev),
      false
    )
    try {
      await removePost(post.id)
    } catch {
      mutate(POSTS_KEY, undefined, true)
    }
  }

  async function handleShare() {
    setMenuOpen(false)
    const data = { title: `${post.user.name} en Topikgram`, text: post.body }
    if (navigator.share) await navigator.share(data).catch(() => {})
    else await navigator.clipboard.writeText(post.body).catch(() => {})
  }

  async function handleCopy() {
    setMenuOpen(false)
    await navigator.clipboard.writeText(post.body).catch(() => {})
  }

  function handleDownload() {
    setMenuOpen(false)
    const blob = new Blob([post.body, '\nposted by ', post.user.username, ' in Topikgram'], { type: 'text/plain' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `post-${post.id}.txt`
    a.click()
    URL.revokeObjectURL(url)
  }

  const topicLabels = getTopicsLabels(post.topics_id)

  return (
    <article class={styles.post}>
      <Link to={`/u/${post.user.username}`} class={styles.avatar}>
        {getInitials(post.user.name)}
      </Link>
      <div class={styles.content}>
        <div class={styles.top}>
          <div class={styles.head}>
            <span class={styles.name}>{post.user.name}</span>
            <span class={styles.handle}>@{post.user.username}</span>
            <span class={styles.dot}>·</span>
            <span class={styles.time}>{getAgeAndDaysSince(post.created_at)}</span>
            {!isOwn && (
              <button class={`${styles.follow} ${following ? styles.followingBtn : ''}`} onClick={toggleFollow}>
                {following ? 'Siguiendo' : 'Seguir'}
              </button>
            )}
          </div>
          <div class={styles.tools}>
            <div class={styles.topicsList}>
              {topicLabels.map((label) => (
                <span key={label} class={styles.topic}>
                  #{label}
                </span>
              ))}
            </div>
            <div class={styles.menuWrap} ref={menuRef}>
              <button
                class={styles.menuBtn}
                aria-label="Más opciones"
                aria-haspopup="menu"
                aria-expanded={menuOpen}
                onClick={() => setMenuOpen((v) => !v)}
              >
                <Icon path={icons.more} />
              </button>
              {menuOpen && (
                <div class={styles.menu} role="menu">
                  {[
                    { icon: icons.share, label: 'Compartir', action: handleShare },
                    { icon: icons.copy, label: 'Copiar', action: handleCopy },
                    { icon: icons.download, label: 'Descargar', action: handleDownload },
                    ...(isOwn ? [{ icon: icons.trash, label: 'Eliminar', action: handleDelete, danger: true }] : [])
                  ].map((item) => (
                    <button
                      key={item.label}
                      class={`${styles.menuItem} ${item.danger ? styles.menuDanger : ''}`}
                      role="menuitem"
                      onClick={item.action}
                    >
                      <Icon path={item.icon} />
                      <span>{item.label}</span>
                    </button>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
        <p class={styles.body}>{post.body}</p>
        <div class={styles.actions}>
          <button
            class={`${styles.action} ${styles.like} ${liked ? styles.liked : ''}`}
            onClick={toggleLike}
            disabled={!isAuth}
          >
            <Icon path={icons.like} fill={liked ? '#ff0000' : ''} />
            <span>{likesCount}</span>
          </button>
          <Link to={`/p/${post.id}`} class={`${styles.action} ${styles.link}`}>
            <Icon path={icons.reply} />
            <span>{post.comments}</span>
          </Link>
          <button
            class={`${styles.action} ${styles.bookmark} ${saved ? styles.saved : ''}`}
            onClick={toggleSaved}
            disabled={!isAuth}
          >
            <Icon path={icons.bookmark} fill={saved ? '#fbbf24' : ''} />
            <span>{saved ? 'Guardado' : 'Guardar'}</span>
          </button>
        </div>
      </div>
    </article>
  )
}
