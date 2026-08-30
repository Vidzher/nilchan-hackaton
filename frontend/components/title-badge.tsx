import { cn } from '@/lib/utils'

function rarityStyle(price: number) {
  if (price <= 50) return 'border-border bg-muted text-muted-foreground'
  if (price <= 100) return 'border-[#91b5d8] bg-[#eaf3fb] text-[#34688f]'
  if (price <= 150) return 'border-[#b9a2d5] bg-[#f1ebf8] text-[#704d98]'
  return 'border-[#d6b768] bg-[#f8efd3] text-[#806117] shadow-[0_0_12px_rgba(183,138,43,0.14)]'
}

export function TitleBadge({
  name,
  price,
  className,
}: {
  name: string
  price: number
  className?: string
}) {
  return (
    <span
      className={cn(
        'inline-flex max-w-full items-center rounded-md border px-2 py-0.5 font-mono text-[11px] font-medium leading-4',
        rarityStyle(price),
        className,
      )}
    >
      <span className="truncate">〔 {name} 〕</span>
    </span>
  )
}
