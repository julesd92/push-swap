// Package solver computes a short list of push-swap instructions.
package solver

import "push-swap/stack"

type solver struct {
	st  *stack.State
	ops []string
}

// Solve returns the instructions that sort nums (top = first element).
func Solve(nums []int) []string {
	s := &solver{st: stack.New(nums)}
	if s.st.Sorted() {
		return nil
	}
	switch n := len(nums); {
	case n == 2:
		s.do("sa", 1)
	case n == 3:
		s.sort3()
	case n <= 5:
		s.sortSmall()
	default:
		s.sortLarge()
	}
	return s.ops
}

// do applies op n times and records it.
func (s *solver) do(op string, n int) {
	for i := 0; i < n; i++ {
		s.st.Apply(op)
		s.ops = append(s.ops, op)
	}
}

// sort3 sorts a stack a of exactly 3 elements.
func (s *solver) sort3() {
	a := s.st.A
	x, y, z := a[0], a[1], a[2]
	switch {
	case x < y && y < z: // 1 2 3
	case x < z && z < y: // 1 3 2
		s.do("sa", 1)
		s.do("ra", 1)
	case y < x && x < z: // 2 1 3
		s.do("sa", 1)
	case y < z && z < x: // 3 1 2
		s.do("ra", 1)
	case z < x && x < y: // 2 3 1
		s.do("rra", 1)
	default: // 3 2 1
		s.do("sa", 1)
		s.do("rra", 1)
	}
}

// sortSmall handles 4 or 5 elements: push the smallest ones to b,
// sort the 3 left in a, then bring everything back.
func (s *solver) sortSmall() {
	for len(s.st.A) > 3 {
		s.rotateA(indexOf(s.st.A, minOf(s.st.A)))
		s.do("pb", 1)
	}
	s.sort3()
	s.do("pa", len(s.st.B))
}

// sortLarge is the cost-based algorithm: keep b sorted in descending
// order while pushing the cheapest element of a each time, then push
// everything back into a at the right place.
func (s *solver) sortLarge() {
	toPush := len(s.st.A) - 3
	for i := 0; i < toPush; i++ {
		if i < 2 {
			s.do("pb", 1)
			continue
		}
		s.pushCheapest()
	}
	s.sort3()
	for len(s.st.B) > 0 {
		s.rotateA(targetInA(s.st.A, s.st.B[0]))
		s.do("pa", 1)
	}
	s.rotateA(indexOf(s.st.A, minOf(s.st.A)))
}

const (
	bothUp = iota
	bothDown
	aUpBDown
	aDownBUp
)

// cost returns the cheapest way to bring a[i] and b[j] to the top.
func cost(i, j, la, lb int) (best, mode int) {
	upA, upB := i, j
	downA, downB := (la-i)%la, (lb-j)%lb
	options := [4]int{
		max(upA, upB),
		max(downA, downB),
		upA + downB,
		downA + upB,
	}
	best = options[0]
	for m := 1; m < 4; m++ {
		if options[m] < best {
			best, mode = options[m], m
		}
	}
	return best, mode
}

func (s *solver) pushCheapest() {
	a, b := s.st.A, s.st.B
	la, lb := len(a), len(b)
	bestCost, bestI, bestJ, bestMode := -1, 0, 0, 0
	for i := range a {
		j := targetInB(b, a[i])
		c, m := cost(i, j, la, lb)
		if bestCost == -1 || c < bestCost {
			bestCost, bestI, bestJ, bestMode = c, i, j, m
		}
	}
	upA, upB := bestI, bestJ
	downA, downB := (la-bestI)%la, (lb-bestJ)%lb
	switch bestMode {
	case bothUp:
		r := min(upA, upB)
		s.do("rr", r)
		s.do("ra", upA-r)
		s.do("rb", upB-r)
	case bothDown:
		r := min(downA, downB)
		s.do("rrr", r)
		s.do("rra", downA-r)
		s.do("rrb", downB-r)
	case aUpBDown:
		s.do("ra", upA)
		s.do("rrb", downB)
	case aDownBUp:
		s.do("rra", downA)
		s.do("rb", upB)
	}
	s.do("pb", 1)
}

// rotateA brings a[j] to the top of a using the shortest rotation.
func (s *solver) rotateA(j int) {
	la := len(s.st.A)
	if j <= la-j {
		s.do("ra", j)
	} else {
		s.do("rra", la-j)
	}
}

// targetInB: index of the largest value smaller than x in b (b is
// descending); if none, the index of the maximum.
func targetInB(b []int, x int) int {
	best := -1
	for i, v := range b {
		if v < x && (best == -1 || v > b[best]) {
			best = i
		}
	}
	if best == -1 {
		return indexOf(b, maxOf(b))
	}
	return best
}

// targetInA: index of the smallest value greater than x in a (a is
// ascending, cyclically); if none, the index of the minimum.
func targetInA(a []int, x int) int {
	best := -1
	for i, v := range a {
		if v > x && (best == -1 || v < a[best]) {
			best = i
		}
	}
	if best == -1 {
		return indexOf(a, minOf(a))
	}
	return best
}

func indexOf(s []int, v int) int {
	for i, x := range s {
		if x == v {
			return i
		}
	}
	return -1
}

func minOf(s []int) int {
	m := s[0]
	for _, v := range s {
		if v < m {
			m = v
		}
	}
	return m
}

func maxOf(s []int) int {
	m := s[0]
	for _, v := range s {
		if v > m {
			m = v
		}
	}
	return m
}
