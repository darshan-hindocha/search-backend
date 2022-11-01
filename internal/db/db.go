package db

import (
	"time"
)

type SearchEngine interface {
	Setup() error
	Search(SearchRequest) (SearchResponse, error)
}

// NilDB is a nil Implementation of SearchEngine
type NilDB struct {
}

func (i *NilDB) Setup() error {
	return nil
}

func (i *NilDB) Search(_ SearchRequest) (SearchResponse, error) {
	return SearchResponse{}, nil
}

type SearchRequest struct {
	Query     string     `json:"query"`
	QueryType QueryTypes `json:"query_type"`
}

type QueryTypes string

const (
	ExactMatch       QueryTypes = "exact"
	AnalyseThenMatch            = "preprocess"
	Fuzzy                       = "fuzzy"
	Default                     = ""
)

// SearchResponse defines the response schema that the frontend receives when searching for documents
type SearchResponse struct {
	SearchParameters SearchRequest `json:"search_parameters"`
	SearchResults    string        `json:"search_results"`
	MaxScore         float64       `json:"max_score,omitempty"`
	NumberOfMatches  int32         `json:"number_of_matches"`
	TimeTaken        time.Duration `json:"time_taken,omitempty"`
}

type SearchResult struct {
	Field   string `json:"field"`
	Content string `json:"content"`
}
