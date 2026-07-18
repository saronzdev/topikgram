import useSWR, { mutate } from 'swr'
import { API_BASE, apiFetch, fetcher } from './fetch'

export const POSTS_KEY = `${API_BASE}/posts`

export function usePosts(cursor?: string) {
  const url = cursor ? `${POSTS_KEY}?cursor=${encodeURIComponent(cursor)}` : POSTS_KEY
  const { data, error, isLoading } = useSWR<PostListResponse>(url, fetcher, {
    revalidateOnFocus: true,
    dedupingInterval: 5000
  })
  return { data, posts: data?.posts, isLoading, error, hasMore: data?.has_more ?? false, nextCursor: data?.next_cursor }
}

export async function createPost(body: string, topics: number[]) {
  await apiFetch(POSTS_KEY, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ body, topics })
  })
  await mutate(POSTS_KEY)
}

export async function updatePost(id: number, body: string) {
  const res = await apiFetch(`${POSTS_KEY}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ body })
  })
  await mutate(POSTS_KEY)
  return res
}

export async function removePost(id: number) {
  await apiFetch(`${POSTS_KEY}/${id}`, { method: 'DELETE' })
  await mutate(POSTS_KEY)
}

export async function likePost(id: number) {
  await apiFetch(`${POSTS_KEY}/${id}/like`, { method: 'POST' })
}

export async function unlikePost(id: number) {
  await apiFetch(`${POSTS_KEY}/${id}/like`, { method: 'DELETE' })
}

export async function savePost(id: number) {
  await apiFetch(`${POSTS_KEY}/${id}/save`, { method: 'POST' })
}

export async function unsavePost(id: number) {
  await apiFetch(`${POSTS_KEY}/${id}/save`, { method: 'DELETE' })
}

export async function getPostLikes(id: number, page = 1, limit = 20): Promise<PaginatedUsers> {
  return await apiFetch(`${POSTS_KEY}/${id}/likes?page=${page}&limit=${limit}`)
}

export async function getPostSaves(id: number, page = 1, limit = 20): Promise<PaginatedUsers> {
  return await apiFetch(`${POSTS_KEY}/${id}/saves?page=${page}&limit=${limit}`)
}

export async function createComment(postID: number, content: string) {
  await apiFetch(`${API_BASE}/comments`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ post_id: postID, content })
  })
}

export async function getComments(postID: number): Promise<CommentInterface[]> {
  return await apiFetch(`${API_BASE}/comments/${postID}`)
}

export async function getCommentsByUserID(userID: number): Promise<CommentInterface[]> {
  return await apiFetch(`${API_BASE}/comments/user/${userID}`)
}
