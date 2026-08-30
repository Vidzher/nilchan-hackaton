package cosmetics

import "testing"

func TestCatalog(t *testing.T) {
	items := All()
	if len(items) < 10 || len(items) > 20 {
		t.Fatalf("catalog size = %d, want 10..20", len(items))
	}
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if item.ID == "" {
			t.Error("item has empty ID")
		}
		if _, exists := seen[item.ID]; exists {
			t.Errorf("duplicate item ID %q", item.ID)
		}
		seen[item.ID] = struct{}{}
		if !IsValidType(item.Type) {
			t.Errorf("item %q has invalid type %q", item.ID, item.Type)
		}
		if item.Price < 0 {
			t.Errorf("item %q has negative price", item.ID)
		}
		if item.DisplayName == "" {
			t.Errorf("item %q has empty display name", item.ID)
		}
		if item.Free && item.Price != 0 {
			t.Errorf("free item %q must have zero price", item.ID)
		}
		if item.Presentation.Emoji == "" &&
			item.Presentation.AssetKey == "" &&
			item.Presentation.CSSClass == "" {
			t.Errorf("item %q has empty presentation data", item.ID)
		}
	}
}
func TestDefaultItems(t *testing.T) {
	avatar, ok := GetByID(DefaultAvatarID)
	if !ok {
		t.Fatal("default avatar is missing")
	}
	if avatar.Type != ItemTypeAvatar {
		t.Fatalf("default avatar has type %q", avatar.Type)
	}
	if !avatar.Free || avatar.Price != 0 {
		t.Fatal("default avatar must be free")
	}
	frame, ok := GetByID(DefaultFrameID)
	if !ok {
		t.Fatal("default frame is missing")
	}
	if frame.Type != ItemTypeFrame {
		t.Fatalf("default frame has type %q", frame.Type)
	}
	if !frame.Free || frame.Price != 0 {
		t.Fatal("default frame must be free")
	}
}
