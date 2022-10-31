package internal

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var searchEngineType = flag.String("engine", "bleve", "search engine type 'bleve'")

func SetupRouter() *gin.Engine {
	flag.Parse()
	r := gin.Default()

	var db SearchEngine
	if *searchEngineType == "bleve" {
		db = &BleveDB{bleveIndex: nil}
		err := db.Setup()
		if err != nil {
			log.Fatal("error initialising bleve search index")
			return nil
		}
	} else {
		db = &NilDB{}
	}

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Search using query
	r.GET("/search/:query", func(c *gin.Context) {
		query := c.Params.ByName("query")
		log.Print("here is the query: ", query)
		value, err := db.Search(query)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"query": query, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"query": query, "status": "nothing found"})
		}
	})

	return r
}
