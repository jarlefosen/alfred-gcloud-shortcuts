package alfred

import (
	"encoding/xml"

	"github.com/jarlefosen/alfred-gcloud-shortcuts/gcloud-butler/google"
)

type Items struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
}

type Item struct {
	XMLName  xml.Name `xml:"item"`
	Title    string   `xml:"title"`
	Subtitle string   `xml:"subtitle"`
	Argument string   `xml:"arg"`
	Icon     string   `xml:"icon"`
}

// ProjectToAlfredModel turns a google.Project specification into a usable Alfred 3 model
func ProjectToAlfredModel(projects ...google.Project) Items {
	var ap []Item
	for _, p := range projects {
		ap = append(ap, Item{
			Title:    p.Name,
			Subtitle: p.ProjectID,
			Argument: p.ProjectID,
		})
	}

	return Items{Items: ap}
}

const (
	// RefreshMsgXML is an Alfred message when there are no projects to be listed
	RefreshMsgXML = `
<items>
  <item uid="none" valid="NO">
		<title>No projects found</title>
		<subtitle>Refresh with "g-refresh"</subtitle>
  </item>
</items>`

	// InvalidCmdXML is a predefined message for when you try to run an unknown command
	InvalidCmdXML = `
<items>
	<item uid="none" valid="NO">
		<title>Invalid command</title>
	</item>
</items>`
)
