package solver

import (
	"math/rand"
	"testing"

	"push-swap/stack"
)

func check(t *testing.T, nums []int) int {
	t.Helper()
	ops := Solve(nums)
	s := stack.New(nums)
	for _, op := range ops {
		if err := s.Apply(op); err != nil {
			t.Fatalf("bad op %q", op)
		}
	}
	if !s.Sorted() {
		t.Fatalf("not sorted for %v: a=%v b=%v", nums, s.A, s.B)
	}
	return len(ops)
}

func TestAllPermutationsUpTo6(t *testing.T) {
	var perm func(cur, rest []int)
	perm = func(cur, rest []int) {
		if len(rest) == 0 {
			check(t, append([]int(nil), cur...))
			return
		}
		for i := range rest {
			next := append([]int(nil), rest[:i]...)
			next = append(next, rest[i+1:]...)
			perm(append(cur, rest[i]), next)
		}
	}
	for n := 0; n <= 6; n++ {
		base := make([]int, n)
		for i := range base {
			base[i] = i*3 - 4 // include negatives
		}
		perm(nil, base)
	}
}

func TestRandom(t *testing.T) {
	r := rand.New(rand.NewSource(42))
	for _, n := range []int{7, 10, 50, 100, 500} {
		worst := 0
		for k := 0; k < 20; k++ {
			nums := r.Perm(n)
			for i := range nums {
				nums[i] -= n / 2
			}
			if c := check(t, nums); c > worst {
				worst = c
			}
		}
		t.Logf("n=%d worst ops over 20 runs: %d", n, worst)
	}
}
