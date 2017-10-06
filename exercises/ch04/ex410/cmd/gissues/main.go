// Exercise 4.10: Modify issues to report the results in age categories, say less than a month old, less than a year old, and more than a year old.
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github"
)

const (
	Day   = 24 * time.Hour
	Month = 30 * Day  // duration in nanosecods
	Year  = 365 * Day // duration in nanosecods
)

func display(item *github.Issue) {
	fmt.Printf("#%-5d %9.9s %.55s\n", item.Number, item.User.Login, item.Title)
}

func main() {
	result, err := github.SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d issuses:\n", result.TotalCount)

	groups := make(map[string][]*github.Issue)

	for _, item := range result.Items {
		age := time.Now().Sub(item.CreatedAt)
		if age < Month {
			groups["month"] = append(groups["month"], item)
		} else if age <= Year {
			groups["year"] = append(groups["year"], item)
		} else {
			groups["year+"] = append(groups["year+"], item)
		}
	}

	for _, k := range [...]string{"month", "year", "year+"} {
		fmt.Printf("%q old issues:\n", k)
		for _, item := range groups[k] {
			display(item)
		}
	}
}
