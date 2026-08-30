import { cn } from '@/lib/utils'

const frameRing: Record<string, string> = {
  Neon: 'ring-[color:var(--brand)] shadow-[0_0_0_3px_var(--brand-soft)]',
  Fire: 'ring-[#c9541f] shadow-[0_0_0_3px_#f4e0d3]',
  Gold: 'ring-[#b78a2b] shadow-[0_0_0_3px_#f2ead4]',
}

const sizes = {
  sm: 'size-9 text-lg',
  md: 'size-12 text-2xl',
  lg: 'size-20 text-4xl',
}

export function AvatarFrame({
  emoji,
  frame,
  size = 'md',
  className,
}: {
  emoji: string
  frame?: string
  size?: keyof typeof sizes
  className?: string
}) {
  return (
    <span
      className={cn(
        'grid shrink-0 place-items-center rounded-full bg-brand-soft ring-2',
        sizes[size],
        frame ? frameRing[frame] : 'ring-border',
        className,
      )}
      aria-hidden="true"
    >
      {emoji}
    </span>
  )
}
