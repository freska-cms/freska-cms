package pageactions

import (
	"net/http"

	"github.com/freska-cms/auth/can"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"
	"github.com/freska-cms/view"

	"github.com/freska-cms/freska-cms/src/lib/session"
	"github.com/freska-cms/freska-cms/src/pages"
	"github.com/freska-cms/freska-cms/src/users"
)

// HandleCreateShow serves the create form via GET for pages.
func HandleCreateShow(w http.ResponseWriter, r *http.Request) error {

	page := pages.New()

	// Authorise
	user := session.CurrentUser(w, r)
	err := can.Create(page, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Fetch the users
	authors, err := users.FindAll(users.Where("role=?", users.Admin))
	if err != nil {
		return server.InternalError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("page", page)
	view.AddKey("authors", authors)
	view.AddKey("currentUser", user)
	return view.Render()
}

// HandleCreate handles the POST of the create form for pages
func HandleCreate(w http.ResponseWriter, r *http.Request) error {

	page := pages.New()

	// Check the authenticity token
	err := session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise
	user := session.CurrentUser(w, r)
	err = can.Create(page, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Setup context
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Validate the params, removing any we don't accept
	pageParams := page.ValidateParams(params.Map(), pages.AllowedParams())

	id, err := page.Create(pageParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to the new page
	page, err = pages.Find(id)
	if err != nil {
		return server.InternalError(err)
	}

	return server.Redirect(w, r, page.IndexURL())
}
