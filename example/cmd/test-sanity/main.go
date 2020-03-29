package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mjm/mpsanity"
)

var (
	projectID = flag.String("project", "", "Sanity project ID")
	dataset   = flag.String("dataset", "production", "Sanity dataset name")
	docID     = flag.String("doc", "", "Sanity document ID to fetch")
	query     = flag.String("query", "", "Sanity query to run")
)

func main() {
	flag.Parse()

	sanity, err := mpsanity.New(*projectID, mpsanity.WithDataset(*dataset))
	if err != nil {
		log.Fatal(err)
	}

	if *docID != "" {
		var doc struct {
			Body        []map[string]interface{} `json:"body"`
			PublishedAt time.Time                `json:"publishedAt"`
			Slug        mpsanity.Slug            `json:"slug"`
		}
		if err := sanity.Doc(context.Background(), *docID, &doc); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", doc)
	} else if *query != "" {
		var res interface{}
		if err := sanity.Query(context.Background(), *query, &res); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", res)
	}
}
