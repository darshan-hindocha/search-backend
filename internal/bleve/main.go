package bleve

import (
	"flag"
	"github.com/blevesearch/bleve/v2"
	"log"
	"os"
	"runtime/pprof"
)

var batchSize = flag.Int("batchSize", 100, "batch size for indexing")
var jsonDir = flag.String("jsonDir", "data/", "json directory")
var indexPath = flag.String("index", "search-index.bleve", "index path")
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write mem profile to file")

func GetBleveIndex() (bleve.Index, error) {
	flag.Parse()
	// open the index
	booksIndex, err := bleve.Open(*indexPath)
	if err == bleve.ErrorIndexPathDoesNotExist {
		log.Printf("Creating new index...")
		// create a mapping
		indexMapping, err := buildIndexMapping()
		if err != nil {
			log.Fatal(err)
		}
		booksIndex, err = bleve.New(*indexPath, indexMapping)
		if err != nil {
			log.Fatal(err)
		}

		// index in background
		go func() {
			err = indexBooks(booksIndex)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StopCPUProfile()
			if *memprofile != "" {
				f, err := os.Create(*memprofile)
				if err != nil {
					log.Fatal("profiling error ", err)
				}
				pprof.WriteHeapProfile(f)
				f.Close()
			}
		}()

	} else if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		log.Printf("Opening existing index...")
		return booksIndex, err
	}

	return booksIndex, err
}
