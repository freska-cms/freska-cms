package useractions

import (
	"net/http"

	"github.com/freska-cms/auth/can"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"

	"github.com/freska-cms/freska-cms/src/lib/session"
	"github.com/freska-cms/freska-cms/src/users"
)

// HandleDestroy responds to /users/n/destroy by deleting the user.
func HandleDestroy(w http.ResponseWriter, r *http.Request) error {

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

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise destroy user
	currentUser := session.CurrentUser(w, r)
	err = can.Destroy(user, currentUser)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Destroy the user
	user.Destroy()

	// Redirect to users root
	return server.Redirect(w, r, user.IndexURL())

}
