// Package reader provides functionality to read TFRecord files.
package reader

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"

	"github.com/ibakaidov/tfrecord-go/example"
)

// ReadChannel reads a TFRecord file in a streaming fashion using a channel.
//
// It launches a background goroutine that reads Examples from the file at `path`
// and sends them to the returned channel. The process is cancellable using `ctx`.
//
// Parameters:
//   - ctx: context for cancellation (useful for timeouts or manual stop)
//   - path: path to the TFRecord file
//   - bufferSizeMB: size of the internal read buffer in megabytes
//
// Returns:
//   - <-chan *example.Example: channel through which decoded examples are streamed
//   - <-chan error: channel that will emit a single error (or nil) and then close
//
// Example:
//
//	ctx := context.Background()
//	exCh, errCh := reader.ReadChannel(ctx, "data.tfrecord", 4)
//	for ex := range exCh {
//	    fmt.Println(ex)
//	}
//	if err := <-errCh; err != nil {
//	    log.Fatal(err)
//	}
func ReadChannel(ctx context.Context, path string, bufferSizeMB int) (<-chan *example.Example, <-chan error) {
	out := make(chan *example.Example)
	errCh := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errCh)

		file, err := os.Open(path)
		if err != nil {
			errCh <- err
			return
		}
		defer file.Close()

		reader := bufio.NewReaderSize(file, bufferSizeMB*1024*1024)

		for {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			default:
				ex, err := readRecord(reader)
				if errors.Is(err, io.EOF) {
					return
				}
				if err != nil {
					errCh <- err
					return
				}
				out <- ex
			}
		}
	}()

	return out, errCh
}
