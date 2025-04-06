// Package reader provides utilities to read Examples from TFRecord files.
package reader

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"

	"github.com/ibakaidov/tfrecord-go/example"
	"google.golang.org/protobuf/proto"
)

// ReadAll reads all serialized Examples from a TFRecord file.
//
// Parameters:
//   - path: path to the TFRecord file.
//
// Returns:
//   - []*example.Example: slice of all successfully read Examples.
//   - error: if file access, CRC check, or decoding fails.
//
// Example:
//
//	records, err := reader.ReadAll("file.tfrecord")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, ex := range records {
//	    fmt.Println(ex)
//	}
func ReadAll(path string) ([]*example.Example, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	var results []*example.Example

	for {
		ex, err := readRecord(reader)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, ex)
	}

	return results, nil
}

// readRecord reads and deserializes a single Example from the given stream.
//
// It expects the binary layout of a TFRecord block:
// [uint64 length][uint32 masked crc of length][byte[length]][uint32 masked crc of data].
func readRecord(r *bufio.Reader) (*example.Example, error) {
	lengthBytes := make([]byte, 8)
	if _, err := io.ReadFull(r, lengthBytes); err != nil {
		return nil, err
	}
	length := binary.LittleEndian.Uint64(lengthBytes)

	// Validate length CRC
	lengthCRC := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthCRC); err != nil {
		return nil, err
	}
	if maskedCRC(lengthBytes) != binary.LittleEndian.Uint32(lengthCRC) {
		return nil, errors.New("length CRC mismatch")
	}

	// Read and validate data
	data := make([]byte, length)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}

	dataCRC := make([]byte, 4)
	if _, err := io.ReadFull(r, dataCRC); err != nil {
		return nil, err
	}
	if maskedCRC(data) != binary.LittleEndian.Uint32(dataCRC) {
		return nil, errors.New("data CRC mismatch")
	}

	var ex example.Example
	if err := proto.Unmarshal(data, &ex); err != nil {
		return nil, err
	}

	return &ex, nil
}
