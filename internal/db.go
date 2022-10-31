package internal

import (
	"github.com/blevesearch/bleve/v2"
	internalBleve "github.com/darshan-hindocha/search-backend/internal/bleve"
	"log"
)

type SearchEngine interface {
	Setup() error
	Search(string) (string, error)
}

type NilDB struct {
}

func (i *NilDB) Setup() error {
	return nil
}

func (i *NilDB) Search(_ string) (string, error) {
	return "using nil db", nil
}

// BleveDB Implements the SearchEngine interface
type BleveDB struct {
	bleveIndex bleve.Index
}

func (i *BleveDB) Setup() error {
	var err error
	i.bleveIndex, err = internalBleve.GetBleveIndex()
	if err != nil {
		if err != bleve.ErrorIndexPathDoesNotExist {
			log.Fatal(err)
			return err
		}
	}
	return nil
}

func (i *BleveDB) Search(query string) (string, error) {
	bleveQuery := bleve.NewQueryStringQuery(query)
	searchRequest := bleve.NewSearchRequest(bleveQuery)
	searchResult, err := i.bleveIndex.Search(searchRequest)
	if err != nil {
		log.Print("error in searching for query: " + query)
		log.Print(err)
		return "", err
	}
	return searchResult.String(), nil
}
