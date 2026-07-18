export function Icon({ path, fill }: { path: string; fill?: string }) {
  return (
    <svg
      style={{ height: '18px', width: '18px' }}
      viewBox="0 0 24 24"
      fill={fill ? fill : 'none'}
      stroke={fill ? fill : 'currentColor'}
      stroke-width="2"
      stroke-linecap="round"
      stroke-linejoin="round"
    >
      <path d={path} />
    </svg>
  )
}
