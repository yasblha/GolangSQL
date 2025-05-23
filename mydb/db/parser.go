package db

import (
	"regexp"
	"strings"
)

func ParseCreateSQL(sql string) []Column {
	re := regexp.MustCompile(`\((.*)\)`)
	match := re.FindStringSubmatch(sql)
	if len(match) < 2 {
		return nil
	}

	def := match[1]
	parts := strings.Split(def, ",")
	var cols []Column

	for _, part := range parts {
		tokens := strings.Fields(strings.TrimSpace(part))
		if len(tokens) >= 2 {
			cols = append(cols, Column{
				Name: tokens[0],
				Type: tokens[1],
			})
		}
	}

	return cols
}
