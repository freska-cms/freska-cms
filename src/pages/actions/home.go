package pageactions

import (
	"net/http"

	"github.com/freska-cms/server"
	"github.com/freska-cms/server/log"
	"github.com/freska-cms/view"

	"github.com/freska-cms/freska-cms/src/lib/session"
	"github.com/freska-cms/freska-cms/src/pages"
	"github.com/freska-cms/freska-cms/src/users"
)

// HandleShowHome serves our home page with a simple template.
func HandleShowHome(w http.ResponseWriter, r *http.Request) error {

	// Demonstrate tracing in log messages
	log.Info(log.Values{"msg": "Home handler", "trace": log.Trace(r)})

	// If we have no users (first run), redirect to setup
	if users.Count() == 0 {
		return server.Redirect(w, r, "/freska/setup")
	}

	// Home fetches the first page with the url '/' and uses it for the home page of the site
	page, err := pages.FindFirst("url=?", "/")
	if err != nil {
		return server.NotFoundError(nil)
	}

	currentUser := session.CurrentUser(w, r)

	view := view.NewWithPath(r.URL.Path, w)
	view.AddKey("title", "Freska CMS app")
	view.AddKey("page", page)
	view.AddKey("currentUser", currentUser)
	view.Template("pages/views/templates/default.html.got")
	return view.Render()
}
