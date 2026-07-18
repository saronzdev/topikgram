import { Composer } from './Composer'
import { PostCard } from './PostCard'
import { usePosts, POSTS_KEY } from '../services/posts'
import { useAuth } from '../services/auth'
import { mutate } from 'swr'
import { useCallback, useRef } from 'preact/hooks'
import styles from '../styles/Home.module.css'

export function Home() {
  const { posts, error, isLoading, hasMore, nextCursor } = usePosts()
  const { isAuth } = useAuth()
  const loadingMore = useRef(false)

  const loadMore = useCallback(async () => {
    if (!hasMore || !nextCursor || loadingMore.current) return
    loadingMore.current = true
    try {
      const res = await fetch(`${POSTS_KEY}?cursor=${encodeURIComponent(nextCursor)}`, {
        credentials: 'include'
      }).then(r => r.json()) as PostListResponse
      mutate(POSTS_KEY, (prev) => {
        if (!prev) return prev
        return {
          posts: [...prev.posts, ...res.posts],
          next_cursor: res.next_cursor,
          has_more: res.has_more
        }
      }, false)
    } catch {}
    loadingMore.current = false
  }, [hasMore, nextCursor])

  return (
    <main class={styles.home}>
      <header class={styles.header}>
        <h1 class={styles.h1}>Topikgram</h1>
      </header>

      {isAuth && <Composer />}

      <section class={styles.posts}>
        {isLoading && <p class={styles.state}>Cargando…</p>}
        {error && <p class={styles.state}>No se pudieron cargar los posts.</p>}
        {!isLoading && !error && (!posts || posts.length === 0) && <p class={styles.state}>Aún no hay posts.</p>}
        {posts?.map((post: Post) => <PostCard key={post.id} post={post} />)}
        {hasMore && (
          <button class={styles.loadMore} onClick={loadMore}>
            Cargar más
          </button>
        )}
      </section>
    </main>
  )
}
