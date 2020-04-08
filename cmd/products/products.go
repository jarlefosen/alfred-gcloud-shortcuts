package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/url"
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

	queryParams := map[string]string{
		"authuser": wf.Config.Get("authuser"),
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
		urlBytes, err := ioutil.ReadAll(buf)
		if err != nil {
			wf.FatalError(err)
			return
		}

		generatedURL := appendURLParameters(string(urlBytes), queryParams)
		wf.NewItem(p.Name).Arg(generatedURL).Subtitle(project).UID(string(urlBytes)).Valid(true)
	}

	if query != "" {
		wf.Filter(query)
	}

	wf.SendFeedback()
}

func appendURLParameters(urlString string, keyval map[string]string) string {
	u, err := url.ParseRequestURI(urlString)
	if err != nil {
		log.Printf("silently failed to parse url: %v", err)
		return urlString
	}
	q := u.Query()
	for key, val := range keyval {
		if key == "" || val == "" {
			continue
		}
		q.Add(key, val)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func main() {
	wf.Run(run)
}
