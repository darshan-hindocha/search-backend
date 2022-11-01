package internal

import (
	"flag"
	"github.com/darshan-hindocha/search-backend/internal/bleve"
	db2 "github.com/darshan-hindocha/search-backend/internal/db"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var searchEngineType = flag.String("engine", "bleve", "search engine type 'bleve'")

func SetupRouter() *gin.Engine {
	flag.Parse()
	r := gin.Default()

	var db db2.SearchEngine
	if *searchEngineType == "bleve" {
		db = &bleve.DB{BleveIndex: nil}
		err := db.Setup()
		if err != nil {
			log.Fatal("error initialising bleve search index")
			return nil
		}
	} else {
		db = &db2.NilDB{}
	}

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Search using query
	r.GET("/search/:query", func(c *gin.Context) {
		query := c.Params.ByName("query")
		log.Print("here is the query: ", query)
		res, err := db.Search(db2.SearchRequest{
			Query: query,
		})
		if err != nil {
			log.Print("error searching the search engine ", err)
			c.JSON(http.StatusOK, gin.H{"query": query, "status": "nothing found"})
		}

		if err == nil {
			c.JSON(http.StatusOK, gin.H{
				"results":                 res.SearchResults,
				"max_score":               res.MaxScore,
				"number_of_matches":       res.NumberOfMatches,
				"time_taken_microseconds": res.TimeTaken.Microseconds(),
				"query":                   res.SearchParameters.Query,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"query": query, "status": "nothing found"})
		}
	})

	return r
}
