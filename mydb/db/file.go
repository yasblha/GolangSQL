package db

import "os"

func ReadBytes(file *os.File, offset int64, size int) []byte {
	buf := make([]byte, size)
	_, err := file.ReadAt(buf, offset)
	if err != nil {
		panic(err)
	}
	return buf
}
