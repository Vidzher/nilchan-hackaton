"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import {
  ArrowRight,
  Check,
  CircleCheck,
  Coins,
  ExternalLink,
  Flame,
  Layers,
  Quote,
  RotateCw,
  Sparkles,
  Trophy,
  X,
  Zap,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { ProgressBar } from "@/components/primitives";
import { cn } from "@/lib/utils";
import {
  api,
  ApiError,
  type CompletionResponse,
  type Quiz,
  type Resource,
  type SubmittedAnswer,
} from "@/lib/api";

type SavedProgress = {
  answers: SubmittedAnswer[];
  queue?: number[];
  position?: number;
  missed?: number[];
  pendingRound?: number[] | null;
  round?: number;
};

async function verify(
  quizId: number,
  questionIndex: number,
  selectedIndex: number,
  salt: string,
  expected: string,
) {
  const value = `v1:${quizId}:${questionIndex}:${selectedIndex}:${salt}`;
  const digest = await crypto.subtle.digest(
    "SHA-256",
    new TextEncoder().encode(value),
  );
  const actual = Array.from(new Uint8Array(digest), (byte) =>
    byte.toString(16).padStart(2, "0"),
  ).join("");
  return actual === expected.toLowerCase();
}

export function QuizRunner({
  resource,
  quiz,
}: {
  resource: Resource;
  quiz: Quiz;
}) {
  const storageKey = `learning-backlog.quiz.${quiz.id}`;
  const total = quiz.questions.length;
  const [answers, setAnswers] = useState<SubmittedAnswer[]>([]);
  const [queue, setQueue] = useState<number[]>([]);
  const [position, setPosition] = useState(0);
  const [missed, setMissed] = useState<number[]>([]);
  const [pendingRound, setPendingRound] = useState<number[] | null>(null);
  const [round, setRound] = useState(1);
  const [restored, setRestored] = useState(false);
  const [selected, setSelected] = useState<number | null>(null);
  const [result, setResult] = useState<"correct" | "wrong" | null>(null);
  const [completion, setCompletion] = useState<CompletionResponse | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [completionAttempted, setCompletionAttempted] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    async function restore() {
      const restoredAnswers: SubmittedAnswer[] = [];
      const seen = new Set<number>();
      let parsed: SavedProgress = { answers: [] };
      try {
        parsed = JSON.parse(
          localStorage.getItem(storageKey) ?? "{}",
        ) as SavedProgress;
        if (Array.isArray(parsed.answers)) {
          for (const answer of parsed.answers) {
            if (
              !Number.isInteger(answer.questionIndex) ||
              answer.questionIndex < 0 ||
              answer.questionIndex >= total ||
              seen.has(answer.questionIndex) ||
              !Number.isInteger(answer.selectedIndex) ||
              answer.selectedIndex < 0 ||
              answer.selectedIndex > 3
            )
              continue;

            const question = quiz.questions[answer.questionIndex];
            if (
              await verify(
                quiz.id,
                answer.questionIndex,
                answer.selectedIndex,
                question.verificationSalt,
                question.correctAnswerHash,
              )
            ) {
              restoredAnswers.push(answer);
              seen.add(answer.questionIndex);
            }
          }
        }
      } catch {
        localStorage.removeItem(storageKey);
      }

      if (!cancelled) {
        restoredAnswers.sort(
          (left, right) => left.questionIndex - right.questionIndex,
        );
        const unresolved = quiz.questions
          .map((_, index) => index)
          .filter((index) => !seen.has(index));
        const validIndices = (value: unknown) =>
          Array.isArray(value) &&
          value.every(
            (index) => Number.isInteger(index) && index >= 0 && index < total,
          ) &&
          new Set(value).size === value.length;
        const savedQueue = validIndices(parsed.queue)
          ? (parsed.queue as number[])
          : null;
        const savedMissed = validIndices(parsed.missed)
          ? (parsed.missed as number[]).filter((index) => !seen.has(index))
          : [];
        const savedPendingRound =
          validIndices(parsed.pendingRound) &&
          (parsed.pendingRound as number[]).length > 0
            ? (parsed.pendingRound as number[]).filter(
                (index) => !seen.has(index),
              )
            : null;
        const savedPosition = Number.isInteger(parsed.position)
          ? (parsed.position as number)
          : -1;
        const savedRound =
          Number.isInteger(parsed.round) && (parsed.round as number) > 0
            ? (parsed.round as number)
            : 1;

        setAnswers(restoredAnswers);
        setRound(savedRound);
        if (
          savedQueue &&
          savedQueue.length > 0 &&
          savedPosition >= 0 &&
          savedPosition < savedQueue.length
        ) {
          let resumePosition = savedPosition;
          const alreadyAttempted = new Set([...seen, ...savedMissed]);
          while (
            resumePosition < savedQueue.length &&
            alreadyAttempted.has(savedQueue[resumePosition])
          )
            resumePosition += 1;

          setQueue(savedQueue);
          setMissed(savedMissed);
          if (savedPendingRound && savedPendingRound.length > 0) {
            setPosition(savedPosition);
            setPendingRound(savedPendingRound);
          } else if (resumePosition < savedQueue.length) {
            setPosition(resumePosition);
          } else if (savedMissed.length > 0) {
            setPosition(savedQueue.length - 1);
            setPendingRound(savedMissed);
          } else {
            setQueue(unresolved);
            setPosition(0);
          }
        } else {
          setQueue(unresolved);
        }
        setRestored(true);
      }
    }

    void restore();
    return () => {
      cancelled = true;
    };
  }, [quiz.id, quiz.questions, storageKey, total]);

  useEffect(() => {
    if (restored) {
      localStorage.setItem(
        storageKey,
        JSON.stringify({
          answers,
          queue,
          position,
          missed,
          pendingRound,
          round,
        } satisfies SavedProgress),
      );
    }
  }, [
    answers,
    missed,
    pendingRound,
    position,
    queue,
    restored,
    round,
    storageKey,
  ]);

  useEffect(() => {
    if (
      restored &&
      answers.length === total &&
      result === null &&
      !completion &&
      !completionAttempted
    ) {
      void submitCompletion();
    }
  }, [answers, completion, completionAttempted, restored, result, total]);

  const questionIndex = queue[position];
  const question = quiz.questions[questionIndex];

  async function check() {
    if (selected === null || !question) return;
    setError(null);
    const correct = await verify(
      quiz.id,
      questionIndex,
      selected,
      question.verificationSalt,
      question.correctAnswerHash,
    );
    if (!correct) {
      setMissed((current) => [...current, questionIndex]);
      setResult("wrong");
      return;
    }

    setAnswers((current) =>
      [...current, { questionIndex, selectedIndex: selected }].sort(
        (left, right) => left.questionIndex - right.questionIndex,
      ),
    );
    setResult("correct");
  }

  async function submitCompletion() {
    if (submitting || completion) return;
    setCompletionAttempted(true);
    setSubmitting(true);
    setError(null);
    try {
      const completed = await api.completeQuiz(resource.id, answers);
      setCompletion(completed);
      localStorage.removeItem(storageKey);
    } catch (caught) {
      setError(
        caught instanceof ApiError
          ? caught.message
          : "Не удалось завершить quiz.",
      );
    } finally {
      setSubmitting(false);
    }
  }

  function next() {
    if (position + 1 < queue.length) {
      setPosition((current) => current + 1);
      setSelected(null);
      setResult(null);
      return;
    }
    if (missed.length > 0) {
      setPendingRound(missed);
      setSelected(null);
      setResult(null);
      return;
    }
    void submitCompletion();
  }

  function startRetryRound() {
    if (!pendingRound) return;
    setQueue(pendingRound);
    setMissed([]);
    setPendingRound(null);
    setPosition(0);
    setRound((current) => current + 1);
  }

  if (completion)
    return <QuizResult resource={resource} completion={completion} />;

  if (pendingRound) {
    return (
      <main className="grid min-h-screen place-items-center bg-background px-4">
        <div className="w-full max-w-md rounded-2xl border border-border bg-card p-6 text-center">
          <span className="mx-auto grid size-12 place-items-center rounded-full bg-warning-soft">
            <RotateCw
              className="size-6 text-[color:var(--warning)]"
              aria-hidden="true"
            />
          </span>
          <h1 className="mt-4 text-xl font-semibold">
            Повторим сложные вопросы
          </h1>
          <p className="mt-2 text-sm text-muted-foreground">
            {pendingRound.length} из них{" "}
            {pendingRound.length === 1 ? "требует" : "требуют"} ещё одной
            попытки.
          </p>
          <Button className="mt-6 w-full" onClick={startRetryRound}>
            Начать следующий раунд
            <ArrowRight className="size-4" aria-hidden="true" />
          </Button>
        </div>
      </main>
    );
  }

  if (restored && queue.length === 0) {
    return (
      <main className="grid min-h-screen place-items-center bg-background px-4">
        <div className="w-full max-w-md rounded-2xl border border-border bg-card p-6 text-center">
          <h1 className="text-xl font-semibold">
            {submitting ? "Завершаем quiz…" : "Все вопросы отвечены"}
          </h1>
          {error ? (
            <p
              className="mt-3 text-sm text-[color:var(--destructive)]"
              role="alert"
            >
              {error}
            </p>
          ) : null}
          <Button
            className="mt-5 w-full"
            onClick={() => void submitCompletion()}
            disabled={submitting}
          >
            {submitting ? "Завершаем…" : "Повторить завершение"}
          </Button>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-background">
      <header className="sticky top-0 z-20 border-b border-border bg-background/90 backdrop-blur">
        <div className="mx-auto flex max-w-2xl items-center gap-4 px-4 py-3">
          <Link
            href="/"
            aria-label="Выйти из quiz"
            className="grid size-8 place-items-center rounded-md text-muted-foreground hover:bg-secondary"
          >
            <X className="size-4" aria-hidden="true" />
          </Link>
          <div className="flex-1">
            <ProgressBar
              value={answers.length}
              max={total}
              label="Прогресс quiz"
            />
          </div>
          <span className="tabular text-sm font-medium text-muted-foreground">
            {position + 1} / {queue.length}
          </span>
        </div>
      </header>
      <div className="mx-auto max-w-2xl px-4 py-8">
        {question ? (
          <div>
            <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
              Раунд {round} · Вопрос {position + 1}
            </p>
            <h1 className="mt-2 text-xl font-semibold leading-snug text-balance">
              {question.text}
            </h1>
            <fieldset className="mt-6 flex flex-col gap-2.5">
              <legend className="sr-only">Варианты ответа</legend>
              {question.options.map((option, optionIndex) => {
                const chosen = selected === optionIndex;
                const showCorrect = result === "correct" && chosen;
                const showWrong = result === "wrong" && chosen;
                return (
                  <label
                    key={optionIndex}
                    className={cn(
                      "flex items-start gap-3 rounded-xl border p-3.5 text-sm transition-colors",
                      !result &&
                        chosen &&
                        "border-[color:var(--brand)] bg-brand-soft",
                      !result &&
                        !chosen &&
                        "cursor-pointer border-border bg-card hover:border-muted-foreground/40",
                      showCorrect &&
                        "border-[color:var(--success)] bg-success-soft",
                      showWrong &&
                        "border-[color:var(--destructive)] bg-[#f6e3e3]",
                      result && !chosen && "border-border bg-card opacity-70",
                    )}
                  >
                    <input
                      type="radio"
                      name="answer"
                      className="sr-only"
                      checked={chosen}
                      disabled={result !== null}
                      onChange={() => setSelected(optionIndex)}
                    />
                    <span
                      className={cn(
                        "mt-0.5 grid size-5 shrink-0 place-items-center rounded-full border",
                        chosen &&
                          !result &&
                          "border-[color:var(--brand)] bg-[color:var(--brand)] text-white",
                        showCorrect &&
                          "border-[color:var(--success)] bg-[color:var(--success)] text-white",
                        showWrong &&
                          "border-[color:var(--destructive)] bg-[color:var(--destructive)] text-white",
                      )}
                    >
                      {showCorrect ? (
                        <Check className="size-3.5" />
                      ) : showWrong ? (
                        <X className="size-3.5" />
                      ) : null}
                    </span>
                    <span className="pt-0.5 text-pretty">{option}</span>
                  </label>
                );
              })}
            </fieldset>
            {result ? (
              <div
                className={cn(
                  "mt-5 rounded-xl border p-4",
                  result === "correct"
                    ? "border-[color:var(--success)]/30 bg-success-soft"
                    : "border-[color:var(--warning)]/30 bg-warning-soft",
                )}
              >
                <p
                  className={cn(
                    "flex items-center gap-2 text-sm font-semibold",
                    result === "correct"
                      ? "text-[color:var(--success)]"
                      : "text-[color:var(--warning)]",
                  )}
                >
                  <CircleCheck className="size-4" aria-hidden="true" />
                  {result === "correct"
                    ? "Верно!"
                    : "Не совсем. Вернёмся к этому вопросу позже."}
                </p>
                {result === "correct" ? (
                  <>
                    <p className="mt-2 text-sm text-foreground/90 text-pretty">
                      {question.explanation}
                    </p>
                    <figure className="mt-3 rounded-lg border border-border bg-card p-3">
                      <Quote
                        className="size-4 text-muted-foreground"
                        aria-hidden="true"
                      />
                      <blockquote className="mt-1.5 text-sm italic text-muted-foreground text-pretty">
                        {question.evidence}
                      </blockquote>
                    </figure>
                  </>
                ) : null}
              </div>
            ) : null}
            {error ? (
              <p
                className="mt-4 text-sm text-[color:var(--destructive)]"
                role="alert"
              >
                {error}
              </p>
            ) : null}
            <div className="mt-6">
              {result ? (
                <Button className="w-full" onClick={next} disabled={submitting}>
                  {submitting
                    ? "Завершаем…"
                    : position + 1 === queue.length && missed.length === 0
                      ? "Завершить quiz"
                      : "Продолжить"}
                  <ArrowRight className="size-4" aria-hidden="true" />
                </Button>
              ) : (
                <Button
                  className="w-full"
                  onClick={() => void check()}
                  disabled={selected === null}
                >
                  Проверить ответ
                </Button>
              )}
            </div>
          </div>
        ) : null}
      </div>
    </main>
  );
}

export function QuizResult({
  resource,
  completion,
}: {
  resource: Resource;
  completion: CompletionResponse;
}) {
  const { totalQuestions, xpEarned, ePointsEarned } = completion.completion;

  return (
    <main className="relative grid min-h-screen place-items-center overflow-hidden bg-background px-4 py-10">
      <div className="pointer-events-none absolute left-[8%] top-[12%] size-56 rounded-full bg-brand-soft/70 blur-3xl" />
      <div className="pointer-events-none absolute bottom-[8%] right-[10%] size-64 rounded-full bg-success-soft/80 blur-3xl" />
      <Sparkles
        className="pointer-events-none absolute left-[18%] top-[22%] hidden size-6 rotate-12 text-[color:var(--brand)]/40 sm:block"
        aria-hidden="true"
      />
      <Sparkles
        className="pointer-events-none absolute bottom-[24%] right-[18%] hidden size-5 -rotate-12 text-[color:var(--success)]/40 sm:block"
        aria-hidden="true"
      />

      <section className="animate-in fade-in zoom-in-95 relative w-full max-w-2xl overflow-hidden rounded-2xl border border-border bg-card shadow-[0_24px_80px_rgba(36,36,31,0.10)] duration-500">
        <div className="p-6 text-center sm:p-9">
          <div className="relative mx-auto w-fit">
            <span className="grid size-16 place-items-center rounded-2xl bg-brand-soft ring-8 ring-brand-soft/40">
              <Trophy
                className="size-8 text-[color:var(--brand)]"
                aria-hidden="true"
              />
            </span>
            <span className="absolute -right-3 -top-3 grid size-7 place-items-center rounded-full bg-[color:var(--success)] text-white shadow-sm">
              <Check className="size-4" aria-hidden="true" />
            </span>
          </div>

          <p className="mt-6 text-xs font-semibold uppercase tracking-[0.18em] text-[color:var(--success)]">
            Материал освоен
          </p>
          <h1 className="mt-2 text-3xl font-semibold tracking-tight">
            Тест пройден!
          </h1>
          <p className="mx-auto mt-2 max-w-lg text-sm leading-relaxed text-muted-foreground text-balance">
            {resource.title}
          </p>

          <div className="mx-auto mt-6 inline-flex items-center gap-2 rounded-full border border-border bg-background px-4 py-2">
            <CircleCheck
              className="size-4 text-[color:var(--success)]"
              aria-hidden="true"
            />
            <span className="tabular text-sm font-semibold">
              {totalQuestions} из {totalQuestions}
            </span>
            <span className="text-sm text-muted-foreground">
              вопросов закрыто
            </span>
          </div>

          <div className="mt-7 grid grid-cols-2 gap-3">
            <div className="rounded-xl border border-[color:var(--brand)]/20 bg-brand-soft/60 p-4 text-left sm:p-5">
              <span className="grid size-8 place-items-center rounded-lg bg-card text-[color:var(--brand)] shadow-sm">
                <Zap className="size-4" aria-hidden="true" />
              </span>
              <p className="tabular mt-3 text-3xl font-semibold text-[color:var(--brand)]">
                +{xpEarned}
              </p>
              <p className="mt-0.5 text-xs font-medium text-muted-foreground">
                опыта получено
              </p>
            </div>
            <div className="rounded-xl border border-[color:var(--success)]/20 bg-success-soft/70 p-4 text-left sm:p-5">
              <span className="grid size-8 place-items-center rounded-lg bg-card text-[color:var(--success)] shadow-sm">
                <Coins className="size-4" aria-hidden="true" />
              </span>
              <p className="tabular mt-3 text-3xl font-semibold text-[color:var(--success)]">
                +{ePointsEarned}
              </p>
              <p className="mt-0.5 text-xs font-medium text-muted-foreground">
                е-баллов заработано
              </p>
            </div>
          </div>

          <div className="mt-3 grid gap-3 sm:grid-cols-2">
            <div className="flex items-center gap-3 rounded-xl border border-border bg-background p-3.5 text-left">
              <span className="grid size-9 shrink-0 place-items-center rounded-lg bg-warning-soft">
                <Flame
                  className="size-5 text-[color:var(--warning)]"
                  aria-hidden="true"
                />
              </span>
              <div>
                <p className="tabular text-sm font-semibold">
                  Серия: {completion.progress.currentStreak}
                </p>
                <p className="text-xs text-muted-foreground">Так держать</p>
              </div>
            </div>
            <div className="flex items-center gap-3 rounded-xl border border-border bg-background p-3.5 text-left">
              <span className="grid size-9 shrink-0 place-items-center rounded-lg bg-success-soft">
                <Layers
                  className="size-5 text-[color:var(--success)]"
                  aria-hidden="true"
                />
              </span>
              <div>
                <p className="text-sm font-semibold">Слот освобождён</p>
                <p className="text-xs text-muted-foreground">
                  Backlog стал легче
                </p>
              </div>
            </div>
          </div>

          <div className="mt-7 grid grid-cols-1 gap-2 sm:grid-cols-2">
            <Button
              className="w-full"
              nativeButton={false}
              render={<Link href="/" />}
            >
              Вернуться в backlog
              <ArrowRight className="size-4" aria-hidden="true" />
            </Button>
            <Button
              variant="outline"
              className="w-full"
              nativeButton={false}
              render={
                <a
                  href={resource.url}
                  target="_blank"
                  rel="noreferrer noopener"
                />
              }
            >
              Открыть оригинал
              <ExternalLink className="size-4" aria-hidden="true" />
            </Button>
          </div>
        </div>
      </section>
    </main>
  );
}
