import Link from 'next/link'
import { ExternalLink, RotateCw, Play, Calendar } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { StatusBadge } from '@/components/status-badge'
import { Tag } from '@/components/primitives'
import type { Resource } from '@/lib/api'

function domain(url: string) {
  try { return new URL(url).hostname } catch { return url }
}

export function ResourceCard({ resource, onRetry, retrying }: { resource: Resource; onRetry: (resource: Resource) => void; retrying?: boolean }) {
  const { id, title, url, tags, status, createdAt } = resource

  return (
    <article className="flex flex-col rounded-xl border border-border bg-card p-4">
      <div className="flex items-start justify-between gap-3">
        <StatusBadge status={status} />
        <span className="tabular inline-flex items-center gap-1 whitespace-nowrap text-xs text-muted-foreground">
          <Calendar className="size-3.5" aria-hidden="true" />
          {new Date(createdAt).toLocaleDateString('ru-RU')}
        </span>
      </div>
      <Link href={`/resource?id=${id}`} className="mt-3 text-[15px] font-semibold leading-snug text-balance hover:underline">{title}</Link>
      <p className="mt-1 text-sm text-muted-foreground">{domain(url)}</p>
      <div className="mt-3 flex flex-wrap gap-1.5">{tags.map((tag) => <Tag key={tag}>{tag}</Tag>)}</div>
      {status === 'FAILED' ? <p className="mt-3 text-sm text-muted-foreground">Попробуйте запустить обработку ещё раз.</p> : null}
      {status === 'COMPLETED' ? (
        <div className="mt-3 flex flex-wrap gap-2 text-sm">
          <span className="tabular font-medium text-[color:var(--success)]">+{resource.xpEarned ?? 0} XP</span>
          <span className="tabular font-medium text-[color:var(--success)]">+{resource.ePointsEarned ?? 0} е-баллов</span>
        </div>
      ) : null}
      <div className="mt-4 flex flex-wrap items-center gap-2 border-t border-border pt-4">
        {status === 'NOT_COMPLETED' ? <Button size="sm" nativeButton={false} render={<Link href={`/quiz?id=${id}`} />}><Play className="size-3.5" aria-hidden="true" />Начать quiz</Button> : null}
        {status === 'PROCESSING' ? <Button size="sm" disabled>Создаём quiz…</Button> : null}
        {status === 'FAILED' ? <Button size="sm" variant="secondary" onClick={() => onRetry(resource)} disabled={retrying}><RotateCw className="size-3.5" aria-hidden="true" />{retrying ? 'Повторяем…' : 'Повторить'}</Button> : null}
        <Button size="sm" variant="outline" nativeButton={false} render={<a href={url} target="_blank" rel="noreferrer noopener" />}><ExternalLink className="size-3.5" aria-hidden="true" />Открыть оригинал</Button>
      </div>
    </article>
  )
}
