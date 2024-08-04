package searcher

import (
	"fmt"

	"github.com/amscotti/HackerJobs/internal/indexer"
	"github.com/blevesearch/bleve/v2"
)

// Search performs a search on the index and returns the results
func Search(jobPostingId int64, queryText string, searchCount int) (map[string]string, error) {
	indexPath := fmt.Sprintf("hackernews_job_postings_%d.bleve", jobPostingId)

	index, err := indexer.CreateIndex(indexPath, jobPostingId)
	if err != nil {
		return nil, fmt.Errorf("error creating/opening index: %v", err)
	}
	defer index.Close()

	query := bleve.NewQueryStringQuery(queryText)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"text"}
	searchRequest.Size = searchCount
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("error performing search: %v", err)
	}

	results := make(map[string]string)
	for _, result := range searchResults.Hits {
		results[result.ID] = fmt.Sprintf("%v", result.Fields["text"])
	}

	return results, nil
}
