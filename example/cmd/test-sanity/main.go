package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mjm/mpsanity"
	"github.com/mjm/mpsanity/block"
	"github.com/mjm/mpsanity/patch"
)

var (
	projectID = flag.String("project", "", "Sanity project ID")
	dataset   = flag.String("dataset", "production", "Sanity dataset name")
	docID     = flag.String("doc", "", "Sanity document ID to fetch")
	query     = flag.String("query", "", "Sanity query to run")
	mutate    = flag.Bool("mutate", false, "Test mutations")
)

func main() {
	flag.Parse()

	sanity, err := mpsanity.New(*projectID, mpsanity.WithDataset(*dataset))
	if err != nil {
		log.Fatal(err)
	}

	if *docID != "" {
		var doc struct {
			Body        []block.Block `json:"body"`
			PublishedAt time.Time     `json:"publishedAt"`
			Slug        mpsanity.Slug `json:"slug"`
		}
		if err := sanity.Doc(context.Background(), *docID, &doc); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", doc)
	} else if *query != "" {
		var res []struct {
			Body        []block.Block `json:"body"`
			PublishedAt time.Time     `json:"publishedAt"`
			Slug        mpsanity.Slug `json:"slug"`
		}
		if err := sanity.Query(context.Background(), *query, &res); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%+v\n", res)
	} else if *mutate {
		sanity.Token = os.Getenv("SANITY_TOKEN")

		// var doc struct {
		// 	Type        string                   `json:"_type"`
		// 	Body        []map[string]interface{} `json:"body"`
		// 	PublishedAt time.Time                `json:"publishedAt"`
		// 	Slug        mpsanity.Slug            `json:"slug"`
		// }
		// doc.Type = "micropost"
		// doc.Slug = mpsanity.Slug(fmt.Sprintf("test-post-%d", time.Now().Unix()))
		// doc.PublishedAt = time.Now()
		// doc.Body = []map[string]interface{}{
		// 	{
		// 		"_type": "block",
		// 		"style": "normal",
		// 		"children": []map[string]interface{}{
		// 			{
		// 				"_type": "span",
		// 				"text":  "This is some content.",
		// 			},
		// 		},
		// 	},
		// }
		//
		// if err := sanity.Txn().Create(doc).Commit(context.Background()); err != nil {
		// 	log.Fatal(err)
		// }

		if err := sanity.Txn().Patch("Jk1HkXoAIFyd382ZShFZ1s",
			patch.Set("slug.current", "this-is-a-test"),
			patch.InsertAfter("body[0]",
				block.New("normal", block.Text("This is another paragraph")))).
			Commit(context.Background()); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Mutation applied.")
	}
}
