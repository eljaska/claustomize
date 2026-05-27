package blocks

// CWD shows the current working directory in one of several styles.
var CWD = Block{
	ID:          "cwd",
	Name:        "Working directory",
	Description: "Current working directory",
	Styles: []Style{
		{
			ID:      "basename",
			Name:    "basename",
			Snippet: `printf '%s' "$(printf '%s' "$input" | jq -r '.workspace.current_dir // empty' | awk -F/ '{print $NF}')"`,
		},
		{
			ID:      "tilde",
			Name:    "~-relative",
			Snippet: `path=$(printf '%s' "$input" | jq -r '.workspace.current_dir // empty'); printf '%s' "${path/#$HOME/~}"`,
		},
		{
			ID:      "full",
			Name:    "full path",
			Snippet: `printf '%s' "$(printf '%s' "$input" | jq -r '.workspace.current_dir // empty')"`,
		},
	},
}
