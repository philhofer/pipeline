package prototype

import (
	"sync"
	"time"
)

type T struct{}
type J struct{}

// Merge merges multiple channels into one
func MergeT(v ...<-chan T) <-chan T {
	wg := new(sync.WaitGroup)
	out := make(chan T)
	for _, c := range v {
		wg.Add(1)
		go func(c <-chan T) {
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
func FanoutT(v <-chan T, n int) []<-chan T {
	out := make([]<-chan T, 0, n)
	for i := 0; i < n; i++ {
		local := make(chan T)
		go func(c <-chan T, i int, out chan T) {
			for e := range c {
				out <- e
			}
			close(out)
		}(v, i, local)
		out = append(out, local)
	}
	return out
}

// Apply applies 'fs' successively to each element of a channel
// and sends the the object to the output channel.
func ApplyT(v <-chan T, fs ...func(T)) <-chan T {
	out := make(chan T)
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
func PapplyT(v <-chan T, n int, fs ...func(T)) <-chan T {
	out := make(chan T)
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
func MapT(v <-chan T, fs ...func(T) T) <-chan T {
	out := make(chan T)
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
func PmapT(v <-chan T, n int, fs ...func(T) T) <-chan T {
	out := make(chan T)
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

// Filter filters elements from 'v' using the function
// 'f'.
func FilterT(v <-chan T, f func(T) bool) <-chan T {
	out := make(chan T)
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

// Transform transforms elements from 'v' from 'T' to 'J' using
// 'f' and sends them on the output channel.
func TransformTJ(v <-chan T, f func(T) J) <-chan J {
	out := make(chan J)
	go func() {
		for e := range v {
			out <- f(e)
		}
		close(out)
	}()
	return out
}

// Ptransform performs Transform() in parallel using 'n' goroutines.
func PtransformTJ(v <-chan T, f func(T) J, n int) <-chan J {
	out := make(chan J)
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
func SendAllT(in []T) <-chan T {
	out := make(chan T)
	go func() {
		for i := range in {
			out <- in[i]
		}
	}()
	return out
}

// RecvAll collects every element sent on the
// input channel into a slice.
func RecvAllT(in <-chan T) []T {
	out := make([]T, 0, 50)
	for e := range in {
		out = append(out, e)
	}
	return out
}

// RecvN pulls 'n' elements out of 'in' and returns
// the slice and the channel.
func RecvNT(in <-chan T, n int) []T {
	ot := make([]T, 0, n)
	for i := 0; i < n; i++ {
		ot = append(ot, <-in)
	}
	return ot
}

// Buffer reads elements from 'in' and attempts to send them
// to 'out'. If the send would block, the messages are buffered
// internally. Buffer closes 'out' after 't' is closed. Buffer() blocks
// until 'in' is closed, so in most cases it should be run asynchronously. Buffer
// uses a LIFO queue internally, so it should only be used when ordering
// doesn't matter.
func BufferT(in <-chan T, out chan<- T) {
	var buf []T
	for {
		if len(buf) > 0 {
			select {
			case v, ok := <-in:
				if !ok {
					// flush everything
					for i := range buf {
						out <- buf[i]
					}
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
			case v, ok := in:
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
