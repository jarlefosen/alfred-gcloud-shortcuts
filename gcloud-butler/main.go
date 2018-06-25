package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/oauth2/google"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
)

const (
	// ScopeListProjects is a scope for listing google cloud projects
	ScopeListProjects = "https://www.googleapis.com/auth/cloudplatformprojects.readonly"
)

var (
	projectsCacheFile = ".gcloud-projects"
	butlerlog         = ".butlerlog"
)

type items struct {
	Items []item `xml:"item"`
}

type item struct {
	Name      string `json:"name" xml:"title"`
	ProjectID string `json:"project_id" xml:"arg"`
}

func filter(f string) []item {
	data, err := readFile(projectsCacheFile)
	if err != nil {
		panic(err)
	}

	var projects []item
	err = json.Unmarshal(data, &projects)
	if err != nil {
		panic(err)
	}

	var filtered []item
	for _, p := range projects {
		cmpString := fmt.Sprintf("%s %s", p.Name, p.ProjectID)
		if strings.Contains(cmpString, f) {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func refresh() error {
	ctx := context.Background()
	c, err := google.DefaultClient(ctx, ScopeListProjects)
	if err != nil {
		return errors.Wrapf(err, "could not create default google client")
	}
	s, err := cloudresourcemanager.New(c)
	if err != nil {
		return errors.Wrapf(err, "failed to create cloudresourcemanager")
	}
	res, err := s.Projects.List().Do()
	if err != nil {
		return errors.Wrapf(err, "failed to list google cloud projects")
	}

	var projects []item
	for _, project := range res.Projects {
		p := item{
			Name:      project.Name,
			ProjectID: project.ProjectId,
		}
		projects = append(projects, p)
	}

	f, err := openAbsoluteFile(projectsCacheFile, true)
	if err != nil {
		panic(errors.Wrapf(err, "failed to open gcloud file cache"))
	}

	byt, err := json.Marshal(projects)
	if err != nil {
		panic(err)
	}
	_, err = fmt.Fprintf(f, "%s", byt)
	if err != nil {
		panic(err)
	}

	return nil
}

func accessLog() error {
	f, err := openAbsoluteFile(butlerlog, false)
	if err != nil {
		return err
	}
	fmt.Fprintf(f, "Butler called: %+v\n", os.Args)
	f.Close()
	return nil
}

func getProgramDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func readFile(filename string) ([]byte, error) {
	dir := getProgramDir()
	fullPath := fmt.Sprintf("%s/%s", dir, filename)
	return ioutil.ReadFile(fullPath)
}

func openAbsoluteFile(filename string, truncate bool) (*os.File, error) {
	dir := getProgramDir()
	mode := os.O_APPEND | os.O_RDWR | os.O_CREATE
	if truncate {
		mode = mode | os.O_TRUNC
	}
	return os.OpenFile(fmt.Sprintf("%s/%s", dir, filename), mode, 0660)
}

func printRefreshMessage() {
	fmt.Fprintf(os.Stdout, `
		<?xml version="1.0"?>
		<items>
			<item uid="example" arg="NOTHING" valid="YES" autocomplete="example" type="file">
			<title>Run: "g-refresh"</title>
			</item>
		</items>
		`)
}

func printProjects(filterStr string) {
	projects := filter(filterStr)
	if len(projects) == 0 {
		printRefreshMessage()
		return
	}
	items := items{projects}
	b, err := xml.Marshal(items)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stdout, "%s", b)
}

func main() {

	// Log queried data in .butlerlog for debug purposes
	if os.Getenv("DEBUG") != "" {
		accessLog()
	}

	args := os.Args[1:]
	switch args[0] {
	case "filter":
		printProjects(strings.Join(args[1:], " "))
	case "refresh":
		err := refresh()
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("Unknown command"))
	}
}
