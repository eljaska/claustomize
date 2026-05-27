// Package blocks defines the built-in statusline segments and their
// rendering styles. Each block lives in its own file; this file holds the
// shared types and the registry that pulls them together.
package blocks

// Style is one rendering variant of a Block. Snippet is shell code that
// runs inside the generated statusline script; $input holds the JSON
// Claude Code pipes to the statusline command on stdin. The snippet
// echoes the segment text (or nothing to omit the segment).
type Style struct {
	ID      string
	Name    string
	Snippet string
}

// Block is a composable statusline segment with one or more styles.
type Block struct {
	ID          string
	Name        string
	Description string
	Styles      []Style
}

// All returns the built-in blocks in palette order.
func All() []Block {
	return []Block{
		Model,
		CWD,
		GitBranch,
		Time,
		Pipe,
		Dot,
		Space,
	}
}

// ByID returns the block with the given ID, or nil if not found.
func ByID(id string) *Block {
	all := All()
	for i := range all {
		if all[i].ID == id {
			return &all[i]
		}
	}
	return nil
}
