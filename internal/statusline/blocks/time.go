package blocks

// Time shows the current local time.
var Time = Block{
	ID:          "time",
	Name:        "Time",
	Description: "Current local time",
	Styles: []Style{
		{ID: "hm", Name: "HH:MM", Snippet: `printf '%s' "$(date +%H:%M)"`},
		{ID: "hms", Name: "HH:MM:SS", Snippet: `printf '%s' "$(date +%H:%M:%S)"`},
	},
}
