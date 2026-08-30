import { CircleCheck, CircleX, Loader, Sparkles } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { ResourceStatus } from '@/lib/api'

const statusMeta: Record<ResourceStatus, { label: string; tone: 'brand' | 'warning' | 'error' | 'success' }> = {
  NOT_COMPLETED: { label: 'Quiz готов', tone: 'brand' },
  PROCESSING: { label: 'Создаём quiz', tone: 'warning' },
  FAILED: { label: 'Не удалось создать quiz', tone: 'error' },
  COMPLETED: { label: 'Завершено', tone: 'success' },
}

const toneStyles: Record<string, string> = {
  brand: 'bg-brand-soft text-[color:var(--brand)] border-[color:var(--brand)]/25',
  warning: 'bg-warning-soft text-[color:var(--warning)] border-[color:var(--warning)]/25',
  error: 'bg-[#f6e3e3] text-[color:var(--destructive)] border-[color:var(--destructive)]/25',
  success: 'bg-success-soft text-[color:var(--success)] border-[color:var(--success)]/25',
}

const toneIcon = {
  brand: Sparkles,
  warning: Loader,
  error: CircleX,
  success: CircleCheck,
}

export function StatusBadge({
  status,
  className,
}: {
  status: ResourceStatus
  className?: string
}) {
  const meta = statusMeta[status]
  const Icon = toneIcon[meta.tone]
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1.5 rounded-md border px-2 py-1 text-xs font-medium',
        toneStyles[meta.tone],
        className,
      )}
    >
      <Icon
        className={cn('size-3.5', meta.tone === 'warning' && 'animate-spin')}
        aria-hidden="true"
      />
      {meta.label}
    </span>
  )
}
