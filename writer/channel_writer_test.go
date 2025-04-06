package writer_test

import (
	"os"
	"testing"

	"github.com/ibakaidov/tfrecord-go/example"
	"github.com/ibakaidov/tfrecord-go/reader"
	"github.com/ibakaidov/tfrecord-go/writer"
)

func TestTFRecordChannelWriter(t *testing.T) {
	path := "testdata/channel_writer.tfrecord"
	_ = os.MkdirAll("testdata", 0755)
	defer os.Remove(path)

	ch := make(chan *example.Example)

	go func() {
		defer close(ch)
		ch <- example.NewExample(
			func() (string, *example.Feature) {
				return example.NewBytesFeature("user", []byte("Anna"))
			},
			func() (string, *example.Feature) {
				return example.NewFloatFeature("score", 4.9)
			},
		)
	}()

	w, err := writer.NewTFRecordChannelWriter(path, ch, 4)
	if err != nil {
		t.Fatalf("failed to create writer: %v", err)
	}
	if err := w.Wait(); err != nil {
		t.Fatalf("write process failed: %v", err)
	}

	// Read back
	read, err := reader.ReadAll(path)
	if err != nil {
		t.Fatalf("failed to read back: %v", err)
	}
	if len(read) != 1 {
		t.Fatalf("expected 1 example, got %d", len(read))
	}

	got := string(read[0].GetFeatures().Feature["user"].GetBytesList().Value[0])
	if got != "Anna" {
		t.Errorf("expected user name 'Anna', got %s", got)
	}
}
