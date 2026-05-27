package blocks

// GitBranch shows the current git branch in the working directory.
var GitBranch = Block{
	ID:          "git_branch",
	Name:        "Git branch",
	Description: "Current git branch in the working directory",
	Styles: []Style{
		{
			ID:      "name",
			Name:    "branch name",
			Snippet: `branch=$(git -C "$(printf '%s' "$input" | jq -r '.workspace.current_dir // empty')" rev-parse --abbrev-ref HEAD 2>/dev/null); [ -n "$branch" ] && printf '%s' "$branch"`,
		},
		{
			ID:   "name_dirty",
			Name: "branch + dirty marker",
			Snippet: `dir=$(printf '%s' "$input" | jq -r '.workspace.current_dir // empty'); ` +
				`branch=$(git -C "$dir" rev-parse --abbrev-ref HEAD 2>/dev/null); ` +
				`if [ -n "$branch" ]; then ` +
				`  if [ -n "$(git -C "$dir" status --porcelain 2>/dev/null)" ]; then printf '%s*' "$branch"; else printf '%s' "$branch"; fi; ` +
				`fi`,
		},
	},
}
