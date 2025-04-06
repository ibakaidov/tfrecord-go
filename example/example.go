// Package example provides helper functions and type aliases
// for working with TensorFlow Example protocol buffers.
package example

import (
	proto "github.com/ibakaidov/tfrecord-go/proto/github.com/tensorflow/tensorflow/tensorflow/go/core/example/example_protos_go_proto"
)

// Type aliases for improved readability and shorter code.
type (
	Example           = proto.Example
	Feature           = proto.Feature
	Features          = proto.Features
	Feature_BytesList = proto.Feature_BytesList
	Feature_FloatList = proto.Feature_FloatList
	Feature_Int64List = proto.Feature_Int64List
	BytesList         = proto.BytesList
	FloatList         = proto.FloatList
	Int64List         = proto.Int64List
)

// NewBytesFeature creates a Feature containing a BytesList.
//
// Example:
//
//	key, feature := example.NewBytesFeature("user_id", []byte("123"))
func NewBytesFeature(key string, values ...[]byte) (string, *Feature) {
	return key, &Feature{
		Kind: &Feature_BytesList{
			BytesList: &BytesList{Value: values},
		},
	}
}

// NewFloatFeature creates a Feature containing a FloatList.
//
// Example:
//
//	key, feature := example.NewFloatFeature("confidence", 0.99)
func NewFloatFeature(key string, values ...float32) (string, *Feature) {
	return key, &Feature{
		Kind: &Feature_FloatList{
			FloatList: &FloatList{Value: values},
		},
	}
}

// NewIntFeature creates a Feature containing an Int64List.
//
// Example:
//
//	key, feature := example.NewIntFeature("id", 42)
func NewIntFeature(key string, values ...int64) (string, *Feature) {
	return key, &Feature{
		Kind: &Feature_Int64List{
			Int64List: &Int64List{Value: values},
		},
	}
}

// NewExample constructs a new Example from a list of functions
// that each return a key and a Feature.
//
// Example:
//
//	ex := example.NewExample(
//	    func() (string, *example.Feature) {
//	        return example.NewBytesFeature("name", []byte("Alice"))
//	    },
//	    func() (string, *example.Feature) {
//	        return example.NewFloatFeature("score", 4.7)
//	    },
//	)
//
// Returns:
//   - *Example ready to be written to a TFRecord.
func NewExample(featureFuncs ...func() (string, *Feature)) *Example {
	featureMap := make(map[string]*Feature)
	for _, f := range featureFuncs {
		k, v := f()
		featureMap[k] = v
	}
	return &Example{
		Features: &Features{
			Feature: featureMap,
		},
	}
}
