package blocks

// Model shows the currently active Claude model.
var Model = Block{
	ID:          "model",
	Name:        "Model",
	Description: "Current Claude model display name",
	Styles: []Style{
		{
			ID:      "display",
			Name:    "display name",
			Snippet: `printf '%s' "$(printf '%s' "$input" | jq -r '.model.display_name // empty')"`,
		},
	},
}
