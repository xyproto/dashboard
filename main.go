package main

import (
	"fmt"
	"net/http"

	"github.com/Unknwon/macaron"
	"github.com/unrolled/render"
	"github.com/xyproto/permissions2"
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
	m := macaron.Classic()

	perm := permissions.New()

    // Set up a middleware handler for Macaron, with a custom "permission denied" message.
    permissionHandler := func(ctx *macaron.Context) {
        // Check if the user has the right admin/user rights
        if perm.Rejected(ctx.Resp, ctx.Req.Request) {
            fmt.Fprintf(ctx.Resp, "Permission denied!")
            // Deny the request
            ctx.Error(http.StatusForbidden)
            // Don't call other middleware handlers
            return
        }
        // Call the next middleware handler
        ctx.Next()
    }

    m.Use(macaron.Logger())
    m.Use(permissionHandler)
    m.Use(macaron.Gziper())
    m.Use(macaron.Static("public"))

	// Dashboard
	m.Get("/", func(ctx *macaron.Context) {
		var page PageData
		page.Text = map[string]string{
			"description": "Dashboard",
			"title":       "Dashboard",
			"subtitle":    "Charts and graphs",
		}

		// Generate a menu where item 0 is active
		page.Menu = GenerateMenu(0)

		// !! Reload template !!
		r = render.New(render.Options{})

		// Render the specified templates/.tmpl file as HTML and return
		r.HTML(ctx.Resp, http.StatusOK, "dashboard", page)
	})

	m.Get("/mirrors", func(ctx *macaron.Context) {
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
		r.HTML(ctx.Resp, http.StatusOK, "mirrors", page)
	})

	// TODO: Admin panel with simple user management

    // Recovery middleware goes last
    m.Use(macaron.Recovery())

	// Serve
	m.Run(3000)
}
