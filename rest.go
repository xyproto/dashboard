package main

import (
	"github.com/go-martini/martini"
	"github.com/xyproto/permissions"
)

func Rest(m *martini.ClassicMartini, perm *permissions.Permissions) {
	perm.AddPublicPath("/books")
	perm.AddPublicPath("/new")
	perm.AddPublicPath("/newquick")
	perm.AddPublicPath("/delete")
	m.Group("/books", func(r martini.Router) {
		r.Get("/:id", GetBooks)
		r.Post("/new", NewBook)
		r.Post("/newquick/:id", NewQuickBook)
		r.Put("/update/:id", UpdateBook)
		r.Delete("/delete/:id", DeleteBook)
	})
}

var books []string

func GetBooks(params martini.Params) string {
	book := params["id"]
	found := -1
	for i, id := range books {
		if id == book {
			found = i
		}
	}
	if found != -1 {
		return book
	}
	return ""
}

func NewBook(params martini.Params) {
	// TODO: Use a form to create a new book
	//book := "New Book"
	//books = append(books, book)
}

func NewQuickBook(params martini.Params) {
	book := params["id"]
	books = append(books, book)
}

func UpdateBook(params martini.Params) {
	// TODO: Update book timestamp?
	//book := params["id"]
}

func DeleteBook(params martini.Params) {
	book := params["id"]
	found := -1
	for i, id := range books {
		if id == book {
			found = i
		}
	}
	if found != -1 {
		books = append(books[:found], books[found+1:]...)
	}
}
