'use client'

import Image from 'next/image'
import { useEffect, useState } from 'react'
import { cn } from '@/lib/utils'

const frameStyles: Record<string, { shell: string; avatar: string }> = {
  'frame-default': {
    shell: 'bg-border p-[2px]',
    avatar: 'ring-1 ring-inset ring-white/70',
  },
  'frame-neon': {
    shell: 'bg-[conic-gradient(from_210deg,#22d3ee,#6366f1,#d946ef,#22d3ee)] p-[3px] shadow-[0_0_14px_rgba(99,102,241,0.38)]',
    avatar: 'ring-1 ring-inset ring-white/50',
  },
  'frame-fire': {
    shell: 'bg-[conic-gradient(from_200deg,#7f1d1d,#ef4444,#f59e0b,#fde68a,#ef4444,#7f1d1d)] p-[3px] shadow-[0_0_14px_rgba(234,88,12,0.38)]',
    avatar: 'ring-1 ring-inset ring-[#fff1d6]/70',
  },
  'frame-gold': {
    shell: 'bg-[conic-gradient(from_215deg,#6b4f12,#f6d365,#fff4b0,#b7791f,#f6d365,#6b4f12)] p-[4px] shadow-[0_0_0_1px_#8a6518,0_3px_12px_rgba(120,83,15,0.28)]',
    avatar: 'ring-2 ring-inset ring-[#fff1a8]/80',
  },
}

const defaultFrame = frameStyles['frame-default']

const sizes = {
  sm: { className: 'size-9 text-lg', pixels: 36 },
  md: { className: 'size-12 text-2xl', pixels: 48 },
  lg: { className: 'size-20 text-4xl', pixels: 80 },
}

export function avatarImagePath(assetKey?: string) {
  return assetKey ? `/avatars/${assetKey}.png` : undefined
}

export function AvatarFrame({
  src,
  fallback = '🙂',
  frame,
  size = 'md',
  className,
}: {
  src?: string
  fallback?: string
  frame?: string
  size?: keyof typeof sizes
  className?: string
}) {
  const [imageFailed, setImageFailed] = useState(false)
  const dimensions = sizes[size]
  const style = (frame && frameStyles[frame]) || defaultFrame

  useEffect(() => setImageFailed(false), [src])

  return (
    <span
      className={cn(
        'relative grid shrink-0 place-items-center rounded-full',
        dimensions.className,
        style.shell,
        className,
      )}
      aria-hidden="true"
    >
      <span
        className={cn(
          'grid size-full place-items-center overflow-hidden rounded-full bg-brand-soft',
          style.avatar,
        )}
      >
        {src && !imageFailed ? (
          <Image
            src={src}
            alt=""
            width={dimensions.pixels}
            height={dimensions.pixels}
            className="size-full object-cover"
            onError={() => setImageFailed(true)}
          />
        ) : (
          fallback
        )}
      </span>
    </span>
  )
}
