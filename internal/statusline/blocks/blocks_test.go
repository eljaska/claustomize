package blocks

import "testing"

func TestAll_BlocksAreWellFormed(t *testing.T) {
	seen := map[string]bool{}
	for _, b := range All() {
		if b.ID == "" {
			t.Errorf("block %q has empty ID", b.Name)
		}
		if seen[b.ID] {
			t.Errorf("duplicate block ID: %s", b.ID)
		}
		seen[b.ID] = true
		if len(b.Styles) == 0 {
			t.Errorf("block %s has no styles", b.ID)
		}
		styleIDs := map[string]bool{}
		for _, s := range b.Styles {
			if s.ID == "" {
				t.Errorf("block %s has style with empty ID", b.ID)
			}
			if styleIDs[s.ID] {
				t.Errorf("block %s has duplicate style ID: %s", b.ID, s.ID)
			}
			styleIDs[s.ID] = true
			if s.Snippet == "" {
				t.Errorf("block %s style %s has empty snippet", b.ID, s.ID)
			}
		}
	}
}

func TestByID(t *testing.T) {
	if ByID("model") == nil {
		t.Error("expected to find 'model' block")
	}
	if ByID("does-not-exist") != nil {
		t.Error("expected nil for unknown block")
	}
}
