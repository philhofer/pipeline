package test

import (
	"sync"
)

// Merge merges multiple channels into one
func Mergeint(v ...<-chan int) <-chan int {
	wg := new(sync.WaitGroup)
	out := make(chan int)
	for _, c := range v {
		wg.Add(1)
		go func(c <-chan int) {
			for e := range c {
				out <- e
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// Fanout splits one channel into 'n' channels
func Fanoutint(v <-chan int, n int) []<-chan int {
	out := make([]<-chan int, 0, n)
	for i := 0; i < n; i++ {
		local := make(chan int)
		go func(c <-chan int, i int, out chan int) {
			for e := range c {
				out <- e
			}
			close(out)
		}(v, i, local)
		out = append(out, local)
	}
	return out
}

// Transform transforms elements from 'v' from int to *string using
// 'f' and sends them on the output channel.
func Transformintptrstring(v <-chan int, f func(int) *string) <-chan *string {
	out := make(chan *string)
	go func() {
		for e := range v {
			out <- f(e)
		}
		close(out)
	}()
	return out
}

// Apply applies 'fs' successively to each element of a channel
// and sends the the object to the output channel.
func Applyint(v <-chan int, fs ...func(int)) <-chan int {
	out := make(chan int)
	go func() {
		for e := range v {
			for _, f := range fs {
				f(e)
			}
			out <- e
		}
		close(out)
	}()
	return out
}

// Papply applies 'fs' successively to each element received on 'v' and
// sends the object along to the output channel. 'n' goroutines
// are used for processing.
func Papplyint(v <-chan int, n int, fs ...func(int)) <-chan int {
	out := make(chan int)
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for e := range v {
				for _, f := range fs {
					f(e)
				}
				out <- e
			}
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// Map applies 'fs' successively to each element read from
// 'v' and sends them to the output channel.
func Mapint(v <-chan int, fs ...func(int) int) <-chan int {
	out := make(chan int)
	go func() {
		for e := range v {
			for _, f := range fs {
				e = f(e)
			}
			out <- e
		}
		close(out)
	}()
	return out
}

// Pmap applies 'fs' successively to elements read from 'v' and
// sends the result to the output channel. 'n' goroutines are used
// for processing.
func Pmapint(v <-chan int, n int, fs ...func(int) int) <-chan int {
	out := make(chan int)
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for e := range v {
				for _, f := range fs {
					e = f(e)
				}
				out <- e
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// Filter filters the input channel using 'f'. Only elements
// for which f(int) evaluates to 'true' will appear on the output channel.
func Filterint(v <-chan int, f func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		for e := range v {
			if f(e) {
				out <- e
			}
		}
		close(out)
	}()
	return out
}

// Ptransform performs Transformintptrstring() in parallel using 'n' goroutines.
func Ptransformintptrstring(v <-chan int, f func(int) *string, n int) <-chan *string {
	out := make(chan *string)
	wg := new(sync.WaitGroup)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			for e := range v {
				out <- f(e)
			}
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// SendAll sends a slice sequentially over a channel
func SendAllint(in []int) <-chan int {
	out := make(chan int)
	go func() {
		for i := range in {
			out <- in[i]
		}
	}()
	return out
}// RecvAll collects every element sent on the
// input channel into a slice. It blocks until the input channel is closed.
func RecvAllint(in <-chan int) []int {
	out := make([]int, 0, 50)
	for e := range in {
		out = append(out, e)
	}
	return out
}

// RecvN pulls 'n' elements out of 'in' and returns a slice
// of those elements. The length of the output slice may be less than 'n' if
// the channel was closed before 'n' elements were collected.
func RecvNint(in <-chan int, n int) []int {
	ot := make([]int, 0, n)
	for i := 0; i < n; i++ {
		ot = append(ot, <-in)
	}
	return ot
}

// Buffer reads elements from 'in' and attempts to send them
// to 'out'. If the send would block, the messages are buffered
// internally. Buffer closes 'out' after 't' is closed. Buffer() blocks
// until 'in' is closed, so in most cases it should be run asynchronously. Buffer uses
// a LIFO queue, so it should only be used in cases where ordering doesn't matter.
func Bufferint(in <-chan int, out chan<- int) {
	var buf []int
	for {
		if len(buf) > 0 {
			select {
			case v, ok := <-in:
				if !ok {
					break
				}
				select {
				case out <- v:
					continue
				default:
					buf = append(buf, v)
					continue
				}
			case out <- buf[len(buf)-1]:
				buf = buf[:len(buf)-1]
				continue
			}
		} else {
			select {
			case v, ok := <-in:
				if !ok {
					break
				}
				select {
				case out <- v:
					continue
				default:
					buf = append(buf, v)
					continue
				}
			}
		}
	}
	close(out)
}

