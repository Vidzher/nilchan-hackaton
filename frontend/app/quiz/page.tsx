'use client'

import Link from 'next/link'
import { useSearchParams } from 'next/navigation'
import { Suspense, useEffect, useState } from 'react'
import { ArrowLeft } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { ProtectedRoute } from '@/components/protected-route'
import { QuizResult, QuizRunner } from '@/components/quiz-runner'
import { api, ApiError, type CompletionResponse, type Quiz, type Resource } from '@/lib/api'

const previewResource: Resource = {
  id: 123,
  url: 'https://example.com/article',
  title: 'Как вернуть билет на самолёт в 2026 году: инструкция по возврату денег',
  tags: ['путешествия', 'авиабилеты'],
  status: 'COMPLETED',
  createdAt: new Date().toISOString(),
  completedAt: new Date().toISOString(),
  xpEarned: 55,
  ePointsEarned: 7,
}

const previewCompletion: CompletionResponse = {
  completion: {
    completedAt: new Date().toISOString(),
    totalQuestions: 7,
    xpEarned: 55,
    ePointsEarned: 7,
  },
  progress: {
    xp: 355,
    ePoints: 42,
    currentStreak: 4,
    level: 3,
    activeBacklogLimit: 7,
  },
}

function QuizContent() {
  const rawId = useSearchParams().get('id')
  const id = rawId && /^\d+$/.test(rawId) ? Number(rawId) : null
  const [resource, setResource] = useState<Resource | null>(null)
  const [quiz, setQuiz] = useState<Quiz | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (id === null) { setError('Quiz не найден.'); return }
    Promise.all([api.resource(id), api.quiz(id)])
      .then(([foundResource, foundQuiz]) => {
        if (foundResource.status !== 'NOT_COMPLETED') throw new ApiError('Этот quiz недоступен или уже завершён.', 409)
        setResource(foundResource)
        setQuiz(foundQuiz)
      })
      .catch((caught) => setError(caught instanceof ApiError ? caught.message : 'Не удалось загрузить quiz.'))
  }, [id])

  if (resource && quiz) return <QuizRunner resource={resource} quiz={quiz} />
  return <main className="grid min-h-screen place-items-center bg-background px-4"><div className="w-full max-w-md rounded-2xl border border-border bg-card p-6 text-center"><h1 className="text-xl font-semibold">{error ?? 'Загружаем quiz…'}</h1>{error ? <Button className="mt-5" nativeButton={false} render={<Link href="/" />}><ArrowLeft className="size-4" aria-hidden="true" />Вернуться в backlog</Button> : null}</div></main>
}

function QuizRoute() {
  const preview = useSearchParams().get('preview')
  if (process.env.NODE_ENV === 'development' && preview === 'result') {
    return <QuizResult resource={previewResource} completion={previewCompletion} />
  }
  return <ProtectedRoute><QuizContent /></ProtectedRoute>
}

export default function QuizPage() {
  return <Suspense fallback={<main className="min-h-screen bg-background" />}><QuizRoute /></Suspense>
}
