package blocks

// Space is a single literal space block.
var Space = Block{
	ID:          "space",
	Name:        "Space",
	Description: "A single space",
	Styles: []Style{
		{ID: "one", Name: "' '", Snippet: `printf ' '`},
	},
}
