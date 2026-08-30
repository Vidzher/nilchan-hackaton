import Link from 'next/link'
import { AuthCard } from '@/components/auth-card'
import { LoginForm } from '@/components/login-form'

export default function LoginPage() {
  return (
    <AuthCard
      title="С возвращением"
      subtitle="Войдите, чтобы продолжить разгребать backlog."
      footer={
        <>
          Ещё нет аккаунта?{' '}
          <Link href="/register" className="font-medium text-foreground hover:underline">
            Зарегистрироваться
          </Link>
        </>
      }
    >
      <LoginForm />
    </AuthCard>
  )
}
