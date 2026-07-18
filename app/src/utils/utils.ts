export function getAgeAndDaysSince(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()

  const msDiff = Number(now) - Number(date)
  const daysDiff = Math.floor(msDiff / (1000 * 60 * 60 * 24))
  const hoursDiff = Math.floor((msDiff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minDiff = Math.floor((msDiff % (1000 * 60 * 60 * 24)) / (1000 * 60))
  const secDiff = Math.floor((msDiff % (1000 * 60 * 60 * 24)) / 1000)

  const isYearPassed = daysDiff >= 365

  if (daysDiff < 1 && hoursDiff < 1 && minDiff < 1) {
    return `${secDiff}s`
  } else if (daysDiff < 1 && hoursDiff < 1) {
    return `${minDiff}m`
  } else if (hoursDiff >= 1 && daysDiff < 1) {
    return `${hoursDiff}h`
  } else if (daysDiff <= 7) {
    return `${daysDiff}d`
  } else if (isYearPassed) {
    return `${date.getDate()} ${date.toLocaleString('default', {
      month: 'short'
    })} ${date.getFullYear()}`
  } else {
    return `${date.getDate()} ${date.toLocaleString('default', { month: 'short' })}`
  }
}

export function getInitials(name: string) {
  if (!name) return '?'
  const parts = name.trim().split(/\s+/)
  const first = parts[0]?.[0] ?? ''
  const last = parts.length > 1 ? parts[parts.length - 1][0] : ''
  return (first + last).toUpperCase() || '?'
}
