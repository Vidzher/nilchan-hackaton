'use client'

import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Field, TextInput } from '@/components/field'
import { api, ApiError, setToken } from '@/lib/api'

export function RegisterForm() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)
  const [errors, setErrors] = useState<{ email?: string; username?: string; password?: string; confirm?: string; form?: string }>({})

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    const email = String(form.get('email') ?? '').trim()
    const username = String(form.get('username') ?? '').trim()
    const password = String(form.get('password') ?? '')
    const confirm = String(form.get('confirm') ?? '')
    const next: typeof errors = {}
    if (!email) next.email = 'Введите email.'
    if (username.length < 3 || username.length > 32) next.username = 'От 3 до 32 символов.'
    if (password.length < 8 || password.length > 72) next.password = 'От 8 до 72 символов.'
    if (confirm !== password) next.confirm = 'Пароли не совпадают.'
    setErrors(next)
    if (Object.keys(next).length) return

    setLoading(true)
    try {
      const result = await api.register(email, username, password)
      setToken(result.token)
      router.replace('/')
    } catch (error) {
      setErrors({ form: error instanceof ApiError ? error.message : 'Не удалось создать аккаунт.' })
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4" noValidate>
      <Field label="Email" htmlFor="email" error={errors.email}>
        <TextInput id="email" name="email" type="email" autoComplete="email" placeholder="you@example.com" aria-invalid={errors.email ? true : undefined} />
      </Field>
      <Field label="Имя пользователя" htmlFor="username" error={errors.username} hint={errors.username ? undefined : 'Так вас будут видеть в рейтинге.'}>
        <TextInput id="username" name="username" autoComplete="username" placeholder="vasya" aria-invalid={errors.username ? true : undefined} />
      </Field>
      <Field label="Пароль" htmlFor="password" error={errors.password}>
        <TextInput id="password" name="password" type="password" autoComplete="new-password" placeholder="••••••••" aria-invalid={errors.password ? true : undefined} />
      </Field>
      <Field label="Повторите пароль" htmlFor="confirm" error={errors.confirm}>
        <TextInput id="confirm" name="confirm" type="password" autoComplete="new-password" placeholder="••••••••" aria-invalid={errors.confirm ? true : undefined} />
      </Field>
      {errors.form ? <p className="text-sm text-[color:var(--destructive)]" role="alert">{errors.form}</p> : null}
      <Button type="submit" className="mt-1 w-full" disabled={loading}>{loading ? 'Создаём аккаунт…' : 'Создать аккаунт'}</Button>
    </form>
  )
}
