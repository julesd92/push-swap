package main

import (
	"bufio"
	"fmt"
	"os"

	"push-swap/solver"
	"push-swap/stack"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}

	nums, err := stack.Parse(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		os.Exit(1)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for _, op := range solver.Solve(nums) {
		fmt.Fprintln(out, op)
	}
}
