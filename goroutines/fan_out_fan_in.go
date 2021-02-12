package goroutines

import (
	"fmt"
	"sync"
	"testing"
)

func worker(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			out <- fmt.Sprintf("out%s", v)
		}
	}()

	return out
}

func merge(cs ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan string) {
		defer wg.Done()
		for n := range c {
			out <- n
		}

	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
func TestGoroutinesPatterns(t *testing.T) {
	t.Run("fun-out and fun-in", func(t *testing.T) {
		//generate
		var in = make(chan string)

		// producer
		go func() {
			defer close(in)
			for i := 0; i < 100; i++ {
				in <- fmt.Sprintf("%d", i)
			}
		}()

		// fun-out
		c1 := worker(in)
		c2 := worker(in)

		for n := range merge(c1, c2) {
			fmt.Println(n) // 4 then 9, or 9 then 4
		}
	})
}

