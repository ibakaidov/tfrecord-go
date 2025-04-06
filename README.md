# tfrecord-go

**tfrecord-go** is a minimal, idiomatic Go library for working with [TensorFlow's TFRecord format](https://www.tensorflow.org/tutorials/load_data/tfrecord).  
It supports reading and writing serialized `Example` protos with full CRC32C checks.

---

## ✨ Features

- ✅ Write and read TFRecord files with `Example` protobuf messages
- ✅ Streaming support via channels and `context.Context`
- ✅ Compatible with Python `tf.train.Example`
- ✅ Zero dependencies beyond `protobuf`

---

## 📦 Installation

```bash
go get github.com/ibakaidov/tfrecord-go
```

---

## 📝 Writing Examples

```go
import (
	"github.com/ibakaidov/tfrecord-go/example"
	"github.com/ibakaidov/tfrecord-go/writer"
)

ex := example.NewExample(
	func() (string, *example.Feature) {
		return example.NewBytesFeature("user", []byte("Alice"))
	},
	func() (string, *example.Feature) {
		return example.NewFloatFeature("score", 4.95)
	},
)

err := writer.WriteArray("output.tfrecord", []*example.Example{ex}, 4)
```

### Streaming via channel

```go
ch := make(chan *example.Example)

w, _ := writer.NewTFRecordChannelWriter("stream.tfrecord", ch, 4)
go func() {
	defer close(ch)
	ch <- ex
}()
_ = w.Wait()
```

---

## 📖 Reading Examples

```go
import (
	"github.com/ibakaidov/tfrecord-go/reader"
)

examples, err := reader.ReadAll("output.tfrecord")
for _, ex := range examples {
	fmt.Println(ex.GetFeatures().Feature["user"])
}
```

### Streaming via channel

```go
ctx := context.Background()
ch, errCh := reader.ReadChannel(ctx, "stream.tfrecord", 4)

for ex := range ch {
	fmt.Println(ex)
}
if err := <-errCh; err != nil {
	log.Fatal(err)
}
```

---

## 🧪 Interoperability with Python

```python
import tensorflow as tf

for r in tf.data.TFRecordDataset("output.tfrecord"):
    ex = tf.train.Example()
    ex.ParseFromString(r.numpy())
    print(ex)
```

---

## 📚 Docs

Go reference: [pkg.go.dev/github.com/ibakaidov/tfrecord-go](https://pkg.go.dev/github.com/ibakaidov/tfrecord-go)

---

## 🛠️ License

MIT — feel free to use in personal or commercial projects.
