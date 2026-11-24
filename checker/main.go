package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"push-swap/stack"
)

func fail() {
	fmt.Fprintln(os.Stderr, "Error")
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}

	nums, err := stack.Parse(args)
	if err != nil {
		fail()
	}
	if len(nums) == 0 {
		return
	}

	// Read and validate every instruction before executing any.
	var ops []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSuffix(scanner.Text(), "\r")
		if line == "" { // tolerate trailing/blank lines
			continue
		}
		if !stack.Valid(line) {
			fail()
		}
		ops = append(ops, line)
	}
	if scanner.Err() != nil {
		fail()
	}

	s := stack.New(nums)
	for _, op := range ops {
		if err := s.Apply(op); err != nil {
			fail()
		}
	}

	if s.Sorted() {
		fmt.Println("OK")
	} else {
		fmt.Println("KO")
	}
}
