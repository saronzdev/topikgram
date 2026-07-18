interface User {
  id: number
  name: string
  username: string
  email: string
  birthday: string
  created_at: string
}

interface UserPublic {
  id: number
  name: string
  username: string
}

interface Post {
  id: number
  user_id: number
  body: string
  topics_id: number[]
  user: UserPublic
  likes: number
  liked: boolean
  saves: number
  saved: boolean
  comments: number
  created_at: string
}

interface PostListResponse {
  posts: Post[]
  next_cursor?: string
  has_more: boolean
}

interface CommentInterface {
  id: number
  user_id: number
  post_id: number
  content: string
  user: UserPublic
  created_at: string
}

interface PaginatedUsers {
  users: UserPublic[]
  total: number
  page: number
  limit: number
}

interface TopicID {
  id: number
  name: string
}
