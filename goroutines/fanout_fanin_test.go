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

func doWork(in <-chan string) <-chan string {
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
	// copies values from channel to out until channel is closed, then calls wg.Done.
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
		out := doWork(c)

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
			channels[cpu] = doWork(in)
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

	t.Run("fan-out fan-in in order using promises", func(t *testing.T) {
		type slowFunc func() string
		type future struct {
			channel chan string
			fn      slowFunc
		}

		newSlowFunc := func(i int) slowFunc {
			return func() string {
				time.Sleep(100 * time.Millisecond)
				return fmt.Sprintf("out-%d", i)
			}
		}

		var futures []*future
		newFuture := func(fn slowFunc) *future {
			f := &future{
				channel: make(chan string),
				fn:      fn,
			}

			go func() {
				f.channel <- f.fn()
			}()

			return f
		}

		n := 100
		// produce and fan out
		for i := 0; i < n; i++ {
			futures = append(futures, newFuture(newSlowFunc(i)))
		}

		now := time.Now()
		// fan-in
		var result []string
		for i := range futures {
			result = append(result, <-futures[i].channel)
		}

		expected := expectedSortedOutput(n)
		require.Equal(t, len(expected), len(result))
		require.Equal(t, expectedSortedOutput(n), result)
		require.LessOrEqual(t, time.Since(now).Milliseconds(), workTime.Milliseconds()+10)
	})
}
