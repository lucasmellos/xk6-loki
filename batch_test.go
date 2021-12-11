package loki

import (
	"context"
	"testing"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	"go.k6.io/k6/lib"
	"go.k6.io/k6/stats"
)

func BenchmarkNewBatch(b *testing.B) {
	samples := make(chan stats.SampleContainer)
	state := &lib.State{
		Samples: samples,
		VUID:    15,
	}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = lib.WithState(ctx, state)

	defer cancel()
	defer close(samples)
	go func() { // this is so that we read the send samples
		for range samples {
		}
	}()
	faker := gofakeit.New(12345)
	cardinalities := map[string]int{
		"app":       5,
		"namespace": 10,
		"pod":       100,
	}
	streams, minBatchSize, maxBatchSize := 5, 500, 1000
	labels := newLabelPool(faker, cardinalities)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = newBatch(ctx, labels, streams, minBatchSize, maxBatchSize)
	}
}

func BenchmarkEncode(b *testing.B) {
	samples := make(chan stats.SampleContainer)
	state := &lib.State{
		Samples: samples,
		VUID:    15,
	}
	ctx, cancel := context.WithCancel(context.Background())
	ctx = lib.WithState(ctx, state)

	defer cancel()
	defer close(samples)
	go func() { // this is so that we read the send samples
		for range samples {
		}
	}()
	faker := gofakeit.New(12345)
	cardinalities := map[string]int{
		"app":       5,
		"namespace": 10,
		"pod":       100,
	}
	streams, minBatchSize, maxBatchSize := 5, 500, 1000
	labels := newLabelPool(faker, cardinalities)

	b.ReportAllocs()
	batch := newBatch(ctx, labels, streams, minBatchSize, maxBatchSize)

	b.Run("encode protobuf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = batch.createPushRequest()
		}
	})

	b.Run("encode json", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = batch.createJSONPushRequest()
		}
	})
}
