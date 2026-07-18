export interface TopicDef {
  id: number
  name: string
  label: string
}

export const TOPICS: TopicDef[] = [
  { id: 0, name: 'General', label: 'General' },
  { id: 1, name: 'Programming', label: 'Programación' },
  { id: 2, name: 'Cybersecurity', label: 'Ciberseguridad' },
  { id: 3, name: 'Entertainment', label: 'Entretenimiento' },
  { id: 4, name: 'Funny', label: 'Humor' },
  { id: 5, name: 'Art', label: 'Arte' },
  { id: 6, name: 'Sports', label: 'Deportes' },
  { id: 7, name: 'Politics', label: 'Política' },
  { id: 8, name: 'Science', label: 'Ciencia' },
  { id: 9, name: 'News', label: 'Noticias' },
  { id: 10, name: 'Cinema', label: 'Cine' },
  { id: 11, name: 'Games', label: 'Videojuegos' },
  { id: 12, name: 'Literature', label: 'Literatura' },
  { id: 13, name: 'Travel', label: 'Viajes' },
  { id: 14, name: 'Cuisive', label: 'Gastronomía' },
  { id: 15, name: 'Tech', label: 'Tecnología' },
  { id: 16, name: 'Economy', label: 'Economía' },
  { id: 17, name: 'Health', label: 'Salud' },
  { id: 18, name: 'Philosophy', label: 'Filosofía' },
  { id: 19, name: 'Opinion', label: 'Opinión' },
  { id: 20, name: 'Ad', label: 'Anuncio' },
] as const

export function getTopicById(id: number): TopicDef | undefined {
  return TOPICS.find(t => t.id === id)
}

export function getTopicsLabels(ids: number[]): string[] {
  return ids.map(id => getTopicById(id)?.label ?? `#${id}`)
}
