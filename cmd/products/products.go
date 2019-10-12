package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/deanishe/awgo/fuzzy"

	aw "github.com/deanishe/awgo"
)

var (
	logger = log.New(os.Stderr, "[products] ", log.LstdFlags)
	wf     *aw.Workflow

	query   string
	project string
)

func init() {
	flag.StringVar(&query, "query", "", "search query")
	flag.StringVar(&project, "project", "", "google cloud project")

	sopts := []fuzzy.Option{
		fuzzy.AdjacencyBonus(10.0),
		fuzzy.LeadingLetterPenalty(-0.1),
		fuzzy.MaxLeadingLetterPenalty(-3.0),
		fuzzy.UnmatchedLetterPenalty(-0.5),
	}
	wf = aw.New(aw.SortOptions(sopts...))

}

type ProductTemplate struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type urlTemplate struct {
	ProjectID string
}

var productTemplates = []ProductTemplate{
	{
		Name: "BigQuery",
		URL:  "https://console.cloud.google.com/bigquery?project={{.ProjectID}}",
	},
	{
		Name: "Dataflow",
		URL:  "https://console.cloud.google.com/dataflow?project={{.ProjectID}}",
	},
	{
		Name: "Console",
		URL:  "https://console.cloud.google.com/home/dashboard?project={{.ProjectID}}",
	},
	{
		Name: "Logs",
		URL:  "https://console.cloud.google.com/logs/viewer?project={{.ProjectID}}",
	},
	{
		Name: "Storage",
		URL:  "https://console.cloud.google.com/storage/browser?project={{.ProjectID}}",
	},
}

func readProducts() ([]ProductTemplate, error) {
	f, err := os.Open("./products.json")
	if err != nil {
		return nil, err
	}
	byt, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var products []ProductTemplate
	if err := json.Unmarshal(byt, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func run() {
	wf.Args()
	flag.Parse()

	templateArgs := urlTemplate{
		ProjectID: project,
	}

	products, err := readProducts()
	if err != nil {
		wf.FatalError(err)
		return
	}
	for _, p := range products {
		urlTemplate, err := template.New("").Parse(p.URL)
		if err != nil {
			wf.FatalError(err)
		}
		buf := &bytes.Buffer{}
		if err := urlTemplate.Execute(buf, templateArgs); err != nil {
			wf.FatalError(err)
		}
		parsedURLBytes, err := ioutil.ReadAll(buf)
		if err != nil {
			wf.FatalError(err)
		}
		parsedURL := string(parsedURLBytes)
		wf.NewItem(p.Name).Arg(parsedURL).Subtitle(project).UID(parsedURL).Valid(true)
	}

	if query != "" {
		wf.Filter(query)
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
