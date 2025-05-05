package gpool_test

import (
	"testing"
	"time"

	"github.com/TDroyal/gpool"
)

func TestGpool(t *testing.T) {
	pool, _ := gpool.NewPool(2)
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
	pool.Submit(t1)
	pool.Submit(t2)
	pool.Submit(t3)
	pool.Submit(t4)
}
