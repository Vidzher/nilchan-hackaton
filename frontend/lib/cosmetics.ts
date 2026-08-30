const cosmetics: Record<string, { name: string; preview?: string }> = {
  default_avatar: { name: 'Default Avatar', preview: '🙂' },
  default_frame: { name: 'Default Frame' },
  avatar_frog: { name: 'Frog', preview: '🐸' },
  avatar_robot: { name: 'Robot', preview: '🤖' },
  avatar_wizard: { name: 'Wizard', preview: '🧙' },
  frame_neon: { name: 'Neon' },
  frame_fire: { name: 'Fire' },
  frame_gold: { name: 'Gold' },
  title_razgreshatel: { name: 'Разгребатель' },
  title_knowledge_goblin: { name: 'Knowledge Goblin' },
  title_backlog_destroyer: { name: 'Backlog Destroyer' },
  showcase_rubber_duck: { name: 'Senior Rubber Duck', preview: '🦆' },
  showcase_cactus: { name: 'Кактус прокрастинации', preview: '🌵' },
  showcase_cat: { name: 'Кот', preview: '🐈' },
  showcase_golden_duck: { name: 'Golden Duck', preview: '👑' },
}

export function cosmeticName(id?: string) {
  return id ? (cosmetics[id]?.name ?? id) : 'Не выбрано'
}

export function cosmeticPreview(id?: string) {
  return id ? (cosmetics[id]?.preview ?? '🙂') : '🙂'
}
