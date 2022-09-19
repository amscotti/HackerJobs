package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/blevesearch/bleve/v2"
)

type Posting struct {
	By          string  `json:"by"`
	Descendants int64   `json:"descendants"`
	ID          int64   `json:"id"`
	Kids        []int64 `json:"kids"`
	Score       int64   `json:"score"`
	Text        string  `json:"text"`
	Time        int64   `json:"time"`
	Title       string  `json:"title"`
	Type        string  `json:"type"`
}

const (
	ITEM_URL_BASE = "https://hacker-news.firebaseio.com/v0/item"
)

func fetch(itemId int64) (Posting, error) {
	rsp, err := http.Get(fmt.Sprintf("%s/%d.json", ITEM_URL_BASE, itemId))
	if err != nil {
		return Posting{}, err
	}
	defer rsp.Body.Close()

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return Posting{}, err
	}

	var posting Posting
	if err := json.Unmarshal(data, &posting); err != nil {
		return Posting{}, err
	}

	return posting, nil
}

func createIndex(index bleve.Index, postingId int64) {
	posting, err := fetch(postingId)
	log.Printf("Indexing job postings from %s \n", posting.Title)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	wg.Add(len(posting.Kids))

	for _, item := range posting.Kids {
		go func(i int64, index bleve.Index) {
			job, err := fetch(i)
			if err != nil {
				log.Fatal(err)
			}
			err = index.Index(strconv.FormatInt(i, 10), job)
			if err != nil {
				log.Fatal(err)
			}
			wg.Done()
		}(item, index)
	}

	wg.Wait()
	log.Println("Done indexing")
}

func RemoveHtmlTag(in string) string {
	const pattern = `(<\/?[a-zA-A]+?[^>]*\/?>)*`
	r := regexp.MustCompile(pattern)
	groups := r.FindAllString(in, -1)

	sort.Slice(groups, func(i, j int) bool {
		return len(groups[i]) > len(groups[j])
	})

	for _, group := range groups {
		if strings.TrimSpace(group) != "" {
			in = strings.ReplaceAll(in, group, "")
		}
	}

	return in
}

func search(jobPostingId int64, queryText string, searchCount int) map[string]string {
	indexPath := fmt.Sprintf("hackernews_job_postings_%d.bleve", jobPostingId)
	index, err := bleve.Open(indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index...")
		enFieldMapping := bleve.NewTextFieldMapping()
		enFieldMapping.Analyzer = "en"

		postingMapping := bleve.NewDocumentMapping()
		postingMapping.AddFieldMappingsAt("text", enFieldMapping)

		indexMapping := bleve.NewIndexMapping()
		indexMapping.DefaultMapping = postingMapping

		index, err = bleve.New(indexPath, indexMapping)
		if err != nil {
			log.Fatal(err)
		}
		createIndex(index, jobPostingId)
	} else if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Opening existing index...")
	}

	query := bleve.NewQueryStringQuery(queryText)
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"text"}
	searchRequest.Size = searchCount
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	results := make(map[string]string)
	for _, result := range searchResults.Hits {
		results[result.ID] = fmt.Sprintf("%v", result.Fields["text"])
	}

	return results
}

func main() {
	var (
		jobPostingId int64
		queryText    string
		searchCount  int
	)

	flag.Int64Var(&jobPostingId, "j", 32677265, "Job posting ID from HackerNews")
	flag.StringVar(&queryText, "q", "+text:golang +text:remote", "Text to search for in postings")
	flag.IntVar(&searchCount, "c", 100, "Count of posting to be return")
	flag.Parse()

	links := search(jobPostingId, queryText, searchCount)

	fmt.Printf("Found %d postings for %s\n\n", len(links), queryText)
	fmt.Printf("%-50s %-20s\n", "Link", "Posting Snippet")
	fmt.Printf("%-50s %-20s\n", strings.Repeat("-", 25), strings.Repeat("-", 25))
	for id, text := range links {
		posting := RemoveHtmlTag(html.UnescapeString(strings.Replace(text, "\n", " ", -1)))
		fmt.Printf("%-50s %.75s\n", fmt.Sprintf("https://news.ycombinator.com/item?id=%s", id), posting)
	}
}
