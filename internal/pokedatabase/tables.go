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
				"id": field{
					dataType: "integer primary key autoincrement",
				},
				"name": field{
					dataType: "text not null",
				},
				"caught": field{
					dataType: "boolean not null",
				},
				"caughtAt": field{
					dataType: "text not null",
				},
				"createdAt": field{
					dataType: "datetime default current_timestamp",
				},
			},
		},
	}
}
