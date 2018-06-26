package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jarlefosen/alfred-gcloud-shortcuts/gcloud-butler/alfred"

	"github.com/jarlefosen/alfred-gcloud-shortcuts/gcloud-butler/filehelp"
	goog "github.com/jarlefosen/alfred-gcloud-shortcuts/gcloud-butler/google"
)

var (
	butlerlog       = ".butlerlog"
	gcloudCacheFile = ".gcloud-projects"
)

func filter(projects []goog.Project, filterArgs []string) []goog.Project {
	for _, arg := range filterArgs {
		var filtered []goog.Project
		arg = strings.ToLower(arg)
		for _, p := range projects {
			cmpString := strings.ToLower(fmt.Sprintf("%s %s", p.Name, p.ProjectID))
			if strings.Contains(cmpString, arg) {
				filtered = append(filtered, p)
			}
		}
		projects = filtered
	}
	return projects
}

func refresh() error {
	ctx := context.Background()
	p, err := goog.FetchProjects(ctx)
	if err != nil {
		panic(err)
	}
	return goog.SaveProjects(ctx, filehelp.RelativePath(gcloudCacheFile), p)
}

func accessLog() error {
	buff := bytes.NewBuffer([]byte{})
	fmt.Fprintf(buff, "Butler called: %+v\n", os.Args)
	return ioutil.WriteFile(filehelp.RelativePath(butlerlog), buff.Bytes(), os.ModeAppend)
}

func printProjects(filterArgs []string) {
	projects, err := goog.ListProjects(filehelp.RelativePath(gcloudCacheFile))
	if err != nil {
		panic(err)
	}
	projects = filter(projects, filterArgs)
	if len(projects) == 0 {
		fmt.Fprint(os.Stdout, alfred.RefreshMsgXML)
		return
	}
	items := alfred.ProjectToAlfredModel(projects...)
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
		printProjects(args[1:])
	case "refresh":
		err := refresh()
		if err != nil {
			panic(err)
		}
	default:
		fmt.Fprint(os.Stdout, alfred.InvalidCmdXML)
		panic(fmt.Errorf("Unknown command"))
	}
}
