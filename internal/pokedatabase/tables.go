package pokedatabase

type table struct {
	fields 		map[string]field
}

type field struct {
    dataType   string
}

func getTables() map[string]table {
	return map[string]table{
		"pokemon": {
			fields: map[string]field{
				"id": 			  {dataType: "integer primary key autoincrement",},
				"number":         { dataType: "integer unique" },
				"name":           { dataType: "text unique" },
				"gender":         { dataType: "integer" },
				"capture_rate":   { dataType: "integer" },
				"base_happiness": { dataType: "integer" },
				"is_baby":        { dataType: "boolean" },
				"is_legendary":   { dataType: "boolean" },
				"is_mythical":    { dataType: "boolean" },
				"is_shiny":       { dataType: "boolean" },
				"sprite_path":     { dataType: "text" },
				"caught":         { dataType: "boolean default 0" },
				"caught_at":   	  { dataType: "datetime" },
				"registered_at":   { dataType: "datetime DEFAULT current_timestamp not null" },
			},
		},
	}
}
