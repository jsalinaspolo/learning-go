package goroutines

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func randomNumber() time.Duration {
	return time.Duration(rand.Intn(100-1) + 1)
}

func generate(n int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			out <- fmt.Sprintf("%d", i)
		}
	}()
	return out
}

func worker(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			time.Sleep(randomNumber() * time.Millisecond)
			out <- fmt.Sprintf("out%s", v)
		}
	}()

	return out
}

func fanIn(cs ...<-chan string) <-chan string {
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
	t.Run("pipelines", func(t *testing.T) {
		c := generate(10)
		out := worker(c)

		for v := range out {
			fmt.Println(v)
		}
	})

	t.Run("fun-out and fun-in", func(t *testing.T) {
		in := generate(100)

		// fan-out
		maxCPU := runtime.NumCPU()
		channels := make([]<-chan string, maxCPU)
		for cpu := 0; cpu < maxCPU; cpu++ {
			channels[cpu] = worker(in)
		}

		// fan-in
		for n := range fanIn(channels...) {
			fmt.Println(n)
		}
	})

	t.Run("fun-out and fun-in in order?", func(t *testing.T) {
		// Produce input
		in := generate(100)

		// fun-out
		c1 := worker(in)
		c2 := worker(in)

		var seq []string
		for n := range fanIn(c1, c2) {
			fmt.Println(n)
			seq = append(seq, n)
		}

		fmt.Println(len(seq))
	})
}
