package imageactions

import (
	"net/http"

	"github.com/freska-cms/auth/can"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"

	"github.com/freska-cms/freska-cms/src/images"
	"github.com/freska-cms/freska-cms/src/lib/session"
)

// HandleDestroy responds to /images/n/destroy by deleting the image.
func HandleDestroy(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the image
	image, err := images.Find(params.GetInt(images.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Check the authenticity token
	err = session.CheckAuthenticity(w, r)
	if err != nil {
		return err
	}

	// Authorise destroy image
	user := session.CurrentUser(w, r)
	err = can.Destroy(image, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Destroy the image
	image.Destroy()

	// Redirect to images root
	return server.Redirect(w, r, image.IndexURL())

}
