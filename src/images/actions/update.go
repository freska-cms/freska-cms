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

// HandleUpdateShow renders the form to update a image.
func HandleUpdateShow(w http.ResponseWriter, r *http.Request) error {

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

	// Authorise update image
	user := session.CurrentUser(w, r)
	err = can.Update(image, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.AddKey("image", image)
	return view.Render()
}

// HandleUpdate handles the POST of the form to update a image
func HandleUpdate(w http.ResponseWriter, r *http.Request) error {

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

	// Authorise update image
	user := session.CurrentUser(w, r)
	err = can.Update(image, user)
	if err != nil {
		return server.NotAuthorizedError(err)
	}

	// Validate the params, removing any we don't accept
	imageParams := image.ValidateParams(params.Map(), images.AllowedParams())

	err = image.Update(imageParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to image
	return server.Redirect(w, r, image.ShowURL())
}
