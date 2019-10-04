package main

import (
	"log"
	"net/http"
	"os"

	"github.com/smartystreets/processor"
)

func main() {

	client := processor.NewAuthenticationClient(http.DefaultClient,
		"https", "us-street-api.smartystreets.com", 
		"9bd279bf-edd1-7564-5d2d-61f58512f0e5", "8qy1ZRT32jG27mLYYtrs")

	pipeline := processor.NewPipeline(os.Stdin, os.Stdout, client, 8)

	if err := pipeline.Process(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}
