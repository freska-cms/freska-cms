package tagactions

import (
	"net/http"

	"github.com/freska-cms/auth/can"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"

	"github.com/freska-cms/freska-cms/src/lib/session"
	"github.com/freska-cms/freska-cms/src/tags"
)

// HandleDestroy responds to /tags/n/destroy by deleting the tag.
func HandleDestroy(w http.ResponseWriter, r *http.Request) error {

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

	// Authorise destroy tag
	user := session.CurrentUser(w, r)
	err = can.Destroy(tag, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Destroy the tag
	tag.Destroy()

	// Redirect to tags root
	return server.Redirect(w, r, tag.IndexURL())

}
