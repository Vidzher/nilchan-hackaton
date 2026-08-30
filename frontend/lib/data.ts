export type ResourceStatus =
  | 'NOT_COMPLETED'
  | 'PROCESSING'
  | 'FAILED'
  | 'COMPLETED'

export type Resource = {
  id: string
  title: string
  domain: string
  url: string
  tags: string[]
  status: ResourceStatus
  createdAt: string
  questions?: number
  completedAt?: string
  rewardXp?: number
  rewardPoints?: number
}

export const user = {
  username: 'vasya',
  avatar: '🐸',
  title: 'Backlog Destroyer',
  showcase: '🦆 Senior Rubber Duck',
  frame: 'Neon',
  level: 4,
  xp: 820,
  nextLevel: 1000,
  points: 74,
  streak: 9,
  slotsUsed: 6,
  slotsTotal: 8,
}

export const resources: Resource[] = [
  {
    id: 'go-concurrency',
    title: 'Go Concurrency Patterns: Pipelines and cancellation',
    domain: 'go.dev',
    url: 'https://go.dev/blog/pipelines',
    tags: ['go', 'concurrency', 'goroutines'],
    status: 'NOT_COMPLETED',
    createdAt: '12 авг',
    questions: 8,
  },
  {
    id: 'css-grid',
    title: 'Understanding CSS Grid',
    domain: 'web.dev',
    url: 'https://web.dev/learn/css/grid',
    tags: ['css', 'layout'],
    status: 'PROCESSING',
    createdAt: '14 авг',
  },
  {
    id: 'sqlite-tx',
    title: 'SQLite Transaction Behavior',
    domain: 'sqlite.org',
    url: 'https://www.sqlite.org/lang_transaction.html',
    tags: ['database', 'sqlite'],
    status: 'FAILED',
    createdAt: '11 авг',
  },
  {
    id: 'rust-ownership',
    title: 'Ownership and Borrowing in Rust',
    domain: 'doc.rust-lang.org',
    url: 'https://doc.rust-lang.org/book/ch04-00-understanding-ownership.html',
    tags: ['rust', 'memory'],
    status: 'NOT_COMPLETED',
    createdAt: '10 авг',
    questions: 6,
  },
  {
    id: 'http-caching',
    title: 'HTTP Caching Explained',
    domain: 'developer.mozilla.org',
    url: 'https://developer.mozilla.org/en-US/docs/Web/HTTP/Caching',
    tags: ['http', 'web', 'performance'],
    status: 'COMPLETED',
    createdAt: '5 авг',
    completedAt: '8 авг',
    questions: 7,
    rewardXp: 60,
    rewardPoints: 8,
  },
  {
    id: 'postgres-indexes',
    title: 'How Postgres Chooses an Index',
    domain: 'postgresql.org',
    url: 'https://www.postgresql.org/docs/current/indexes.html',
    tags: ['database', 'postgres'],
    status: 'COMPLETED',
    createdAt: '2 авг',
    completedAt: '4 авг',
    questions: 9,
    rewardXp: 75,
    rewardPoints: 10,
  },
]

export const statusMeta: Record<
  ResourceStatus,
  { label: string; tone: 'brand' | 'warning' | 'error' | 'success' }
> = {
  NOT_COMPLETED: { label: 'Quiz готов', tone: 'brand' },
  PROCESSING: { label: 'Создаём quiz', tone: 'warning' },
  FAILED: { label: 'Не удалось создать quiz', tone: 'error' },
  COMPLETED: { label: 'Завершено', tone: 'success' },
}

export type QuizQuestion = {
  prompt: string
  options: string[]
  correct: number
  explanation: string
  evidence: string
}

export const quizQuestions: QuizQuestion[] = [
  {
    prompt:
      'Что произойдёт, если получатель канала перестанет читать значения в пайплайне без отмены?',
    options: [
      'Горутины-отправители завершатся автоматически',
      'Отправляющие горутины заблокируются и произойдёт утечка',
      'Канал закроется сам через таймаут',
      'Данные будут отброшены без блокировки',
    ],
    correct: 1,
    explanation:
      'Без явного сигнала отмены отправители остаются заблокированными на записи в канал, что приводит к утечке горутин.',
    evidence:
      '«If a receiver stops reading, the senders block forever — we need a way to tell them to stop, even when the receiver hasn’t consumed all values.»',
  },
  {
    prompt: 'Для чего используется канал done в паттерне отмены?',
    options: [
      'Для передачи результатов вычислений',
      'Для буферизации входных данных',
      'Для широковещательного сигнала об остановке через закрытие',
      'Для подсчёта количества активных горутин',
    ],
    correct: 2,
    explanation:
      'Закрытие канала done служит широковещательным сигналом: все горутины, читающие из него, одновременно разблокируются.',
    evidence:
      '«Closing a channel that serves as a done signal broadcasts to all goroutines that they should abandon their work and return.»',
  },
  {
    prompt: 'Почему select с каналом done предпочтительнее прямой записи?',
    options: [
      'Он быстрее выполняется',
      'Он позволяет горутине выйти вместо вечной блокировки',
      'Он гарантирует порядок сообщений',
      'Он автоматически закрывает выходной канал',
    ],
    correct: 1,
    explanation:
      'select даёт горутине альтернативную ветку: при закрытии done она может завершиться, а не блокироваться на отправке.',
    evidence:
      '«Using select with a done case lets a goroutine give up on a send rather than blocking indefinitely.»',
  },
]

export type ShopItem = {
  id: string
  name: string
  type: 'Аватар' | 'Рамка' | 'Титул' | 'Витрина'
  price: number
  preview: string
  state: 'buy' | 'owned' | 'equip' | 'equipped'
}

export const shopItems: ShopItem[] = [
  { id: 'a-frog', name: 'Frog', type: 'Аватар', price: 15, preview: '🐸', state: 'equipped' },
  { id: 'a-robot', name: 'Robot', type: 'Аватар', price: 25, preview: '🤖', state: 'buy' },
  { id: 'a-wizard', name: 'Wizard', type: 'Аватар', price: 40, preview: '🧙', state: 'buy' },
  { id: 'f-neon', name: 'Neon', type: 'Рамка', price: 30, preview: 'Neon', state: 'equipped' },
  { id: 'f-fire', name: 'Fire', type: 'Рамка', price: 50, preview: 'Fire', state: 'equip' },
  { id: 'f-gold', name: 'Gold', type: 'Рамка', price: 120, preview: 'Gold', state: 'buy' },
  { id: 't-razgreb', name: 'Разгребатель', type: 'Титул', price: 20, preview: 'Разгребатель', state: 'equip' },
  { id: 't-goblin', name: 'Knowledge Goblin', type: 'Титул', price: 40, preview: 'Knowledge Goblin', state: 'buy' },
  { id: 't-destroyer', name: 'Backlog Destroyer', type: 'Титул', price: 70, preview: 'Backlog Destroyer', state: 'equipped' },
  { id: 's-duck', name: 'Senior Rubber Duck', type: 'Витрина', price: 40, preview: '🦆', state: 'equipped' },
  { id: 's-cactus', name: 'Кактус прокрастинации', type: 'Витрина', price: 55, preview: '🌵', state: 'buy' },
  { id: 's-cat', name: 'Кот', type: 'Витрина', price: 120, preview: '🐈', state: 'buy' },
  { id: 's-golden-duck', name: 'Golden Duck', type: 'Витрина', price: 250, preview: '👑', state: 'buy' },
]

export type LeaderRow = {
  rank: number
  username: string
  avatar: string
  frame: string
  level: number
  xp: number
  isCurrent?: boolean
}

export const leaderboard: LeaderRow[] = [
  { rank: 1, username: 'vasya', avatar: '🐸', frame: 'Neon', level: 5, xp: 3420, isCurrent: true },
  { rank: 2, username: 'lena', avatar: '🦊', frame: 'Fire', level: 5, xp: 3180 },
  { rank: 3, username: 'max', avatar: '🤖', frame: 'Gold', level: 5, xp: 2940 },
  { rank: 4, username: 'nina', avatar: '🧙', frame: 'Neon', level: 5, xp: 2710 },
  { rank: 5, username: 'oleg', avatar: '🦉', frame: 'Fire', level: 4, xp: 2480 },
  { rank: 6, username: 'kira', avatar: '🐱', frame: 'Gold', level: 4, xp: 2255 },
  { rank: 7, username: 'petya', avatar: '🐸', frame: 'Neon', level: 4, xp: 2090 },
  { rank: 8, username: 'sonya', avatar: '🦆', frame: 'Fire', level: 4, xp: 1970 },
  { rank: 9, username: 'dima', avatar: '🐢', frame: 'Neon', level: 4, xp: 1840 },
  { rank: 10, username: 'anya', avatar: '🐙', frame: 'Gold', level: 4, xp: 1725 },
  { rank: 11, username: 'roma', avatar: '🐝', frame: 'Fire', level: 3, xp: 1610 },
  { rank: 12, username: 'yulia', avatar: '🦋', frame: 'Neon', level: 3, xp: 1495 },
  { rank: 13, username: 'stas', avatar: '🐧', frame: 'Gold', level: 3, xp: 1380 },
  { rank: 14, username: 'vera', avatar: '🐬', frame: 'Fire', level: 3, xp: 1265 },
  { rank: 15, username: 'gleb', avatar: '🦖', frame: 'Neon', level: 3, xp: 1150 },
  { rank: 16, username: 'mila', avatar: '🐰', frame: 'Gold', level: 2, xp: 1035 },
  { rank: 17, username: 'egor', avatar: '🦔', frame: 'Fire', level: 2, xp: 920 },
  { rank: 18, username: 'zoya', avatar: '🐳', frame: 'Neon', level: 2, xp: 805 },
  { rank: 19, username: 'artem', avatar: '🦫', frame: 'Gold', level: 2, xp: 690 },
  { rank: 20, username: 'olya', avatar: '🐌', frame: 'Fire', level: 2, xp: 575 },
]
