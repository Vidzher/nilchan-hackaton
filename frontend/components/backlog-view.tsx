'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import { Flame, Layers, Plus, Search, Coins } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ProgressBar } from '@/components/primitives'
import { ResourceCard } from '@/components/resource-card'
import { AddResourceModal } from '@/components/add-resource-modal'
import { inputClass } from '@/components/field'
import { cn } from '@/lib/utils'
import { api, ApiError, type Profile, type Resource, type ResourceStatus } from '@/lib/api'

type Filter = 'ALL' | ResourceStatus
const filters: { key: Filter; label: string }[] = [
  { key: 'ALL', label: 'Все' }, { key: 'PROCESSING', label: 'В работе' },
  { key: 'NOT_COMPLETED', label: 'Готовы к quiz' }, { key: 'FAILED', label: 'Ошибка' },
  { key: 'COMPLETED', label: 'Завершённые' },
]
const levelCeilings: Record<number, number> = { 1: 120, 2: 300, 3: 600, 4: 1000 }

function SummaryBlock({ label, children }: { label: string; children: React.ReactNode }) {
  return <div className="rounded-xl border border-border bg-card p-3.5"><p className="text-xs text-muted-foreground">{label}</p>{children}</div>
}

export function BacklogView() {
  const [open, setOpen] = useState(false)
  const [query, setQuery] = useState('')
  const [filter, setFilter] = useState<Filter>('ALL')
  const [tag, setTag] = useState('ALL')
  const [resources, setResources] = useState<Resource[]>([])
  const [profile, setProfile] = useState<Profile | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [retrying, setRetrying] = useState<number | null>(null)

  const load = useCallback(async () => {
    try {
      const [nextResources, nextProfile] = await Promise.all([api.resources(), api.profile()])
      setResources(nextResources)
      setProfile(nextProfile)
      setError(null)
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось загрузить backlog.')
    }
  }, [])

  useEffect(() => { void load() }, [load])
  useEffect(() => {
    if (!resources.some((resource) => resource.status === 'PROCESSING')) return
    const timer = window.setInterval(() => void load(), 3000)
    return () => window.clearInterval(timer)
  }, [load, resources])

  const used = resources.filter((resource) => resource.status === 'PROCESSING' || resource.status === 'NOT_COMPLETED').length
  const progress = profile?.progress

  async function retry(resource: Resource) {
    if (!progress) return
    if (used > progress.activeBacklogLimit) {
      setError('Сначала завершите один из материалов: временный слот уже используется.')
      return
    }

    const purchaseOverflowSlot = used === progress.activeBacklogLimit
    if (purchaseOverflowSlot) {
      if (progress.ePoints < 25) {
        setError('Для повторной обработки нужен временный слот за 25 е-баллов.')
        return
      }
      if (!window.confirm('Backlog заполнен. Купить временный слот за 25 е-баллов и повторить обработку?')) return
    }

    setRetrying(resource.id)
    try {
      await api.createResource(resource.url, purchaseOverflowSlot)
      await load()
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось повторить обработку.')
    } finally { setRetrying(null) }
  }

  const allTags = useMemo(() => Array.from(new Set(resources.flatMap((resource) => resource.tags))).sort(), [resources])
  const filtered = resources.filter((resource) => resource.title.toLowerCase().includes(query.toLowerCase()) && (filter === 'ALL' || resource.status === filter) && (tag === 'ALL' || resource.tags.includes(tag)))
  const ceiling = progress ? (levelCeilings[progress.level] ?? Math.max(progress.xp, 1000)) : 1

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-wrap items-start justify-between gap-4"><div><h1 className="text-2xl font-semibold tracking-tight">Backlog</h1><p className="mt-1 text-sm text-muted-foreground">Материалы, которые вы решили закончить.</p></div><Button onClick={() => setOpen(true)}><Plus className="size-4" aria-hidden="true" />Добавить материал</Button></div>
      {error ? <p className="rounded-lg bg-warning-soft p-3 text-sm text-[color:var(--destructive)]" role="alert">{error}</p> : null}
      <div className="grid grid-cols-2 gap-3 lg:grid-cols-4">
        <SummaryBlock label={`Уровень ${progress?.level ?? '…'}`}><p className="tabular mt-1 text-sm font-semibold">{progress?.xp ?? 0} / {ceiling} XP</p><ProgressBar className="mt-2" value={progress?.xp ?? 0} max={ceiling} label="Прогресс до следующего уровня" /></SummaryBlock>
        <SummaryBlock label="Баланс"><p className="mt-1 flex items-center gap-1.5"><Coins className="size-4 text-muted-foreground" aria-hidden="true" /><span className="tabular text-lg font-semibold">{progress?.ePoints ?? 0}</span><span className="text-sm text-muted-foreground">е-баллов</span></p></SummaryBlock>
        <SummaryBlock label="Серия"><p className="mt-1 flex items-center gap-1.5"><Flame className="size-4 text-[color:var(--brand)]" aria-hidden="true" /><span className="tabular text-lg font-semibold">{progress?.currentStreak ?? 0}</span><span className="text-sm text-muted-foreground">дней</span></p></SummaryBlock>
        <SummaryBlock label="Backlog"><p className="mt-1 flex items-center gap-1.5"><Layers className="size-4 text-muted-foreground" aria-hidden="true" /><span className="tabular text-lg font-semibold">{used} из {progress?.activeBacklogLimit ?? '…'}</span><span className="text-sm text-muted-foreground">слотов</span></p></SummaryBlock>
      </div>
      <div className="flex flex-col gap-3"><div className="flex flex-col gap-3 sm:flex-row"><div className="relative flex-1"><Search className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" aria-hidden="true" /><label htmlFor="backlog-search" className="sr-only">Поиск по названию</label><input id="backlog-search" type="search" value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Поиск по названию" className={cn(inputClass, 'pl-9')} /></div><select value={tag} onChange={(event) => setTag(event.target.value)} className={cn(inputClass, 'cursor-pointer sm:w-52')} aria-label="Фильтр по тегу"><option value="ALL">Все теги</option>{allTags.map((value) => <option key={value} value={value}>{value}</option>)}</select></div>
        <div className="flex flex-wrap gap-2" role="group" aria-label="Фильтр по статусу">{filters.map((item) => <button key={item.key} type="button" onClick={() => setFilter(item.key)} aria-pressed={filter === item.key} className={cn('rounded-full border px-3 py-1.5 text-sm font-medium transition-colors', filter === item.key ? 'border-primary bg-primary text-primary-foreground' : 'border-border bg-card text-muted-foreground hover:text-foreground')}>{item.label}</button>)}</div></div>
      {filtered.length ? (
        <div className="grid grid-cols-1 gap-4 xl:grid-cols-2">
          {filtered.map((resource) => <ResourceCard key={resource.id} resource={resource} onRetry={retry} retrying={retrying === resource.id} />)}
        </div>
      ) : resources.length ? (
        <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-border bg-card px-6 py-16 text-center">
          <Search className="size-8 text-muted-foreground" aria-hidden="true" />
          <h2 className="mt-4 text-lg font-semibold">Ничего не найдено</h2>
          <p className="mt-1 text-sm text-muted-foreground">Попробуйте изменить поиск или фильтры.</p>
          <Button
            className="mt-5"
            variant="outline"
            onClick={() => {
              setQuery('')
              setFilter('ALL')
              setTag('ALL')
            }}
          >
            Сбросить фильтры
          </Button>
        </div>
      ) : (
        <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-border bg-card px-6 py-16 text-center">
          <Layers className="size-8 text-muted-foreground" aria-hidden="true" />
          <h2 className="mt-4 text-lg font-semibold">Backlog пуст</h2>
          <Button className="mt-5" onClick={() => setOpen(true)}><Plus className="size-4" aria-hidden="true" />Добавить материал</Button>
        </div>
      )}
      <AddResourceModal open={open} onClose={() => setOpen(false)} onCreated={load} slotsUsed={used} slotsTotal={progress?.activeBacklogLimit ?? 0} balance={progress?.ePoints ?? 0} />
    </div>
  )
}
