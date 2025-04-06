package writer

import (
	"bufio"
	"encoding/binary"

	"github.com/ibakaidov/tfrecord-go/example"
	"google.golang.org/protobuf/proto"
)

// writeRecord writes a single Example in TFRecord format to the stream.
func writeRecord(w *bufio.Writer, ex *example.Example) error {
	data, err := proto.Marshal(ex) // Serialize the Example to a byte array.
	if err != nil {
		return err
	}

	length := uint64(len(data)) // Get the length of the serialized data.
	lengthBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lengthBytes, length) // Encode the length as 8 bytes in little-endian format.

	lengthCRC := maskedCRC(lengthBytes) // Compute the masked CRC for the length.
	dataCRC := maskedCRC(data)          // Compute the masked CRC for the data.

	// Write the length of the data to the stream.
	if _, err := w.Write(lengthBytes); err != nil {
		return err
	}
	// Write the CRC for the length to the stream.
	if err := binary.Write(w, binary.LittleEndian, lengthCRC); err != nil {
		return err
	}
	// Write the serialized data to the stream.
	if _, err := w.Write(data); err != nil {
		return err
	}
	// Write the CRC for the data to the stream.
	if err := binary.Write(w, binary.LittleEndian, dataCRC); err != nil {
		return err
	}

	return nil
}
