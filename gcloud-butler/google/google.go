package google

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

const (
	// ScopeListProjects is a scope for listing google cloud projects
	ScopeListProjects = "https://www.googleapis.com/auth/cloudplatformprojects.readonly"
)

// Project is a container for Google Cloud project specifications
type Project struct {
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
}

// ListProjects returns list of cached projects
func ListProjects(filepath string) ([]Project, error) {
	byt, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var projects []Project
	err = json.Unmarshal(byt, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

// SaveProjects saves a list of projects to a cache file
func SaveProjects(ctx context.Context, filepath string, projects []Project) error {
	byt, err := json.Marshal(projects)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, byt, 0660)
}

// FetchProjects fetches all projects associated with your Google account
func FetchProjects(ctx context.Context) ([]Project, error) {
	c, err := google.DefaultClient(ctx, ScopeListProjects)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create google client")
	}
	s, err := cloudresourcemanager.New(c)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cloudresourcemanager")
	}
	res, err := s.Projects.List().Do()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list google cloud projects")
	}

	var projects []Project
	for _, project := range res.Projects {
		p := Project{
			Name:      project.Name,
			ProjectID: project.ProjectId,
		}
		projects = append(projects, p)
	}

	return projects, nil
}
