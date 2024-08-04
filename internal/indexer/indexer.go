package indexer

import (
	"html"
	"log"
	"strconv"
	"sync"

	"github.com/blevesearch/bleve/v2"
	"jaytaylor.com/html2text"

	"github.com/amscotti/HackerJobs/internal/fetcher"
)

// CreateIndex creates a new index or opens an existing one and populates it with job postings
func CreateIndex(indexPath string, postingId int64) (bleve.Index, error) {
	index, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index...")
		index, err = createNewIndex(indexPath)
		if err != nil {
			return nil, err
		}
		err = populateIndex(index, postingId)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		log.Printf("Opening existing index...")
	}
	return index, nil
}

func createNewIndex(indexPath string) (bleve.Index, error) {
	enFieldMapping := bleve.NewTextFieldMapping()
	enFieldMapping.Analyzer = "en"

	byFieldMapping := bleve.NewTextFieldMapping()
	byFieldMapping.Analyzer = "keyword"

	postingMapping := bleve.NewDocumentMapping()
	postingMapping.AddFieldMappingsAt("text", enFieldMapping)
	postingMapping.AddFieldMappingsAt("by", byFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.DefaultMapping = postingMapping

	return bleve.New(indexPath, indexMapping)
}

func populateIndex(index bleve.Index, postingId int64) error {
	posting, err := fetcher.Fetch(postingId)
	if err != nil {
		return err
	}

	log.Printf("Indexing job postings from %s \n", posting.Title)
	log.Printf("Indexing %d postings", len(posting.Kids))

	var wg sync.WaitGroup
	wg.Add(len(posting.Kids))

	for _, item := range posting.Kids {
		go func(i int64, index bleve.Index) {
			defer wg.Done()
			job, err := fetcher.Fetch(i)
			if err != nil {
				log.Printf("Error fetching job %d: %v", i, err)
				return
			}

			text, err := html2text.FromString(html.UnescapeString(job.Text))
			if err != nil {
				log.Printf("Error converting HTML to text for job %d: %v", i, err)
				return
			}
			job.Text = text

			err = index.Index(strconv.FormatInt(i, 10), job)
			if err != nil {
				log.Printf("Error indexing job %d: %v", i, err)
			}
		}(item, index)
	}

	wg.Wait()
	log.Println("Done indexing")
	return nil
}
