import { cn } from '@/lib/utils'

export function PageHeader({
  title,
  description,
  action,
  meta,
  divided = true,
}: {
  title: string
  description?: string
  action?: React.ReactNode
  meta?: React.ReactNode
  divided?: boolean
}) {
  return (
    <header className={cn('flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between', divided && 'border-b border-border pb-6')}>
      <div className="min-w-0 max-w-2xl">
        {meta ? <div className="mb-2 text-xs font-medium uppercase tracking-[0.14em] text-muted-foreground">{meta}</div> : null}
        <h1 className="text-2xl font-semibold tracking-[-0.025em] text-balance sm:text-3xl">{title}</h1>
        {description ? <p className="mt-1.5 text-sm leading-6 text-muted-foreground">{description}</p> : null}
      </div>
      {action ? <div className="shrink-0">{action}</div> : null}
    </header>
  )
}

export function SectionHeader({
  title,
  description,
  aside,
}: {
  title: string
  description?: string
  aside?: React.ReactNode
}) {
  return (
    <div className="flex items-end justify-between gap-4">
      <div>
        <h2 className="text-base font-semibold tracking-tight">{title}</h2>
        {description ? <p className="mt-1 text-sm text-muted-foreground">{description}</p> : null}
      </div>
      {aside}
    </div>
  )
}

export function StatCard({
  icon: Icon,
  label,
  value,
  suffix,
  accent,
  className,
}: {
  icon: React.ComponentType<{ className?: string; 'aria-hidden'?: boolean }>
  label: string
  value: React.ReactNode
  suffix?: string
  accent?: boolean
  className?: string
}) {
  return (
    <div className={cn('rounded-xl bg-card/70 p-4', className)}>
      <div className="flex items-center gap-2 text-xs font-medium text-muted-foreground">
        <Icon className={cn('size-4', accent && 'text-[color:var(--brand)]')} aria-hidden={true} />
        {label}
      </div>
      <p className="tabular mt-2 text-2xl font-semibold tracking-tight">
        {value}
        {suffix ? <span className="ml-1.5 text-sm font-normal text-muted-foreground">{suffix}</span> : null}
      </p>
    </div>
  )
}

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
  const safeMax = Math.max(1, max)
  const safeValue = Math.max(0, Math.min(value, safeMax))
  const pct = Math.round((safeValue / safeMax) * 100)
  return (
    <div
      className={cn('h-1.5 w-full overflow-hidden rounded-full bg-border', className)}
      role="progressbar"
      aria-valuenow={safeValue}
      aria-valuemin={0}
      aria-valuemax={safeMax}
      aria-label={label}
    >
      <div
        className="h-full rounded-full bg-[color:var(--brand)] transition-all"
        style={{ width: `${pct}%` }}
      />
    </div>
  )
}
