import { useEffect, useState } from 'preact/hooks'
import { useAuth } from '../services/auth'
import { useParams, Link } from 'wouter'
import { getCommentsByUserID, usePosts } from '../services/posts'
import { getUser, followUser, unfollowUser, getFollowers, getFollowing } from '../services/fetch'
import { PostCard } from './PostCard'
import { Composer } from './Composer'
import { Icon } from './Icon'
import { icons } from '../utils/svgs'
import { getInitials } from '../utils/utils'
import styles from '../styles/UserProfile.module.css'
import { CommentCard } from './CommentCard'

interface ListModal {
  type: 'followers' | 'following'
  users: UserPublic[]
}

export function UserProfile() {
  const { username } = useParams() as { username: string }
  const { user: currentUser, isAuth } = useAuth()
  const { posts } = usePosts()
  const [user, setUser] = useState<UserPublic | null>(null)
  const [userComments, setUserComments] = useState<CommentInterface[]>([])
  const [isLoadingComments, setIsLoadingComments] = useState(false)
  const [tabIdx, setTabIdx] = useState(0)
  const [following, setFollowing] = useState(false)
  const [followersCount, setFollowersCount] = useState(0)
  const [followingCount, setFollowingCount] = useState(0)
  const [listModal, setListModal] = useState<ListModal | null>(null)

  useEffect(() => {
    if (!username) return
    getUser(username).then((u) => {
      if (!u) return
      setUser(u)
      getFollowers(u.id)
        .then((f) => setFollowersCount(f.length))
        .catch(() => {})
      getFollowing(u.id)
        .then((f) => setFollowingCount(f.length))
        .catch(() => {})
      if (currentUser) {
        getFollowing(currentUser.id)
          .then((f) => {
            setFollowing(f.some((uu) => uu.id === user!.id))
          })
          .catch(() => {})
      }
    })
  }, [username, currentUser])

  useEffect(() => {
    if (!user || tabIdx !== 1) return
    setIsLoadingComments(true)
    getCommentsByUserID(user.id)
      .then((c) => {
        console.log(c)
        setUserComments(c)
      })
      .finally(() => setIsLoadingComments(false))
  }, [tabIdx])

  async function toggleFollow() {
    if (!isAuth || !user) return
    const prev = following
    setFollowing((v) => !v)
    setFollowersCount((n) => (prev ? n - 1 : n + 1))
    try {
      prev ? await unfollowUser(user.id) : await followUser(user.id)
    } catch {
      setFollowing(prev)
      setFollowersCount((n) => (prev ? n + 1 : n - 1))
    }
  }

  async function openList(type: 'followers' | 'following') {
    if (!user) return
    try {
      const users = type === 'followers' ? await getFollowers(user.id) : await getFollowing(user.id)
      setListModal({ type, users })
    } catch {}
  }

  const userPosts = Array.isArray(posts) && user ? posts.filter((p) => p.user.username === username) : []
  const userLikedPosts = Array.isArray(posts) && user ? posts.filter((p) => p.liked) : []

  if (!user) return <p class={styles.state}>Cargando…</p>

  const isOwn = isAuth && currentUser?.username === username

  return (
    <>
      <section class={styles.section}>
        <Link to="/" class={styles.back}>
          <Icon path={icons.arrowLeft} />
        </Link>
        <div class={styles.info}>
          <span class={styles.avatar}>{getInitials(user.name)}</span>
          <div class={styles.data}>
            <p class={styles.name}>{user.name}</p>
            <p class={styles.username}>@{user.username}</p>
          </div>
        </div>
        <div class={styles.stats}>
          <p class={styles.statItem}>{userPosts.length} Posts</p>
          <p class={`${styles.statItem} ${styles.statClickable}`} onClick={() => openList('followers')}>
            {followersCount} Seguidores
          </p>
          <p class={`${styles.statItem} ${styles.statClickable}`} onClick={() => openList('following')}>
            {followingCount} Seguidos
          </p>
        </div>
        {!isOwn && (
          <button class={`${styles.followBtn} ${following ? styles.followingBtn : ''}`} onClick={toggleFollow}>
            {following ? 'Siguiendo' : 'Seguir'}
          </button>
        )}
      </section>
      <section class={styles.content}>
        <header class={styles.tabs}>
          {['Publicaciones', 'Respuestas', 'Me gustas'].map((label, i) => (
            <p key={label} class={`${styles.tab} ${tabIdx === i ? styles.active : ''}`} onClick={() => setTabIdx(i)}>
              {label}
            </p>
          ))}
        </header>
        <div class={styles.posts}>
          {tabIdx === 0 && isOwn && <Composer />}
          {tabIdx === 0 &&
            (userPosts.length === 0 ? (
              <p class={styles.state}>Aún no hay posts.</p>
            ) : (
              userPosts.map((p) => <PostCard key={p.id} post={p} />)
            ))}
          {tabIdx === 1 && isLoadingComments && <p class={styles.state}>Cargando...</p>}
          {tabIdx === 1 &&
            (!Array.isArray(userComments) || userComments.length === 0 ? (
              <p class={styles.state}>No se ha hecho ningun comentario.</p>
            ) : (
              userComments.map((c) => <CommentCard key={c.id} c={c} />)
            ))}
          {tabIdx === 2 &&
            (userLikedPosts.length === 0 ? (
              <p class={styles.state}>No ha reaccionado a ninguna publicación.</p>
            ) : (
              userLikedPosts.map((p) => <PostCard key={p.id} post={p} />)
            ))}
        </div>
      </section>

      {listModal && (
        <div class={styles.overlay} onClick={() => setListModal(null)}>
          <div class={styles.modal} onClick={(e) => e.stopPropagation()}>
            <header class={styles.modalHeader}>
              <h3 class={styles.modalTitle}>{listModal.type === 'followers' ? 'Seguidores' : 'Seguidos'}</h3>
              <button class={styles.modalClose} onClick={() => setListModal(null)}>
                <Icon path={icons.close} />
              </button>
            </header>
            <div class={styles.modalList}>
              {listModal.users.length === 0 && (
                <p class={styles.state}>
                  {listModal.type === 'followers' ? 'Aún no tiene seguidores.' : 'No sigue a nadie aún.'}
                </p>
              )}
              {listModal.users.map((u) => (
                <Link key={u.id} to={`/u/${u.username}`} class={styles.userRow} onClick={() => setListModal(null)}>
                  <span class={styles.userAvatar}>{getInitials(u.name)}</span>
                  <div class={styles.userData}>
                    <span class={styles.userName}>{u.name}</span>
                    <span class={styles.userHandle}>@{u.username}</span>
                  </div>
                </Link>
              ))}
            </div>
          </div>
        </div>
      )}
    </>
  )
}
