'use client'

import { useRouter } from 'next/navigation'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Field, TextInput } from '@/components/field'
import { api, ApiError, setToken } from '@/lib/api'

export function LoginForm() {
  const router = useRouter()
  const [loading, setLoading] = useState(false)
  const [errors, setErrors] = useState<{ email?: string; password?: string; form?: string }>({})

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    const form = new FormData(e.currentTarget)
    const email = String(form.get('email') ?? '').trim()
    const password = String(form.get('password') ?? '')
    const next: typeof errors = {}
    if (!email) next.email = 'Введите email.'
    if (!password) next.password = 'Введите пароль.'
    setErrors(next)
    if (Object.keys(next).length) return

    setLoading(true)
    try {
      const result = await api.login(email, password)
      setToken(result.token)
      router.replace('/')
    } catch (error) {
      setErrors({ form: error instanceof ApiError ? error.message : 'Не удалось войти.' })
    } finally {
      setLoading(false)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-4" noValidate>
      <Field label="Email" htmlFor="email" error={errors.email}>
        <TextInput id="email" name="email" type="email" autoComplete="email" placeholder="you@example.com" aria-invalid={errors.email ? true : undefined} />
      </Field>
      <Field label="Пароль" htmlFor="password" error={errors.password}>
        <TextInput id="password" name="password" type="password" autoComplete="current-password" placeholder="••••••••" aria-invalid={errors.password ? true : undefined} />
      </Field>
      {errors.form ? <p className="text-sm text-[color:var(--destructive)]" role="alert">{errors.form}</p> : null}
      <Button type="submit" className="mt-1 w-full" disabled={loading}>
        {loading ? 'Входим…' : 'Войти'}
      </Button>
    </form>
  )
}
