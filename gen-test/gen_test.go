package test

import (
	"sort"
	"testing"
	"time"
)

func TestBuild(t *testing.T) {}

func TestMerge(t *testing.T) {
	one := make(chan int, 3)
	two := make(chan int, 3)

	one <- 1
	one <- 2
	one <- 3
	two <- 4
	two <- 5
	two <- 6

	out := Mergeint(one, two)
	close(one)
	close(two)

	slc := make([]int, 6)
	for i := 0; i < 6; i++ {
		select {
		case slc[i] = <-out:
		case <-time.After(1 * time.Millisecond):
			t.Fatal("didn't get 6 elements")
		}
	}
	sort.Ints(slc)
	for i := range slc {
		if slc[i] != i+1 {
			t.Errorf("expected %d at index %d", i+1, i)
		}
	}
}

func TestFanout(t *testing.T) {
	in := make(chan int, 3)

	in <- 1
	in <- 2
	in <- 3
	close(in)

	out := Fanoutint(in, 2)

	if len(out) != 2 {
		t.Fatalf("wanted 2 channels, got %d", len(out))
	}

	slc := make([]int, 3)
	for i := 0; i < 3; i++ {
		select {
		case v := <-out[0]:
			slc[i] = v
		case v := <-out[1]:
			slc[i] = v
		case <-time.After(1 * time.Millisecond):
			t.Fatal("didn't get three elements.")
		}
	}

	select {
	case <-out[0]:
		t.Error("can still receive from channel")
	case <-out[1]:
		t.Error("can still receive from channel")
	default:
	}

	sort.Ints(slc)
	for i := range slc {
		if slc[i] != i+1 {
			t.Errorf("expected %d at index %d", i+1, i)
		}
	}
}
