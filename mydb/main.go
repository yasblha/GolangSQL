package main

import (
	"fmt"
	"mydb/db"
	"os"
)

func main() {
	file, err := os.Open("sample.db") // <-- place ici un fichier SQLite valide
	if err != nil {
		panic(err)
	}
	defer file.Close()

	info := db.ParseHeader(file)
	fmt.Printf("Page size: %d bytes\n", info.PageSize)

	tables := db.ReadMasterTable(file, info)
	for _, t := range tables {
		fmt.Printf("Table: %s\n", t.Name)
		for _, col := range t.Columns {
			fmt.Printf(" - %s (%s)\n", col.Name, col.Type)
		}
	}
}
