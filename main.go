package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/xyproto/fizz"
	"github.com/martini-contrib/auth"

)

const (
	Version = "0.1"
)

type MenuItem struct {
	Url      string
	Text     string
	Selected bool
}

type Menu struct {
	Title     string
	MenuItems []MenuItem
}

type Post struct {
	Title string
	Body  string
}

type PageData struct {
	Text  map[string]string
	Menu  Menu
	Posts []Post
}

func GenerateMenu(active int) Menu {

	// Create a menu
	var menu Menu
	menu.Title = "Menu"

	var menuitem1 MenuItem
	menuitem1.Url = "/"
	menuitem1.Text = "Hi"
	if active == 0 {
		menuitem1.Selected = true
	}

	var menuitem2 MenuItem
	menuitem2.Url = "/admin"
	menuitem2.Text = "Admin"
	if active == 1 {
		menuitem2.Selected = true
	}

	var menuitem3 MenuItem
	menuitem3.Url = "/forum"
	menuitem3.Text = "Forum"
	if active == 2 {
		menuitem3.Selected = true
	}

	var menuitem4 MenuItem
	menuitem4.Url = "/logout"
	menuitem4.Text = "Logout"
	if active == 3 {
		menuitem4.Selected = true
	}

	menu.MenuItems = append(menu.MenuItems, menuitem1)
	menu.MenuItems = append(menu.MenuItems, menuitem2)
	menu.MenuItems = append(menu.MenuItems, menuitem3)
	menu.MenuItems = append(menu.MenuItems, menuitem4)

	return menu
}

func main() {
	fmt.Println("bumpfriend ", Version)

	r := render.Renderer(render.Options{})

	m := martini.Classic()

    m.Use(r)

	fizz := fizz.NewWithRedisConf(7, "")
	userstate := fizz.UserState()
	perm := fizz.Perm()

	m.Get("/login/:password", func(w http.ResponseWriter, req *http.Request, params martini.Params) string {
		if userstate.CorrectPassword("admin", params["password"]) {
			userstate.Login(w, "admin")
			return "Wrong password"
		}
		return fmt.Sprintf("Logged in as administrator: %v\n", userstate.AdminRights(req))
	})

	m.Get("/logout", func(w http.ResponseWriter) string {
		return "Logged out"
	})

	// Public page
	m.Get("/", func() string {
		return "hi"
	})

	// Admin panel
	m.Get("/admin", func(r render.Render) {
		var page PageData
		page.Text = map[string]string{
			"title":    "Bumpfriend",
			"subtitle": "Admin panel",
		}

		// Generate a menu where item 0 is active
		page.Menu = GenerateMenu(0)

		// !! Reload template !!
		//r = render.New(render.Options{})

		// Render the specified templates/.tmpl file as HTML and return
		r.HTML(http.StatusOK, "admin", page)
	})

	// For testing the API
	perm.AddPublicPath("/api/1.0/")
	m.Post("/api/1.0/", func(r render.Render) {
		r.JSON(200, map[string]interface{}{"hello": "fjaselus"})
	})

	m.Get("/api/1.0/", func(r render.Render) {
		r.JSON(200, map[string]interface{}{"hello": "dolly"})
	})

	// For adding users
	m.Post("/api/1.0/create/:username/:password", func(params martini.Params, r render.Render) {
		username := params["username"]
		password := params["password"]
		if userstate.HasUser(username) {
			r.JSON(200, map[string]interface{}{"error": "user " + username + " already exists"})
			return
		}
		userstate.AddUser(username, password, "")
		if userstate.HasUser(username) {
			r.JSON(200, map[string]interface{}{"create": "true"})
		} else {
			r.JSON(200, map[string]interface{}{"error": "user " + username + " was not created"})
		}
	})

	// For logging in
	m.Post("/api/1.0/login/:username/:password", func(w http.ResponseWriter, params martini.Params, r render.Render) {
		username := params["username"]
		password := params["password"]
		if userstate.CorrectPassword(username, password) {
			userstate.SetLoggedIn(username)
		}
		if !userstate.IsLoggedIn(username) {
			r.JSON(200, map[string]interface{}{"error": "could not log in " + username})
			return
		}
		r.JSON(200, map[string]interface{}{"login": "true"})
	})

	// For logging out
	m.Post("/api/1.0/logout/:username", func(params martini.Params, r render.Render) {
		username := params["username"]
		userstate.Logout(username)
		if userstate.IsLoggedIn(username) {
			r.JSON(200, map[string]interface{}{"error": "user " + username + " is still logged in!"})
			return
		}
		r.JSON(200, map[string]interface{}{"logout": "true"})
	})

	// For login status
	m.Post("/api/1.0/status/:username", func(params martini.Params, r render.Render) {
		username := params["username"]
		if userstate.IsLoggedIn(username) {
			r.JSON(200, map[string]interface{}{"login": "true"})
		} else {
			r.JSON(200, map[string]interface{}{"login": "false"})
		}
	})

	// Score POST og GET + timestamp
	m.Post("/api/1.0/score/:username/:score", func(params martini.Params, r render.Render) {
		username := params["username"]
		score := params["score"]

		if !userstate.HasUser(username) {
			r.JSON(200, map[string]interface{}{"error": "no such user " + username})
			return
		}

		users := userstate.GetUsers()
		users.Set(username, "score", score)

		r.JSON(200, map[string]interface{}{"score set": "true"})
	})
	m.Get("/api/1.0/score/:username", func(params martini.Params, r render.Render) {
		username := params["username"]

		if !userstate.HasUser(username) {
			r.JSON(200, map[string]interface{}{"error": "no such user " + username})
			return
		}

		users := userstate.GetUsers()
		score, err := users.Get(username, "score")
		if err != nil {
			r.JSON(200, map[string]interface{}{"error": "could not get score for " + username})
			return
		}

		r.JSON(200, map[string]interface{}{"score": score})
	})


	m.Get("/mirrors", func(r render.Render) {
		var page PageData
		page.Text = map[string]string{
			"title":    "Mirrors",
			"subtitle": "List of mirrors",
		}

		// Generate a menu where item 1 is active
		page.Menu = GenerateMenu(1)

		// !! Reload template !!
		//r = render.New(render.Options{})

		// Render the specified templates/.tmpl file as HTML and return
		r.HTML(http.StatusOK, "mirrors", page)
	})

	// Rest
	Rest(m, fizz.Perm())

	// Activate the permission middleware
	m.Use(fizz.All())

	// Share the files in static
	m.Use(martini.Static("static"))

	// HTTP Basic Auth
	m.Use(auth.BasicFunc(func(username, password string) bool {

		return auth.SecureCompare(username, "admin") && auth.SecureCompare(password, "testfest")
	}))

	m.Run() // port 3000 by default
}
