'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import {
  Layers,
  LogOut,
  ShoppingBag,
  Trophy,
  User,
  type LucideIcon,
} from 'lucide-react'
import { AvatarFrame, avatarImagePath } from '@/components/avatar-frame'
import { BrandMark } from '@/components/brand-assets'
import { cn } from '@/lib/utils'
import { api, clearToken, type Profile, type ShopItem } from '@/lib/api'

type NavItem = { href: string; label: string; icon: LucideIcon }

const navItems: NavItem[] = [
  { href: '/', label: 'Backlog', icon: Layers },
  { href: '/profile', label: 'Профиль', icon: User },
  { href: '/shop', label: 'Магазин', icon: ShoppingBag },
  { href: '/leaderboard', label: 'Рейтинг', icon: Trophy },
]

function Wordmark() {
  return (
    <Link href="/" className="group flex items-center gap-2.5" aria-label="Learning Backlog">
      <BrandMark className="size-7 transition-transform group-hover:-rotate-3" />
      <span className="leading-none">
        <span className="block text-[10px] font-semibold uppercase tracking-[0.18em] text-[color:var(--brand)]">Learning</span>
        <span className="mt-0.5 block text-sm font-bold tracking-tight">Backlog</span>
      </span>
    </Link>
  )
}

function isActive(pathname: string, href: string) {
  return href === '/' ? pathname === '/' : pathname.startsWith(href)
}

export function AppShell({ children }: { children: React.ReactNode }) {
  const pathname = usePathname()
  const router = useRouter()
  const [profile, setProfile] = useState<Profile | null>(null)
  const [avatar, setAvatar] = useState<ShopItem | null>(null)
  const [frame, setFrame] = useState<ShopItem | null>(null)
  const [slotsUsed, setSlotsUsed] = useState(0)

  useEffect(() => {
    const loadProfile = () => {
      Promise.all([api.profile(), api.shop(), api.resources()])
        .then(([nextProfile, catalog, resources]) => {
          setProfile(nextProfile)
          setAvatar(catalog.find((item) => item.id === nextProfile.cosmetics.avatarId) ?? null)
          setFrame(catalog.find((item) => item.id === nextProfile.cosmetics.frameId) ?? null)
          setSlotsUsed(resources.filter((resource) => resource.status === 'PROCESSING' || resource.status === 'NOT_COMPLETED').length)
        })
        .catch(() => undefined)
    }

    loadProfile()
    window.addEventListener('cosmeticschange', loadProfile)
    window.addEventListener('resourceschange', loadProfile)
    return () => {
      window.removeEventListener('cosmeticschange', loadProfile)
      window.removeEventListener('resourceschange', loadProfile)
    }
  }, [pathname])

  function logout() {
    clearToken()
    router.replace('/login')
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Desktop sidebar */}
      <aside className="fixed inset-y-0 left-0 z-30 hidden w-[220px] flex-col border-r border-border bg-sidebar px-4 py-6 lg:flex">
        <div className="px-1">
          <Wordmark />
        </div>
        <nav className="mt-8 flex flex-1 flex-col gap-1" aria-label="Основная навигация">
          {navItems.map(({ href, label, icon: Icon }) => {
            const active = isActive(pathname, href)
            return (
              <Link
                key={href}
                href={href}
                aria-current={active ? 'page' : undefined}
                className={cn(
                  'flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors',
                  active
                    ? 'bg-brand-soft text-[color:var(--brand)]'
                    : 'text-muted-foreground hover:bg-secondary/60 hover:text-foreground',
                )}
              >
                <Icon className="size-4" aria-hidden="true" />
                {label}
              </Link>
            )
          })}
        </nav>

        <div className="mb-3 rounded-lg border border-border bg-card/60 p-3">
          <div className="flex items-center justify-between gap-3">
            <span className="text-xs font-medium text-muted-foreground">Активные слоты</span>
            <span className="tabular text-xs font-semibold">{slotsUsed} / {profile?.progress.activeBacklogLimit ?? '…'}</span>
          </div>
          <div className="mt-2 flex gap-1.5" aria-label={`Занято слотов: ${slotsUsed} из ${profile?.progress.activeBacklogLimit ?? 0}`}>
            {Array.from({ length: Math.max(slotsUsed, profile?.progress.activeBacklogLimit ?? 0) }, (_, index) => (
              <span
                key={index}
                className={cn(
                  'size-2.5 rounded-full border',
                  index < slotsUsed
                    ? index < (profile?.progress.activeBacklogLimit ?? 0) ? 'border-[color:var(--brand)] bg-[color:var(--brand)]' : 'border-[color:var(--warning)] bg-[color:var(--warning)]'
                    : 'border-border bg-background',
                )}
                aria-hidden="true"
              />
            ))}
          </div>
        </div>

        <div className="flex items-center gap-3 rounded-lg border border-border p-2.5">
          <AvatarFrame
            src={avatarImagePath(avatar?.presentation.assetKey)}
            fallback={avatar?.presentation.emoji}
            frame={frame?.presentation.cssClass}
            size="sm"
          />
          <div className="min-w-0 flex-1">
            <p className="truncate text-sm font-medium">{profile?.user.username ?? '…'}</p>
            <p className="text-xs text-muted-foreground">Уровень {profile?.progress.level ?? '…'}</p>
          </div>
          <button
            type="button"
            onClick={logout}
            aria-label="Выйти"
            className="grid size-8 shrink-0 place-items-center rounded-md text-muted-foreground transition-colors hover:bg-secondary hover:text-foreground"
          >
            <LogOut className="size-4" aria-hidden="true" />
          </button>
        </div>
      </aside>

      {/* Mobile top bar */}
      <header className="sticky top-0 z-30 flex items-center justify-between border-b border-border bg-background/90 px-4 py-3 backdrop-blur lg:hidden">
        <Wordmark />
        <span className="tabular inline-flex items-center gap-1 rounded-md border border-border bg-card px-2 py-1 text-sm font-medium">
          {profile?.progress.ePoints ?? '…'}
          <span className="text-xs text-muted-foreground">е-баллов</span>
        </span>
      </header>

      {/* Main content */}
      <div className="lg:pl-[220px]">
        <main className="mx-auto w-full max-w-[1180px] px-4 pt-6 pb-28 lg:px-8 lg:pt-8 lg:pb-12">
          {children}
        </main>
      </div>

      {/* Mobile bottom nav */}
      <nav
        aria-label="Основная навигация"
        className="fixed inset-x-0 bottom-0 z-30 border-t border-border bg-card lg:hidden"
        style={{ paddingBottom: 'env(safe-area-inset-bottom)' }}
      >
        <div className="mx-auto grid max-w-md grid-cols-4">
          {navItems.map(({ href, label, icon: Icon }) => {
            const active = isActive(pathname, href)
            return (
              <Link
                key={href}
                href={href}
                aria-current={active ? 'page' : undefined}
                className={cn(
                  'flex min-h-[56px] flex-col items-center justify-center gap-1 text-xs font-medium transition-colors',
                  active ? 'text-[color:var(--brand)]' : 'text-muted-foreground',
                )}
              >
                <Icon className="size-5" aria-hidden="true" />
                {label}
              </Link>
            )
          })}
        </div>
      </nav>
    </div>
  )
}
