import Link from 'next/link'
import { ArrowRight, Calendar, Coins, ExternalLink, Globe2, Layers, Play, RotateCw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { StatusBadge } from '@/components/status-badge'
import { Tag } from '@/components/primitives'
import type { Resource } from '@/lib/api'

function domain(url: string) {
  try {
    return new URL(url).hostname.replace(/^www\./, '')
  } catch {
    return url
  }
}

function SourceIdentity({ url }: { url: string }) {
  return (
    <span className="inline-flex min-w-0 items-center gap-1.5 text-xs font-medium text-muted-foreground">
      <Globe2 className="size-3.5 shrink-0" aria-hidden="true" />
      <span className="truncate">{domain(url)}</span>
    </span>
  )
}

function ResourceTags({ tags }: { tags: string[] }) {
  return <div className="flex flex-wrap gap-1.5">{tags.map((tag) => <Tag key={tag}>{tag}</Tag>)}</div>
}

export function ContinueLearningCard({ resource, remaining }: { resource: Resource; remaining: number }) {
  return (
    <section className="relative overflow-hidden rounded-2xl border border-[color:var(--brand)]/20 bg-gradient-to-br from-brand-soft via-card to-card p-5 shadow-[0_16px_40px_rgba(86,54,37,0.08)] sm:p-6" aria-labelledby="continue-learning-title">
      <div className="pointer-events-none absolute -right-12 -top-16 size-48 rounded-full bg-[color:var(--brand)]/8 blur-2xl" />
      <div className="relative flex flex-col gap-5 sm:flex-row sm:items-center sm:justify-between">
        <div className="min-w-0 max-w-2xl">
          <div className="flex flex-wrap items-center gap-3">
            <p className="text-xs font-semibold uppercase tracking-[0.14em] text-[color:var(--brand)]">Продолжить обучение</p>
            {remaining > 1 ? <span className="rounded-full bg-card/80 px-2 py-0.5 text-xs text-muted-foreground">Ещё {remaining - 1}</span> : null}
          </div>
          <h2 id="continue-learning-title" className="mt-3 text-xl font-semibold leading-snug tracking-tight text-balance">{resource.title}</h2>
          <div className="mt-3 flex flex-wrap items-center gap-x-4 gap-y-2">
            <SourceIdentity url={resource.url} />
            <ResourceTags tags={resource.tags} />
          </div>
        </div>
        <Button className="w-full shrink-0 sm:w-auto" nativeButton={false} render={<Link href={`/quiz?id=${resource.id}`} />}>
          Продолжить квиз
          <ArrowRight className="size-4" aria-hidden="true" />
        </Button>
      </div>
    </section>
  )
}

export function ResourceCard({ resource, onRetry, retrying }: { resource: Resource; onRetry: (resource: Resource) => void; retrying?: boolean }) {
  const { id, title, url, tags, status, createdAt, completedAt } = resource

  return (
    <article className="group flex h-full w-full min-w-0 flex-col rounded-2xl border border-border/90 bg-card p-5 shadow-[0_8px_30px_rgba(86,54,37,0.055)] transition-[border-color,box-shadow] hover:border-[color:var(--brand)]/25 hover:shadow-[0_14px_34px_rgba(86,54,37,0.08)]">
      <div className="flex items-center justify-between gap-3">
        <SourceIdentity url={url} />
        <Button size="sm" variant="ghost" nativeButton={false} render={<a href={url} target="_blank" rel="noreferrer noopener" aria-label={`Открыть оригинал: ${title}`} />}>
          Оригинал
          <ExternalLink className="size-3.5" aria-hidden="true" />
        </Button>
      </div>

      <Link href={`/resource?id=${id}`} className="mt-3 text-base font-semibold leading-snug tracking-tight text-balance hover:text-[color:var(--brand)] xl:min-h-16">{title}</Link>

      <div className="mt-3 flex flex-wrap items-center gap-2">
        <StatusBadge status={status} />
        {status === 'PROCESSING' || status === 'NOT_COMPLETED' ? (
          <span className="inline-flex items-center gap-1.5 rounded-md border border-border bg-secondary px-2 py-1 text-xs font-medium text-muted-foreground" title="Этот материал занимает активный слот">
            <Layers className="size-3.5" aria-hidden="true" />
            Активный слот
          </span>
        ) : null}
        {status === 'COMPLETED' ? (
          <span className="tabular inline-flex items-center gap-1.5 rounded-md border border-[color:var(--success)]/20 bg-success-soft px-2 py-1 text-xs font-semibold text-[color:var(--success)]">
            +{resource.xpEarned ?? 0} XP
            <span aria-hidden="true">·</span>
            <Coins className="size-3.5" aria-hidden="true" />+{resource.ePointsEarned ?? 0}
            <span className="sr-only">е-баллов</span>
          </span>
        ) : null}
      </div>

      <p className="tabular mt-2 inline-flex items-center gap-1.5 text-xs text-muted-foreground">
        <Calendar className="size-3.5" aria-hidden="true" />
        {status === 'COMPLETED' && completedAt ? 'Завершено' : 'Добавлено'} {new Date(status === 'COMPLETED' && completedAt ? completedAt : createdAt).toLocaleDateString('ru-RU')}
      </p>

      {status === 'PROCESSING' ? (
        <div className="mt-4 flex items-center gap-3 text-xs text-muted-foreground">
          <div className="h-1.5 flex-1 overflow-hidden rounded-full bg-border">
            <div className="h-full w-1/2 animate-pulse rounded-full bg-[color:var(--brand)]" />
          </div>
          <span className="shrink-0">Несколько минут</span>
        </div>
      ) : null}

      {status === 'FAILED' ? <p className="mt-4 text-sm text-muted-foreground">Обработка не удалась. Можно попробовать ещё раз.</p> : null}

      <div className="mt-auto flex flex-wrap items-end justify-between gap-3 pt-4">
        <div className="min-h-6">{tags.length ? <ResourceTags tags={tags} /> : null}</div>
        {status !== 'PROCESSING' ? (
          <div className="ml-auto">
            {status === 'NOT_COMPLETED' ? <Button size="sm" nativeButton={false} render={<Link href={`/quiz?id=${id}`} />}><Play className="size-3.5" aria-hidden="true" />Начать квиз</Button> : null}
            {status === 'FAILED' ? <Button size="sm" onClick={() => onRetry(resource)} disabled={retrying}><RotateCw className="size-3.5" aria-hidden="true" />{retrying ? 'Повторяем…' : 'Повторить'}</Button> : null}
            {status === 'COMPLETED' ? <Button size="sm" variant="secondary" nativeButton={false} render={<Link href={`/resource?id=${id}`} />}>Посмотреть<ArrowRight className="size-3.5" aria-hidden="true" /></Button> : null}
          </div>
        ) : null}
      </div>
    </article>
  )
}
