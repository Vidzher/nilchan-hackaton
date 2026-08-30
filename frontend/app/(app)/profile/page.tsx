'use client'

import Link from 'next/link'
import { useEffect, useState } from 'react'
import { CircleCheck, Coins, Flame, Layers } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { AvatarFrame, avatarImagePath } from '@/components/avatar-frame'
import { ProgressBar } from '@/components/primitives'
import { TitleBadge } from '@/components/title-badge'
import { api, ApiError, type Profile, type ShopItem } from '@/lib/api'

const ceilings: Record<number, number> = { 1: 120, 2: 300, 3: 600, 4: 1000 }

export default function ProfilePage() {
  const [profile, setProfile] = useState<Profile | null>(null)
  const [avatar, setAvatar] = useState<ShopItem | null>(null)
  const [frame, setFrame] = useState<ShopItem | null>(null)
  const [title, setTitle] = useState<ShopItem | null>(null)
  const [showcase, setShowcase] = useState<ShopItem | null>(null)
  const [completed, setCompleted] = useState(0)
  const [used, setUsed] = useState(0)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    Promise.all([api.profile(), api.resources(), api.shop()]).then(([found, resources, catalog]) => {
      setProfile(found)
      setAvatar(catalog.find((item) => item.id === found.cosmetics.avatarId) ?? null)
      setFrame(catalog.find((item) => item.id === found.cosmetics.frameId) ?? null)
      setTitle(catalog.find((item) => item.id === found.cosmetics.titleId) ?? null)
      setShowcase(catalog.find((item) => item.id === found.cosmetics.showcaseItemId) ?? null)
      setCompleted(resources.filter((resource) => resource.status === 'COMPLETED').length)
      setUsed(resources.filter((resource) => resource.status === 'PROCESSING' || resource.status === 'NOT_COMPLETED').length)
    }).catch((caught) => setError(caught instanceof ApiError ? caught.message : 'Не удалось загрузить профиль.'))
  }, [])

  if (!profile) return <div className="rounded-2xl border border-border bg-card p-6"><h1 className="text-xl font-semibold">{error ?? 'Загружаем профиль…'}</h1></div>
  const { progress, user } = profile
  const ceiling = ceilings[progress.level] ?? Math.max(progress.xp, 1000)
  const equipped = [
    { label: 'Аватар', value: avatar?.displayName ?? 'Не выбрано' },
    { label: 'Рамка', value: frame?.displayName ?? 'Без рамки' }, { label: 'Титул', value: title?.displayName ?? 'Не выбрано' },
    { label: 'Витрина', value: showcase?.displayName ?? 'Не выбрано' },
  ]

  return <div className="flex flex-col gap-6"><h1 className="text-2xl font-semibold tracking-tight">Профиль</h1>
    <div className="rounded-2xl border border-border bg-card p-6"><div className="flex flex-col gap-5 sm:flex-row sm:items-center"><AvatarFrame src={avatarImagePath(avatar?.presentation.assetKey)} fallback={avatar?.presentation.emoji} frame={frame?.presentation.cssClass} size="lg" /><div className="min-w-0 flex-1"><div className="flex flex-wrap items-center gap-2"><h2 className="text-xl font-semibold">{user.username}</h2><span className="rounded-md bg-primary px-2 py-0.5 text-xs font-medium text-primary-foreground">Уровень {progress.level}</span></div>{title ? <TitleBadge className="mt-1.5" name={title.displayName} price={title.price} /> : null}{showcase ? <p className="mt-0.5 text-sm text-muted-foreground">{showcase.presentation.emoji} {showcase.displayName}</p> : null}</div><Button variant="outline" nativeButton={false} render={<Link href="/shop" />}>Кастомизация</Button></div><div className="mt-6"><div className="mb-2 flex items-center justify-between text-sm"><span className="text-muted-foreground">Прогресс уровня</span><span className="tabular font-medium">{progress.xp} / {ceiling} XP</span></div><ProgressBar value={progress.xp} max={ceiling} label="Прогресс уровня" /></div></div>
    <div className="grid grid-cols-2 gap-3 lg:grid-cols-4"><StatCard icon={CircleCheck} label="Завершено" value={completed} suffix="материалов" /><StatCard icon={Coins} label="Баланс" value={progress.ePoints} suffix="е-баллов" /><StatCard icon={Flame} label="Серия" value={progress.currentStreak} suffix="дней" accent /><StatCard icon={Layers} label="Backlog" value={`${used}/${progress.activeBacklogLimit}`} suffix="слотов" /></div>
    <div className="rounded-2xl border border-border bg-card p-6"><h2 className="text-base font-semibold">Экипировка</h2><dl className="mt-5 grid grid-cols-1 gap-3 sm:grid-cols-2">{equipped.map((item) => <div key={item.label} className="flex items-center justify-between rounded-xl border border-border bg-background px-4 py-3"><dt className="text-sm text-muted-foreground">{item.label}</dt><dd className="text-sm font-medium">{item.value}</dd></div>)}</dl></div>
  </div>
}

function StatCard({ icon: Icon, label, value, suffix, accent }: { icon: React.ComponentType<{ className?: string; 'aria-hidden'?: boolean }>; label: string; value: string | number; suffix: string; accent?: boolean }) {
  return <div className="rounded-xl border border-border bg-card p-4"><Icon className={accent ? 'size-4 text-[color:var(--brand)]' : 'size-4 text-muted-foreground'} aria-hidden={true} /><p className="tabular mt-2 text-2xl font-semibold">{value}</p><p className="text-xs text-muted-foreground">{label} · {suffix}</p></div>
}
