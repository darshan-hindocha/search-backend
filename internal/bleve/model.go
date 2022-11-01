package bleve

import (
	"github.com/blevesearch/bleve/v2"
	"github.com/darshan-hindocha/search-backend/internal/db"
	"log"
)

type Book struct {
	Title    string    `json:"title"`
	Author   string    `json:"author"`
	Chapters []Chapter `json:"chapters"`
	Text     string    `json:"text"`
}

type Chapter struct {
	ChapterTitle string      `json:"chapter_title"`
	Paragraphs   []Paragraph `json:"paragraphs"`
}

type Paragraph struct {
	Text string `json:"text"`
}

// DB Implements the db interface
type DB struct {
	BleveIndex bleve.Index
}

func (i *DB) Setup() error {
	var err error
	i.BleveIndex, err = GetBleveIndex()
	if err != nil {
		if err != bleve.ErrorIndexPathDoesNotExist {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func (i *DB) Search(query db.SearchRequest) (db.SearchResponse, error) {
	bleveQuery := bleve.NewQueryStringQuery(query.Query)
	searchRequest := bleve.NewSearchRequest(bleveQuery)
	searchResult, err := i.BleveIndex.Search(searchRequest)
	if err != nil {
		log.Print("error in searching for query: " + query.Query)
		log.Print(err)
		return db.SearchResponse{}, err
	}
	log.Printf("here is the string result: " + searchResult.String())
	res := db.SearchResponse{
		SearchParameters: query,
		SearchResults:    searchResult.String(),
		MaxScore:         searchResult.MaxScore,
		NumberOfMatches:  int32(searchResult.Hits.Len()),
		TimeTaken:        searchResult.Took,
	}
	return res, nil
}
