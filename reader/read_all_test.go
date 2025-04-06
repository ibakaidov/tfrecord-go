package reader_test

import (
	"os"
	"testing"

	"github.com/ibakaidov/tfrecord-go/example"
	"github.com/ibakaidov/tfrecord-go/reader"
	"github.com/ibakaidov/tfrecord-go/writer"
)

func TestReadAll(t *testing.T) {
	path := "testdata/read_all.tfrecord"
	_ = os.MkdirAll("testdata", 0755)
	defer os.Remove(path)

	// Prepare test data
	ex := example.NewExample(
		func() (string, *example.Feature) {
			return example.NewBytesFeature("city", []byte("Paris"))
		},
		func() (string, *example.Feature) {
			return example.NewFloatFeature("temp", 18.3)
		},
	)

	// Write the example
	err := writer.WriteArray(path, []*example.Example{ex}, 2)
	if err != nil {
		t.Fatalf("failed to write: %v", err)
	}

	// Read the example back
	result, err := reader.ReadAll(path)
	if err != nil {
		t.Fatalf("failed to read: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 example, got %d", len(result))
	}

	got := string(result[0].GetFeatures().Feature["city"].GetBytesList().Value[0])
	want := "Paris"
	if got != want {
		t.Errorf("expected city %q, got %q", want, got)
	}
}
