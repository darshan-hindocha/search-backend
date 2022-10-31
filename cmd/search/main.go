package main

import "github.com/darshan-hindocha/search-backend/internal"

// Initialise the app
// // connect to db

func main() {
	r := internal.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
