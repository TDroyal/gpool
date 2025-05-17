package gpool_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/TDroyal/gpool"
)

func TestGpool(t *testing.T) {
	pool, _ := gpool.NewPool(4)
	t1 := func() {
		t.Log("[task]===1111\n")
		time.Sleep(time.Second * 3)
	}
	t2 := func() {
		t.Log("[task]===2222\n")
		time.Sleep(time.Second * 3)
	}
	t3 := func() {
		t.Log("[task]===3333\n")
		time.Sleep(time.Second * 3)
	}
	t4 := func() {
		t.Log("[task]===4444\n")
		time.Sleep(time.Second * 3)
	}
	t5 := func() {
		t.Log("[task]===5555\n")
		time.Sleep(time.Second * 3)
	}
	t6 := func() {
		t.Log("[task]===6666\n")
		time.Sleep(time.Second * 3)
	}
	t7 := func() {
		t.Log("[task]===7777\n")
		time.Sleep(time.Second * 3)
	}
	t8 := func() {
		t.Log("[task]===8888\n")
		time.Sleep(time.Second * 3)
	}
	t9 := func() {
		t.Log("[task]===9999\n")
		time.Sleep(time.Second * 3)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println("[info]: the number of goroutines: ", runtime.NumGoroutine())
		}
	}()

	ts := []gpool.Task{t1, t2, t3, t4, t5, t6, t7, t8, t9}
	for i := range ts {
		pool.Submit(ts[i])
	}

	select {}
}
