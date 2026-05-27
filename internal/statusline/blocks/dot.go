package blocks

// Dot is a middle-dot separator block.
var Dot = Block{
	ID:          "dot",
	Name:        "Dot separator",
	Description: "A middle-dot separator",
	Styles: []Style{
		{ID: "padded", Name: "' · '", Snippet: `printf ' · '`},
		{ID: "tight", Name: "'·'", Snippet: `printf '·'`},
	},
}
