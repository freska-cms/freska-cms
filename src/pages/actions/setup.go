package pageactions

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/freska-cms/auth"
	"github.com/freska-cms/mux"
	"github.com/freska-cms/server"
	"github.com/freska-cms/server/log"
	"github.com/freska-cms/view"

	"github.com/freska-cms/freska-cms/src/pages"
	"github.com/freska-cms/freska-cms/src/users"
)

// HandleSetupShow responds to GET /freska/setup
func HandleSetupShow(w http.ResponseWriter, r *http.Request) error {

	// Check we have no users, if not bail out
	if users.Count() != 0 {
		return server.NotAuthorizedError(nil, "Users already exist")
	}

	// Render the template
	view := view.NewRenderer(w, r)
	view.Template("pages/views/setup.html.got")
	return view.Render()
}

// HandleSetup responds to POST /freska/setup
func HandleSetup(w http.ResponseWriter, r *http.Request) error {

	// Check we have no users, if not bail out
	if users.Count() != 0 {
		return server.NotAuthorizedError(nil, "Users already exist")
	}

	// Fetch the  params
	params, err := mux.Params(r)
	if err != nil {
		return server.InternalError(err)
	}

	// Convert the password param to a password_hash
	hash, err := auth.HashPassword(params.Get("password"))
	if err != nil {
		return server.InternalError(err, "Problem hashing password")
	}

	// Take the details given and create the first user
	userParams := map[string]string{
		"email":         params.Get("email"),
		"password_hash": hash,
		"name":          nameFromEmail(params.Get("email")),
		"status":        "100",
		"role":          "100",
		"title":         "Administrator",
	}

	user := users.New()
	uid, err := user.Create(userParams)
	if err != nil {
		return server.InternalError(err)
	}

	user, err = users.Find(uid)
	if err != nil {
		return server.InternalError(err, "Error creating user")
	}
	// Login this user automatically - save cookie
	session, err := auth.Session(w, r)
	if err != nil {
		log.Info(log.V{"msg": "login failed", "user_id": user.ID, "status": http.StatusInternalServerError})
	}

	// Success, log it and set the cookie with user id
	session.Set(auth.SessionUserKey, fmt.Sprintf("%d", user.ID))
	session.Save(w)

	// Log action
	log.Info(log.V{"msg": "login", "user_email": user.Email, "user_id": user.ID})

	// Create a welcome home page
	pageParams := map[string]string{
		"status": "100",
		"name":   "Freska",
		"url":    "/",
		"text":   "<section class=\"padded\"><h1>Welcome to Freska</h1><p><a href=\"/pages/1/update\">Edit this page</a></p></section>",
	}
	_, err = pages.New().Create(pageParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Create another couple of simple pages as examples (about and privacy)
	pageParams = map[string]string{
		"status": "100",
		"name":   "About Us",
		"url":    "/about",
		"text":   "<section class=\"narrow\"><h1>About us</h1><p>About us</p></section>",
	}
	_, err = pages.New().Create(pageParams)
	if err != nil {
		return server.InternalError(err)
	}
	pageParams = map[string]string{
		"status": "100",
		"name":   "Privacy Policy",
		"url":    "/privacy",
		"text":   "<section class=\"narrow\"><h1>Privacy Policy</h1><p>We respect your privacy.</p></section>",
	}
	_, err = pages.New().Create(pageParams)
	if err != nil {
		return server.InternalError(err)
	}

	// Redirect to the home page (newly set up we hope)
	return server.Redirect(w, r, "/")
}

// nameFromEmail grabs a name from an email address
func nameFromEmail(e string) string {
	// Split email on @, and separate by removing . or _
	parts := strings.Split(e, "@")
	if len(parts) > 0 {
		n := strings.Replace(parts[0], ".", " ", -1)
		n = strings.Replace(n, "_", " ", -1)
		return n
	}

	return e
}
