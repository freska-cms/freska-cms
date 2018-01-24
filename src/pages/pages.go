// Package pages represents the page resource
package pages

import (
	"github.com/freska-cms/view/helpers"

	"github.com/freska-cms/freska-cms/src/lib/resource"
	"github.com/freska-cms/freska-cms/src/lib/status"
)

// Page handles saving and retreiving pages from the database
type Page struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	AuthorID int64
	Keywords string
	Name     string
	Summary  string
	Template string
	Text     string
	URL      string
}

// ShowURL returns our canonical url for showing the page
func (p *Page) ShowURL() string {
	return p.URL
}

// ShowTemplate returns the default template if none is set, or the template selected
func (p *Page) ShowTemplate() string {
	if p.Template == "" {
		return "pages/views/templates/default.html.got"
	}
	return p.Template
}

// TemplateOptions provides a set of options for the templates menu
// ids are indexes into the templates array above
func (p *Page) TemplateOptions() []helpers.Selectable {
	var options []helpers.Selectable

	options = append(options, helpers.SelectableOption{Value: "pages/views/templates/default.html.got", Name: "Default"})

	return options
}
