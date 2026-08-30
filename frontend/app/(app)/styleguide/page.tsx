'use client'

import { useState } from 'react'
import {
  AlertTriangle,
  CheckCircle2,
  Inbox,
  RotateCw,
} from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Field, TextInput } from '@/components/field'
import { StatusBadge } from '@/components/status-badge'
import { PageHeader, Tag, ProgressBar } from '@/components/primitives'
import { AddResourceModal } from '@/components/add-resource-modal'
import { cn } from '@/lib/utils'
import type { ResourceStatus } from '@/lib/data'

function Section({
  title,
  description,
  children,
}: {
  title: string
  description?: string
  children: React.ReactNode
}) {
  return (
    <section className="rounded-2xl border border-border bg-card p-6">
      <h2 className="text-base font-semibold">{title}</h2>
      {description ? (
        <p className="mt-1 text-sm text-muted-foreground">{description}</p>
      ) : null}
      <div className="mt-5">{children}</div>
    </section>
  )
}

const statuses: ResourceStatus[] = [
  'NOT_COMPLETED',
  'PROCESSING',
  'FAILED',
  'COMPLETED',
]

export default function StyleguidePage() {
  const [modal, setModal] = useState<null | 'available' | 'full'>(null)
  const [toast, setToast] = useState(false)
  const [radio, setRadio] = useState(0)
  const [loading, setLoading] = useState(false)

  function showToast() {
    setToast(true)
    setTimeout(() => setToast(false), 3000)
  }

  return (
    <div className="flex flex-col gap-8">
      <PageHeader
        title="Компоненты"
        description="Состояния интерфейса и переиспользуемые элементы."
      />

      <Section title="Кнопки" description="Варианты, размеры и состояния.">
        <div className="flex flex-col gap-4">
          <div className="flex flex-wrap items-center gap-2">
            <Button>Основная</Button>
            <Button variant="secondary">Вторичная</Button>
            <Button variant="outline">Контурная</Button>
            <Button variant="ghost">Призрак</Button>
            <Button variant="destructive">Опасная</Button>
            <Button variant="link">Ссылка</Button>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <Button size="sm">Малая</Button>
            <Button>Обычная</Button>
            <Button size="lg">Большая</Button>
            <Button disabled>Отключена</Button>
            <Button
              disabled={loading}
              onClick={() => {
                setLoading(true)
                setTimeout(() => setLoading(false), 1500)
              }}
            >
              {loading ? 'Загрузка…' : 'С загрузкой'}
            </Button>
          </div>
        </div>
      </Section>

      <Section title="Поля ввода" description="Обычное, с подсказкой, с ошибкой и отключённое.">
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label="Обычное" htmlFor="sg-1">
            <TextInput id="sg-1" placeholder="Введите значение" />
          </Field>
          <Field label="С подсказкой" htmlFor="sg-2" hint="Так вас увидят в рейтинге.">
            <TextInput id="sg-2" placeholder="vasya" />
          </Field>
          <Field label="С ошибкой" htmlFor="sg-3" error="Не похоже на ссылку.">
            <TextInput id="sg-3" defaultValue="htp://oops" aria-invalid />
          </Field>
          <Field label="Отключённое" htmlFor="sg-4">
            <TextInput id="sg-4" placeholder="Недоступно" disabled />
          </Field>
        </div>
      </Section>

      <Section title="Статусы материалов" description="Статус всегда подкреплён иконкой и текстом.">
        <div className="flex flex-wrap gap-2">
          {statuses.map((s) => (
            <StatusBadge key={s} status={s} />
          ))}
        </div>
      </Section>

      <Section title="Теги">
        <div className="flex flex-wrap gap-1.5">
          {['go', 'concurrency', 'database', 'css', 'rust'].map((t) => (
            <Tag key={t}>{t}</Tag>
          ))}
        </div>
      </Section>

      <Section title="Варианты ответа (radio)" description="Используются в квизе.">
        <fieldset className="flex flex-col gap-2.5">
          <legend className="sr-only">Пример вопроса</legend>
          {['Первый вариант', 'Второй вариант', 'Третий вариант'].map((opt, i) => {
            const active = radio === i
            return (
              <label
                key={i}
                className={cn(
                  'flex cursor-pointer items-center gap-3 rounded-xl border p-3.5 text-sm transition-colors',
                  active
                    ? 'border-[color:var(--brand)] bg-brand-soft'
                    : 'border-border bg-card hover:border-muted-foreground/40',
                )}
              >
                <input
                  type="radio"
                  name="sg-radio"
                  className="sr-only"
                  checked={active}
                  onChange={() => setRadio(i)}
                />
                <span
                  className={cn(
                    'grid size-5 place-items-center rounded-full border',
                    active
                      ? 'border-[color:var(--brand)] bg-[color:var(--brand)]'
                      : 'border-muted-foreground/40',
                  )}
                  aria-hidden="true"
                >
                  {active ? <span className="size-2 rounded-full bg-white" /> : null}
                </span>
                {opt}
              </label>
            )
          })}
        </fieldset>
      </Section>

      <Section title="Прогресс">
        <div className="flex flex-col gap-4">
          <ProgressBar value={25} max={100} label="25 процентов" />
          <ProgressBar value={60} max={100} label="60 процентов" />
          <ProgressBar value={100} max={100} label="100 процентов" />
        </div>
      </Section>

      <Section title="Модальные окна и уведомления">
        <div className="flex flex-wrap gap-2">
          <Button onClick={() => setModal('available')}>Модалка: есть слоты</Button>
          <Button variant="secondary" onClick={() => setModal('full')}>
            Модалка: backlog полон
          </Button>
          <Button variant="outline" onClick={showToast}>
            Показать toast
          </Button>
        </div>
      </Section>

      <Section title="Пустое состояние">
        <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-border bg-background px-6 py-12 text-center">
          <Inbox className="size-8 text-muted-foreground" aria-hidden="true" />
          <p className="mt-3 text-sm font-semibold">Пока ничего нет</p>
          <p className="mt-1 max-w-xs text-sm text-muted-foreground">
            Добавьте первый материал, чтобы начать.
          </p>
        </div>
      </Section>

      <Section title="Скелетон загрузки">
        <div className="grid gap-4 sm:grid-cols-2">
          {[0, 1].map((i) => (
            <div key={i} className="rounded-xl border border-border bg-card p-4">
              <div className="h-6 w-28 animate-pulse rounded-md bg-muted" />
              <div className="mt-3 h-5 w-full animate-pulse rounded-md bg-muted" />
              <div className="mt-2 h-5 w-2/3 animate-pulse rounded-md bg-muted" />
              <div className="mt-4 flex gap-2">
                <div className="h-6 w-14 animate-pulse rounded-md bg-muted" />
                <div className="h-6 w-14 animate-pulse rounded-md bg-muted" />
              </div>
            </div>
          ))}
        </div>
      </Section>

      <Section title="Состояние ошибки">
        <div className="flex flex-col items-center justify-center rounded-xl border border-[color:var(--destructive)]/25 bg-[#f6e3e3] px-6 py-12 text-center">
          <AlertTriangle
            className="size-8 text-[color:var(--destructive)]"
            aria-hidden="true"
          />
          <p className="mt-3 text-sm font-semibold text-[color:var(--destructive)]">
            Не удалось загрузить данные
          </p>
          <p className="mt-1 max-w-xs text-sm text-muted-foreground">
            Проверьте подключение и попробуйте снова.
          </p>
          <Button variant="outline" className="mt-4">
            <RotateCw className="size-4" aria-hidden="true" />
            Повторить
          </Button>
        </div>
      </Section>

      <AddResourceModal
        open={modal === 'available'}
        onClose={() => setModal(null)}
      />
      <AddResourceModal
        open={modal === 'full'}
        onClose={() => setModal(null)}
        slotsUsed={8}
        slotsTotal={8}
        balance={74}
      />

      {/* Toast */}
      {toast ? (
        <div
          role="status"
          className="fixed bottom-24 left-1/2 z-50 flex -translate-x-1/2 items-center gap-2 rounded-xl border border-border bg-card px-4 py-3 text-sm shadow-lg lg:bottom-8"
        >
          <CheckCircle2 className="size-4 text-[color:var(--success)]" aria-hidden="true" />
          Материал добавлен в backlog
        </div>
      ) : null}
    </div>
  )
}
