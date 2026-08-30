'use client'

import { useEffect, useState } from 'react'
import { Check, Coins, Lock } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { TitleBadge } from '@/components/title-badge'
import { cn } from '@/lib/utils'
import {
  api,
  ApiError,
  type CosmeticsUpdate,
  type CosmeticType,
  type Profile,
  type ShopItem,
} from '@/lib/api'

type Category = 'Все' | CosmeticType

const categories: { value: Category; label: string }[] = [
  { value: 'Все', label: 'Все' },
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
  const [category, setCategory] = useState<Category>('Все')
  const [busyItemID, setBusyItemID] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    Promise.all([api.shop(), api.profile()])
      .then(([catalog, foundProfile]) => {
        setItems(catalog)
        setProfile(foundProfile)
      })
      .catch((caught) => {
        setError(caught instanceof ApiError ? caught.message : 'Не удалось загрузить магазин.')
      })
  }, [])

  const visible = items.filter((item) => category === 'Все' || item.type === category)
  const balance = profile?.progress.ePoints ?? 0
  const owned = new Set(profile?.cosmetics.ownedCosmeticIds ?? [])

  async function buy(item: ShopItem) {
    setBusyItemID(item.id)
    setError(null)
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
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось купить косметику.')
    } finally {
      setBusyItemID(null)
    }
  }

  async function equip(item: ShopItem) {
    setBusyItemID(item.id)
    setError(null)
    try {
      setProfile(await api.updateCosmetics(equipUpdate(item)))
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось надеть косметику.')
    } finally {
      setBusyItemID(null)
    }
  }

  if (!profile) {
    return <div className="rounded-2xl border border-border bg-card p-6"><h1 className="text-xl font-semibold">{error ?? 'Загружаем магазин…'}</h1></div>
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-semibold tracking-tight">Магазин</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Тратьте е-баллы на косметику — она видна в рейтинге.
          </p>
        </div>
        <span className="inline-flex items-center gap-2 rounded-xl border border-border bg-card px-4 py-2.5">
          <Coins className="size-4 text-[color:var(--brand)]" aria-hidden="true" />
          <span className="tabular text-lg font-semibold">{balance}</span>
          <span className="text-sm text-muted-foreground">е-баллов</span>
        </span>
      </div>

      {error ? <p className="rounded-lg bg-warning-soft p-3 text-sm text-[color:var(--destructive)]" role="alert">{error}</p> : null}

      <div className="flex flex-wrap gap-2" role="group" aria-label="Категории">
        {categories.map((item) => {
          const active = category === item.value
          return (
            <button
              key={item.value}
              type="button"
              onClick={() => setCategory(item.value)}
              aria-pressed={active}
              className={cn(
                'rounded-full border px-3 py-1.5 text-sm font-medium transition-colors',
                active
                  ? 'border-primary bg-primary text-primary-foreground'
                  : 'border-border bg-card text-muted-foreground hover:text-foreground',
              )}
            >
              {item.label}
            </button>
          )
        })}
      </div>

      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
        {visible.map((item) => {
          const isOwned = owned.has(item.id)
          const isEquipped = equippedItemID(profile, item.type) === item.id
          const affordable = balance >= item.price
          const busy = busyItemID === item.id
          return (
            <div
              key={item.id}
              className={cn(
                'flex flex-col rounded-xl border bg-card p-4',
                isEquipped ? 'border-[color:var(--brand)]' : 'border-border',
              )}
            >
              <div className="flex items-start justify-between">
                <span className="text-xs text-muted-foreground">{typeLabels[item.type]}</span>
                {isEquipped ? (
                  <span className="inline-flex items-center gap-1 rounded-md bg-brand-soft px-1.5 py-0.5 text-xs font-medium text-[color:var(--brand)]">
                    <Check className="size-3" aria-hidden="true" />
                    Надето
                  </span>
                ) : null}
              </div>

              <div className="mt-3 grid h-20 place-items-center rounded-lg bg-background">
                {item.type === 'title' ? (
                  <div className="flex max-w-full flex-col items-center gap-1.5 px-2">
                    <span className="max-w-full truncate text-xs font-medium">{profile.user.username}</span>
                    <TitleBadge name={item.displayName} price={item.price} />
                  </div>
                ) : item.presentation.emoji ? (
                  <span className="text-4xl" aria-hidden="true">{item.presentation.emoji}</span>
                ) : (
                  <span className="px-2 text-center text-sm font-semibold text-pretty">{item.displayName}</span>
                )}
              </div>

              <p className="mt-3 text-sm font-medium text-pretty">{item.displayName}</p>

              <div className="mt-3">
                {!isOwned ? (
                  <Button
                    size="sm"
                    className="w-full"
                    disabled={!affordable || busyItemID !== null}
                    onClick={() => void buy(item)}
                  >
                    {affordable ? <Coins className="size-3.5" aria-hidden="true" /> : <Lock className="size-3.5" aria-hidden="true" />}
                    {busy ? 'Покупаем…' : item.price}
                  </Button>
                ) : isEquipped ? (
                  <Button size="sm" variant="outline" className="w-full" disabled>Надето</Button>
                ) : (
                  <Button
                    size="sm"
                    variant="secondary"
                    className="w-full"
                    disabled={busyItemID !== null}
                    onClick={() => void equip(item)}
                  >
                    {busy ? 'Надеваем…' : 'Надеть'}
                  </Button>
                )}
              </div>
            </div>
          )
        })}
      </div>
    </div>
  )
}
