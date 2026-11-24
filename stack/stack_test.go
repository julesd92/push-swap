package stack

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	good := []struct {
		in   []string
		want []int
	}{
		{[]string{"2 1 3"}, []int{2, 1, 3}},
		{[]string{"2", "1", "3"}, []int{2, 1, 3}},
		{[]string{"+5 -7"}, []int{5, -7}},
		{[]string{""}, nil},
	}
	for _, c := range good {
		got, err := Parse(c.in)
		if err != nil || !reflect.DeepEqual(got, c.want) {
			t.Errorf("Parse(%q) = %v, %v; want %v", c.in, got, err, c.want)
		}
	}
	for _, in := range [][]string{{"1 one"}, {"1 2 1"}, {"1.5"}, {"99999999999999999999"}, {"0 -0"}} {
		if _, err := Parse(in); err == nil {
			t.Errorf("Parse(%q) expected error", in)
		}
	}
}

func TestSubjectExample(t *testing.T) {
	s := New([]int{2, 1, 3, 6, 8, 5})
	for _, op := range []string{"sa", "pb", "pb", "pb", "rb", "rra", "rrb", "pa", "pa", "pa"} {
		if err := s.Apply(op); err != nil {
			t.Fatal(err)
		}
	}
	if !s.Sorted() {
		t.Errorf("expected sorted, got a=%v b=%v", s.A, s.B)
	}
}

func TestOps(t *testing.T) {
	s := New([]int{1, 2, 3})
	s.Apply("ra")
	if !reflect.DeepEqual(s.A, []int{2, 3, 1}) {
		t.Errorf("ra: %v", s.A)
	}
	s.Apply("rra")
	if !reflect.DeepEqual(s.A, []int{1, 2, 3}) {
		t.Errorf("rra: %v", s.A)
	}
	s.Apply("pa") // empty b: no-op
	if len(s.A) != 3 {
		t.Errorf("pa on empty b changed a: %v", s.A)
	}
	if s.Apply("xx") == nil {
		t.Error("unknown instruction accepted")
	}
}
