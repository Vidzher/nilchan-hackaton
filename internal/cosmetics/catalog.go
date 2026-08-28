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
		ID:          "title_razgreshatel",
		Type:        ItemTypeTitle,
		DisplayName: "Разгребатель",
		Price:       20,
		Presentation: Presentation{
			AssetKey: "title-razgreshatel",
		},
	},
	{
		ID:          "title_knowledge_goblin",
		Type:        ItemTypeTitle,
		DisplayName: "Knowledge Goblin",
		Price:       40,
		Presentation: Presentation{
			AssetKey: "title-knowledge-goblin",
		},
	},
	{
		ID:          "title_backlog_destroyer",
		Type:        ItemTypeTitle,
		DisplayName: "Backlog Destroyer",
		Price:       70,
		Presentation: Presentation{
			AssetKey: "title-backlog-destroyer",
		},
	},
	{
		ID:          "showcase_rubber_duck",
		Type:        ItemTypeShowcase,
		DisplayName: "Senior Rubber Duck",
		Price:       40,
		Presentation: Presentation{
			Emoji:    "🦆",
			AssetKey: "showcase-rubber-duck",
		},
	},
	{
		ID:          "showcase_cactus",
		Type:        ItemTypeShowcase,
		DisplayName: "Кактус прокрастинации",
		Price:       55,
		Presentation: Presentation{
			Emoji:    "🌵",
			AssetKey: "showcase-cactus",
		},
	},
	{
		ID:          "showcase_cat",
		Type:        ItemTypeShowcase,
		DisplayName: "Кот",
		Price:       120,
		Presentation: Presentation{
			Emoji:    "🐈",
			AssetKey: "showcase-cat",
		},
	},
	{
		ID:          "showcase_golden_duck",
		Type:        ItemTypeShowcase,
		DisplayName: "Golden Duck",
		Price:       250,
		Presentation: Presentation{
			Emoji:    "👑",
			AssetKey: "showcase-golden-duck",
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
