package reader_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/ibakaidov/tfrecord-go/example"
	"github.com/ibakaidov/tfrecord-go/reader"
	"github.com/ibakaidov/tfrecord-go/writer"
)

func TestReadChannel(t *testing.T) {
	path := "testdata/read_channel.tfrecord"
	_ = os.MkdirAll("testdata", 0755)
	defer os.Remove(path)

	// Write test file
	data := []*example.Example{
		example.NewExample(
			func() (string, *example.Feature) {
				return example.NewBytesFeature("user", []byte("Charlie"))
			},
		),
	}
	if err := writer.WriteArray(path, data, 2); err != nil {
		t.Fatalf("WriteArray failed: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	exCh, errCh := reader.ReadChannel(ctx, path, 2)

	var count int
	for ex := range exCh {
		if ex == nil {
			t.Error("received nil example")
		}
		count++
	}

	if count != 1 {
		t.Errorf("expected 1 example, got %d", count)
	}

	if err := <-errCh; err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
