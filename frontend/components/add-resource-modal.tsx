'use client'

import { useState } from 'react'
import { AlertCircle, Clock, Info } from 'lucide-react'
import { Modal } from '@/components/modal'
import { Field, TextInput } from '@/components/field'
import { Button } from '@/components/ui/button'
import { api, ApiError } from '@/lib/api'

const SLOT_COST = 25

export function AddResourceModal({
  open,
  onClose,
  onCreated = () => undefined,
  slotsUsed = 0,
  slotsTotal = 0,
  balance = 0,
}: {
  open: boolean
  onClose: () => void
  onCreated?: () => void | Promise<void>
  slotsUsed?: number
  slotsTotal?: number
  balance?: number
}) {
  const [url, setUrl] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [buySlot, setBuySlot] = useState(false)

  const free = slotsTotal - slotsUsed
  const isOverCapacity = slotsUsed > slotsTotal
  const needsPurchase = slotsUsed === slotsTotal
  const canAfford = balance >= SLOT_COST

  function reset() {
    setUrl('')
    setError(null)
    setLoading(false)
    setBuySlot(false)
  }

  function handleClose() {
    if (loading) return
    reset()
    onClose()
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (loading) return

    let parsed: URL | null = null
    try {
      parsed = new URL(url)
    } catch {
      setError('Не похоже на ссылку. Проверьте адрес статьи.')
      return
    }
    if (!/^https?:$/.test(parsed.protocol)) {
      setError('Поддерживаются только ссылки http и https.')
      return
    }
    if (isOverCapacity) {
      setError('Сначала завершите один из материалов: временный слот уже используется.')
      return
    }
    if (needsPurchase && !buySlot) {
      setError('Backlog заполнен. Отметьте покупку слота, чтобы продолжить.')
      return
    }

    setError(null)
    setLoading(true)
    try {
      await api.createResource(url, needsPurchase && buySlot)
      await onCreated()
      reset()
      onClose()
    } catch (caught) {
      setError(caught instanceof ApiError ? caught.message : 'Не удалось добавить материал.')
      setLoading(false)
    }
  }

  const submitLabel = loading
    ? 'Загружаем материал…'
    : isOverCapacity
      ? 'Временный слот занят'
      : needsPurchase
        ? 'Купить слот и добавить'
        : 'Добавить в backlog'

  return (
    <Modal
      open={open}
      onClose={handleClose}
      title="Добавить материал"
      description="Вставьте ссылку на статью, гайд или документацию."
    >
      <form onSubmit={handleSubmit} className="flex flex-col gap-4">
        <Field
          label="Ссылка на материал"
          htmlFor="resource-url"
          error={error ?? undefined}
        >
          <TextInput
            id="resource-url"
            type="url"
            inputMode="url"
            placeholder="https://example.com/article"
            value={url}
            onChange={(e) => {
              setUrl(e.target.value)
              if (error) setError(null)
            }}
            aria-invalid={error ? true : undefined}
            disabled={loading}
            autoFocus
          />
        </Field>

        {/* Capacity / purchase */}
        {free > 0 ? (
          <p className="tabular text-sm text-muted-foreground">
            Свободно {free} из {slotsTotal} слотов
          </p>
        ) : isOverCapacity ? (
          <div className="rounded-lg border border-[color:var(--warning)]/30 bg-warning-soft p-3">
            <p className="flex items-center gap-2 text-sm font-medium text-[color:var(--warning)]">
              <AlertCircle className="size-4" aria-hidden="true" />
              Временный слот уже используется.
            </p>
            <p className="mt-2 text-xs text-muted-foreground">
              Завершите один из материалов, прежде чем добавлять следующий.
            </p>
          </div>
        ) : (
          <div className="rounded-lg border border-[color:var(--warning)]/30 bg-warning-soft p-3">
            <p className="flex items-center gap-2 text-sm font-medium text-[color:var(--warning)]">
              <AlertCircle className="size-4" aria-hidden="true" />
              Все {slotsTotal} слотов заняты.
            </p>
            <label className="mt-3 flex items-start gap-2.5 text-sm">
              <input
                type="checkbox"
                checked={buySlot}
                disabled={!canAfford || loading}
                onChange={(e) => setBuySlot(e.target.checked)}
                className="mt-0.5 size-4 accent-[color:var(--brand)]"
              />
              <span>
                Купить временный слот за{' '}
                <span className="tabular font-medium">{SLOT_COST}</span> е-баллов
              </span>
            </label>
            <p className="tabular mt-2 text-xs text-muted-foreground">
              Баланс: {balance} е-балла
            </p>
            {!canAfford ? (
              <p className="mt-1 text-xs text-[color:var(--destructive)]">
                Недостаточно е-баллов для покупки слота.
              </p>
            ) : null}
          </div>
        )}

        {loading ? (
          <p className="flex items-center gap-2 text-xs text-muted-foreground">
            <Clock className="size-3.5" aria-hidden="true" />
            Загрузка материала может занять до 30 секунд.
          </p>
        ) : (
          <p className="flex items-center gap-2 text-xs text-muted-foreground">
            <Info className="size-3.5" aria-hidden="true" />
            После добавления мы автоматически создадим квиз по материалу.
          </p>
        )}

        <div className="mt-1 flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <Button
            type="button"
            variant="outline"
            onClick={handleClose}
            disabled={loading}
          >
            Отмена
          </Button>
          <Button
            type="submit"
            disabled={loading || isOverCapacity || (needsPurchase && (!canAfford || !buySlot))}
          >
            {submitLabel}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
