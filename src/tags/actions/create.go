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

// HandleCreateShow serves the create form via GET for tags.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	tag := tags.New()

	// Authorise
	user := session.CurrentUser(w, r)
	err := can.Create(tag, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", user)
	view.AddKey("tag", tag)
	return view.Render()
}

// HandleCreate handles the POST of the create form for tags
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	tag := tags.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise
	user := session.CurrentUser(w, r)
	err = can.Create(tag, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Setup context
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Validate the params, removing any we don't accept
	tagParams := tag.ValidateParams(params.Map(), tags.AllowedParams())

	id, err := tag.Create(tagParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to the new tag
	tag, err = tags.Find(id)
	if err != nil {
		return server.InternalError(err)
	}

	return server.Redirect(w, r, tag.IndexURL())
}
