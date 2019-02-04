// Package radix implements three-way radix Quicksort for strings.
// https://en.wikipedia.org/wiki/Multi-key_quicksort

package radix

func Sort(a []string) {
	qsort3(a, 0, maxDepth(len(a)))
}

// 3-way radix Quicksort
// Based on Sedgewick, Algorithms (4th ed.) p 720;
// Algorithm 5.3: "Three-way string quicksort"
func qsort3(a []string, d, maxDepth int) {
	if maxDepth == 0 {
		heapSort(a, 0, len(a))
		return
	}
	maxDepth--
	for len(a) > 1 {
		p := pivot(a, d)
		eq, x, gt := 0, 0, len(a)
		// Use Dijkstra's DNF algorithm to perform a 3-way partition of
		// the slice based on the bytes at index d.
		// Dijkstra, A Discipline of Programming, Ch. 14: "The problem
		// of the Dutch national flag"
		//
		// Invariants are:
		//  data[:eq]  < pivot
		//  data[eq:x] = pivot
		//  data[x:gt] unexamined
		//  data[gt:]  > pivot
		//
		// Each iteration, the length of the unexamined slice decreases by 1.
		for x < gt {
			cur := byteAt(a[x], d)
			if cur < p {
				a[x], a[eq] = a[eq], a[x]
				x++
				eq++
			} else if cur == p {
				x++
			} else {
				gt--
				a[x], a[gt] = a[gt], a[x]
			}
		}
		// Avoid recursion on the larger subproblem. This guarantees a
		// stack depth of at most lg(len(a)).
		// If we didn't care about worst-case stack depth, the outer
		// for loop would become an if statement, and the code below would
		// simply be:
		//   qsort3(a[:eq], d)
		//   if p > -1 {
		//     qsort3(a[eq:gt], d+1)
		//   }
		//   qsort3(a[gt:], d)
		if p < 0 {
			if eq < len(a)-gt {
				qsort3(a[:eq], d, maxDepth)
				a = a[gt:]
			} else {
				qsort3(a[gt:], d, maxDepth)
				a = a[:eq]
			}
		} else {
			max := max3(eq, gt-eq, len(a)-gt)
			if eq == max {
				qsort3(a[eq:gt], d+1, maxDepth)
				qsort3(a[gt:], d, maxDepth)
				a = a[:eq]
			} else if gt-eq == max {
				qsort3(a[:eq], d, maxDepth)
				qsort3(a[gt:], d, maxDepth)
				a = a[eq:gt]
				d++
			} else {
				qsort3(a[:eq], d, maxDepth)
				qsort3(a[eq:gt], d+1, maxDepth)
				a = a[gt:]
			}
		}
	}
}

func byteAt(s string, d int) int {
	if d < len(s) {
		return int(s[d])
	}
	return -1
}

func max3(i, j, k int) int {
	max := i
	if max < j {
		max = j
	}
	if max < k {
		max = k
	}
	return max
}

func pivot(a []string, d int) int {
	m := len(a) >> 1
	if len(a) > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := len(a) / 8
		return median3(
			median3(byteAt(a[0], d), byteAt(a[s], d), byteAt(a[2*s], d)),
			median3(byteAt(a[m], d), byteAt(a[m-s], d), byteAt(a[m+s], d)),
			median3(byteAt(a[len(a)-1], d), byteAt(a[len(a)-1-s], d), byteAt(a[len(a)-1-2*s], d)))
	}
	return median3(byteAt(a[0], d), byteAt(a[m], d), byteAt(a[len(a)-1], d))
}

func median3(i, j, k int) int {
	if i < j {
		if j < k {
			return j
		}
		if i < k {
			return k
		}
		return i
	}
	if k < j {
		return j
	}
	if k < i {
		return k
	}
	return i
}

// siftDown_ implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
// Copied from Go's built-in sort library, and specialized to strings.
func siftDown(data []string, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data[first+child] < data[first+child+1] {
			child++
		}
		if !(data[first+root] < data[first+child]) {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

// Copied from Go's built-in sort library, and specialized to strings.
func heapSort(data []string, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDown(data, lo, i, first)
	}
}

// maxDepth returns a threshold at which quicksort should switch
// to heapsort. It returns 2*ceil(lg(n+1)).
// COpied from Go's built-in sort library.
func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}
