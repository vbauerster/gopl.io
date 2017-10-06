// Exercise 5.11: The instructor of the linear algebra course decides that
// calculus is now a prerequisite. Extend the topoSort function to report
// cycles.
package main

import (
	"fmt"
	"sort"
)

// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"algorithms":     {"data structures"},
	"calculus":       {"linear algebra"},
	"linear algebra": {"calculus"},

	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},

	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"discrete math":         {"intro to programming"},
	"formal languages":      {"discrete math"},
	"networks":              {"operating systems"},
	"operating systems":     {"data structures", "computer organization"},
	"programming languages": {"data structures", "computer organization"},
}

func main() {
	fmt.Println(prereqs["intro to programming"])
	seq, err := topoSort(prereqs)
	if err != nil {
		fmt.Println(err)
	}
	for i, course := range seq {
		fmt.Printf("%.2d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) ([]string, error) {
	var (
		order    []string
		visitAll func([]string)

		seen        = make(map[string]bool)
		visitedPath = make(map[string]bool)

		err error
	)

	visitAll = func(items []string) {
		for _, item := range items {
			if visitedPath[item] {
				err = fmt.Errorf("Cycle at: %s", item)
				return
			}
			if !seen[item] {
				seen[item] = true
				visitedPath[item] = true
				visitAll(m[item])
				delete(visitedPath, item)
				order = append(order, item)
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	visitAll(keys)
	return order, err
}

func contains(items []string, dep string) bool {
	for _, item := range items {
		if item == dep {
			return true
		}
	}
	return false
}
