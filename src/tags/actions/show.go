package tagactions

import (
	"net/http"

	"github.com/freska-cms/auth/can"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"
	"github.com/freska-cms/view"

	"github.com/freska-cms/freska-cms/src/lib/session"
	"github.com/freska-cms/freska-cms/src/tags"
)

// HandleShow displays a single tag.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the tag
	tag, err := tags.Find(params.GetInt(tags.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise access
	user := session.CurrentUser(w, r)
	err = can.Show(tag, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(tag.CacheKey())
	view.AddKey("currentUser", user)
	view.AddKey("tag", tag)
	return view.Render()
}
