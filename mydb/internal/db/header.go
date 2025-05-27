package db

import (
	"encoding/binary"
	"os"
)

type DBInfo struct {
	PageSize uint16
}

func ParseHeader(file *os.File) DBInfo {
	header := ReadBytes(file, 0, 100)
	if string(header[:16]) != "SQLite format 3\x00" {
		panic("Not a valid SQLite file")
	}

	pageSize := binary.BigEndian.Uint16(header[16:18])
	return DBInfo{PageSize: pageSize}
}
