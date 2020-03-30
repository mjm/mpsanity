package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/exporters/trace/stdout"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/mjm/mpsanity"
	"github.com/mjm/mpsanity/mpapi"
)

var (
	projectID = flag.String("project", "", "Sanity project ID")
	dataset   = flag.String("dataset", "production", "Sanity dataset name")

	port = flag.String("port", "9090", "Port to listen on for HTTP")
)

func main() {
	flag.Parse()

	exporter, err := stdout.NewExporter(stdout.Options{PrettyPrint: true})
	if err != nil {
		log.Fatal(err)
	}
	tp, err := sdktrace.NewProvider(sdktrace.WithConfig(sdktrace.Config{DefaultSampler: sdktrace.AlwaysSample()}),
		sdktrace.WithSyncer(exporter))
	if err != nil {
		log.Fatal(err)
	}
	global.SetTraceProvider(tp)

	sanity, err := mpsanity.New(*projectID,
		mpsanity.WithDataset(*dataset),
		mpsanity.WithToken(os.Getenv("SANITY_TOKEN")))
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", mpapi.New(sanity))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
