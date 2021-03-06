package sort

import (
	"github.com/shady831213/algorithms/heap"
	"testing"
)

func countingSortWrapper(arr []int) {
	_arr := make([]int, len(arr), cap(arr))
	copy(_arr, arr)
	h := heap.NewHeapIntArray(arr)
	max := h.Pop().(int)
	sortedArry := countingSort(arr, max)
	copy(arr, sortedArry)
}

func Test_countingSort(t *testing.T) {
	testSort(t, countingSortWrapper)
}
