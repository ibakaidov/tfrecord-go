package writer

import (
	"hash/crc32"
)

func maskCRC(crc uint32) uint32 {
	return ((crc >> 15) | (crc << 17)) + 0xa282ead8
}

func maskedCRC(data []byte) uint32 {
	table := crc32.MakeTable(crc32.Castagnoli)
	return maskCRC(crc32.Checksum(data, table))
}
