'use client'

import { useRouter } from 'next/navigation'
import { useEffect, useState } from 'react'
import { getToken } from '@/lib/api'

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const router = useRouter()
  const [allowed, setAllowed] = useState(false)

  useEffect(() => {
    const check = () => {
      if (!getToken()) {
        setAllowed(false)
        router.replace('/login')
      } else {
        setAllowed(true)
      }
    }
    check()
    window.addEventListener('authchange', check)
    return () => window.removeEventListener('authchange', check)
  }, [router])

  return allowed ? children : <div className="min-h-screen bg-background" />
}
