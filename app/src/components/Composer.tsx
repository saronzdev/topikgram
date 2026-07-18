import { useState } from 'preact/hooks'
import { usePosts, createPost, POSTS_KEY } from '../services/posts'
import { TOPICS } from '../utils/topics'
import { useAuth } from '../services/auth'
import { mutate } from 'swr'
import styles from '../styles/Composer.module.css'

export function Composer() {
  const { user, isAuth } = useAuth()
  const { posts } = usePosts()
  const [text, setText] = useState('')
  const [selectedTopics, setSelectedTopics] = useState<number[]>([0])
  const [submitting, setSubmitting] = useState(false)

  const max = 5000
  const len = text.length
  const canPost = len > 0 && len <= max && selectedTopics.length >= 1 && selectedTopics.length <= 3 && isAuth && !submitting

  function toggleTopic(id: number) {
    setSelectedTopics(prev => {
      if (id === 0) {
        if (prev.includes(0)) {
          return prev.length > 1 ? prev.filter(t => t !== 0) : prev
        }
        if (prev.length >= 3) return prev
        return [...prev, 0]
      }
      if (prev.includes(id)) return prev.filter(t => t !== id)
      if (prev.length >= 3) return prev
      return [...prev, id]
    })
  }

  async function submit() {
    if (!isAuth || !user || submitting) return
    setSubmitting(true)
    const previous = posts
    const newPost: Post = {
      id: Date.now(), user_id: user.id, user,
      body: text, topics_id: selectedTopics,
      likes: 0, liked: false, saved: false,
      created_at: new Date().toISOString()
    }
    setText('')
    setSelectedTopics([0])
    mutate(POSTS_KEY, (prev) => prev ? { ...prev, posts: [...prev.posts, newPost] } : prev, false)
    try {
      await createPost(text, selectedTopics)
    } catch {
      mutate(POSTS_KEY, prev => prev ? { ...prev, posts: previous || [] } : prev, false)
    }
    setSubmitting(false)
  }

  return (
    <section class={styles.composer}>
      <div class={styles.row}>
        <div class={styles.avatar}><span>Tú</span></div>
        <textarea class={styles.input} placeholder="¿Qué está pasando?" maxLength={max}
          value={text} onInput={e => setText(e.currentTarget.value)} rows={3} />
      </div>
      <div class={styles.topics}>
        <span class={styles.topicLabel}>Temas ({selectedTopics.length}/3)</span>
        <div class={styles.topicGrid}>
          {TOPICS.filter(t => t.id !== 20).map(t => (
            <button
              key={t.id}
              type="button"
              class={`${styles.topicChip} ${selectedTopics.includes(t.id) ? styles.topicActive : ''}`}
              onClick={() => toggleTopic(t.id)}
              disabled={t.id !== 0 && !selectedTopics.includes(t.id) && selectedTopics.length >= 3}
            >
              {t.label}
            </button>
          ))}
        </div>
      </div>
      <div class={styles.actions}>
        <div class={styles.meta}>
          <span class={`${styles.counter} ${len > max * 0.9 ? styles.counterWarn : ''}`}>{len}/{max}</span>
          <button class={styles.btnPrimary} disabled={!canPost} onClick={submit}>
            {submitting ? 'Publicando…' : 'Postear'}
          </button>
        </div>
      </div>
    </section>
  )
}
