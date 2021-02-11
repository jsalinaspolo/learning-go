package goroutines

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLearningGoRoutines(t *testing.T) {
	t.Run("goroutines are async", func(t *testing.T) {
		elements := make([]string, 100)
		for i := 0; i < len(elements); i++ {
			elements[i] = fmt.Sprintf("el%d", i)
		}

		r, w, _ := os.Pipe()
		os.Stdout = w
		var goroutine = func() {
			fmt.Fprint(w, "1")
			// goroutines are async and will trigger eventually
			go func() {
				for _, id := range elements {
					fmt.Fprint(w, id)
				}
			}()

			fmt.Fprint(w, "2")
		}

		goroutine()
		w.Close()
		out, _ := ioutil.ReadAll(r)

		require.Equal(t, "12", string(out))
	})

	t.Run("goroutines async wait for trigger", func(t *testing.T) {
		elements := []string{"el1", "el2", "el3"}
			r, w, _ := os.Pipe()
			os.Stdout = w
		var goroutine = func() {

			fmt.Fprint(w, "1")

			// goroutines are async and will trigger eventually
			go func() {
				// loop is blocking
				for _, id := range elements {
					fmt.Fprint(w,  id)
				}
			}()
			fmt.Fprint(w, "2")

			// gives time to triggers the goroutine
			time.Sleep(100 * time.Millisecond)
		}

		goroutine()
		w.Close()
		out, _ := ioutil.ReadAll(r)

		require.Equal(t, "12el1el2el3", string(out))
	})

	t.Run("goroutines using loop make it blocking", func(t *testing.T) {
		elements := []string{"el1", "el2", "el3"}

		var goroutine = func() []string {
			var sequence []string
			ch := make(chan string, 5)

			sequence = append(sequence, "1")
			go func() {
				defer close(ch)
				//blocking in order
				for _, id := range elements {
					ch <- func() string {
						// add sleep to get mix between adds and reads
						time.Sleep(100 * time.Millisecond)
						return id
					}()
				}
			}()

			sequence = append(sequence, "2")

			for range elements {
				select {
				case res := <-ch:
					sequence = append(sequence, res)
				}
			}
			sequence = append(sequence, "3")
			return sequence
		}

		now := time.Now()
		r := goroutine()
		var expected = append(append([]string{"1", "2"}, elements...), "3")
		require.Equal(t, expected, r)
		// time is more than 100ms (is the sum of 100*elements)
		require.GreaterOrEqual(t, time.Since(now).Milliseconds(), 100*int64(len(elements))*time.Millisecond.Milliseconds())
	})

	t.Run("goroutines run in parallel", func(t *testing.T) {
		elements := []string{"el1", "el2", "el3"}

		var goroutine = func() []string {
			var sequence []string
			ch := make(chan string, 5)

			sleep := func(id string) string {
				time.Sleep(100 * time.Millisecond)
				return id
			}

			sequence = append(sequence, "1")
			for _, id := range elements {
				go func(identifier string) {
					ch <- sleep(identifier)
				}(id)
			}

			sequence = append(sequence, "2")

			for range elements {
				select {
				case res := <-ch:
					sequence = append(sequence, res)
				}
			}
			sequence = append(sequence, "3")
			return sequence
		}

		now := time.Now()
		r := goroutine()
		// async and non sequential
		require.Equal(t, "1", r[0])
		require.Equal(t, "2", r[1])
		require.Equal(t, "3", r[len(elements)+2])
		fmt.Println(r)
		// time is around 100ms
		require.LessOrEqual(t, time.Since(now).Milliseconds(), 110*time.Millisecond.Milliseconds())
	})

	t.Run("goroutines limited channel", func(t *testing.T) {
		elements := make([]string, 100)
		for i := 0; i < len(elements); i++ {
			elements[i] = fmt.Sprintf("el%d", i)
		}

		var goroutine = func() []string {
			var sequence []string
			ch := make(chan string, 5)
			defer close(ch)

			sleep := func(id string) string {
				time.Sleep(100 * time.Millisecond)
				return id
			}

			sequence = append(sequence, "1")
			for _, id := range elements {
				go func(identifier string) {
					ch <- sleep(identifier)
				}(id)
			}

			sequence = append(sequence, "2")

			for range elements {
				select {
				case res := <-ch:
					sequence = append(sequence, res)
				}
			}
			sequence = append(sequence, "3")
			return sequence
		}

		now := time.Now()
		r := goroutine()
		fmt.Println(r)
		require.Equal(t, "1", r[0])
		require.Equal(t, "3", r[len(r)-1])
		require.LessOrEqual(t, time.Since(now).Milliseconds(), 120*time.Millisecond.Milliseconds())
		require.Equal(t, len(elements)+3, len(r))
	})

	t.Run("goroutines use waitGroups", func(t *testing.T) {
		elements := make([]string, 100)

		for i := 0; i < len(elements); i++ {
			elements[i] = fmt.Sprintf("el%d", i)
		}

		var goroutine = func() []string {
			var sequence []string
			var wg sync.WaitGroup

			sleep := func(id string, wg *sync.WaitGroup) string {
				defer wg.Done()
				time.Sleep(100 * time.Millisecond)
				return id
			}

			sequence = append(sequence, "1")
			mu := &sync.Mutex{}
			for _, id := range elements {
				wg.Add(1)
				go func(identifier string) {
					mu.Lock()
					// sequence is shared so need a lock
					sequence = append(sequence, identifier)
					mu.Unlock()
					sleep(identifier, &wg)
				}(id)
			}
			wg.Wait()
			sequence = append(sequence, "2")
			return sequence
		}

		now := time.Now()
		r := goroutine()
		fmt.Println(r)
		require.Equal(t, "1", r[0])
		require.Equal(t, "2", r[len(r)-1])
		require.LessOrEqual(t, time.Since(now).Milliseconds(), 120*time.Millisecond.Milliseconds())
		require.Equal(t, len(elements)+2, len(r))
	})
}
