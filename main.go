package main

import (
	"flag"
	"log"

	"github.com/amscotti/HackerJobs/internal/display"
	"github.com/amscotti/HackerJobs/internal/searcher"
)

func main() {
	var (
		jobPostingId int64
		queryText    string
		searchCount  int
	)

	flag.Int64Var(&jobPostingId, "j", 41129813, "Job posting ID from HackerNews")
	flag.StringVar(&queryText, "q", "+text:golang +text:remote", "Text to search for in postings")
	flag.IntVar(&searchCount, "c", 100, "Count of posting to be return")
	flag.Parse()

	links, err := searcher.Search(jobPostingId, queryText, searchCount)
	if err != nil {
		log.Fatalf("Error searching: %v", err)
	}

	display.DisplayResults(links)
}
