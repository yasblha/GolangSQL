package db

import "os"

func ReadPage(file *os.File, pageNum int, pageSize uint16) []byte {
	offset := int64((pageNum - 1) * int(pageSize))
	return ReadBytes(file, offset, int(pageSize))
}
