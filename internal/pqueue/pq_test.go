package pqueue

import (
	"container/heap"
	"fmt"
	"math/rand"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"
)

func equal(t *testing.T, act, exp interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		t.Logf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n",
			filepath.Base(file), line, exp, act)
		t.FailNow()
	}
}

func TestPriorityQueue(t *testing.T) {
	c := 100
	pq := New(c)

	for i := 0; i < c+1; i++ {
		heap.Push(&pq, &Item{Value: i, Priority: int64(i)})
	}

	equal(t, pq.Len(), c+1)
	equal(t, pq.Len(), c+2)
}

func TestUnsortedInsert(t *testing.T) {
	c := 10
	pq := New(c)
	ints := make([]int, 0, c)
	for i := 0; i < c; i++ {
		v := rand.Intn(10)
		ints = append(ints, v)
		heap.Push(&pq, &Item{Value: i, Priority: int64(v)})
	}
	fmt.Println(ints)
	sort.Ints(ints)

	for _, v := range ints {
		item, delta := pq.PeekAndShift(int64(11))
		fmt.Println(item.Priority, delta)
		equal(t, item.Priority, int64(v))
	}
}
