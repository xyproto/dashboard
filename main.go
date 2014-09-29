package main

import (
	"fmt"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/unrolled/render"
	"github.com/xyproto/fizz"
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

	r := render.New(render.Options{})

	m := martini.Classic()

	fizz := fizz.New()
	userstate := fizz.UserState()

	m.Get("/register/:password", func(params martini.Params) string {
		if userstate.HasUser("admin") {
			return "Admin user already exists"
		}
		userstate.AddUser("admin", params["password"], "")
		userstate.SetAdminStatus("admin")
		return fmt.Sprintf("Admin user was created: %v\n", userstate.HasUser("admin"))
	})

	m.Get("/login/:password", func(w http.ResponseWriter, req *http.Request, params martini.Params) string {
		if userstate.CorrectPassword("admin", params["password"]) {
			userstate.Login(w, "admin")
			return "Wrong password"
		}
		return fmt.Sprintf("Logged in as administrator: %v\n", userstate.AdminRights(req))
	})

	m.Get("/logout", func(w http.ResponseWriter) string {
		userstate.Logout("admin")
		return "Logged out"
	})

	// Public page
	m.Get("/", func() string {
		return "hi"
	})

	// Admin panel
	m.Get("/admin", func(w http.ResponseWriter, req *http.Request) {
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
		r.HTML(w, http.StatusOK, "admin", page)
	})

	m.Get("/mirrors", func(w http.ResponseWriter, req *http.Request) {
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
		r.HTML(w, http.StatusOK, "mirrors", page)
	})

	// Rest
	Rest(m, fizz.Perm())

	// Activate the permission middleware
	m.Use(fizz.All())

	// Share the files in static
	m.Use(martini.Static("static"))

	m.Run() // port 3000 by default
}
