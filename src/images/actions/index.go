package imageactions

import (
	"net/http"

	"github.com/freska-cms/auth/can"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"
	"github.com/freska-cms/view"

	"github.com/freska-cms/freska-cms/src/images"
	"github.com/freska-cms/freska-cms/src/lib/session"
)

// HandleIndex displays a list of images.
func HandleIndex(w http.ResponseWriter, r *http.Request) error {

	// Authorise list image
	user := session.CurrentUser(w, r)
	err := can.List(images.New(), user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Get the params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Build a query
	q := images.Query()

	// Order by required order, or default to id asc
	switch params.Get("order") {

	case "1":
		q.Order("created_at desc")

	case "2":
		q.Order("updated_at desc")

	default:
		q.Order("id asc")
	}

	// Filter if requested
	filter := params.Get("filter")
	if len(filter) > 0 {
		q.Where("name ILIKE ?", filter)
	}

	// Fetch the images
	results, err := images.FindAll(q)
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("filter", filter)
	view.AddKey("images", results)
	return view.Render()
}
