package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudresourcemanager/v1"

	"github.com/deanishe/awgo/fuzzy"

	aw "github.com/deanishe/awgo"
)

const (
	cacheGoogleProjects = "google-projects"
	scopeListProjects   = "https://www.googleapis.com/auth/cloudplatformprojects.readonly"
)

var (
	logger = log.New(os.Stderr, "[projects] ", log.LstdFlags)
	wf     *aw.Workflow

	argRefreshProjects bool
)

type ProjectDescription struct {
	Name   string `json:"name"`
	ID     string `json:"id"`
	Number int64  `json:"number"`
}

func init() {
	sopts := []fuzzy.Option{
		fuzzy.AdjacencyBonus(10.0),
		fuzzy.LeadingLetterPenalty(-0.1),
		fuzzy.MaxLeadingLetterPenalty(-3.0),
		fuzzy.UnmatchedLetterPenalty(-0.5),
	}
	wf = aw.New(aw.SortOptions(sopts...))

	flag.BoolVar(&argRefreshProjects, "refresh", false, "refresh authenticated projects")
}

func FetchGoogleProjects(ctx context.Context) ([]ProjectDescription, error) {
	goog, err := google.DefaultClient(ctx, scopeListProjects)
	if err != nil {
		return nil, fmt.Errorf("failed to create google client: %w", err)
	}

	s, err := cloudresourcemanager.New(goog)
	if err != nil {
		return nil, fmt.Errorf("google resource manager: %w", err)
	}

	allProjects := make([]ProjectDescription, 0)

	nextPageToken := ""
	for {
		res, err := s.Projects.List().Context(ctx).PageToken(nextPageToken).PageSize(200).Do()
		if err != nil {
			return nil, err
		}
		nextPageToken = res.NextPageToken
		for _, project := range res.Projects {
			allProjects = append(allProjects, ProjectDescription{
				Name:   project.Name,
				ID:     project.ProjectId,
				Number: project.ProjectNumber,
			})
		}
		if nextPageToken == "" {
			break
		}
	}

	return allProjects, nil
}

func run() {
	args := wf.Args()
	flag.Parse()
	ctx := context.Background()

	if argRefreshProjects {
		wf.Configure(aw.TextErrors(true))
		logger.Printf("refreshing projects")
		projects, err := FetchGoogleProjects(ctx)
		if err != nil {
			wf.FatalError(err)
		}
		err = wf.Data.StoreJSON(cacheGoogleProjects, projects)
		if err != nil {
			wf.FatalError(err)
		}
		logger.Printf("refresh done")
		wf.SendFeedback()
		return
	}
	var query string
	if len(args) > 0 {
		query = args[0]
	}

	if strings.HasPrefix(query, "-") {
		wf.NewItem("Refresh projects").Arg("-refresh").Autocomplete("-refresh").Valid(false)
		wf.SendFeedback()
		return
	}

	logger.Printf("query=%s", query)
	var projects []ProjectDescription
	if !wf.Data.Exists(cacheGoogleProjects) {
		wf.Fatal(`No projects cached, please run "-refresh"`)
	}
	if err := wf.Data.LoadJSON(cacheGoogleProjects, &projects); err != nil {
		wf.FatalError(err)
	}

	for _, p := range projects {
		wf.NewItem(p.Name).Arg(p.ID).Subtitle(p.ID).UID(p.ID).Valid(true)
	}

	if query != "" {
		wf.Filter(query)
	}

	wf.SendFeedback()
}

func main() {
	wf.Run(run)
}
