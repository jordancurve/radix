package radix

import (
	"reflect"
	"sort"
	"testing"
	"testing/quick"
)

func TestMax3(t *testing.T) {
	f := func(a, b, c int) bool {
		abc := []int{a, b, c}
		sort.Ints(abc)
		return max3(a, b, c) == abc[2]
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestMedian(t *testing.T) {
	f := func(a, b, c int) bool {
		abc := []int{a, b, c}
		sort.Ints(abc)
		return median3(a, b, c) == abc[1]
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestQuickSort3Radix(t *testing.T) {
	f := func(a []string) bool {
		b := make([]string, len(a))
		copy(b, a)
		Sort(a)
		sort.Strings(b)
		return reflect.DeepEqual(a, b)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
