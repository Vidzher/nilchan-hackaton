'use client'

import Link from 'next/link'
import { useEffect, useState } from 'react'
import { CircleCheck, Coins, Flame, Layers } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { AvatarFrame, avatarImagePath } from '@/components/avatar-frame'
import { PageHeader, ProgressBar } from '@/components/primitives'
import { TitleBadge } from '@/components/title-badge'
import { getLevelProgress } from '@/lib/progress'
import { api, ApiError, type Profile, type ShopItem } from '@/lib/api'

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

  if (!profile) return <div className="border-y border-border py-12"><h1 className="text-xl font-semibold">{error ?? 'Загружаем профиль…'}</h1></div>

  const { progress, user, cosmetics } = profile
  const levelProgress = getLevelProgress(progress.xp, progress.level)
  const equipped = [
    { label: 'Аватар', value: avatar?.displayName ?? (cosmetics.avatarId === 'default_avatar' ? 'Нормис' : 'Не выбрано') },
    { label: 'Рамка', value: frame?.displayName ?? (cosmetics.frameId === 'default_frame' ? 'Без рамки' : 'Не выбрано') },
    { label: 'Титул', value: title?.displayName ?? 'Не выбрано' },
    { label: 'Витрина', value: showcase?.displayName ?? 'Не выбрано' },
  ]

  return (
    <div className="flex flex-col gap-7">
      <PageHeader title="Профиль" description="Ваш прогресс и оформление, которое видят другие." />

      <section className="grid overflow-hidden rounded-xl border border-border bg-card lg:grid-cols-[1fr_1.25fr]">
        <div className="flex flex-col gap-5 p-5 sm:flex-row sm:items-center sm:p-6">
          <AvatarFrame src={avatarImagePath(avatar?.presentation.assetKey)} fallback={avatar?.presentation.emoji} frame={frame?.presentation.cssClass} size="lg" />
          <div className="min-w-0 flex-1">
            <div className="flex flex-wrap items-center gap-2">
              <h1 className="truncate text-xl font-semibold">{user.username}</h1>
              <span className="border-l border-border pl-2 text-xs font-medium text-muted-foreground">Уровень {progress.level}</span>
            </div>
            {title ? <TitleBadge className="mt-2" name={title.displayName} price={title.price} /> : null}
            {showcase ? <p className="mt-2 text-sm text-muted-foreground">{showcase.presentation.emoji} {showcase.displayName}</p> : null}
            <Button className="mt-5" size="sm" variant="outline" nativeButton={false} render={<Link href="/shop" />}>Изменить оформление</Button>
          </div>
        </div>

        <div className="border-t border-border p-5 sm:p-6 lg:border-l lg:border-t-0">
          <div className="flex items-start justify-between gap-4">
            <div>
              <p className="text-xs font-medium uppercase tracking-[0.14em] text-muted-foreground">Опыт</p>
              <p className="tabular mt-2 text-3xl font-semibold tracking-tight">{progress.xp} XP</p>
            </div>
            <p className="tabular text-right text-sm text-muted-foreground">
              {levelProgress.isMaxLevel ? 'Максимальный уровень' : `${levelProgress.remaining} XP до уровня ${progress.level + 1}`}
            </p>
          </div>
          <ProgressBar className="mt-5 h-2" value={levelProgress.current} max={levelProgress.required} label="Прогресс текущего уровня" />
          {!levelProgress.isMaxLevel ? <p className="tabular mt-2 text-xs text-muted-foreground">{levelProgress.current} из {levelProgress.required} XP на этом уровне</p> : null}
        </div>
      </section>

      <dl className="grid grid-cols-2 border-y border-border lg:grid-cols-4">
        <ProfileStat icon={CircleCheck} label="Завершено" value={completed} suffix="мат." />
        <ProfileStat icon={Coins} label="Баланс" value={progress.ePoints} suffix="е-баллов" />
        <ProfileStat icon={Flame} label="Серия" value={progress.currentStreak} suffix="дн." accent />
        <ProfileStat icon={Layers} label="Backlog" value={`${used}/${progress.activeBacklogLimit}`} suffix="слотов" last />
      </dl>

      <section>
        <h2 className="text-base font-semibold tracking-tight">Текущее оформление</h2>
        <dl className="mt-4 divide-y divide-border border-y border-border">
          {equipped.map((item) => (
            <div key={item.label} className="flex items-center justify-between gap-6 py-3.5">
              <dt className="text-sm text-muted-foreground">{item.label}</dt>
              <dd className="truncate text-right text-sm font-medium">{item.value}</dd>
            </div>
          ))}
        </dl>
      </section>
    </div>
  )
}

function ProfileStat({
  icon: Icon,
  label,
  value,
  suffix,
  accent,
  last,
}: {
  icon: React.ComponentType<{ className?: string; 'aria-hidden'?: boolean }>
  label: string
  value: string | number
  suffix: string
  accent?: boolean
  last?: boolean
}) {
  return (
    <div className={`px-4 py-5 ${last ? '' : 'border-r border-border'} max-lg:[&:nth-child(2)]:border-r-0 max-lg:[&:nth-child(-n+2)]:border-b`}>
      <dt className="flex items-center gap-1.5 text-xs text-muted-foreground"><Icon className={accent ? 'size-3.5 text-[color:var(--brand)]' : 'size-3.5'} aria-hidden={true} />{label}</dt>
      <dd className="tabular mt-2 text-2xl font-semibold">{value}<span className="ml-1.5 text-xs font-normal text-muted-foreground">{suffix}</span></dd>
    </div>
  )
}
