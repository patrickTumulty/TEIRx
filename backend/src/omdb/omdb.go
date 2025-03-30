package omdb

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"teirxserver/src/cfg"
	"teirxserver/src/txlog"
)

const OMDB_BASE_URL string = "http://www.omdbapi.com/"

type OmdbSearchResponse struct {
	Response     string           `json:"Response"`
	Search       []OmdbSearchItem `json:"Search"`
	TotalResults string           `json:"totalResults"`
}

type OmdbSearchItem struct {
	Title  string `json:"Title"`
	Year   string `json:"Year"`
	ImdbID string `json:"imdbID"`
	Type   string `json:"Type"`
	Poster string `json:"Poster"`
}

func (item *OmdbSearchItem) ToJson() map[string]any {
	m := make(map[string]any)
	m["title"] = item.Title
	m["year"] = item.Year
	m["imdb_id"] = item.ImdbID
	m["type"] = item.Type
	m["poster"] = item.Poster
	return m
}

func OmdbSearch(query string) ([]OmdbSearchItem, error) {
	var searchResponse OmdbSearchResponse
	searchItems := &searchResponse.Search

	appConfig := cfg.GetAppConfig()

	parsedURL, err := url.Parse(OMDB_BASE_URL)
	if err != nil {
		txlog.TxLogError("Error parsing OMDb URL: '%s'", err)
		return *searchItems, err
	}

	params := parsedURL.Query()
	params.Set("type", "movie")
	params.Set("s", query)
	params.Set("apiKey", appConfig.Keys.Omdb)

	parsedURL.RawQuery = params.Encode()

	url := parsedURL.String()

	resp, err := http.Get(url)
	if err != nil {
		txlog.TxLogError("Error performing get call: '%s'", err)
		return *searchItems, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		txlog.TxLogError("Error processing OMDb response: '%s'", err)
		return *searchItems, err
	}

	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		txlog.TxLogError("Error unmarshalling OMDb response: '%s'", err)
		return *searchItems, err
	}

	return *searchItems, nil
}
