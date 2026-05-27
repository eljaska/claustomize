package slot

import "github.com/eljaska/claustomize/internal/statusline/blocks"

// Slot is one position in the statusline. An empty slot has Block == nil.
// A filled slot points at a Block and selects one of its styles.
type Slot struct {
	Block    *blocks.Block
	StyleIdx int
}

// IsEmpty reports whether this slot has no block.
func (s Slot) IsEmpty() bool { return s.Block == nil }

// List is an ordered sequence of slots that always satisfies the layout
// invariant: empty slots at both ends, and exactly one empty slot between
// every pair of filled slots. Equivalently, slots strictly alternate
// empty/filled/empty/... starting and ending with empty.
type List []Slot

// New returns the initial list: a single empty slot.
func New() List {
	return List{{}}
}

// Fill replaces the empty slot at i with [empty, filled, empty], preserving
// the invariant. Returns the updated list and the new cursor index (the
// newly-filled slot). Panics if l[i] is not empty.
func (l List) Fill(i int, b *blocks.Block, styleIdx int) (List, int) {
	if !l[i].IsEmpty() {
		panic("slot.Fill: target slot is not empty")
	}
	out := make(List, 0, len(l)+2)
	out = append(out, l[:i]...)
	out = append(out, Slot{})
	out = append(out, Slot{Block: b, StyleIdx: styleIdx})
	out = append(out, Slot{})
	out = append(out, l[i+1:]...)
	return out, i + 1
}

// Empty replaces [empty, filled, empty] at i-1..i+1 with a single empty,
// preserving the invariant. Returns the updated list and the new cursor
// index (the collapsed empty). Panics if l[i] is not filled.
func (l List) Empty(i int) (List, int) {
	if l[i].IsEmpty() {
		panic("slot.Empty: target slot is not filled")
	}
	out := make(List, 0, len(l)-2)
	out = append(out, l[:i-1]...)
	out = append(out, Slot{})
	out = append(out, l[i+2:]...)
	return out, i - 1
}

// Replace swaps the filled slot at i for a different block and/or style.
// Returns the updated list; cursor index is unchanged. Panics if l[i] is
// empty (use Fill instead).
func (l List) Replace(i int, b *blocks.Block, styleIdx int) List {
	if l[i].IsEmpty() {
		panic("slot.Replace: target slot is empty; use Fill")
	}
	out := make(List, len(l))
	copy(out, l)
	out[i] = Slot{Block: b, StyleIdx: styleIdx}
	return out
}
