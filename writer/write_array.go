// Package writer provides utilities for writing data in the TFRecord format.
package writer

import (
	"bufio"
	"os"

	"github.com/ibakaidov/tfrecord-go/example"
)

// WriteArray writes a slice of Examples to a TFRecord file.
//
// Parameters:
//   - path: the output file path.
//   - data: a slice of pointers to Examples that will be serialized and written.
//   - bufferSizeMB: size of the internal write buffer in megabytes.
//
// Returns an error if file creation, buffering, or record writing fails.
//
// Example:
//
//	examples := []*example.Example{
//		example.NewExample(
//			func() (string, *example.Feature) {
//				return example.NewBytesFeature("name", []byte("Alice"))
//			},
//		),
//	}
//
//	err := writer.WriteArray("output.tfrecord", examples, 4)
//	if err != nil {
//		log.Fatal(err)
//	}
func WriteArray(path string, data []*example.Example, bufferSizeMB int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := bufio.NewWriterSize(file, bufferSizeMB*1024*1024)
	defer buffer.Flush()

	for _, ex := range data {
		if err := writeRecord(buffer, ex); err != nil {
			return err
		}
	}

	return nil
}
