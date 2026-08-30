'use client'

import { useEffect, useMemo, useState } from 'react'
import { Check, Coins, Lock, SlidersHorizontal } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { AvatarFrame, avatarImagePath } from '@/components/avatar-frame'
import { PageHeader, SectionHeader } from '@/components/primitives'
import { TitleBadge } from '@/components/title-badge'
import { inputClass } from '@/components/field'
import { cn } from '@/lib/utils'
import {
  api,
  ApiError,
  type CosmeticsUpdate,
  type CosmeticType,
  type Profile,
  type ShopItem,
} from '@/lib/api'

type View = 'catalog' | 'owned'
type Category = 'all' | CosmeticType
type Sort = 'featured' | 'price-asc' | 'price-desc'

const categories: { value: Category; label: string }[] = [
  { value: 'all', label: 'Все' },
  { value: 'avatar', label: 'Аватары' },
  { value: 'frame', label: 'Рамки' },
  { value: 'title', label: 'Титулы' },
  { value: 'showcase', label: 'Витрина' },
]

const typeLabels: Record<CosmeticType, string> = {
  avatar: 'Аватар',
  frame: 'Рамка',
  title: 'Титул',
  showcase: 'Витрина',
}

const defaultItems: ShopItem[] = [
  {
    id: 'default_avatar',
    type: 'avatar',
    displayName: 'Нормис',
    price: 0,
    presentation: { emoji: '🙂' },
  },
  {
    id: 'default_frame',
    type: 'frame',
    displayName: 'Без рамки',
    price: 0,
    presentation: { cssClass: 'frame-default' },
  },
]

function equippedItemID(profile: Profile, type: CosmeticType) {
  if (type === 'avatar') return profile.cosmetics.avatarId
  if (type === 'frame') return profile.cosmetics.frameId
  if (type === 'title') return profile.cosmetics.titleId
  return profile.cosmetics.showcaseItemId
}

function equipUpdate(item: ShopItem): CosmeticsUpdate {
  if (item.type === 'avatar') return { avatarId: item.id }
  if (item.type === 'frame') return { frameId: item.id }
  if (item.type === 'title') return { titleId: item.id }
  return { showcaseItemId: item.id }
}

export function ShopView() {
  const [items, setItems] = useState<ShopItem[]>([])
  const [profile, setProfile] = useState<Profile | null>(null)
  const [view, setView] = useState<View>('catalog')
  const [category, setCategory] = useState<Category>('all')
  const [sort, setSort] = useState<Sort>('featured')
  const [busyItemID, setBusyItemID] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  useEffect(() => {
    Promise.all([api.shop(), api.profile()])
      .then(([catalog, foundProfile]) => {
        setItems(Array.from(new Map([...defaultItems, ...catalog].map((item) => [item.id, item])).values()))
        setProfile(foundProfile)
        setError(null)
      })
      .catch((caught) => {
        setError(caught instanceof ApiError ? caught.message : 'Не удалось загрузить магазин.')
      })
  }, [])

  const balance = profile?.progress.ePoints ?? 0
  const owned = useMemo(() => new Set([...defaultItems.map((item) => item.id), ...(profile?.cosmetics.ownedCosmeticIds ?? [])]), [profile])
  const equippedAvatar = items.find((item) => item.id === profile?.cosmetics.avatarId)
  const equippedFrame = items.find((item) => item.id === profile?.cosmetics.frameId)

  const visible = useMemo(() => {
    const filtered = items.filter((item) => {
      const matchesOwnership = view === 'owned' ? owned.has(item.id) : true
      const matchesCategory = category === 'all' || item.type === category
      return matchesOwnership && matchesCategory
    })
    if (sort === 'price-asc') return [...filtered].sort((left, right) => left.price - right.price)
    if (sort === 'price-desc') return [...filtered].sort((left, right) => right.price - left.price)
    if (view === 'owned' && profile) {
      return [...filtered].sort((left, right) => Number(equippedItemID(profile, right.type) === right.id) - Number(equippedItemID(profile, left.type) === left.id))
    }
    return filtered
  }, [category, items, owned, profile, sort, view])

  async function buy(item: ShopItem) {
    if (item.price >= 100 && !window.confirm(`Купить «${item.displayName}» за ${item.price} е-баллов?`)) return
    setBusyItemID(item.id)
    setError(null)
    setNotice(null)
    try {
      const result = await api.purchaseCosmetic(item.id)
      setProfile((current) => current && ({
        ...current,
        progress: { ...current.progress, ePoints: result.ePoints },
        cosmetics: {
          ...current.cosmetics,
          ownedCosmeticIds: Array.from(new Set([...current.cosmetics.ownedCosmeticIds, result.item.id])),
        },
      }))
      setNotice(`«${item.displayName}» добавлено в коллекцию.`)
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось купить косметику.')
    } finally {
      setBusyItemID(null)
    }
  }

  async function equip(item: ShopItem) {
    setBusyItemID(item.id)
    setError(null)
    setNotice(null)
    try {
      setProfile(await api.updateCosmetics(equipUpdate(item)))
      window.dispatchEvent(new Event('cosmeticschange'))
      setNotice(`Теперь надето: «${item.displayName}».`)
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось надеть косметику.')
    } finally {
      setBusyItemID(null)
    }
  }

  if (!profile) {
    return <div className="border-y border-border py-12"><h1 className="text-xl font-semibold">{error ?? 'Загружаем магазин…'}</h1></div>
  }

  return (
    <div className="flex flex-col gap-7">
      <PageHeader
        title="Магазин"
        description="Оформление профиля, которое видно в рейтинге. Без бустеров и pay-to-win."
        action={
          <div className="flex items-baseline gap-2 border-l border-border pl-4">
            <Coins className="size-4 text-[color:var(--brand)]" aria-hidden="true" />
            <span className="tabular text-2xl font-semibold">{balance}</span>
            <span className="text-sm text-muted-foreground">е-баллов</span>
          </div>
        }
      />

      {error ? <p className="border-l-2 border-[color:var(--destructive)] bg-card px-4 py-3 text-sm text-[color:var(--destructive)]" role="alert">{error}</p> : null}
      {notice ? <p className="border-l-2 border-[color:var(--success)] bg-card px-4 py-3 text-sm text-[color:var(--success)]" role="status">{notice}</p> : null}

      <div className="flex w-fit rounded-lg bg-secondary p-1" role="tablist" aria-label="Раздел магазина">
        {([
          { value: 'catalog', label: 'Все предметы', count: items.length },
          { value: 'owned', label: 'Моя коллекция', count: items.filter((item) => owned.has(item.id)).length },
        ] as const).map((item) => (
          <button
            key={item.value}
            type="button"
            role="tab"
            aria-selected={view === item.value}
            onClick={() => setView(item.value)}
            className={cn(
              'rounded-md px-4 py-2 text-sm font-semibold transition-colors',
              view === item.value ? 'bg-card text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground',
            )}
          >
            {item.label}
            <span className="tabular ml-2 text-xs font-normal text-muted-foreground">{item.count}</span>
          </button>
        ))}
      </div>

      <div className="flex flex-col gap-3 rounded-xl border border-border bg-card p-3 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex gap-1 overflow-x-auto" role="group" aria-label="Категории предметов">
          {categories.map((item) => {
            const active = category === item.value
            return (
              <button
                key={item.value}
                type="button"
                onClick={() => setCategory(item.value)}
                aria-pressed={active}
                className={cn(
                  'shrink-0 rounded-md px-3 py-2 text-sm font-medium transition-colors',
                  active ? 'bg-brand-soft text-[color:var(--brand)]' : 'text-muted-foreground hover:bg-secondary hover:text-foreground',
                )}
              >
                {item.label}
              </button>
            )
          })}
        </div>
        <label className="flex shrink-0 items-center gap-2 text-sm text-muted-foreground">
          <SlidersHorizontal className="size-4" aria-hidden="true" />
          <span className="sr-only">Сортировка</span>
          <select value={sort} onChange={(event) => setSort(event.target.value as Sort)} className={cn(inputClass, 'h-9 w-full cursor-pointer bg-background sm:w-auto sm:min-w-40')}>
            <option value="featured">По умолчанию</option>
            <option value="price-asc">Сначала дешевле</option>
            <option value="price-desc">Сначала дороже</option>
          </select>
        </label>
      </div>

      {visible.length ? (
        <section className="flex flex-col gap-4">
          <SectionHeader
            title={view === 'owned' ? 'Моя коллекция' : 'Все предметы'}
            description={view === 'owned' ? 'Выберите предмет и примените его к профилю.' : 'Купленные предметы не исчезают — их можно применить прямо здесь.'}
            aside={<span className="tabular rounded-full bg-secondary px-2.5 py-1 text-xs text-muted-foreground">{visible.length}</span>}
          />
          <div className="grid items-stretch gap-4 sm:grid-cols-2 xl:grid-cols-3">
            {visible.map((item) => (
              <ShopItemCard
                key={item.id}
                item={item}
                profile={profile}
                equippedAvatar={equippedAvatar}
                equippedFrame={equippedFrame}
                isOwned={owned.has(item.id)}
                isEquipped={equippedItemID(profile, item.type) === item.id}
                balance={balance}
                busyItemID={busyItemID}
                onBuy={buy}
                onEquip={equip}
              />
            ))}
          </div>
        </section>
      ) : <p className="border-y border-border py-12 text-sm text-muted-foreground">В этой категории ничего нет.</p>}
    </div>
  )
}

function ShopItemCard({
  item,
  profile,
  equippedAvatar,
  equippedFrame,
  isOwned,
  isEquipped,
  balance,
  busyItemID,
  onBuy,
  onEquip,
}: {
  item: ShopItem
  profile: Profile
  equippedAvatar?: ShopItem
  equippedFrame?: ShopItem
  isOwned: boolean
  isEquipped: boolean
  balance: number
  busyItemID: string | null
  onBuy: (item: ShopItem) => Promise<void>
  onEquip: (item: ShopItem) => Promise<void>
}) {
  const affordable = balance >= item.price
  const busy = busyItemID === item.id
  const shortfall = Math.max(0, item.price - balance)

  return (
    <article className={cn(
      'group flex h-full min-w-0 flex-col rounded-xl border bg-card p-4 shadow-[0_6px_22px_rgba(86,54,37,0.04)] transition-[border-color,box-shadow]',
      isEquipped ? 'border-[color:var(--brand)]/35 ring-1 ring-[color:var(--brand)]/10' : 'border-border hover:border-[color:var(--brand)]/25 hover:shadow-[0_10px_28px_rgba(86,54,37,0.07)]',
    )}>
      <div className="flex items-center justify-between gap-3">
        <span className="rounded-md bg-secondary px-2 py-1 text-[11px] font-semibold uppercase tracking-[0.1em] text-muted-foreground">{typeLabels[item.type]}</span>
        {isEquipped ? <span className="inline-flex items-center gap-1 rounded-full bg-brand-soft px-2 py-1 text-xs font-semibold text-[color:var(--brand)]"><Check className="size-3.5" aria-hidden="true" />Используется</span> : isOwned ? <span className="inline-flex items-center gap-1 text-xs font-medium text-[color:var(--success)]"><Check className="size-3.5" aria-hidden="true" />Куплено</span> : null}
      </div>

      <div className="mt-3 grid h-36 place-items-center overflow-hidden rounded-lg border border-border/60 bg-secondary/45 transition-colors group-hover:bg-secondary/70">
        {item.type === 'title' ? (
          <div className="flex max-w-full flex-col items-center gap-2 px-4">
            <span className="max-w-full truncate text-sm font-medium">{profile.user.username}</span>
            <TitleBadge name={item.displayName} price={item.price} />
          </div>
        ) : item.type === 'avatar' ? (
          <AvatarFrame src={avatarImagePath(item.presentation.assetKey)} fallback={item.presentation.emoji} frame={equippedFrame?.presentation.cssClass} size="lg" />
        ) : item.type === 'frame' ? (
          <AvatarFrame src={avatarImagePath(equippedAvatar?.presentation.assetKey)} fallback={equippedAvatar?.presentation.emoji} frame={item.presentation.cssClass} size="lg" />
        ) : item.presentation.emoji ? (
          <span className="text-5xl" aria-hidden="true">{item.presentation.emoji}</span>
        ) : (
          <span className="px-4 text-center text-sm font-semibold text-pretty">{item.displayName}</span>
        )}
      </div>

      <div className="mt-4 flex items-start justify-between gap-3">
        <div className="min-w-0">
          <h3 className="truncate text-sm font-semibold">{item.displayName}</h3>
          <p className="mt-1 text-xs text-muted-foreground">
            {isEquipped ? 'Сейчас отображается в профиле' : isOwned ? 'Доступно в вашей коллекции' : 'Останется в коллекции навсегда'}
          </p>
        </div>
        {!isOwned ? <span className="tabular flex shrink-0 items-center gap-1 rounded-md bg-secondary px-2 py-1 text-sm font-semibold"><Coins className="size-3.5 text-[color:var(--brand)]" aria-hidden="true" />{item.price}</span> : null}
      </div>

      <div className="mt-auto pt-4">
        {!isOwned ? (
          <>
            <Button size="sm" className="w-full" disabled={!affordable || busyItemID !== null} onClick={() => void onBuy(item)}>
              {affordable ? <Coins className="size-3.5" aria-hidden="true" /> : <Lock className="size-3.5" aria-hidden="true" />}
              {busy ? 'Покупаем…' : `Купить за ${item.price}`}
            </Button>
            {!affordable ? <p className="tabular mt-2 text-center text-xs text-[color:var(--destructive)]">Не хватает {shortfall} е-баллов</p> : null}
          </>
        ) : isEquipped ? (
          <Button size="sm" variant="outline" className="w-full" disabled><Check className="size-3.5" aria-hidden="true" />Выбрано</Button>
        ) : (
          <Button size="sm" variant="secondary" className="w-full" disabled={busyItemID !== null} onClick={() => void onEquip(item)}>
            {busy ? 'Применяем…' : 'Применить к профилю'}
          </Button>
        )}
      </div>
    </article>
  )
}
