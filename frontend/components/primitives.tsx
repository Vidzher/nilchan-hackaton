import { cn } from '@/lib/utils'

export function Tag({ children }: { children: React.ReactNode }) {
  return (
    <span className="inline-flex items-center rounded-md border border-border bg-muted px-1.5 py-0.5 text-xs text-muted-foreground">
      {children}
    </span>
  )
}

export function ProgressBar({
  value,
  max,
  className,
  label,
}: {
  value: number
  max: number
  className?: string
  label?: string
}) {
  const pct = Math.min(100, Math.round((value / max) * 100))
  return (
    <div
      className={cn('h-1.5 w-full overflow-hidden rounded-full bg-border', className)}
      role="progressbar"
      aria-valuenow={value}
      aria-valuemin={0}
      aria-valuemax={max}
      aria-label={label}
    >
      <div
        className="h-full rounded-full bg-[color:var(--brand)] transition-all"
        style={{ width: `${pct}%` }}
      />
    </div>
  )
}
