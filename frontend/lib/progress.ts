const levelMinimumXP: Record<number, number> = {
  1: 0,
  2: 120,
  3: 300,
  4: 600,
  5: 1000,
}

export function getLevelProgress(xp: number, level: number) {
  const levelStart = levelMinimumXP[level] ?? 0
  const nextLevelXP = levelMinimumXP[level + 1]

  if (nextLevelXP === undefined) {
    return {
      current: 1,
      required: 1,
      remaining: 0,
      nextLevel: null,
      isMaxLevel: true,
    }
  }

  const required = nextLevelXP - levelStart
  const current = Math.min(required, Math.max(0, xp - levelStart))

  return {
    current,
    required,
    remaining: Math.max(0, required - current),
    nextLevel: level + 1,
    isMaxLevel: false,
  }
}
