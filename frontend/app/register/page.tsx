import Link from 'next/link'
import { AuthCard } from '@/components/auth-card'
import { RegisterForm } from '@/components/register-form'

export default function RegisterPage() {
  return (
    <AuthCard
      title="Создать аккаунт"
      subtitle="Начните дочитывать то, что сохранили на потом."
      footer={
        <>
          Уже есть аккаунт?{' '}
          <Link href="/login" className="font-medium text-foreground hover:underline">
            Войти
          </Link>
        </>
      }
    >
      <RegisterForm />
    </AuthCard>
  )
}
