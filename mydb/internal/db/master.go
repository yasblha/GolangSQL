package db

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func ReadMasterTable(file *os.File, info DBInfo) []Table {
	page := ReadPage(file, 1, info.PageSize)

	var tables []Table

	if bytes.Contains(page, []byte("CREATE TABLE")) {
		start := bytes.Index(page, []byte("CREATE TABLE"))
		if start == -1 {
			return tables
		}

		// Essaye de lire jusqu’à un caractère de fin, mais sans sortir de la page
		end := bytes.IndexByte(page[start:], 0x00)
		if end == -1 || start+end > len(page) {
			end = len(page) - start
		}

		sql := string(page[start : start+end])
		fmt.Println("Found SQL:", sql)

		name := extractTableName(sql)
		cols := ParseCreateSQL(sql)

		tables = append(tables, Table{
			Name:    name,
			Columns: cols,
		})
	}

	return tables
}

func extractTableName(sql string) string {
	parts := strings.Split(sql, " ")
	for i, p := range parts {
		if strings.ToUpper(p) == "TABLE" && i+1 < len(parts) {
			return strings.TrimSpace(parts[i+1])
		}
	}
	return "unknown"
}
