package slot

import (
	"testing"

	"github.com/eljaska/claustomize/internal/statusline/blocks"
)

func TestNew_SingleEmpty(t *testing.T) {
	l := New()
	if len(l) != 1 || !l[0].IsEmpty() {
		t.Errorf("expected single empty slot, got %v", l)
	}
}

func assertInvariant(t *testing.T, l List) {
	t.Helper()
	if len(l) == 0 {
		t.Fatal("list is empty (must have at least one slot)")
	}
	if len(l)%2 == 0 {
		t.Fatalf("expected odd length, got %d", len(l))
	}
	for i, s := range l {
		wantEmpty := i%2 == 0
		if s.IsEmpty() != wantEmpty {
			t.Errorf("slot %d: expected empty=%v, got empty=%v", i, wantEmpty, s.IsEmpty())
		}
	}
}

func TestFill_PreservesInvariant(t *testing.T) {
	b := blocks.ByID("model")
	l := New()

	l, cursor := l.Fill(0, b, 0)
	assertInvariant(t, l)
	if cursor != 1 {
		t.Errorf("expected cursor at 1, got %d", cursor)
	}
	if len(l) != 3 {
		t.Errorf("expected 3 slots, got %d", len(l))
	}

	// Fill the rightmost empty.
	l, cursor = l.Fill(2, b, 0)
	assertInvariant(t, l)
	if cursor != 3 {
		t.Errorf("expected cursor at 3, got %d", cursor)
	}
	if len(l) != 5 {
		t.Errorf("expected 5 slots, got %d", len(l))
	}
}

func TestEmpty_PreservesInvariant(t *testing.T) {
	b := blocks.ByID("model")
	l := New()
	l, _ = l.Fill(0, b, 0) // [E, F, E]
	l, cursor := l.Empty(1)
	assertInvariant(t, l)
	if cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", cursor)
	}
	if len(l) != 1 {
		t.Errorf("expected single empty slot, got %d", len(l))
	}
}

func TestEmpty_CollapsesNeighborsCorrectly(t *testing.T) {
	b := blocks.ByID("model")
	l := New()
	l, _ = l.Fill(0, b, 0) // [E, F, E]
	l, _ = l.Fill(2, b, 0) // [E, F, E, F, E]
	l, cursor := l.Empty(1)
	assertInvariant(t, l)
	if cursor != 0 {
		t.Errorf("expected cursor at 0, got %d", cursor)
	}
	if len(l) != 3 {
		t.Errorf("expected 3 slots after emptying first, got %d", len(l))
	}
	if l[1].Block == nil || l[1].Block.ID != "model" {
		t.Errorf("remaining filled slot should still be model, got %+v", l[1])
	}
}

func TestReplace_KeepsLayout(t *testing.T) {
	b1 := blocks.ByID("model")
	b2 := blocks.ByID("cwd")
	l := New()
	l, cursor := l.Fill(0, b1, 0)

	l2 := l.Replace(cursor, b2, 1)
	assertInvariant(t, l2)
	if len(l2) != len(l) {
		t.Errorf("Replace changed length: %d -> %d", len(l), len(l2))
	}
	if l2[cursor].Block.ID != "cwd" || l2[cursor].StyleIdx != 1 {
		t.Errorf("replace did not stick: %+v", l2[cursor])
	}
}

func TestFill_PanicsOnFilledTarget(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	b := blocks.ByID("model")
	l := New()
	l, _ = l.Fill(0, b, 0)
	_, _ = l.Fill(1, b, 0) // index 1 is filled, must panic
}

func TestEmpty_PanicsOnEmptyTarget(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	l := New()
	_, _ = l.Empty(0) // index 0 is empty, must panic
}
