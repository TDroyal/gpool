package gpool_test

import (
	"sync"
	"testing"
	"time"

	"github.com/TDroyal/gpool"
)

const (
	epoch    = 100_000 // 10w
	poolSize = 1_000   // 1k
)

func demoFunc() {
	time.Sleep(time.Millisecond * 10) // task cost time
}

func BenchmarkNoPool(b *testing.B) {
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(epoch)
		for j := 0; j < epoch; j++ {
			go func() {
				defer wg.Done()
				demoFunc()
			}()
		}
		wg.Wait()
	}
}

func BenchmarkGPool(b *testing.B) {
	pool, _ := gpool.NewPool(poolSize)
	var wg sync.WaitGroup
	for i := 0; i < b.N; i++ {
		wg.Add(epoch)
		for j := 0; j < epoch; j++ {
			pool.Submit(
				func() {
					defer wg.Done()
					demoFunc()
				})
		}
		wg.Wait()
	}
}

/*
PS > go test -benchmem -run=^$ -bench ^BenchmarkNoPool$ github.com/TDroyal/gpool
goos: windows
goarch: amd64
pkg: github.com/TDroyal/gpool
cpu: AMD Ryzen 5 5600H with Radeon Graphics
BenchmarkNoPool-12            18          61349978 ns/op        10157213 B/op     203133 allocs/op
PASS
ok      github.com/TDroyal/gpool        1.867s
PS > go test -benchmem -run=^$ -bench ^BenchmarkGPool$ github.com/TDroyal/gpool
goos: windows
goarch: amd64
pkg: github.com/TDroyal/gpool
cpu: AMD Ryzen 5 5600H with Radeon Graphics
BenchmarkGPool-12              1        1813833100 ns/op         4310552 B/op     110399 allocs/op
PASS
ok      github.com/TDroyal/gpool        1.891s

*/
