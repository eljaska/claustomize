package blocks

// Pipe is a pipe-character separator block.
var Pipe = Block{
	ID:          "pipe",
	Name:        "Pipe separator",
	Description: "A pipe separator",
	Styles: []Style{
		{ID: "padded", Name: "' | '", Snippet: `printf ' | '`},
		{ID: "tight", Name: "'|'", Snippet: `printf '|'`},
	},
}
