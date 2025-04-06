// Package writer provides utilities for writing data in the TFRecord format.
package writer

import (
	"bufio"
	"os"
	"sync"

	"github.com/ibakaidov/tfrecord-go/example"
)

// TFRecordChannelWriter is a streaming TFRecord writer that consumes examples from a channel.
type TFRecordChannelWriter struct {
	wg  sync.WaitGroup
	err error
}

// NewTFRecordChannelWriter creates a TFRecord writer that reads examples from a channel and writes
// them to a TFRecord file. The writing process is performed in a separate goroutine.
//
// Parameters:
//   - path: the output TFRecord file path.
//   - ch: a channel from which *example.Example instances are received.
//   - bufferSizeMB: the internal write buffer size in megabytes.
//
// Returns:
//   - *TFRecordChannelWriter: an object to wait for the write process.
//   - error: if file creation fails.
//
// Usage:
//
//	ch := make(chan *example.Example)
//	writer, err := writer.NewTFRecordChannelWriter("file.tfrecord", ch, 4)
//	if err != nil { log.Fatal(err) }
//	go func() {
//	    defer close(ch)
//	    ch <- example.NewExample(...)
//	}()
//	if err := writer.Wait(); err != nil {
//	    log.Fatal(err)
//	}
func NewTFRecordChannelWriter(path string, ch <-chan *example.Example, bufferSizeMB int) (*TFRecordChannelWriter, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	bufferSize := bufferSizeMB * 1024 * 1024
	writer := bufio.NewWriterSize(file, bufferSize)

	w := &TFRecordChannelWriter{}
	w.wg.Add(1)

	go func() {
		defer w.wg.Done()
		defer file.Close()

		for ex := range ch {
			if err := writeRecord(writer, ex); err != nil {
				w.err = err
				return
			}
		}

		if err := writer.Flush(); err != nil {
			w.err = err
		}
	}()

	return w, nil
}

// Wait blocks until the background writing process finishes.
// It returns any error that occurred during the writing.
func (w *TFRecordChannelWriter) Wait() error {
	w.wg.Wait()
	return w.err
}
