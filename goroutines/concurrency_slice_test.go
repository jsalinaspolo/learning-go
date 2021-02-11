// +build !race

package goroutines

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

func TestSliceGoRoutines(t *testing.T) {

	// running with -race, will fail
	t.Run("slice is not thread safe", func(t *testing.T) {
		var s []string
		limit := 10
		ch := make(chan string, 5)
		defer close(ch)
		for i := 0; i < limit; i++ {
			go func(id int) {
				identifier := strconv.Itoa(id)
				s = append(s, identifier)
				ch <- identifier
			}(i)
		}

		var expected []string
		for i := 0; i < limit; i++ {
			expected = append(expected, strconv.Itoa(i))
		}

		for i := 0; i < limit; i++ {
			select {
			case _ = <-ch:
				// do nothing
			}
		}

		t.Skip("Expected to fail with race condition")
	})

	t.Run("use channels to get results from goroutines", func(t *testing.T) {
		var s []string
		limit := 10
		ch := make(chan string, 5)
		defer close(ch)
		for i := 0; i < limit; i++ {
			go func(id int) {
				ch <- strconv.Itoa(id)
			}(i)
		}

		var expected []string
		for i := 0; i < limit; i++ {
			expected = append(expected, strconv.Itoa(i))
		}

		for i := 0; i < limit; i++ {
			select {
			case res := <-ch:
				s = append(s, res)
			}
		}

		require.Equal(t, len(expected), len(s))
	})

}

