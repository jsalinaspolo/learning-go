package arrays_slices

import (
	"testing"
)

func TestSum(t *testing.T) {

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		assertSum(t, got, want, numbers)
	})
}

func assertSum(t *testing.T, got int, want int, numbers []int) {
	if got != want {
		t.Errorf("got %d want %d given, %v", got, want, numbers)
	}
}
