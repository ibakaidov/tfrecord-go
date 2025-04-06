package reader

import (
	"hash/crc32"
)

// maskedCRC возвращает masked CRC32C, как требует формат TFRecord.
func maskedCRC(data []byte) uint32 {
	return maskCRC(crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli)))
}

// maskCRC применяет маску к CRC по спецификации TFRecord.
func maskCRC(crc uint32) uint32 {
	return ((crc >> 15) | (crc << 17)) + 0xa282ead8
}
