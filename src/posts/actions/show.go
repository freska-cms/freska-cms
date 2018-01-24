package postactions

import (
	"net/http"

	"github.com/freska-cms/auth/can"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"
	"github.com/freska-cms/view"

	"github.com/freska-cms/freska-cms/src/lib/session"
	"github.com/freska-cms/freska-cms/src/posts"
)

// HandleShow displays a single post.
func HandleShow(w http.ResponseWriter, r *http.Request) error {

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Find the post
	post, err := posts.Find(params.GetInt(posts.KeyName))
	if err != nil {
		return server.NotFoundError(err)
	}

	// Authorise access
	user := session.CurrentUser(w, r)

	if !post.IsPublished() {
		err = can.Show(post, user)
		if err != nil {
			return server.NotAuthorizedError(err)
		}
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.CacheKey(post.CacheKey())
	view.AddKey("currentUser", user)
	view.AddKey("post", post)
	view.Template("posts/views/templates/default.html.got")
	return view.Render()
}
