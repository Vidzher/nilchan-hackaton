'use client'

import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { Suspense, useCallback, useEffect, useState } from 'react'
import { ArrowLeft, ExternalLink, ListChecks, Play, RotateCw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { StatusBadge } from '@/components/status-badge'
import { Tag } from '@/components/primitives'
import { api, ApiError, type Resource } from '@/lib/api'

function ResourceContent() {
  const rawId = useSearchParams().get('id')
  const id = rawId && /^\d+$/.test(rawId) ? Number(rawId) : null
  const [resource, setResource] = useState<Resource | null>(null)
  const [questions, setQuestions] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [retrying, setRetrying] = useState(false)

  const load = useCallback(async () => {
    if (id === null) { setError('Материал не найден.'); return }
    try {
      const found = await api.resource(id)
      setResource(found)
      setError(null)
      if (found.status === 'NOT_COMPLETED') {
        const quiz = await api.quiz(id)
        setQuestions(quiz.questions.length)
      }
    } catch (caught) { setError(caught instanceof ApiError ? caught.message : 'Не удалось загрузить материал.') }
  }, [id])

  useEffect(() => { void load() }, [load])
  useEffect(() => {
    if (resource?.status !== 'PROCESSING') return
    const timer = window.setInterval(() => void load(), 3000)
    return () => window.clearInterval(timer)
  }, [load, resource?.status])

  async function retry() {
    if (!resource) return
    setRetrying(true)
    setError(null)
    try {
      const [profile, resources] = await Promise.all([api.profile(), api.resources()])
      const used = resources.filter((item) => item.status === 'PROCESSING' || item.status === 'NOT_COMPLETED').length
      if (used > profile.progress.activeBacklogLimit) {
        setError('Сначала завершите один из материалов: временный слот уже используется.')
        return
      }

      const purchaseOverflowSlot = used === profile.progress.activeBacklogLimit
      if (purchaseOverflowSlot) {
        if (profile.progress.ePoints < 25) {
          setError('Для повторной обработки нужен временный слот за 25 е-баллов.')
          return
        }
        if (!window.confirm('Backlog заполнен. Купить временный слот за 25 е-баллов и повторить обработку?')) return
      }

      await api.createResource(resource.url, purchaseOverflowSlot)
      await load()
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось повторить обработку.')
    } finally { setRetrying(false) }
  }

  if (!resource) return <div className="mx-auto max-w-3xl rounded-2xl border border-border bg-card p-6"><h1 className="text-xl font-semibold">{error ?? 'Загружаем материал…'}</h1><Button className="mt-5" nativeButton={false} render={<Link href="/" />}><ArrowLeft className="size-4" aria-hidden="true" />Вернуться в backlog</Button></div>
  let hostname = resource.url
  try { hostname = new URL(resource.url).hostname } catch { /* keep URL */ }

  return <div className="mx-auto flex max-w-3xl flex-col gap-6">
    <Link href="/" className="inline-flex w-fit items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground"><ArrowLeft className="size-4" aria-hidden="true" />Назад в backlog</Link>
    <div className="rounded-2xl border border-border bg-card p-6 sm:p-7">
      <StatusBadge status={resource.status} /><h1 className="mt-4 text-2xl font-semibold tracking-tight text-balance">{resource.title}</h1>
      <a href={resource.url} target="_blank" rel="noreferrer noopener" className="mt-1 inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground">{hostname}<ExternalLink className="size-3.5" aria-hidden="true" /></a>
      <div className="mt-4 flex flex-wrap gap-1.5">{resource.tags.map((tag) => <Tag key={tag}>{tag}</Tag>)}</div>
      {error ? <p className="mt-4 text-sm text-[color:var(--destructive)]">{error}</p> : null}
      <dl className="mt-6 grid grid-cols-2 gap-4 border-t border-border pt-6 sm:grid-cols-3"><div><dt className="text-xs text-muted-foreground">Добавлено</dt><dd className="mt-0.5 text-sm font-medium">{new Date(resource.createdAt).toLocaleDateString('ru-RU')}</dd></div>{questions !== null ? <div><dt className="text-xs text-muted-foreground">Вопросов в квизе</dt><dd className="tabular mt-0.5 text-sm font-medium">{questions}</dd></div> : null}</dl>
      {resource.status === 'COMPLETED' ? <div className="mt-6 flex gap-3 rounded-xl bg-success-soft p-4"><span className="text-sm font-semibold text-[color:var(--success)]">+{resource.xpEarned ?? 0} XP</span><span className="text-sm font-semibold text-[color:var(--success)]">+{resource.ePointsEarned ?? 0} е-баллов</span></div> : null}
      <div className="mt-7 flex flex-col gap-2 sm:flex-row">{resource.status === 'NOT_COMPLETED' ? <Button nativeButton={false} render={<Link href={`/quiz?id=${resource.id}`} />}><Play className="size-4" aria-hidden="true" />Начать квиз</Button> : null}{resource.status === 'PROCESSING' ? <Button disabled><ListChecks className="size-4" aria-hidden="true" />Квиз создаётся…</Button> : null}{resource.status === 'FAILED' ? <Button variant="secondary" onClick={retry} disabled={retrying}><RotateCw className="size-4" aria-hidden="true" />{retrying ? 'Повторяем…' : 'Повторить создание квиза'}</Button> : null}<Button variant="outline" nativeButton={false} render={<a href={resource.url} target="_blank" rel="noreferrer noopener" />}><ExternalLink className="size-4" aria-hidden="true" />Открыть оригинал</Button></div>
    </div>
  </div>
}

export default function ResourcePage() { return <Suspense fallback={<div className="h-64 animate-pulse rounded-2xl bg-card" />}><ResourceContent /></Suspense> }
