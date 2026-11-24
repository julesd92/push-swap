// Package stack implements the two stacks and the push-swap instruction set.
// In both stacks, index 0 is the top.
package stack

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Parse turns command-line arguments into a list of integers.
// Each argument may hold several numbers separated by whitespace.
func Parse(args []string) ([]int, error) {
	var nums []int
	seen := make(map[int]struct{})
	for _, arg := range args {
		for _, field := range strings.Fields(arg) {
			n, err := strconv.Atoi(field)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", field)
			}
			if _, dup := seen[n]; dup {
				return nil, fmt.Errorf("duplicate value %d", n)
			}
			seen[n] = struct{}{}
			nums = append(nums, n)
		}
	}
	return nums, nil
}

// State holds stacks a and b.
type State struct {
	A, B []int
}

// New creates a State with the given numbers in stack a and an empty b.
func New(nums []int) *State {
	a := make([]int, len(nums))
	copy(a, nums)
	return &State{A: a, B: []int{}}
}

// ErrUnknown is returned for an instruction that does not exist.
var ErrUnknown = errors.New("unknown instruction")

// Valid reports whether op is one of the 11 instructions.
func Valid(op string) bool {
	switch op {
	case "pa", "pb", "sa", "sb", "ss", "ra", "rb", "rr", "rra", "rrb", "rrr":
		return true
	}
	return false
}

// Apply executes one instruction. Instructions that cannot act
// (e.g. pa on an empty b) are no-ops.
func (s *State) Apply(op string) error {
	switch op {
	case "pa":
		push(&s.B, &s.A)
	case "pb":
		push(&s.A, &s.B)
	case "sa":
		swap(s.A)
	case "sb":
		swap(s.B)
	case "ss":
		swap(s.A)
		swap(s.B)
	case "ra":
		rotate(s.A)
	case "rb":
		rotate(s.B)
	case "rr":
		rotate(s.A)
		rotate(s.B)
	case "rra":
		reverse(s.A)
	case "rrb":
		reverse(s.B)
	case "rrr":
		reverse(s.A)
		reverse(s.B)
	default:
		return ErrUnknown
	}
	return nil
}

// Sorted reports whether a is ascending (top smallest) and b is empty.
func (s *State) Sorted() bool {
	if len(s.B) != 0 {
		return false
	}
	for i := 1; i < len(s.A); i++ {
		if s.A[i-1] > s.A[i] {
			return false
		}
	}
	return true
}

func swap(s []int) {
	if len(s) >= 2 {
		s[0], s[1] = s[1], s[0]
	}
}

// push moves the top of from onto the top of to.
func push(from, to *[]int) {
	if len(*from) == 0 {
		return
	}
	top := (*from)[0]
	*from = (*from)[1:]
	*to = append([]int{top}, *to...)
}

func rotate(s []int) {
	if len(s) >= 2 {
		first := s[0]
		copy(s, s[1:])
		s[len(s)-1] = first
	}
}

func reverse(s []int) {
	if len(s) >= 2 {
		last := s[len(s)-1]
		copy(s[1:], s[:len(s)-1])
		s[0] = last
	}
}
