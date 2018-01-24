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

// HandleUpdateShow renders the form to update a tag.
func HandleUpdateShow(w http.ResponseWriter, r *http.Request) error {

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

	// Authorise update tag
	user := session.CurrentUser(w, r)
	err = can.Update(tag, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("currentUser", user)
	view.AddKey("tag", tag)
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a tag
func HandleUpdate(w http.ResponseWriter, r *http.Request) error {

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

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise update tag
	user := session.CurrentUser(w, r)
	err = can.Update(tag, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Validate the params, removing any we don't accept
	tagParams := tag.ValidateParams(params.Map(), tags.AllowedParams())

	err = tag.Update(tagParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to tag
	return server.Redirect(w, r, tag.ShowURL())
}
