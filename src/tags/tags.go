// Package tags represents the tag resource
package tags

import (
	"github.com/freska-cms/freska-cms/src/lib/resource"
	"github.com/freska-cms/freska-cms/src/lib/status"
)

// Tag handles saving and retreiving tags from the database
type Tag struct {
	// resource.Base defines behaviour and fields shared between all resources
	resource.Base

	// status.ResourceStatus defines a status field and associated behaviour
	status.ResourceStatus

	DottedIDs string
	Name      string
	ParentID  int64
	Sort      int64
	Summary   string
	URL       string
}
