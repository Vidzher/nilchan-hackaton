# Learning Backlog — MVP PRD

## 0. Статус документа

Этот файл — единственный source of truth для MVP product behavior, backend invariants и API contract. Таблицы, constants и явно описанные rules являются нормативными; примеры UI и payloads — иллюстративными и не могут переопределять rules. Любое изменение поведения сначала фиксируется здесь.

---

## 1. Идея

**Learning Backlog** — приложение для разбора сохранённых «на потом» образовательных материалов.

Пользователь добавляет URL статьи, гайда или документации. HTTP request ждёт только **Firecrawl**. После получения content backend сохраняет Resource с title и tags, возвращает его frontend и генерирует quiz через **LLM** в фоне. Пользователь сразу видит материал и может открыть оригинал; quiz становится доступен после завершения генерации.

Quiz проходится на frontend без запроса к backend после каждого ответа. Вопрос считается закрытым только после правильного ответа. После правильного ответа на все вопросы frontend одним запросом отправляет итоговые ответы на backend. Backend проверяет их, начисляет **XP** и **е-баллы**, помечает материал выполненным и освобождает backlog slot.

XP отвечает за постоянный прогресс и leaderboard. Е-баллы тратятся на оформление профиля и покупку временного overflow slot при заполненном backlog.

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
Quiz вместе с массивом Questions сохраняется атомарно в одном JSON-поле
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

Если `usedCapacity` достиг Active Backlog Limit, новый Resource добавить нельзя, кроме атомарной покупки одного временного overflow slot вместе с созданием Resource. Resource со статусом `FAILED` или `COMPLETED` slot не занимает.

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
5. получение title из Firecrawl metadata и tags из metadata keywords или локально из content;
6. повторный duplicate и capacity check в database transaction;
7. при необходимости списание стоимости временного overflow slot и создание Resource со статусом `PROCESSING` в той же transaction;
8. запуск goroutine для генерации quiz;
9. ответ `202 Accepted` с Resource, включая title и tags.

Повторная проверка перед сохранением обязательна: пока Firecrawl выполнялся, другой request мог занять последний slot или добавить тот же URL. Если Firecrawl или content validation завершились ошибкой, Resource не создаётся и backend возвращает ошибку сразу.

Goroutine использует собственный timeout, не привязанный к завершившемуся HTTP request, и:

1. создаёт и валидирует quiz через LLM;
2. в одной transaction сохраняет Quiz со всем массивом Questions в `questions_json`;
3. переводит Resource в `NOT_COMPLETED` в той же transaction.

Frontend опрашивает `GET /api/resources/:id`. Пока Resource находится в `PROCESSING`, frontend уже показывает title, tags и кнопку **Открыть оригинал**, но не позволяет открыть quiz.

### Duplicate и retry rules

Для `(userId, url)` действует database unique constraint. В `url` хранится нормализованный URL. Повторный `POST /api/resources` обрабатывается так:

| Текущий статус  | Поведение                                                                                                                                                                                   |
| --------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `PROCESSING`    | вернуть существующий Resource с `202`, новую goroutine не запускать                                                                                                                         |
| `NOT_COMPLETED` | вернуть duplicate error                                                                                                                                                                     |
| `COMPLETED`     | вернуть duplicate error                                                                                                                                                                     |
| `FAILED`        | повторно проверить capacity и при необходимости заново купить overflow slot, использовать сохранённый content, атомарно перевести Resource в `PROCESSING` и повторить только LLM generation |

Одновременные retry используют условный переход `FAILED → PROCESSING`, поэтому обработку сможет запустить только один запрос.

Firecrawl выполняет одну автоматическую повторную попытку для network error, timeout, `429` и `5xx`, но обе попытки должны укладываться в общий Firecrawl timeout HTTP request.

LLM goroutine выполняет одну автоматическую повторную попытку для network error, timeout, `429`, `5xx`, невалидного ответа LLM или quiz validation error.

```text
MAX_FIRECRAWL_ATTEMPTS = 2
FIRECRAWL_TOTAL_TIMEOUT = 30s
MAX_LLM_ATTEMPTS = 2
LLM_ATTEMPT_TIMEOUT = 120s
```

После последней LLM ошибки Resource переходит в `FAILED` и освобождает slot. Backend не отдаёт frontend сырой ответ провайдера. Если для Resource был куплен overflow slot, его стоимость возвращается в е-баллах ровно один раз в той же transaction, а `purchasedOverflowSlot` сбрасывается в `false`.

При старте backend переводит оставшиеся после сбоя `PROCESSING` Resources в `FAILED` через тот же failure transition: slot освобождается, а стоимость купленного overflow slot возвращается. Пользователь может повторить только LLM generation обычным `POST /api/resources`.

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

Храним один нормализованный `url`. Duplicate определяется по `userId + url`; этот же URL используется для открытия оригинала.

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

Title берётся только из Firecrawl metadata. Если metadata не содержит непустой title, Resource не создаётся. Tags берутся из metadata keywords; если их нет, backend локально выделяет ключевые слова из title, headings и content без дополнительного LLM request. Если подходящих слов нет, используется hostname. Tags приводятся к lowercase, очищаются от дублей и используются для фильтрации Resources.

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

Backend отклоняет весь сгенерированный quiz, если вопросов не 5–10, есть дубли вопросов или хотя бы один Question не проходит эти проверки. Порядок Questions определяется их позицией в JSON-массиве, а `totalQuestions` вычисляется как длина массива и отдельно не хранится.

Не использовать вопросы про название статьи, автора или механическое запоминание формулировок.

Источник считается **untrusted content**. LLM не должна выполнять инструкции, найденные внутри страницы, и должна использовать только информацию из SOURCE.

LLM возвращает только объект с массивом `questions`. `Quiz.id`, `Quiz.title`, `verificationSalt` и `correctAnswerHash` не генерируются LLM: title копируется из сохранённого Resource, остальные значения создаются backend после валидации ответа.

OpenRouter вызывается со следующим structured output contract:

```json
{
  "type": "json_schema",
  "json_schema": {
    "name": "quiz_generation",
    "strict": true,
    "schema": {
      "type": "object",
      "additionalProperties": false,
      "required": ["questions"],
      "properties": {
        "questions": {
          "type": "array",
          "minItems": 5,
          "maxItems": 10,
          "items": {
            "type": "object",
            "additionalProperties": false,
            "required": [
              "text",
              "options",
              "correctIndex",
              "explanation",
              "evidence"
            ],
            "properties": {
              "text": {
                "type": "string",
                "minLength": 1,
                "description": "Question testing understanding of the source material"
              },
              "options": {
                "type": "array",
                "minItems": 4,
                "maxItems": 4,
                "uniqueItems": true,
                "items": {
                  "type": "string",
                  "minLength": 1
                },
                "description": "Exactly four unique answer options"
              },
              "correctIndex": {
                "type": "integer",
                "minimum": 0,
                "maximum": 3,
                "description": "Zero-based index of the only correct option"
              },
              "explanation": {
                "type": "string",
                "minLength": 1,
                "description": "Short explanation of why the answer is correct"
              },
              "evidence": {
                "type": "string",
                "minLength": 1,
                "description": "Exact excerpt from SOURCE supporting the correct answer"
              }
            }
          }
        }
      }
    }
  }
}
```

Сокращённый фрагмент ответа LLM; полный ответ обязан содержать 5–10 вопросов:

```json
{
  "questions": [
    {
      "text": "Why are goroutines considered lightweight compared with operating-system threads?",
      "options": [
        "They use dynamically growing stacks managed by the Go runtime",
        "They always execute without operating-system threads",
        "They cannot perform blocking operations",
        "They share a single fixed-size stack"
      ],
      "correctIndex": 0,
      "explanation": "Goroutine stacks start small and grow when necessary, reducing their initial memory cost.",
      "evidence": "Goroutines have dynamically growing stacks that begin with a small amount of memory."
    }
  ]
}
```

JSON Schema не заменяет backend validation. После parsing backend нормализует whitespace, проверяет непустые строки и уникальность options, отсутствие дублирующихся вопросов, диапазон `correctIndex`, а также точное присутствие нормализованного `evidence` в сохранённом `Resource.content`. Любое нарушение отклоняет весь quiz и считается quiz validation error.

---

## 7. Client-side проверка ответов

Каноническая модель Question определена в разделе 14. Frontend получает все её поля, кроме `correctIndex`.

Для каждого Question backend генерирует случайный `verificationSalt` минимум из 16 random bytes и вычисляет lowercase hex digest от UTF-8 строки. Question идентифицируется позицией `questionIndex` в JSON-массиве:

```text
correctAnswerHash = SHA-256(
  "v1:" + quizId + ":" + questionIndex + ":" + correctIndex + ":" + verificationSalt
)
```

`verificationSalt` использует безопасный алфавит без символа `:`. Entity IDs являются числовыми.

Frontend получает `correctAnswerHash` и `verificationSalt`, вычисляет hash выбранного `selectedIndex` и проверяет ответ локально.

Если ответ неправильный, вопрос остаётся незакрытым и пользователь может попробовать ещё раз.

Если правильный — вопрос закрывается, показываются `explanation` и `evidence`.

Промежуточные ответы на backend не отправляются. Для сохранения progress после reload можно использовать `localStorage`.

Hash не является защитой от намеренного cheating. Финальным источником истины остаётся backend.

---

## 8. Completion

После правильного ответа на все вопросы frontend делает один запрос:

```http
POST /api/resources/:id/quiz/complete
```

```json
{
  "answers": [
    { "questionIndex": 0, "selectedIndex": 2 },
    { "questionIndex": 1, "selectedIndex": 0 }
  ]
}
```

Backend:

1. проверяет ownership Quiz и Resource;
2. требует ровно один answer для каждой позиции Question без неизвестных или повторяющихся `questionIndex`;
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

Reward:

```text
ePointsEarned = totalQuestions
```

| Questions | Е-баллы |
| --------: | ------: |
|         5 |       5 |
|         6 |       6 |
|         8 |       8 |
|        10 |      10 |

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

## 11. Е-магазин

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

Cosmetics покупаются навсегда и могут свободно экипироваться. Purchase выполняется в одной transaction: backend берёт цену из hardcoded catalog, проверяет `ePoints`, списывает е-баллы и добавляет cosmetic. Повторная покупка owned cosmetic возвращает `409 ALREADY_OWNED`. Equip разрешён только для owned cosmetic соответствующего type.

### Временный overflow slot

Временный overflow slot покупается только вместе с добавлением Resource при заполненном Active Backlog. Отдельно купить или хранить его нельзя.

```text
OVERFLOW_SLOT_PRICE = 25 е-баллов
```

```text
Used capacity: 5 / 5
↓
Купить overflow slot за 25 е-баллов и добавить Resource
↓
Used capacity: 6 / 5
```

Ограничение:

```text
notCompletedResources + processingResources <= activeBacklogLimit + 1
```

Если used capacity ниже обычного лимита, `purchaseOverflowSlot` игнорируется и е-баллы не списываются. Если capacity равна обычному лимиту, request с `purchaseOverflowSlot: false` возвращает `BACKLOG_FULL`, а с `purchaseOverflowSlot: true` требует минимум 25 е-баллов и атомарно списывает их при создании Resource.

Если обработка завершилась с `FAILED`, 25 е-баллов возвращаются ровно один раз в transaction перевода Resource в `FAILED`.

Нельзя купить второй overflow slot, пока used capacity уже находится выше обычного лимита. После возвращения к обычному лимиту следующий slot снова можно купить.

Utility items в Shop в MVP нет.

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
- карточка `FAILED` Resource с общим сообщением об ошибке, **Открыть оригинал** и **Повторить**;
- Add Resource.

### Resource

- title;
- tags;
- processing/error state;
- количество вопросов;
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
+8 е-баллов

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

Все entity IDs являются целыми числами: в Go используется `int64`, в JSON — number.

### User

```ts
type User = {
  id: number;
  email: string;
  username: string;
  passwordHash: string;
  createdAt: Date;
};
```

`email` и `username` сравниваются case-insensitive средствами SQLite `COLLATE NOCASE` и уникальны. Username должен иметь длину `3..32` символа, password — `8..72` символа согласно строковой валидации Go. Password хранится только как bcrypt hash.

### UserProgress

```ts
type UserProgress = {
  userId: number;

  xp: number;
  ePoints: number;

  currentStreak: number;
  lastCompletionAt?: Date;

  avatarId: string;
  frameId: string;
  titleId?: string;
  showcaseItemId?: string;

  ownedCosmeticIds: string[];
};
```

Новый пользователь получает `xp = 0`, `ePoints = 0`, `currentStreak = 0`, а также бесплатные default avatar и frame, которые считаются owned. `level` и `activeBacklogLimit` добавляются в API response как вычисляемые из XP значения.

### Resource

```ts
type Resource = {
  id: number;
  userId: number;
  url: string;
  title: string;
  tags: string[];
  content: string;
  status: ResourceStatus;
  purchasedOverflowSlot: boolean;
  createdAt: Date;
  completedAt?: Date;
  xpEarned?: number;
  ePointsEarned?: number;
};
```

Resource создаётся только после успешного Firecrawl request, поэтому `title`, `tags` и `content` доступны уже в `PROCESSING`. Повторное добавление URL со статусом `FAILED` использует сохранённый content и перезапускает только quiz generation.

### Quiz

```ts
type Quiz = {
  id: number;
  resourceId: number;
  title: string;
  questions: Question[];
};
```

### Question

```ts
type Question = {
  text: string;
  options: [string, string, string, string];
  correctIndex: number;
  explanation: string;
  evidence: string;
  verificationSalt: string;
  correctAnswerHash: string;
};
```

Обязательные database invariants:

- `user_progress.userId` — primary key и foreign key на User;
- unique `users.email` и `users.username` с SQLite `COLLATE NOCASE`;
- unique `(resources.userId, resources.url)`;
- unique `quizzes.resourceId` — один Quiz на Resource;
- массив Questions хранится сериализованным в `quizzes.questions_json`; отдельной таблицы `questions` нет;
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
POST   /api/resources/:id/quiz/complete

GET    /api/profile
PATCH  /api/profile/cosmetics

GET    /api/shop
POST   /api/shop/purchase

GET    /api/leaderboard
```

При покупке временного overflow slot:

```json
{
  "url": "https://example.com/article",
  "purchaseOverflowSlot": true
}
```

После успешного Firecrawl request backend возвращает созданный Resource:

```http
HTTP/1.1 202 Accepted
```

```json
{
  "resourceId": 123,
  "status": "PROCESSING",
  "title": "Understanding Go Concurrency",
  "tags": ["go", "concurrency", "goroutines"]
}
```

Frontend опрашивает `GET /api/resources/:id` до состояния `NOT_COMPLETED` или `FAILED`. `GET /api/resources` без filters возвращает все Resources текущего пользователя по `createdAt DESC`; `status` и `tag` — optional exact-match filters.

Register принимает `email`, `username` и `password`. Register и login возвращают JWT. Все остальные endpoints требуют JWT.

Все Resource endpoints проверяют ownership. Для чужого или несуществующего Resource backend возвращает `404`. `GET /api/resources/:id/quiz` возвращает `409`, если Resource находится в `PROCESSING` или `FAILED`.

API errors используют HTTP status и человекочитаемое сообщение без отдельного машинного `code`:

```json
{
  "status": "Error",
  "error": "human-readable message"
}
```

Основные HTTP statuses:

| HTTP | Когда                                                                                                |
| ---: | ---------------------------------------------------------------------------------------------------- |
|  400 | невалидный request или URL                                                                           |
|  401 | отсутствует или невалиден JWT, либо неверны credentials                                              |
|  404 | entity не существует или принадлежит другому user                                                    |
|  409 | конфликт состояния, duplicate email/username/resource, заполненный backlog или недостаточно е-баллов |
|  422 | request валиден, но content или answers не могут быть приняты                                        |
|  502 | Firecrawl завершился ошибкой                                                                         |
|  504 | общий Firecrawl timeout исчерпан                                                                     |

Auth middleware для отсутствующего или невалидного JWT возвращает plain-text сообщение с HTTP `401`; login/register handlers используют JSON envelope выше.

---

## 16. Технический стек

- **Backend:** Go + Chi
- **Database:** SQLite + встроенный mutable `internal/storage/schema.sql` через `CREATE TABLE IF NOT EXISTS` + foreign keys + WAL + busy timeout; versioned migrations не используются
- **Authentication:** JWT HS256 с hardcoded package-level signing key, token lifetime `24h`, refresh tokens не входят в MVP
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
Quiz с массивом Questions в `questions_json` сохраняется атомарно
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
е-баллы → покупка временного overflow slot при добавлении Resource
```

**Главный критерий:** core learning loop должен работать независимо от Shop и Leaderboard. Геймификация должна усиливать completion, но не усложнять прохождение quiz.
