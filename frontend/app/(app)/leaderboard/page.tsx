'use client'

import { useEffect, useState } from 'react'
import { Trophy } from 'lucide-react'
import { cn } from '@/lib/utils'
import { cosmeticName, cosmeticPreview } from '@/lib/cosmetics'
import { AvatarFrame } from '@/components/avatar-frame'
import { api, ApiError, type LeaderboardEntry } from '@/lib/api'

function rankStyle(rank: number) {
  if (rank === 1) return 'bg-[#f2ead4] text-[#8a6a12]'
  if (rank === 2) return 'bg-[#ececec] text-[#5f5f5f]'
  if (rank === 3) return 'bg-[#f0ddd0] text-[#94582f]'
  return 'bg-muted text-muted-foreground'
}

export default function LeaderboardPage() {
  const [entries, setEntries] = useState<LeaderboardEntry[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let active = true
    api.leaderboard()
      .then((result) => {
        if (!active) return
        setEntries(result)
        setError(null)
      })
      .catch((caught) => {
        if (!active) return
        setError(caught instanceof ApiError ? caught.message : 'Не удалось загрузить рейтинг.')
      })
      .finally(() => {
        if (active) setLoading(false)
      })
    return () => { active = false }
  }, [])

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Рейтинг</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Топ разгребателей backlog по накопленному XP.
        </p>
      </div>

      {error ? <p className="rounded-lg bg-warning-soft p-3 text-sm text-[color:var(--destructive)]" role="alert">{error}</p> : null}

      <div className="overflow-hidden rounded-2xl border border-border bg-card">
        <div className="hidden grid-cols-[3rem_1fr_5rem_6rem] gap-4 border-b border-border px-4 py-3 text-xs font-medium uppercase tracking-wide text-muted-foreground sm:grid">
          <span>Место</span>
          <span>Игрок</span>
          <span className="text-right">Уровень</span>
          <span className="text-right">XP</span>
        </div>

        {loading ? (
          <p className="px-4 py-12 text-center text-sm text-muted-foreground">Загружаем рейтинг…</p>
        ) : entries.length ? (
          <ul>
            {entries.map((row) => (
              <li
                key={row.userId}
                className={cn(
                  'grid grid-cols-[2.25rem_1fr_auto] items-center gap-3 border-b border-border px-4 py-3 last:border-b-0 sm:grid-cols-[3rem_1fr_5rem_6rem] sm:gap-4',
                  row.isCurrent && 'bg-brand-soft/50',
                )}
              >
                <span className={cn('tabular grid size-8 place-items-center rounded-full text-sm font-semibold', rankStyle(row.rank))}>
                  {row.rank}
                </span>

                <div className="flex min-w-0 items-center gap-3">
                  <AvatarFrame emoji={cosmeticPreview(row.avatarId)} frame={cosmeticName(row.frameId)} size="sm" />
                  <div className="min-w-0">
                    <p className="flex items-center gap-2 truncate text-sm font-medium">
                      {row.username}
                      {row.isCurrent ? (
                        <span className="rounded bg-[color:var(--brand)] px-1.5 py-0.5 text-[10px] font-semibold text-white">Вы</span>
                      ) : null}
                    </p>
                    <p className="text-xs text-muted-foreground sm:hidden">
                      Ур. {row.level} · {row.xp.toLocaleString('ru-RU')} XP
                    </p>
                  </div>
                </div>

                <span className="tabular hidden text-right text-sm text-muted-foreground sm:block">{row.level}</span>
                <span className="tabular hidden text-right text-sm font-semibold sm:block">{row.xp.toLocaleString('ru-RU')}</span>
              </li>
            ))}
          </ul>
        ) : (
          <div className="flex flex-col items-center px-4 py-12 text-center text-sm text-muted-foreground">
            <Trophy className="mb-3 size-7" aria-hidden="true" />
            В рейтинге пока никого нет.
          </div>
        )}
      </div>
    </div>
  )
}
