'use client'

import { useState } from 'react'
import { Check, Coins, Lock } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'
import { shopItems, user, type ShopItem } from '@/lib/data'

type Category = 'Все' | ShopItem['type']
const categories: Category[] = ['Все', 'Аватар', 'Рамка', 'Титул', 'Витрина']

const emojiTypes: ShopItem['type'][] = ['Аватар', 'Витрина']

export function ShopView() {
  const [items, setItems] = useState<ShopItem[]>(shopItems)
  const [balance, setBalance] = useState(user.points)
  const [category, setCategory] = useState<Category>('Все')

  const visible = items.filter((i) => category === 'Все' || i.type === category)

  function buy(item: ShopItem) {
    if (balance < item.price) return
    setBalance((b) => b - item.price)
    setItems((prev) =>
      prev.map((i) => (i.id === item.id ? { ...i, state: 'equip' } : i)),
    )
  }

  function equip(item: ShopItem) {
    setItems((prev) =>
      prev.map((i) => {
        if (i.type !== item.type) return i
        if (i.id === item.id) return { ...i, state: 'equipped' }
        if (i.state === 'equipped') return { ...i, state: 'equip' }
        return i
      }),
    )
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

      <div className="flex flex-wrap gap-2" role="group" aria-label="Категории">
        {categories.map((c) => {
          const active = category === c
          return (
            <button
              key={c}
              type="button"
              onClick={() => setCategory(c)}
              aria-pressed={active}
              className={cn(
                'rounded-full border px-3 py-1.5 text-sm font-medium transition-colors',
                active
                  ? 'border-primary bg-primary text-primary-foreground'
                  : 'border-border bg-card text-muted-foreground hover:text-foreground',
              )}
            >
              {c}
            </button>
          )
        })}
      </div>

      <div className="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-4">
        {visible.map((item) => {
          const isEmoji = emojiTypes.includes(item.type)
          const affordable = balance >= item.price
          return (
            <div
              key={item.id}
              className={cn(
                'flex flex-col rounded-xl border bg-card p-4',
                item.state === 'equipped'
                  ? 'border-[color:var(--brand)]'
                  : 'border-border',
              )}
            >
              <div className="flex items-start justify-between">
                <span className="text-xs text-muted-foreground">{item.type}</span>
                {item.state === 'equipped' ? (
                  <span className="inline-flex items-center gap-1 rounded-md bg-brand-soft px-1.5 py-0.5 text-xs font-medium text-[color:var(--brand)]">
                    <Check className="size-3" aria-hidden="true" />
                    Надето
                  </span>
                ) : null}
              </div>

              <div className="mt-3 grid h-20 place-items-center rounded-lg bg-background">
                {isEmoji ? (
                  <span className="text-4xl" aria-hidden="true">
                    {item.preview}
                  </span>
                ) : (
                  <span className="px-2 text-center text-sm font-semibold text-pretty">
                    {item.preview}
                  </span>
                )}
              </div>

              <p className="mt-3 text-sm font-medium text-pretty">{item.name}</p>

              <div className="mt-3">
                {item.state === 'buy' ? (
                  <Button
                    size="sm"
                    className="w-full"
                    disabled={!affordable}
                    onClick={() => buy(item)}
                  >
                    {affordable ? (
                      <>
                        <Coins className="size-3.5" aria-hidden="true" />
                        {item.price}
                      </>
                    ) : (
                      <>
                        <Lock className="size-3.5" aria-hidden="true" />
                        {item.price}
                      </>
                    )}
                  </Button>
                ) : item.state === 'equipped' ? (
                  <Button size="sm" variant="outline" className="w-full" disabled>
                    Надето
                  </Button>
                ) : (
                  <Button
                    size="sm"
                    variant="secondary"
                    className="w-full"
                    onClick={() => equip(item)}
                  >
                    Надеть
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
