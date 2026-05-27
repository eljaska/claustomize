package block

import "testing"

func TestAll_HasStableIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, b := range All() {
		if b.ID == "" {
			t.Errorf("block %q has empty ID", b.Name)
		}
		if seen[b.ID] {
			t.Errorf("duplicate block ID: %s", b.ID)
		}
		seen[b.ID] = true
		if b.Snippet == "" {
			t.Errorf("block %s has empty snippet", b.ID)
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
