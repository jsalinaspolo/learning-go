package main

import "testing"

func TestHello(t *testing.T) {
	got := Hello("Javier")
	want := "Hello, Javier"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
