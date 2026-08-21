# Learning Backlog — MVP PRD

## 0. Статус документа

Этот файл — единственный source of truth для MVP product behavior, backend invariants и API contract. Таблицы, constants и явно описанные rules являются нормативными; примеры UI и payloads — иллюстративными и не могут переопределять rules. Любое изменение поведения сначала фиксируется здесь.

---

## 1. Идея

**Learning Backlog** — приложение для разбора сохранённых «на потом» образовательных материалов.

Пользователь добавляет URL статьи, гайда или документации. HTTP request ждёт только **Firecrawl**. После получения content backend сохраняет Resource с title и tags, возвращает его frontend и генерирует quiz через **LLM** в фоне. Пользователь сразу видит материал и может открыть оригинал; quiz становится доступен после завершения генерации.

Quiz проходится на frontend без запроса к backend после каждого ответа. Вопрос считается закрытым только после правильного ответа. После правильного ответа на все вопросы frontend одним запросом отправляет итоговые ответы на backend. Backend проверяет их, начисляет **XP** и **е-баллы**, помечает материал выполненным и освобождает backlog slot.

XP отвечает за постоянный прогресс и leaderboard. Е-баллы тратятся на оформление профиля и `Overflow Pass`.

> Не просто сохранять полезный контент, а реально его завершать.

---

## 2. Проблема и гипотеза

Пользователи регулярно сохраняют полезные материалы, но редко к ним возвращаются. Backlog растёт, а сохранение ссылки создаёт ложное ощущение прогресса.

Гипотеза:

> Ограниченный Active Backlog, короткие quiz и простая игровая прогрессия будут мотивировать пользователя завершать сохранённые материалы и возвращаться в приложение.

---

## 3. Core Loop и статусы

```text
Добавить URL
↓
Backend синхронно получает content через Firecrawl
↓
Backend сохраняет Resource = PROCESSING с title и tags
↓
Frontend получает Resource и может открыть оригинал
↓
Backend в goroutine создаёт quiz через LLM
↓
Quiz и Questions сохраняются атомарно
↓
Resource = NOT_COMPLETED, кнопка quiz становится доступна
↓
Пользователь правильно отвечает на все вопросы
↓
Frontend отправляет один completion request
↓
Backend проверяет ответы
↓
Начисляются XP и е-баллы
↓
Resource = COMPLETED
↓
Backlog slot освобождается
↓
XP → Level и Leaderboard
е-баллы → Shop
```

У Resource четыре статуса:

```ts
enum ResourceStatus {
  PROCESSING = "PROCESSING",
  FAILED = "FAILED",
  NOT_COMPLETED = "NOT_COMPLETED",
  COMPLETED = "COMPLETED",
}
```

Материал считается выполненным только если **все вопросы отвечены правильно** и backend подтвердил ответы.

```text
correctlyAnsweredQuestionsCount == totalQuestionsCount
→ Resource = COMPLETED
→ rewards awarded
→ backlog slot freed
```

`PROCESSING` означает, что content, title и tags уже сохранены, но quiz ещё генерируется. `FAILED` означает окончательную ошибку или прерывание генерации quiz. Quiz доступен frontend только когда Resource перешёл в `NOT_COMPLETED`, поэтому frontend никогда не видит частично готовый quiz.

Разрешены только переходы:

```text
PROCESSING → NOT_COMPLETED
PROCESSING → FAILED
FAILED → PROCESSING
NOT_COMPLETED → COMPLETED
```

---

## 4. Active Backlog и Levels

Slot занимает Resource со статусом `NOT_COMPLETED` или `PROCESSING`, чтобы параллельные добавления не могли превысить лимит.

```text
usedCapacity = notCompletedResources + processingResources
```

Если `usedCapacity` достиг Active Backlog Limit, новый Resource добавить нельзя, кроме случая использования `Overflow Pass`. Resource со статусом `FAILED` или `COMPLETED` slot не занимает.

| Level | Required XP | Active Backlog Limit |
| ----: | ----------: | -------------------: |
|     1 |           0 |                    5 |
|     2 |         120 |                    6 |
|     3 |         300 |                    7 |
|     4 |         600 |                    8 |
|     5 |        1000 |                   10 |

XP никогда не тратится. `xp` — единственный источник истины для Level и Active Backlog Limit: оба значения вычисляются по таблице выше и отдельно не изменяются. В MVP `MAX_LEVEL = 5`; после 1000 XP пользователь остаётся Level 5, а XP продолжает расти.

---

## 5. Добавление и обработка Resource

Пользователь вставляет URL и нажимает **Добавить в backlog**.

Backend выполняет:

1. URL validation и normalization;
2. предварительный duplicate и capacity check;
3. синхронное получение content через Firecrawl;
4. content validation;
5. получение title и tags из Firecrawl metadata или локально из content;
6. повторный duplicate и capacity check в database transaction;
7. при необходимости списание `Overflow Pass` и создание Resource со статусом `PROCESSING` в той же transaction;
8. запуск goroutine для генерации quiz;
9. ответ `202 Accepted` с Resource, включая title и tags.

Повторная проверка перед сохранением обязательна: пока Firecrawl выполнялся, другой request мог занять последний slot или добавить тот же URL. Если Firecrawl или content validation завершились ошибкой, Resource не создаётся и backend возвращает ошибку сразу.

Goroutine использует собственный timeout, не привязанный к завершившемуся HTTP request, и:

1. создаёт и валидирует quiz через LLM;
2. в одной transaction сохраняет Quiz и все Questions;
3. переводит Resource в `NOT_COMPLETED` и впервые устанавливает `activatedAt` в той же transaction.

Frontend опрашивает `GET /api/resources/:id`. Пока Resource находится в `PROCESSING`, frontend уже показывает title, tags и кнопку **Открыть оригинал**, но не позволяет открыть quiz.

### Duplicate и retry rules

Для `(userId, canonicalUrl)` действует database unique constraint. Повторный `POST /api/resources` обрабатывается так:

| Текущий статус | Поведение |
| -------------- | --------- |
| `PROCESSING` | вернуть существующий Resource с `202`, новую goroutine не запускать |
| `NOT_COMPLETED` | вернуть duplicate error |
| `COMPLETED` | вернуть duplicate error |
| `FAILED` | повторно проверить capacity и Overflow Pass, использовать сохранённый content, атомарно перевести Resource в `PROCESSING` и повторить только LLM generation |

Одновременные retry используют условный переход `FAILED → PROCESSING`, поэтому обработку сможет запустить только один запрос.

Firecrawl выполняет одну автоматическую повторную попытку для network error, timeout, `429` и `5xx`, но обе попытки должны укладываться в общий Firecrawl timeout HTTP request.

LLM goroutine выполняет одну автоматическую повторную попытку для network error, timeout, `429`, `5xx`, невалидного ответа LLM или quiz validation error.

```text
MAX_FIRECRAWL_ATTEMPTS = 2
FIRECRAWL_TOTAL_TIMEOUT = 30s
MAX_LLM_ATTEMPTS = 2
LLM_ATTEMPT_TIMEOUT = 60s
```

После последней LLM ошибки Resource переходит в `FAILED` и освобождает slot. Backend хранит безопасный `errorCode`, но не отдаёт frontend сырой ответ провайдера. Если использовался `Overflow Pass`, он возвращается ровно один раз в той же transaction, а `usedOverflowPass` сбрасывается в `false`.

При старте backend переводит оставшиеся после сбоя `PROCESSING` Resources в `FAILED` через тот же failure transition: slot освобождается, а использованный pass возвращается. Пользователь может повторить только LLM generation обычным `POST /api/resources`.

### URL

Разрешены только абсолютные `http://` и `https://` URL.

Блокируются:

- credentials;
- `localhost` и private/internal network;
- `file:`, `javascript:`, `data:`, `ftp:`;
- URL самого Firecrawl.

```text
MAX_URL_LENGTH = 2048
```

Перед duplicate check:

- hostname → lowercase;
- удаляется `#fragment`;
- удаляется default port;
- удаляются `utm_*`, `fbclid`, `gclid`.

Храним `originalUrl` и `canonicalUrl`. Duplicate определяется по `userId + canonicalUrl`.

### Content extraction

```text
URL → Firecrawl → Markdown / readable content
```

Базовые ограничения:

```text
MIN_CONTENT_WORDS = 300
MAX_CONTENT_CHARS = 50_000
MIN_TAGS = 1
MAX_TAGS = 8
MAX_TAG_LENGTH = 32
```

Content короче `MIN_CONTENT_WORDS` отклоняется. Content длиннее `MAX_CONTENT_CHARS` обрезается по последней полной границе paragraph перед лимитом; только сохранённая версия передаётся в LLM и используется как source of truth для evidence.

Title берётся из Firecrawl metadata, затем из первого подходящего heading; последний fallback — hostname. Tags берутся из metadata keywords; если их нет, backend локально выделяет ключевые слова из title, headings и content без дополнительного LLM request. Если подходящих слов нет, используется hostname. Tags приводятся к lowercase, очищаются от дублей и используются для фильтрации Resources.

В MVP не поддерживаются PDF, книги, login pages, CAPTCHA, search results и страницы без достаточного текстового содержимого.

---

## 6. Генерация Quiz

Quiz содержит от **5 до 10 вопросов**.

```text
MIN_QUESTIONS = 5
MAX_QUESTIONS = 10
```

Количество вопросов определяется содержимым материала. Filler-вопросы ради количества не нужны.

Каждый Question должен:

- проверять понимание материала;
- иметь 4 непустых и уникальных options;
- иметь `correctIndex` в диапазоне `0..3` и ровно один правильный ответ;
- иметь правдоподобные distractors;
- не требовать внешних знаний;
- содержать непустое короткое `explanation`;
- содержать непустое `evidence`, которое после нормализации whitespace встречается в source content.

Backend отклоняет весь сгенерированный quiz, если вопросов не 5–10, нарушен порядок, есть дубли вопросов, `totalQuestions` не совпадает с фактическим количеством или хотя бы один Question не проходит эти проверки.

Не использовать вопросы про название статьи, автора или механическое запоминание формулировок.

Источник считается **untrusted content**. LLM не должна выполнять инструкции, найденные внутри страницы, и должна использовать только информацию из SOURCE.

---

## 7. Client-side проверка ответов

Каноническая модель Question определена в разделе 14. Frontend получает все её поля, кроме `correctIndex`.

Для каждого Question backend генерирует случайный `verificationSalt` минимум из 16 random bytes и вычисляет lowercase hex digest от UTF-8 строки:

```text
correctAnswerHash = SHA-256(
  "v1:" + quizId + ":" + questionId + ":" + correctIndex + ":" + verificationSalt
)
```

Генерируемые IDs и salt используют безопасный алфавит без символа `:`.

Frontend получает `correctAnswerHash` и `verificationSalt`, вычисляет hash выбранного `selectedIndex` и проверяет ответ локально.

Если ответ неправильный, вопрос остаётся незакрытым и пользователь может попробовать ещё раз.

Если правильный — вопрос закрывается, показываются `explanation` и `evidence`.

Промежуточные ответы на backend не отправляются. Для сохранения progress после reload можно использовать `localStorage`.

Hash не является защитой от намеренного cheating. Финальным источником истины остаётся backend.

---

## 8. Completion

После правильного ответа на все вопросы frontend делает один запрос:

```http
POST /api/quizzes/:quizId/complete
```

```json
{
  "answers": [
    { "questionId": "q_1", "selectedIndex": 2 },
    { "questionId": "q_2", "selectedIndex": 0 }
  ]
}
```

Backend:

1. проверяет ownership Quiz и Resource;
2. требует ровно один answer для каждого Question без неизвестных или повторяющихся `questionId`;
3. проверяет `selectedIndex` в диапазоне `0..3`;
4. сверяет все answers с `correctIndex`;
5. убеждается, что все ответы правильные;
6. условно переводит Resource из `NOT_COMPLETED` в `COMPLETED`;
7. рассчитывает и начисляет XP и е-баллы;
8. обновляет streak;
9. возвращает новый вычисляемый Level и освобождённый backlog slot.

Проверка статуса, rewards, streak и переход Resource выполняются в одной database transaction. Только запрос, успешно выполнивший переход `NOT_COMPLETED → COMPLETED`, начисляет rewards. Итоговые reward values и breakdown сохраняются в Resource.

Повторный completion не начисляет rewards повторно и возвращает `200 OK` с ранее сохранённым completion result. Это делает endpoint idempotent даже если первый response потерялся.

---

## 9. XP, е-баллы и Streak

### XP

XP — постоянный прогресс пользователя.

```text
XP = 20 + totalQuestions * 5
```

| Questions |  XP |
| --------: | --: |
|         5 |  45 |
|         6 |  50 |
|         8 |  60 |
|        10 |  70 |

XP:

- начисляется только после completion;
- никогда не тратится;
- определяет Level;
- определяет позицию в all-time Leaderboard.

### Е-баллы

Base reward:

```text
baseEPoints = totalQuestions
```

Возраст Resource для rewards считается от `activatedAt` — момента первого перехода в `NOT_COMPLETED`, а не от создания или неудачной генерации quiz. Используется количество полных прошедших 24-часовых периодов.

Дополнительно действует **Old Backlog Bounty**:

| Возраст Resource | Bonus |
| ---------------- | ----: |
| < 7 дней         |    +0 |
| 7–13 дней        |    +1 |
| 14–29 дней       |    +2 |
| 30–59 дней       |    +4 |
| 60+ дней         |    +6 |

Если внутри completion transaction, до перевода завершаемого Resource в `COMPLETED`, используемая capacity была заполнена до обычного лимита, включая Resources в обработке:

```text
notCompletedResources + processingResources >= activeBacklogLimit
→ Full Backlog Bonus = +2 е-балла
```

Итог:

```text
ePointsEarned =
  totalQuestions
  + oldBacklogBonus
  + fullBacklogBonus
```

Пример:

```text
8 вопросов             +8
Resource лежал 37 дней +4
Backlog был заполнен   +2
─────────────────────────
Итого                 +14 е-баллов
```

### Streak

Streak обновляется только после успешного completion. В MVP календарные дни определяются по UTC:

- последний completion в текущий UTC-день → без изменений;
- в предыдущий UTC-день → `streak += 1`;
- раньше или отсутствует → `streak = 1`.

---

## 10. Профиль и Cosmetics

Профиль — простая карточка игрока. Никаких комнат, drag-and-drop и сложного конструктора.

Показываются:

- username;
- avatar;
- frame;
- title;
- один showcase item;
- Level;
- XP;
- streak.

Пример:

```text
┌────────────────────────┐
│       NEON FRAME       │
│                        │
│           🐸           │
│         vasya          │
│                        │
│    Backlog Destroyer   │
│                        │
│           🦆           │
│   Senior Rubber Duck   │
│                        │
│ Level 4 · 820 XP       │
│ 🔥 9 дней              │
└────────────────────────┘
```

Для MVP достаточно emoji, простых SVG и CSS frames.

Храним только выбранные IDs:

```ts
avatarId: string
frameId: string
titleId?: string
showcaseItemId?: string
```

---

## 11. Е-магазин и Overflow Pass

Магазин нужен, чтобы заработок е-баллов имел понятную цель.

Каталог для MVP **hardcoded**, без admin panel и CRUD. Достаточно 10–15 items.

### Пример каталога

| Item                     | Type     | Цена |
| ------------------------ | -------- | ---: |
| 🐸 Frog                  | Avatar   |   15 |
| 🤖 Robot                 | Avatar   |   25 |
| 🧙 Wizard                | Avatar   |   40 |
| Neon                     | Frame    |   30 |
| Fire                     | Frame    |   50 |
| Gold                     | Frame    |  120 |
| Разгребатель             | Title    |   20 |
| Knowledge Goblin         | Title    |   40 |
| Backlog Destroyer        | Title    |   70 |
| 🦆 Senior Rubber Duck    | Showcase |   40 |
| 🌵 Кактус прокрастинации | Showcase |   55 |
| 🐈 Кот                   | Showcase |  120 |
| 👑 Golden Duck           | Showcase |  250 |
| 📦 Overflow Pass         | Utility  |   25 |

Cosmetics покупаются навсегда и могут свободно экипироваться. Purchase выполняется в одной transaction: backend берёт цену из hardcoded catalog, проверяет `ePoints`, списывает е-баллы и добавляет cosmetic или увеличивает `overflowPassCount`. Повторная покупка owned cosmetic возвращает `409 ALREADY_OWNED`; Overflow Pass можно покупать многократно. Equip разрешён только для owned cosmetic соответствующего type.

### Overflow Pass

`Overflow Pass` позволяет один раз добавить Resource при заполненном Active Backlog.

```text
Used capacity: 5 / 5
↓
Использовать Overflow Pass
↓
Used capacity: 6 / 5
```

Ограничение:

```text
notCompletedResources + processingResources <= activeBacklogLimit + 1
```

Если used capacity ниже обычного лимита, `useOverflowPass` игнорируется и pass не списывается. Если capacity равна обычному лимиту, request с `useOverflowPass: false` возвращает `BACKLOG_FULL`, а с `useOverflowPass: true` требует хотя бы один pass и атомарно списывает его при создании Resource.

Если обработка завершилась с `FAILED`, pass возвращается ровно один раз в transaction перевода Resource в `FAILED`.

Нельзя использовать второй pass, пока used capacity уже находится выше обычного лимита. После возвращения к обычному лимиту следующий pass снова можно использовать.

Других utility items в MVP нет.

---

## 12. Leaderboard

В MVP есть только **all-time Leaderboard по XP**.

```text
ORDER BY xp DESC
```

Показываются Top 20 пользователей и текущий пользователь, если он не входит в Top 20.

Пример:

```text
🏆 Leaderboard

1. 🐸 vasya      Lv.5   3420 XP
2. 🤖 lena       Lv.5   3180 XP
3. 🧙 max        Lv.5   2940 XP
...
18. 🐱 you       Lv.5   1320 XP
```

Leaderboard использует avatar и frame пользователя. При открытии профиля можно увидеть title и showcase item. Для стабильного порядка ties сортируются по `xp DESC, userId ASC`.

Leaderboard не выдаёт дополнительных rewards.

---

## 13. Основные экраны

### Home

- Level и XP progress;
- е-баллы;
- streak;
- Active Backlog;
- список незавершённых Resources;
- карточка `PROCESSING` Resource с title, tags, кнопкой **Открыть оригинал** и выключенной кнопкой quiz;
- автоматическое обновление backlog, когда Resource переходит в `NOT_COMPLETED`;
- карточка `FAILED` Resource с безопасной причиной ошибки, **Открыть оригинал** и **Повторить**;
- Old Backlog Bounty на старых материалах;
- Add Resource.

### Resource

- title;
- domain;
- tags;
- processing/error state;
- количество вопросов;
- bounty;
- **Открыть оригинал** доступно в любом статусе;
- **Начать quiz** доступно только при `NOT_COMPLETED`.

### Quiz

- `Вопрос N из M`;
- progress;
- 4 options;
- локальный result;
- explanation и evidence после правильного ответа.

### Result

```text
Quiz выполнен

8 / 8 правильно

+60 XP
+14 е-баллов
  +8 за вопросы
  +4 old backlog bounty
  +2 full backlog bonus

🔥 Streak: 4 дня

Слот в backlog освобождён
```

### Profile

Карточка профиля + выбор купленных cosmetics.

### Shop

Hardcoded список items со статусами:

```text
Купить / Куплено / Надеть / Надето
```

### Leaderboard

All-time Top 20 по XP.

---

## 14. Минимальная модель данных

### User

```ts
type User = {
  id: string;
  email: string;
  username: string;
  passwordHash: string;
  createdAt: Date;
};
```

`email` приводится к lowercase и trim; `username` trim-ится и сравнивается case-insensitive. Оба значения уникальны после normalization. Username должен соответствовать `[A-Za-z0-9_-]{3,32}`, password — иметь длину `8..72` bytes. Password хранится только как bcrypt hash.

### UserProgress

```ts
type UserProgress = {
  userId: string;

  xp: number;
  ePoints: number;

  currentStreak: number;
  lastCompletionAt?: Date;

  overflowPassCount: number;

  avatarId: string;
  frameId: string;
  titleId?: string;
  showcaseItemId?: string;

  ownedCosmeticIds: string[];
};
```

Новый пользователь получает `xp = 0`, `ePoints = 0`, `currentStreak = 0`, `overflowPassCount = 0`, а также бесплатные default avatar и frame, которые считаются owned. `level` и `activeBacklogLimit` добавляются в API response как вычисляемые из XP значения.

### Resource

```ts
type Resource = {
  id: string;
  userId: string;
  originalUrl: string;
  canonicalUrl: string;
  title: string;
  domain: string;
  tags: string[];
  content: string;
  status: ResourceStatus;
  errorCode?: string;
  usedOverflowPass: boolean;
  createdAt: Date;
  updatedAt: Date;
  activatedAt?: Date;
  completedAt?: Date;
  xpEarned?: number;
  ePointsEarned?: number;
  oldBacklogBonus?: number;
  fullBacklogBonus?: number;
};
```

Resource создаётся только после успешного Firecrawl request, поэтому `title`, `domain`, `tags` и `content` доступны уже в `PROCESSING`. При переходе в `FAILED` устанавливается безопасный `errorCode`; retry очищает его. Повторное добавление URL со статусом `FAILED` использует сохранённый content и перезапускает только quiz generation.

### Quiz

```ts
type Quiz = {
  id: string;
  resourceId: string;
  title: string;
  topic?: string;
  totalQuestions: number;
  createdAt: Date;
};
```

### Question

```ts
type Question = {
  id: string;
  quizId: string;
  order: number;
  question: string;
  options: [string, string, string, string];
  correctIndex: number;
  explanation: string;
  evidence: string;
  verificationSalt: string;
  correctAnswerHash: string;
};
```

Отдельные `QuizAttempt` и `QuizAnswer` на backend для MVP не нужны.

`level`, `activeBacklogLimit`, `ownedCosmeticIds` и `tags` в API являются computed representations, а не отдельными изменяемыми источниками истины.

Обязательные database invariants:

- `user_progress.userId` — primary key и foreign key на User;
- unique normalized `users.email` и `users.username`;
- unique `(resources.userId, resources.canonicalUrl)`;
- unique `quizzes.resourceId` — один Quiz на Resource;
- unique `(questions.quizId, questions.order)`;
- unique `(user_cosmetics.userId, user_cosmetics.itemId)`;
- unique `(resource_tags.resourceId, resource_tags.tag)` и index по `resource_tags.tag`;
- допустимы только ResourceStatus из раздела 3;
- foreign keys включены; XP, ePoints, streak и counters не могут быть отрицательными.

Shop catalog хранится в коде. Owned cosmetics хранятся в `user_cosmetics`, tags — в `resource_tags`; массивы в моделях выше являются API representation.

---

## 15. Минимальный Backend API

```http
POST   /api/register
POST   /api/login

POST   /api/resources
GET    /api/resources?status={optionalStatus}&tag={optionalTag}
GET    /api/resources/:id
GET    /api/resources/:id/quiz
POST   /api/quizzes/:quizId/complete

GET    /api/profile
PATCH  /api/profile/cosmetics

GET    /api/shop
POST   /api/shop/purchase

GET    /api/leaderboard
```

При использовании Overflow Pass:

```json
{
  "url": "https://example.com/article",
  "useOverflowPass": true
}
```

После успешного Firecrawl request backend возвращает созданный Resource:

```http
HTTP/1.1 202 Accepted
```

```json
{
  "resourceId": "res_123",
  "status": "PROCESSING",
  "title": "Understanding Go Concurrency",
  "tags": ["go", "concurrency", "goroutines"]
}
```

Frontend опрашивает `GET /api/resources/:id` до состояния `NOT_COMPLETED` или `FAILED`. `GET /api/resources` без filters возвращает все Resources текущего пользователя по `createdAt DESC`; `status` и `tag` — optional exact-match filters.

Register принимает `email`, `username` и `password`. Register и login возвращают JWT. Все остальные endpoints требуют JWT.

Все Resource endpoints проверяют ownership. Для чужого или несуществующего Resource backend возвращает `404`. `GET /api/resources/:id/quiz` возвращает `409 QUIZ_NOT_READY`, если Resource находится в `PROCESSING` или `FAILED`.

Основные API errors:

| HTTP | Code | Когда |
| ---: | ---- | ----- |
| 400 | `VALIDATION_ERROR`, `INVALID_URL` | невалидный request или URL |
| 401 | `UNAUTHORIZED` | отсутствует или невалиден JWT |
| 404 | `NOT_FOUND` | entity не существует или принадлежит другому user |
| 409 | `DUPLICATE_RESOURCE`, `BACKLOG_FULL`, `OVERFLOW_PASS_REQUIRED`, `QUIZ_NOT_READY`, `ALREADY_OWNED`, `INSUFFICIENT_EPOINTS`, `EMAIL_TAKEN`, `USERNAME_TAKEN` | конфликт состояния |
| 422 | `UNSUPPORTED_CONTENT`, `INCORRECT_ANSWERS` | request валиден, но не может быть выполнен |
| 502 | `FIRECRAWL_ERROR` | Firecrawl завершился ошибкой |
| 504 | `FIRECRAWL_TIMEOUT` | общий Firecrawl timeout исчерпан |

---

## 16. Технический стек

- **Backend:** Go + Chi
- **Database:** SQLite + versioned migrations + foreign keys + WAL + busy timeout
- **Authentication:** JWT HS256, signing secret минимум `32` random bytes только из environment, token lifetime `24h`, refresh tokens не входят в MVP
- **Content extraction:** Firecrawl
- **AI:** OpenRouter, `gpt-5-nano`
- **Synchronous extraction:** Firecrawl выполняется внутри `POST /api/resources`
- **Background processing:** только LLM generation, Go goroutine + global semaphore
- **Processing concurrency:** `MAX_CONCURRENT_PROCESSING = 3`
- **Provider calls:** общий Firecrawl request timeout и отдельный OpenRouter timeout
- **HTTP timeout:** server write timeout не меньше `35s`, то есть больше общего Firecrawl timeout
- **Deployment model:** один backend process для MVP; multi-instance deployment не поддерживается
- **Client verification:** SHA-256
- **Local quiz progress:** `localStorage`
- **Cosmetics:** emoji / SVG / CSS
- **Shop catalog:** hardcoded config

---

## 17. Не входит в MVP

- weekly leaderboard;
- leaderboard rewards;
- wishlist;
- Streak Shield;
- ставки на Resource;
- achievements;
- rotating shop;
- loot boxes;
- rarity;
- room builder;
- несколько showcase items;
- server-side сохранение каждого ответа;
- sync quiz progress между устройствами;
- Premium;
- YouTube и PDF;
- bookmark import;
- browser extension;
- notifications;
- adaptive difficulty;
- flashcards;
- spaced repetition;
- social feed и friend system.

---

## 18. Критерий готовности MVP

Core flow:

```text
URL
↓
POST /api/resources ждёт Firecrawl и валидирует content
↓
Backend сохраняет Resource = PROCESSING с title, tags и content
↓
Frontend получает Resource и может открыть оригинал
↓
Goroutine генерирует quiz через LLM
↓
Quiz + Questions сохраняются атомарно
↓
Resource = NOT_COMPLETED, frontend включает кнопку quiz
↓
Ответы проверяются локально
↓
Все вопросы отвечены правильно
↓
Один completion request
↓
Backend повторно проверяет answers
↓
Resource = COMPLETED
↓
XP + е-баллы + streak
↓
Backlog slot освобождён
```

Gamification flow:

```text
XP → Level → больше backlog slots
XP → all-time Leaderboard

е-баллы → Cosmetics → Profile → видны в Leaderboard
е-баллы → Overflow Pass → один временный extra slot

старый Resource → дополнительный bounty → больше мотивации его завершить
```

**Главный критерий:** core learning loop должен работать независимо от Shop и Leaderboard. Геймификация должна усиливать completion, но не усложнять прохождение quiz.
