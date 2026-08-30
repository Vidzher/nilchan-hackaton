package cosmetics

type ItemType string

const (
	ItemTypeAvatar   ItemType = "avatar"
	ItemTypeFrame    ItemType = "frame"
	ItemTypeTitle    ItemType = "title"
	ItemTypeShowcase ItemType = "showcase"
)

const (
	DefaultAvatarID = "default_avatar"
	DefaultFrameID  = "default_frame"
)

type Presentation struct {
	Emoji    string
	AssetKey string
	CSSClass string
}

type Item struct {
	ID           string
	Type         ItemType
	DisplayName  string
	Price        int64
	Free         bool
	Presentation Presentation
}

var catalog = []Item{
	{
		ID:          DefaultAvatarID,
		Type:        ItemTypeAvatar,
		DisplayName: "Default Avatar",
		Price:       0,
		Free:        true,
		Presentation: Presentation{
			Emoji:    "🙂",
			AssetKey: "avatar-default",
		},
	},
	{
		ID:          DefaultFrameID,
		Type:        ItemTypeFrame,
		DisplayName: "Default Frame",
		Price:       0,
		Free:        true,
		Presentation: Presentation{
			AssetKey: "frame-default",
			CSSClass: "frame-default",
		},
	},
	{
		ID:          "avatar_frog",
		Type:        ItemTypeAvatar,
		DisplayName: "Frog",
		Price:       15,
		Presentation: Presentation{
			Emoji:    "🐸",
			AssetKey: "avatar-frog",
		},
	},
	{
		ID:          "avatar_robot",
		Type:        ItemTypeAvatar,
		DisplayName: "Robot",
		Price:       25,
		Presentation: Presentation{
			Emoji:    "🤖",
			AssetKey: "avatar-robot",
		},
	},
	{
		ID:          "avatar_wizard",
		Type:        ItemTypeAvatar,
		DisplayName: "Wizard",
		Price:       40,
		Presentation: Presentation{
			Emoji:    "🧙",
			AssetKey: "avatar-wizard",
		},
	},
	{
		ID:          "frame_neon",
		Type:        ItemTypeFrame,
		DisplayName: "Neon",
		Price:       30,
		Presentation: Presentation{
			AssetKey: "frame-neon",
			CSSClass: "frame-neon",
		},
	},
	{
		ID:          "frame_fire",
		Type:        ItemTypeFrame,
		DisplayName: "Fire",
		Price:       50,
		Presentation: Presentation{
			AssetKey: "frame-fire",
			CSSClass: "frame-fire",
		},
	},
	{
		ID:          "frame_gold",
		Type:        ItemTypeFrame,
		DisplayName: "Gold",
		Price:       120,
		Presentation: Presentation{
			AssetKey: "frame-gold",
			CSSClass: "frame-gold",
		},
	},
	{
		ID:          "title_vkatun",
		Type:        ItemTypeTitle,
		DisplayName: "Вкатун",
		Price:       20,
		Presentation: Presentation{
			AssetKey: "title-vkatun",
		},
	},
	{
		ID:          "title_factory_worker",
		Type:        ItemTypeTitle,
		DisplayName: "Заводчанин",
		Price:       30,
		Presentation: Presentation{
			AssetKey: "title-factory-worker",
		},
	},
	{
		ID:          "title_1c_developer",
		Type:        ItemTypeTitle,
		DisplayName: "1Сник",
		Price:       35,
		Presentation: Presentation{
			AssetKey: "title-1c-developer",
		},
	},
	{
		ID:          "title_insider",
		Type:        ItemTypeTitle,
		DisplayName: "Инсайдер",
		Price:       40,
		Presentation: Presentation{
			AssetKey: "title-insider",
		},
	},
	{
		ID:          "title_vibe_coder",
		Type:        ItemTypeTitle,
		DisplayName: "Вайбкодер",
		Price:       50,
		Presentation: Presentation{
			AssetKey: "title-vibe-coder",
		},
	},
	{
		ID:          "title_top_1_pisa",
		Type:        ItemTypeTitle,
		DisplayName: "Топ-1 в писе",
		Price:       70,
		Presentation: Presentation{
			AssetKey: "title-top-1-pisa",
		},
	},
	{
		ID:          "title_pogromist",
		Type:        ItemTypeTitle,
		DisplayName: "Погромист",
		Price:       90,
		Presentation: Presentation{
			AssetKey: "title-pogromist",
		},
	},
	{
		ID:          "title_middle",
		Type:        ItemTypeTitle,
		DisplayName: "Мидл",
		Price:       120,
		Presentation: Presentation{
			AssetKey: "title-middle",
		},
	},
	{
		ID:          "title_infogypsy",
		Type:        ItemTypeTitle,
		DisplayName: "Инфоцыган",
		Price:       150,
		Presentation: Presentation{
			AssetKey: "title-infogypsy",
		},
	},
	{
		ID:          "title_senior",
		Type:        ItemTypeTitle,
		DisplayName: "Синьор",
		Price:       200,
		Presentation: Presentation{
			AssetKey: "title-senior",
		},
	},
	{
		ID:          "showcase_prod_env",
		Type:        ItemTypeShowcase,
		DisplayName: ".env с прода",
		Price:       25,
		Presentation: Presentation{
			Emoji:    "🔐",
			AssetKey: "showcase-prod-env",
		},
	},
	{
		ID:          "showcase_green_ci",
		Type:        ItemTypeShowcase,
		DisplayName: "Зелёный CI",
		Price:       35,
		Presentation: Presentation{
			Emoji:    "✅",
			AssetKey: "showcase-green-ci",
		},
	},
	{
		ID:          "showcase_stackoverflow_2009",
		Type:        ItemTypeShowcase,
		DisplayName: "Ответ со Stack Overflow, 2009",
		Price:       50,
		Presentation: Presentation{
			Emoji:    "📜",
			AssetKey: "showcase-stackoverflow-2009",
		},
	},
	{
		ID:          "showcase_arch_sticker",
		Type:        ItemTypeShowcase,
		DisplayName: "Стикер «I use Arch btw»",
		Price:       70,
		Presentation: Presentation{
			Emoji:    "👨🏻‍💻",
			AssetKey: "showcase-arch-sticker",
		},
	},
	{
		ID:          "showcase_agents_md",
		Type:        ItemTypeShowcase,
		DisplayName: "AGENTS.md",
		Price:       90,
		Presentation: Presentation{
			Emoji:    "📄",
			AssetKey: "showcase-agents-md",
		},
	},
	{
		ID:          "showcase_monolith",
		Type:        ItemTypeShowcase,
		DisplayName: "Монолит",
		Price:       120,
		Presentation: Presentation{
			Emoji:    "🗿",
			AssetKey: "showcase-monolith",
		},
	},
	{
		ID:          "showcase_pantela_course",
		Type:        ItemTypeShowcase,
		DisplayName: "Курс Евгения Пантелы",
		Price:       150,
		Presentation: Presentation{
			Emoji:    "🤡",
			AssetKey: "showcase-pantela-course",
		},
	},
	{
		ID:          "showcase_winrar_license",
		Type:        ItemTypeShowcase,
		DisplayName: "Лицензия WinRAR",
		Price:       200,
		Presentation: Presentation{
			Emoji:    "📚",
			AssetKey: "showcase-winrar-license",
		},
	},
	{
		ID:          "showcase_general_watch",
		Type:        ItemTypeShowcase,
		DisplayName: "Котлы генеральские",
		Price:       250,
		Presentation: Presentation{
			Emoji:    "⌚",
			AssetKey: "showcase-general-watch",
		},
	},
}

var validTypes = map[ItemType]struct{}{
	ItemTypeAvatar:   {},
	ItemTypeFrame:    {},
	ItemTypeTitle:    {},
	ItemTypeShowcase: {},
}

func All() []Item {
	result := make([]Item, len(catalog))
	copy(result, catalog)
	return result
}

func GetByID(id string) (Item, bool) {
	for _, item := range catalog {
		if item.ID == id {
			return item, true
		}
	}
	return Item{}, false
}

func IsValidType(itemType ItemType) bool {
	_, ok := validTypes[itemType]
	return ok
}
