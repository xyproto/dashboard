package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/unrolled/render"
	"github.com/xyproto/permissions"
	"github.com/xyproto/webhandle"
)

const (
	Version = "0.1"
)

func main() {
	fmt.Println("dashboard ", Version)

	r := render.New(render.Options{})

	n := negroni.Classic()

	perm := permissions.New()
	userstate := perm.UserState()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		data := map[string]string{
			"version_string": "0.1",
			"server":         "server",
			"status":         "up",
			"shortname":      "dx",
			"statuscolor":    "#a0a0ff",
		}

		// !! Reload template !!
		r = render.New(render.Options{})

		// Render and return
		r.HTML(w, http.StatusOK, "index", data)
	})

	mux.HandleFunc("/status", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Has user bob: %v\nLogged in on server: %v\n", userstate.HasUser("bob"), userstate.IsLoggedIn("bob"))
		fmt.Fprintf(w, "Is confirmed: %v\n", userstate.IsConfirmed("bob"))
		fmt.Fprintf(w, "Current user is logged in, has a valid cookie and *user rights*: %v\n", userstate.UserRights(req))
		fmt.Fprintf(w, "Current user is logged in, has a valid cookie and *admin rights*: %v\n", userstate.AdminRights(req))
		fmt.Fprintf(w, "Username stored in cookies (or blank): %v\n", userstate.GetUsername(req))
		fmt.Fprint(w, "\nTry: /register, /confirm, /remove, /login, /logout, /makeadmin and /admin")
	})

	mux.HandleFunc("/hello/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, %v", webhandle.GetLast(req.URL))
	})

	mux.HandleFunc("/register", func(w http.ResponseWriter, req *http.Request) {
		userstate.AddUser("bob", "hunter1", "bob@zombo.com")
		fmt.Fprintf(w, "User bob was created: %v\n", userstate.HasUser("bob"))
	})

	mux.HandleFunc("/confirm", func(w http.ResponseWriter, req *http.Request) {
		userstate.MarkConfirmed("bob")
		fmt.Fprintf(w, "User bob was confirmed: %v\n", userstate.IsConfirmed("bob"))
	})

	mux.HandleFunc("/remove", func(w http.ResponseWriter, req *http.Request) {
		userstate.RemoveUser("bob")
		fmt.Fprintf(w, "User bob was removed: %v\n", !userstate.HasUser("bob"))
	})

	mux.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		userstate.Login(w, "bob")
		fmt.Fprintf(w, "bob is now logged in: %v\n", userstate.IsLoggedIn("bob"))
	})

	mux.HandleFunc("/logout", func(w http.ResponseWriter, req *http.Request) {
		userstate.Logout("bob")
		fmt.Fprintf(w, "bob is now logged out: %v\n", !userstate.IsLoggedIn("bob"))
	})

	mux.HandleFunc("/makeadmin", func(w http.ResponseWriter, req *http.Request) {
		userstate.SetAdminStatus("bob")
		fmt.Fprintf(w, "bob is now administrator: %v\n", userstate.IsAdmin("bob"))
	})

	mux.HandleFunc("/data", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "user page that only logged in users must see!")
	})

	mux.HandleFunc("/admin", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "super secret information that only logged in administrators must see!\n\n")
		usernames, err := userstate.GetAllUsernames()
		if err == nil {
			fmt.Fprintf(w, "list of all users: "+strings.Join(usernames, ", "))
		}
	})

	// Publish the generated SVG as "/circles.svg"
	svgurl := "/circles.svg"
	mux.HandleFunc(svgurl, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprint(w, svgPage().String())
	})

	// Generate a Page that includes the svg image
	page := indexPage(svgurl)
	// Publish the generated Page in a way that connects the HTML and CSS
	page.Publish(mux, "/html", "/style.css", false)

	// Permission middleware

	perm.SetDenyFunction(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Permission denied!")
	})

	// Activate the permission middleware

	n.Use(perm)

	// Already included in negroni.Classic:
	//// Share the files in static
	//n.Use(negroni.NewStatic(http.Dir("static")))

	n.UseHandler(mux)

	n.Run(":3000")
}
