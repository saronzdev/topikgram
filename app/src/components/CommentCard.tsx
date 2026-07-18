import { Link } from 'wouter'
import { getInitials, getAgeAndDaysSince } from '../utils/utils'
import styles from '../styles/PostView.module.css'

export function CommentCard({ c }: { c: CommentInterface }) {
  return (
    <div key={c.id} class={styles.comment}>
      <Link to={`/u/${c.user.username}`} class={styles.commentAvatar}>
        {getInitials(c.user.name)}
      </Link>
      <div class={styles.commentContent}>
        <div class={styles.commentHead}>
          <span class={styles.commentName}>{c.user.name}</span>
          <span class={styles.commentHandle}>@{c.user.username}</span>
          <span class={styles.dot}>·</span>
          <span class={styles.commentTime}>{getAgeAndDaysSince(c.created_at)}</span>
        </div>
        <p class={styles.commentBody}>{c.content}</p>
      </div>
    </div>
  )
}
