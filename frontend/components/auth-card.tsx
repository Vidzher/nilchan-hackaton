import Link from 'next/link'
import { BrandMark } from '@/components/brand-assets'

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
    <main className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden bg-background px-4 py-12">
      <div className="pointer-events-none absolute -left-24 top-[12%] size-72 rounded-full bg-brand-soft/70 blur-3xl" />
      <div className="pointer-events-none absolute -right-24 bottom-[8%] size-72 rounded-full bg-success-soft/70 blur-3xl" />
      <div className="relative w-full max-w-[400px]">
        <Link href="/" className="group flex items-center justify-center gap-2.5" aria-label="Learning Backlog">
          <BrandMark className="transition-transform group-hover:-rotate-3" />
          <span className="leading-none">
            <span className="block text-[10px] font-semibold uppercase tracking-[0.18em] text-[color:var(--brand)]">Learning</span>
            <span className="mt-0.5 block text-base font-bold tracking-tight">Backlog</span>
          </span>
        </Link>

        <div className="mt-8 rounded-2xl border border-border bg-card p-6 shadow-[0_20px_60px_rgba(86,54,37,0.08)] sm:p-7">
          <h1 className="text-xl font-semibold tracking-tight">{title}</h1>
          <p className="mt-1 text-sm text-muted-foreground text-pretty">{subtitle}</p>
          <div className="mt-6">{children}</div>
        </div>

        <p className="mt-5 text-center text-sm text-muted-foreground">{footer}</p>
      </div>
    </main>
  )
}
