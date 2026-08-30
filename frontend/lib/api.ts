'use client'

export type ResourceStatus =
  | 'PROCESSING'
  | 'FAILED'
  | 'NOT_COMPLETED'
  | 'COMPLETED'

export type AuthResponse = {
  userId: number
  email: string
  username: string
  token: string
}

export type Resource = {
  id: number
  url: string
  title: string
  tags: string[]
  status: ResourceStatus
  createdAt: string
  completedAt?: string
  xpEarned?: number
  ePointsEarned?: number
}

export type QuizQuestion = {
  text: string
  options: [string, string, string, string]
  explanation: string
  evidence: string
  verificationSalt: string
  correctAnswerHash: string
}

export type Quiz = {
  id: number
  resourceId: number
  title: string
  questions: QuizQuestion[]
}

export type SubmittedAnswer = {
  questionIndex: number
  selectedIndex: number
}

export type CompletionResponse = {
  completion: {
    completedAt: string
    totalQuestions: number
    xpEarned: number
    ePointsEarned: number
  }
  progress: {
    xp: number
    ePoints: number
    currentStreak: number
    level: number
    activeBacklogLimit: number
  }
}

export type Profile = {
  user: { id: number; email: string; username: string }
  progress: {
    xp: number
    level: number
    activeBacklogLimit: number
    ePoints: number
    currentStreak: number
    lastCompletionAt?: string
  }
  cosmetics: {
    avatarId: string
    frameId: string
    titleId?: string
    showcaseItemId?: string
    ownedCosmeticIds: string[]
  }
}

export type CosmeticType = 'avatar' | 'frame' | 'title' | 'showcase'

export type ShopItem = {
  id: string
  type: CosmeticType
  displayName: string
  price: number
  presentation: {
    emoji?: string
    assetKey?: string
    cssClass?: string
  }
}

const defaultCosmetics: ShopItem[] = [
  {
    id: 'default_avatar',
    type: 'avatar',
    displayName: 'Нормис',
    price: 0,
    presentation: { emoji: '🙂' },
  },
  {
    id: 'default_frame',
    type: 'frame',
    displayName: 'Без рамки',
    price: 0,
    presentation: { assetKey: 'frame-default', cssClass: 'frame-default' },
  },
]

export type PurchaseCosmeticResponse = {
  item: ShopItem
  ePoints: number
}

export type CosmeticsUpdate = {
  avatarId?: string
  frameId?: string
  titleId?: string | null
  showcaseItemId?: string | null
}

export type LeaderboardEntry = {
  rank: number
  userId: number
  username: string
  xp: number
  level: number
  avatarId: string
  frameId: string
  titleId?: string
  showcaseItemId?: string
  isCurrent: boolean
}

type SuccessEnvelope<T> = { status: 'OK'; data: T }
type ErrorEnvelope = { status: 'Error'; error?: string }

const TOKEN_KEY = 'learning-backlog.jwt'
const API_BASE = (process.env.NEXT_PUBLIC_API_BASE_URL ?? '').replace(/\/$/, '')

export class ApiError extends Error {
  constructor(
    message: string,
    readonly status: number,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

export function getToken() {
  return typeof window === 'undefined' ? null : window.localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string) {
  window.localStorage.setItem(TOKEN_KEY, token)
  window.dispatchEvent(new Event('authchange'))
}

export function clearToken() {
  window.localStorage.removeItem(TOKEN_KEY)
  window.dispatchEvent(new Event('authchange'))
}

async function request<T>(path: string, init: RequestInit = {}, authenticated = true): Promise<T> {
  const headers = new Headers(init.headers)
  if (init.body) headers.set('Content-Type', 'application/json')
  if (authenticated) {
    const token = getToken()
    if (token) headers.set('Authorization', `Bearer ${token}`)
  }

  let response: Response
  try {
    response = await fetch(`${API_BASE}${path}`, { ...init, headers })
  } catch {
    throw new ApiError('Не удалось связаться с сервером.', 0)
  }

  const contentType = response.headers.get('content-type') ?? ''
  const body: unknown = contentType.includes('application/json')
    ? await response.json().catch(() => null)
    : await response.text().catch(() => '')

  if (!response.ok) {
    if (response.status === 401 && authenticated) clearToken()
    const message =
      typeof body === 'object' && body !== null && 'error' in body && typeof (body as ErrorEnvelope).error === 'string'
        ? (body as ErrorEnvelope).error!
        : typeof body === 'string' && body.trim()
          ? body.trim()
          : `Ошибка сервера (${response.status}).`
    throw new ApiError(message, response.status)
  }

  if (typeof body !== 'object' || body === null || !('data' in body)) {
    throw new ApiError('Сервер вернул некорректный ответ.', response.status)
  }
  return (body as SuccessEnvelope<T>).data
}

export const api = {
  login: (email: string, password: string) =>
    request<AuthResponse>('/api/login', { method: 'POST', body: JSON.stringify({ email, password }) }, false),
  register: (email: string, username: string, password: string) =>
    request<AuthResponse>(
      '/api/register',
      { method: 'POST', body: JSON.stringify({ email, username, password }) },
      false,
    ),
  profile: () => request<Profile>('/api/profile'),
  updateCosmetics: (update: CosmeticsUpdate) =>
    request<Profile>('/api/profile/cosmetics', {
      method: 'PATCH',
      body: JSON.stringify(update),
    }),
  shop: async () => [...defaultCosmetics, ...(await request<ShopItem[]>('/api/shop'))],
  purchaseCosmetic: (itemId: string) =>
    request<PurchaseCosmeticResponse>('/api/shop/purchase', {
      method: 'POST',
      body: JSON.stringify({ itemId }),
    }),
  leaderboard: () => request<LeaderboardEntry[]>('/api/leaderboard'),
  resources: () => request<Resource[]>('/api/resources'),
  resource: (id: number) => request<Resource>(`/api/resources/${id}`),
  createResource: async (url: string, purchaseOverflowSlot: boolean) => {
    const resource = await request<Resource>('/api/resources', {
      method: 'POST',
      body: JSON.stringify({ url, purchaseOverflowSlot }),
    })
    window.dispatchEvent(new Event('resourceschange'))
    return resource
  },
  quiz: (resourceId: number) => request<Quiz>(`/api/resources/${resourceId}/quiz`),
  completeQuiz: async (resourceId: number, answers: SubmittedAnswer[]) => {
    const completion = await request<CompletionResponse>(`/api/resources/${resourceId}/quiz/complete`, {
      method: 'POST',
      body: JSON.stringify({ answers }),
    })
    window.dispatchEvent(new Event('resourceschange'))
    return completion
  },
}
