package util

import (
	"testing"
)

func TestSimpleArenaAllocator(t *testing.T) {
	arena := NewFastAllocator(1000)
	slice := arena.Alloc(10)
	if arena.off != 10 {
		t.Error("off not match, expect 10 bug got", arena.off)
	}

	if len(slice) != 0 || cap(slice) != 10 {
		t.Error("slice length or cap not match")
	}

	slice = arena.Alloc(20)
	if arena.off != 30 {
		t.Error("off not match, expect 30 bug got", arena.off)
	}

	if len(slice) != 0 || cap(slice) != 20 {
		t.Error("slice length or cap not match")
	}

	slice = arena.Alloc(1024)
	if arena.off != 30 {
		t.Error("off not match, expect 30 bug got", arena.off)
	}

	if len(slice) != 0 || cap(slice) != 1024 {
		t.Error("slice length or cap not match")
	}

	slice = arena.AllocWithLen(2, 10)
	if arena.off != 40 {
		t.Error("off not match, expect 40 bug got", arena.off)
	}

	if len(slice) != 2 || cap(slice) != 10 {
		t.Error("slice length or cap not match")
	}

	arena.Reset()
	if arena.off != 0 || cap(arena.arena) != 1000 {
		t.Error("off or cap not match")
	}
}

func TestStdAllocator(t *testing.T) {
	slice := StdAllocator.Alloc(20)
	if len(slice) != 0 {
		t.Error("length not match")
	}

	if cap(slice) != 20 {
		t.Error("cap not match")
	}

	slice = StdAllocator.AllocWithLen(10, 20)
	if len(slice) != 10 {
		t.Error("length not match")
	}

	if cap(slice) != 20 {
		t.Error("cap not match")
	}
}
