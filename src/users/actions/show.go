package useractions

import (
	"net/http"

	"github.com/freska-cms/auth/can"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"
	"github.com/freska-cms/view"

	"github.com/freska-cms/freska-cms/src/lib/session"
	"github.com/freska-cms/freska-cms/src/users"
)

// HandleShow displays a single user.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

	// Get the user params for id
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the user
	user, err := users.Find(params.GetInt(users.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise access
	currentUser := session.CurrentUser(w, r)
	err = can.Show(user, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(user.CacheKey())
	view.AddKey("user", user)
	view.AddKey("currentUser", currentUser)
	return view.Render()
}
