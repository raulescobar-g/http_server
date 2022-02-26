package main

import (
	"fmt"
)

func printMin(nums []int) {

	if len(nums) == 0 {
		fmt.Printf("Array empty")
	}

	min := nums[0]

	for _, n := range nums {
		if n < min {
			min = n
		}
	}
	fmt.Printf("%d \n", min)
}

func testLinkedList() {
	linky := List{
		size: 0,
	}

	linky.insert(1)
	linky.insert(2)
	linky.insert(10)
	linky.remove(1)
	linky.print()
}

func main() {

	fmt.Printf("hello world \n")
}
