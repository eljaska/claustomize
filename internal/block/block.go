package block

// Block is a single composable segment of the statusline.
//
// Snippet is shell code that runs inside the generated script. The variable
// $input holds the JSON Claude Code pipes to the statusline command on stdin.
// The snippet must echo the segment text on stdout (or nothing to omit).
type Block struct {
	ID          string
	Name        string
	Description string
	Snippet     string
}

// All returns the built-in blocks in the fixed display/render order.
func All() []Block {
	return []Block{
		{
			ID:          "model",
			Name:        "Model",
			Description: "Current Claude model display name",
			Snippet:     `printf '%s' "$(printf '%s' "$input" | jq -r '.model.display_name // empty')"`,
		},
		{
			ID:          "cwd",
			Name:        "Working directory",
			Description: "Basename of the current working directory",
			Snippet:     `printf '%s' "$(printf '%s' "$input" | jq -r '.workspace.current_dir // empty' | awk -F/ '{print $NF}')"`,
		},
		{
			ID:          "git_branch",
			Name:        "Git branch",
			Description: "Current git branch, if inside a repo",
			Snippet:     `branch=$(git -C "$(printf '%s' "$input" | jq -r '.workspace.current_dir // empty')" rev-parse --abbrev-ref HEAD 2>/dev/null); [ -n "$branch" ] && printf '%s' "$branch"`,
		},
		{
			ID:          "time",
			Name:        "Time",
			Description: "Current local time (HH:MM)",
			Snippet:     `printf '%s' "$(date +%H:%M)"`,
		},
	}
}

// ByID returns the block with the given ID, or nil if not found.
func ByID(id string) *Block {
	for _, b := range All() {
		if b.ID == id {
			return &b
		}
	}
	return nil
}
