package pageactions

import (
	"net/http"

	"github.com/freska-cms/auth/can"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"

	"github.com/freska-cms/freska-cms/src/lib/session"
	"github.com/freska-cms/freska-cms/src/pages"
)

// HandleDestroy responds to /pages/n/destroy by deleting the page.
func HandleDestroy(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the page
	page, err := pages.Find(params.GetInt(pages.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise destroy page
	user := session.CurrentUser(w, r)
	err = can.Destroy(page, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Destroy the page
	page.Destroy()

	// Redirect to pages root
	return server.Redirect(w, r, page.IndexURL())

}
