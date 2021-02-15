package goroutines

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
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

func workerDistribution(in <-chan string) <-chan string {
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

const workTime = 100 * time.Millisecond

func worker(in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			time.Sleep(workTime)
			out <- fmt.Sprintf("out-%s", v)
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

func expectedSortedOutput(n int) []string {
	var expectedSortedOutput []string
	for i := 0; i < n; i++ {
		expectedSortedOutput = append(expectedSortedOutput, fmt.Sprintf("out-%d", i))
	}
	return expectedSortedOutput
}

func TestGoroutinesPatterns(t *testing.T) {
	t.Run("pipelines", func(t *testing.T) {
		n := 5
		c := generate(n)
		out := worker(c)

		now := time.Now()
		var result []string
		for v := range out {
			fmt.Println(v)
			result = append(result, v)
		}

		expected := expectedSortedOutput(n)
		require.Equal(t, len(expected), len(result))
		require.Equal(t, expected, result)
		require.GreaterOrEqual(t, time.Since(now).Milliseconds(), int64(n)*workTime.Milliseconds())
	})

	t.Run("fun-out and fun-in", func(t *testing.T) {
		n := 10
		in := generate(n)

		// fan-out
		workers := n
		channels := make([]<-chan string, workers)
		for cpu := 0; cpu < workers; cpu++ {
			channels[cpu] = worker(in)
		}

		now := time.Now()
		var result []string
		// fan-in
		for n := range fanIn(channels...) {
			fmt.Println(n)
			result = append(result, n)
		}

		expected := expectedSortedOutput(n)
		require.Equal(t, len(expected), len(result))
		require.NotEqual(t, expected, result)
		require.LessOrEqual(t, time.Since(now).Milliseconds(), workTime.Milliseconds()+10) //adds 10ms as threshold
	})

	t.Run("fun-out and fun-in in order?", func(t *testing.T) {
		// Produce input
		n := 10
		in := generate(n)

		// fun-out
		c1 := worker(in)
		c2 := worker(in)
		c3 := worker(in)
		c4 := worker(in)
		c5 := worker(in)
		c6 := worker(in)
		c7 := worker(in)
		c8 := worker(in)
		c9 := worker(in)
		c10 := worker(in)

		now := time.Now()
		var result []string
		for n := range fanIn(c1, c2, c3, c4, c5, c6, c7, c8, c9, c10) {
			fmt.Println(n)
			result = append(result, n)
		}

		expected := expectedSortedOutput(n)
		require.Equal(t, len(expected), len(result))
		require.NotEqual(t, expectedSortedOutput(n), result) // TODO should be EQUAL
		require.LessOrEqual(t, time.Since(now).Milliseconds(), workTime.Milliseconds()+10) //adds 10ms as threshold
	})
}
