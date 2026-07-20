package mapreduce

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestWordCountBasic(t *testing.T) {
	docs := []string{
		"the cat sat",
		"the dog sat",
		"cat dog cat",
	}
	got := WordCount(context.Background(), docs, 4)
	want := map[string]int{"the": 2, "cat": 3, "sat": 2, "dog": 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordCount = %v, want %v", got, want)
	}
}

func TestTokenizationLowercasesAndSplitsPunctuation(t *testing.T) {
	docs := []string{
		"Hello, world!",
		"hello HELLO   world",
	}
	got := WordCount(context.Background(), docs, 2)
	want := map[string]int{"hello": 3, "world": 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordCount = %v, want %v", got, want)
	}
}

func TestEmptyInput(t *testing.T) {
	got := WordCount(context.Background(), nil, 4)
	if got == nil {
		t.Fatal("WordCount(nil) = nil map, want an empty (non-nil) map")
	}
	if len(got) != 0 {
		t.Errorf("WordCount(nil) = %v, want empty", got)
	}
}

// bigCorpus builds a deterministic, sizeable set of documents.
func bigCorpus() []string {
	words := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	docs := make([]string, 0, 500)
	for i := 0; i < 500; i++ {
		// Each doc repeats one word (i%5)+1 times plus a constant "zeta".
		w := words[i%len(words)]
		line := "zeta"
		for j := 0; j <= i%5; j++ {
			line += " " + w
		}
		docs = append(docs, line)
	}
	return docs
}

func TestDeterministicAcrossWorkerCounts(t *testing.T) {
	docs := bigCorpus()
	base := WordCount(context.Background(), docs, 1)
	if len(base) == 0 {
		t.Fatal("baseline count is empty")
	}
	for _, w := range []int{2, 3, 8, 16, 64} {
		got := WordCount(context.Background(), docs, w)
		if !reflect.DeepEqual(got, base) {
			t.Errorf("workers=%d gave a different result than workers=1", w)
		}
	}
}

func TestWorkersLessThanOneTreatedAsOne(t *testing.T) {
	docs := []string{"a a b", "b c c c"}
	want := WordCount(context.Background(), docs, 1)
	for _, w := range []int{0, -5} {
		got := WordCount(context.Background(), docs, w)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("workers=%d = %v, want %v", w, got, want)
		}
	}
}

func TestCancellationReturnsPromptly(t *testing.T) {
	// A large corpus so a non-cancelling implementation would take real work.
	docs := make([]string, 100_000)
	for i := range docs {
		docs[i] = fmt.Sprintf("word%d common", i)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled before we start

	done := make(chan map[string]int, 1)
	go func() { done <- WordCount(ctx, docs, 4) }()

	select {
	case <-done: // returned without processing everything — good
	case <-time.After(3 * time.Second):
		t.Fatal("WordCount did not honor context cancellation")
	}
}

func TestConcurrentStress(t *testing.T) {
	docs := bigCorpus()
	// Run many times with varying worker counts under -race to shake out data
	// races in the fan-out/fan-in wiring.
	base := WordCount(context.Background(), docs, 1)
	for i := 0; i < 20; i++ {
		got := WordCount(context.Background(), docs, 8)
		if !reflect.DeepEqual(got, base) {
			t.Fatalf("iteration %d differs from baseline", i)
		}
	}
}
