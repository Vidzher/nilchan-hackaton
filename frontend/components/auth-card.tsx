import Link from 'next/link'
import { Layers } from 'lucide-react'

export function AuthCard({
  title,
  subtitle,
  children,
  footer,
}: {
  title: string
  subtitle: string
  children: React.ReactNode
  footer: React.ReactNode
}) {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-background px-4 py-12">
      <div className="w-full max-w-[400px]">
        <Link href="/" className="flex items-center justify-center gap-2">
          <span className="grid size-8 place-items-center rounded-md bg-primary text-primary-foreground">
            <Layers className="size-4" aria-hidden="true" />
          </span>
          <span className="text-base font-semibold tracking-tight">
            Learning Backlog
          </span>
        </Link>

        <div className="mt-8 rounded-2xl border border-border bg-card p-6 sm:p-7">
          <h1 className="text-xl font-semibold tracking-tight">{title}</h1>
          <p className="mt-1 text-sm text-muted-foreground text-pretty">{subtitle}</p>
          <div className="mt-6">{children}</div>
        </div>

        <p className="mt-5 text-center text-sm text-muted-foreground">{footer}</p>
      </div>
    </main>
  )
}
