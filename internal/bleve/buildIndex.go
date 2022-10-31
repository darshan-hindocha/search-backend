package bleve

import (
	"encoding/json"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/analysis/analyzer/keyword"
	"github.com/blevesearch/bleve/v2/analysis/lang/en"
	"github.com/blevesearch/bleve/v2/mapping"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

func buildIndexMapping() (mapping.IndexMapping, error) {

	// a generic reusable mapping for english text
	englishTextFieldMapping := bleve.NewTextFieldMapping()
	englishTextFieldMapping.Analyzer = en.AnalyzerName

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = keyword.Name

	//bookMapping := bleve.NewDocumentMapping()
	//
	//// title + author
	//bookMapping.AddFieldMappingsAt("title", englishTextFieldMapping)
	//bookMapping.AddFieldMappingsAt("author", englishTextFieldMapping)
	//bookMapping.AddFieldMappingsAt("text", keywordFieldMapping)
	//
	//indexMapping := bleve.NewIndexMapping()
	//indexMapping.AddDocumentMapping("books", bookMapping)

	// TEMPORARY --------------------
	beerMapping := bleve.NewDocumentMapping()

	// name
	beerMapping.AddFieldMappingsAt("name", englishTextFieldMapping)

	// description
	beerMapping.AddFieldMappingsAt("description",
		englishTextFieldMapping)

	beerMapping.AddFieldMappingsAt("type", keywordFieldMapping)
	beerMapping.AddFieldMappingsAt("style", keywordFieldMapping)
	beerMapping.AddFieldMappingsAt("category", keywordFieldMapping)

	breweryMapping := bleve.NewDocumentMapping()
	breweryMapping.AddFieldMappingsAt("name", englishTextFieldMapping)
	breweryMapping.AddFieldMappingsAt("description", englishTextFieldMapping)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("beer", beerMapping)
	indexMapping.AddDocumentMapping("brewery", breweryMapping)
	// TEMPORARY ^^^^^^^^^^^^^^^^^^^^^^^

	indexMapping.TypeField = "type"
	indexMapping.DefaultAnalyzer = "en"

	return indexMapping, nil
}

func indexBooks(i bleve.Index) error {

	// open the directory
	dirEntries, err := ioutil.ReadDir(*jsonDir)
	if err != nil {
		return err
	}

	// walk the directory entries for indexing
	log.Printf("Indexing...")
	count := 0
	startTime := time.Now()
	batch := i.NewBatch()
	batchCount := 0
	for _, dirEntry := range dirEntries {
		filename := dirEntry.Name()
		// read the bytes
		jsonBytes, err := ioutil.ReadFile(*jsonDir + "/" + filename)
		if err != nil {
			return err
		}
		// parse bytes as json
		var jsonDoc interface{}
		err = json.Unmarshal(jsonBytes, &jsonDoc)
		if err != nil {
			return err
		}
		ext := filepath.Ext(filename)
		docID := filename[:(len(filename) - len(ext))]
		batch.Index(docID, jsonDoc)
		batchCount++

		if batchCount >= *batchSize {
			err = i.Batch(batch)
			if err != nil {
				return err
			}
			batch = i.NewBatch()
			batchCount = 0
		}
		count++
		if count%1000 == 0 {
			indexDuration := time.Since(startTime)
			indexDurationSeconds := float64(indexDuration) / float64(time.Second)
			timePerDoc := float64(indexDuration) / float64(count)
			log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
		}
	}
	// flush the last batch
	if batchCount > 0 {
		err = i.Batch(batch)
		if err != nil {
			log.Fatal(err)
		}
	}
	indexDuration := time.Since(startTime)
	indexDurationSeconds := float64(indexDuration) / float64(time.Second)
	timePerDoc := float64(indexDuration) / float64(count)
	log.Printf("Indexed %d documents, in %.2fs (average %.2fms/doc)", count, indexDurationSeconds, timePerDoc/float64(time.Millisecond))
	return nil
}
