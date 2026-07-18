import { useEffect, useState } from 'preact/hooks'
import { useParams, Link } from 'wouter'
import { useAuth } from '../services/auth'
import { apiFetch } from '../services/fetch'
import { POSTS_KEY, getComments, createComment, likePost, unlikePost, savePost, unsavePost } from '../services/posts'
import { Icon } from './Icon'
import { icons } from '../utils/svgs'
import { getInitials, getAgeAndDaysSince } from '../utils/utils'
import { getTopicsLabels } from '../utils/topics'
import styles from '../styles/PostView.module.css'
import { CommentCard } from './CommentCard'

export function PostView() {
  const { id } = useParams() as { id: string }
  const { user: authUser, isAuth } = useAuth()
  const [post, setPost] = useState<Post | null>(null)
  const [comments, setComments] = useState<CommentInterface[]>([])
  const [commentText, setCommentText] = useState('')
  const [sending, setSending] = useState(false)
  const [liked, setLiked] = useState(false)
  const [likesCount, setLikesCount] = useState(0)
  const [saved, setSaved] = useState(false)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(false)

  useEffect(() => {
    setLoading(true)
    setError(false)
    apiFetch(`${POSTS_KEY}/${id}`)
      .then((p: Post) => {
        setPost(p)
        setLiked(p.liked)
        setLikesCount(p.likes)
        setSaved(p.saved)
        return getComments(p.id)
      })
      .then((c) => setComments(c))
      .catch(() => setError(true))
      .finally(() => setLoading(false))
  }, [id])

  async function toggleLike() {
    if (!isAuth || !post) return
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
    if (!isAuth || !post) return
    const prevSaved = saved
    setSaved((v) => !v)
    try {
      prevSaved ? await unsavePost(post.id) : await savePost(post.id)
    } catch {
      setSaved(prevSaved)
    }
  }

  async function submitComment() {
    if (!isAuth || !post || !commentText.trim() || sending) return
    setSending(true)
    const text = commentText.trim()
    setCommentText('')
    const optimistic: CommentInterface = {
      id: Date.now(),
      user_id: authUser!.id,
      post_id: post.id,
      content: text,
      user: { id: authUser!.id, name: authUser!.name, username: authUser!.username },
      created_at: new Date().toISOString()
    }
    setComments((prev) => [...prev, optimistic])
    try {
      await createComment(post.id, text)
      const fresh = await getComments(post.id)
      setComments(fresh)
    } catch {
      setComments((prev) => prev.filter((c) => c.id !== optimistic.id))
    }
    setSending(false)
  }

  if (loading) return <p class={styles.state}>Cargando…</p>
  if (error || !post) return <p class={styles.state}>No se pudo cargar el post.</p>

  const topicLabels = getTopicsLabels(post.topics_id)

  return (
    <main class={styles.view}>
      <header class={styles.header}>
        <Link to="/" class={styles.back}>
          <Icon path={icons.arrowLeft} />
        </Link>
        <h2 class={styles.title}>Post</h2>
      </header>

      <article class={styles.post}>
        <div class={styles.postTop}>
          <Link to={`/u/${post.user.username}`} class={styles.avatar}>
            {getInitials(post.user.name)}
          </Link>
          <div class={styles.postHead}>
            <span class={styles.name}>{post.user.name}</span>
            <span class={styles.handle}>@{post.user.username}</span>
          </div>
          <span class={styles.time}>{getAgeAndDaysSince(post.created_at)}</span>
        </div>
        <p class={styles.body}>{post.body}</p>
        {topicLabels.length > 0 && (
          <div class={styles.topics}>
            {topicLabels.map((label) => (
              <span key={label} class={styles.topic}>
                #{label}
              </span>
            ))}
          </div>
        )}
        <div class={styles.postActions}>
          <button class={`${styles.action} ${liked ? styles.liked : ''}`} onClick={toggleLike} disabled={!isAuth}>
            <Icon path={icons.like} fill={liked ? '#ff0000' : ''} />
            <span>{likesCount}</span>
          </button>
          <span class={styles.actionStatic}>
            <Icon path={icons.reply} />
            <span>{comments.length}</span>
          </span>
          <button class={`${styles.action} ${saved ? styles.saved : ''}`} onClick={toggleSaved} disabled={!isAuth}>
            <Icon path={icons.bookmark} fill={saved ? '#fbbf24' : ''} />
            <span>{saved ? 'Guardado' : 'Guardar'}</span>
          </button>
        </div>
      </article>

      {isAuth && (
        <div class={styles.inputBar}>
          <div class={styles.inputAvatar}>{authUser ? getInitials(authUser.name) : '?'}</div>
          <input
            class={styles.input}
            placeholder="Escribe un comentario…"
            value={commentText}
            onInput={(e) => setCommentText(e.currentTarget.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault()
                submitComment()
              }
            }}
          />
          <button class={styles.sendBtn} disabled={!commentText.trim() || sending} onClick={submitComment}>
            {sending ? '…' : <Icon path={icons.reply} />}
          </button>
        </div>
      )}

      <section class={styles.commentsSection}>
        <h3 class={styles.commentsTitle}>Comentarios ({comments.length})</h3>
        {comments.length === 0 && <p class={styles.state}>Sé el primero en comentar.</p>}
        {comments.map((c) => (
          <CommentCard key={c.id} c={c} />
        ))}
      </section>
    </main>
  )
}
