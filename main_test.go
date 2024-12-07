package main

import "testing"

func TestAdd(t *testing.T) {

	got := Add(2, 3)
	want := 5

	if got != want {
		t.Errorf("got: %d, wanted: %d.", got, want)
	}
}
