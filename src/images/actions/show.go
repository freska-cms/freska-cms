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

// HandleShow displays a single image.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

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

	// Authorise access
	user := session.CurrentUser(w, r)
	err = can.Show(image, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(image.CacheKey())
	view.AddKey("image", image)
	return view.Render()
}
