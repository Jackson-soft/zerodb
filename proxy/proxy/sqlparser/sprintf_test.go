package sqlparser

import (
	"fmt"
	"testing"
)

func BenchmarkFprintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fmt.Sprintf("%s %s", "hello", "world")
	}
}

func BenchmarkDddprintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		buf := NewTrackedBuffer(nil, nil)
		buf.Fprintf("%s %s", "hello", "world")
	}
}
