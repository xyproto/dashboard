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
	menuitem1.Text = "Dashboard"
	if active == 0 {
		menuitem1.Selected = true
	}

	var menuitem2 MenuItem
	menuitem2.Url = "/mirrors"
	menuitem2.Text = "Mirrors"
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
	menuitem4.Url = "/wiki"
	menuitem4.Text = "Wiki"
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
	fmt.Println("dashboard ", Version)

	r := render.New(render.Options{})

	m := martini.Classic()

	fizz := fizz.New()

	// Dashboard
	m.Get("/", func(w http.ResponseWriter, req *http.Request) {
		var page PageData
		page.Text = map[string]string{
			"description": "Dashboard",
			"title":       "Dashboard",
			"subtitle":    "Charts and graphs",
		}

		// Generate a menu where item 0 is active
		page.Menu = GenerateMenu(0)

		// !! Reload template !!
		//r = render.New(render.Options{})

		// Render the specified templates/.tmpl file as HTML and return
		r.HTML(w, http.StatusOK, "dashboard", page)
	})

	m.Get("/mirrors", func(w http.ResponseWriter, req *http.Request) {
		var page PageData
		page.Text = map[string]string{
			"description": "Mirrors",
			"title":       "Mirrors",
			"subtitle":    "Mirror mirror on the wall",
		}

		// Generate a menu where item 1 is active
		page.Menu = GenerateMenu(1)

		// !! Reload template !!
		//r = render.New(render.Options{})

		// Render the specified templates/.tmpl file as HTML and return
		r.HTML(w, http.StatusOK, "mirrors", page)
	})

	// TODO: Admin panel with simple user management

	// Activate the permission middleware
	m.Use(fizz.All())

	// Share the files in static
	m.Use(martini.Static("static"))

	m.Run() // port 3000 by default
}
