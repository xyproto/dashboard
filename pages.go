package main

import (
	"fmt"

	"github.com/xyproto/onthefly"
)

// Generate a new SVG Page
func svgPage() *onthefly.Page {
	page, svg := onthefly.NewTinySVG(0, 0, 128, 64)
	desc := svg.AddNewTag("desc")
	desc.AddContent("Hello SVG")

	// x, y, radius, color
	svg.Circle(30, 10, 5, "red")
	svg.Circle(110, 30, 2, "green")
	svg.Circle(80, 40, 7, "blue")
	return page
}

// Generate a new onthefly Page (HTML5 and CSS combined)
func indexPage(svgurl string) *onthefly.Page {

	// Create a new HTML5 page, with CSS included
	page := onthefly.NewHTML5Page("Demonstration")

	// Add some text
	page.AddContent(fmt.Sprintf("onthefly %.1f", onthefly.Version))

	// Change the margin (em is default)
	page.SetMargin(4)

	// Change the font family
	page.SetFontFamily("serif") // sans serif

	// Change the color scheme
	page.SetColor("black", "#d0d0d0")

	// Include the generated SVG image on the page
	body, err := page.GetTag("body")
	if err == nil {
		// CSS attributes for the body tag
		body.AddStyle("font-size", "2em")

		// Paragraph
		p := body.AddNewTag("p")

		// CSS style
		p.AddStyle("margin-top", "2em")

		// Image tag
		img := p.AddNewTag("img")

		// HTML attributes
		img.AddAttrib("src", svgurl)
		img.AddAttrib("alt", "Three circles")

		// CSS style
		img.AddStyle("width", "60%")
		img.AddStyle("border", "4px solid white")
	}

	return page
}
