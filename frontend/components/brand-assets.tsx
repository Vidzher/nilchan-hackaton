import { Check } from 'lucide-react'
import { cn } from '@/lib/utils'

export function BrandMark({ className }: { className?: string }) {
  return (
    <span
      className={cn('relative block size-8 shrink-0', className)}
      aria-hidden="true"
    >
      <span className="absolute inset-x-1.5 bottom-0 h-5 rounded-md bg-[#9f4528]" />
      <span className="absolute inset-x-0.5 bottom-1 h-5 rotate-[-4deg] rounded-md border border-white/30 bg-[#f4a37e]" />
      <span className="absolute inset-x-1 bottom-2 grid h-5 rotate-[3deg] place-items-center rounded-md bg-[color:var(--brand)] text-white shadow-sm">
        <Check className="size-3.5 stroke-[3]" />
      </span>
    </span>
  )
}

export function BacklogIllustration({ className }: { className?: string }) {
  return (
    <div className={cn('relative h-24 w-32', className)} aria-hidden="true">
      <div className="absolute inset-x-5 bottom-1 h-16 rotate-[-7deg] rounded-xl border border-border bg-[#ede7dc]" />
      <div className="absolute inset-x-3 bottom-3 h-16 rotate-[5deg] rounded-xl border border-border bg-brand-soft" />
      <div className="absolute inset-x-4 bottom-6 flex h-16 flex-col justify-between rounded-xl border border-[color:var(--brand)]/25 bg-card p-3 shadow-[0_8px_24px_rgba(86,54,37,0.10)]">
        <div className="h-2 w-3/4 rounded-full bg-foreground/15" />
        <div className="space-y-1.5">
          <div className="h-1.5 w-full rounded-full bg-foreground/10" />
          <div className="h-1.5 w-2/3 rounded-full bg-foreground/10" />
        </div>
        <span className="absolute -right-2 -top-2 grid size-7 place-items-center rounded-full bg-[color:var(--brand)] text-white shadow-md">
          <Check className="size-4 stroke-[3]" />
        </span>
      </div>
    </div>
  )
}
