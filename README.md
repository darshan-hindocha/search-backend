# search-backend

## handbook

in `cmd/search/main.go` is where the search engine is rooted.

run `go run cmd/search/main.go` to start the service using the bleve search engine

the endpoints are described in `internal/router.go`

## internal

the `internal/*` folder contains the stuff needed to operate this service.

### internal/db

the `internal/db.go` file describes an interface for the db.

the paradigm is that whichever search engine we use, it should implement the db interface.

eventually db will include the users, collections and extracts, and a fine-grained text-corpus

### internal/db/bleve

pycharm very kindly reminds me that bleve is a typo, but it's not.

the `internal/db/bleve/*` has everything we need to implement bleve as a search engine within the db.

### internal/router

`internal/router.go` contains the gin api router that offers endpoints to interact with the db and it's methods
