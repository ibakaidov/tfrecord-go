package writer_test

import (
	"os"
	"testing"

	"github.com/ibakaidov/tfrecord-go/example"
	"github.com/ibakaidov/tfrecord-go/reader"
	"github.com/ibakaidov/tfrecord-go/writer"
)

func TestWriteArray(t *testing.T) {
	path := "testdata/write_array.tfrecord"
	_ = os.MkdirAll("testdata", 0755)
	defer os.Remove(path)

	// Примеры
	data := []*example.Example{
		example.NewExample(
			func() (string, *example.Feature) {
				return example.NewBytesFeature("user", []byte("Alice"))
			},
			func() (string, *example.Feature) {
				return example.NewFloatFeature("score", 3.14)
			},
		),
		example.NewExample(
			func() (string, *example.Feature) {
				return example.NewBytesFeature("user", []byte("Bob"))
			},
			func() (string, *example.Feature) {
				return example.NewFloatFeature("score", 4.2)
			},
		),
	}

	err := writer.WriteArray(path, data, 4)
	if err != nil {
		t.Fatalf("WriteArray failed: %v", err)
	}

	readData, err := reader.ReadAll(path)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}

	if len(readData) != len(data) {
		t.Errorf("expected %d examples, got %d", len(data), len(readData))
	}

	want := "Alice"
	got := string(readData[0].GetFeatures().Feature["user"].GetBytesList().Value[0])
	if got != want {
		t.Errorf("expected user %q, got %q", want, got)
	}
}
