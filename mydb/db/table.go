package db

type Column struct {
	Name string
	Type string
}

type Table struct {
	Name    string
	Columns []Column
}
