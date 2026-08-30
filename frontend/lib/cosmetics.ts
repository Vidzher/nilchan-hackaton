const cosmetics: Record<string, { name: string }> = {
  default_avatar: { name: 'Нормис' },
  default_frame: { name: 'Default Frame' },
  frame_neon: { name: 'Neon' },
  frame_fire: { name: 'Fire' },
  frame_gold: { name: 'Gold' },
}

export function cosmeticName(id?: string) {
  return id ? (cosmetics[id]?.name ?? id) : 'Не выбрано'
}
