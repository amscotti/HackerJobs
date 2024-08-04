package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/amscotti/HackerJobs/internal/models"
)

const ITEM_URL_BASE = "https://hacker-news.firebaseio.com/v0/item"

// FetchURL is a variable that holds the function to construct the URL
var FetchURL = func(itemId int64) (string, error) {
    return fmt.Sprintf("%s/%d.json", ITEM_URL_BASE, itemId), nil
}

// Fetch retrieves a posting from the Hacker News API by its ID
func Fetch(itemId int64) (models.Posting, error) {
    url, err := FetchURL(itemId)
    if err != nil {
        return models.Posting{}, err
    }

    rsp, err := http.Get(url)
    if err != nil {
        return models.Posting{}, err
    }
    defer rsp.Body.Close()

    data, err := io.ReadAll(rsp.Body)
    if err != nil {
        return models.Posting{}, err
    }

    var posting models.Posting
    if err := json.Unmarshal(data, &posting); err != nil {
        return models.Posting{}, err
    }

    return posting, nil
}
