// Exercise 7.10: The sort.Interface type can be adapted to other uses. Write a
// function IsPalindrome(s sort.Interface) bool that reports whether the
// sequence s is a palindrome, in other words, reversing the sequence would
// not change it. Assume that the elements at indices i and j are equal if
// !s.Less(i, j) && !s.Less(j, i).
package main

import (
	"fmt"
	"sort"
)

func IsPalindrome(s sort.Interface) bool {
	for i, j := 0, s.Len()-1; i < j; i, j = i+1, j-1 {
		if !(!s.Less(i, j) && !s.Less(j, i)) {
			return false
		}
	}
	return true
}

func main() {
	data := []string{"a", "b", "b", "a"}
	fmt.Println(IsPalindrome(sort.StringSlice(data)))
	fmt.Println(data)
}
