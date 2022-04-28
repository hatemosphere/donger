package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	dongerListURL = "http://dongerlist.com/" // Fuck HTTP, but HTTPS does not work on this gods forsaken website for some reason
	categoryPath  = "category/"
)

// dongerCategory is the container for a dongers of single category
type dongerCategory struct {
	Name    string
	Dongers []string
}

func main() {
	dongerCategories := make(map[string]dongerCategory)
	c := colly.NewCollector()
	d := c.Clone()

	c.OnHTML(`li[class=list-2-item]`, func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			if strings.HasPrefix(el.Attr("href"), fmt.Sprint(dongerListURL+categoryPath)) {
				dongerCategories[el.Attr("href")] = dongerCategory{
					Name: el.Text,
				}
				d.Visit(e.Request.AbsoluteURL(el.Attr("href")))
			}
		})
	})

	d.OnHTML(`ul[class=list-1]`, func(e *colly.HTMLElement) {
		dongers := []string{}
		e.ForEach(`textarea[class=donger]`, func(_ int, el *colly.HTMLElement) {
			dongers = append(dongers, el.Text)
		})
		// Elegant way to change single value of struct, having map of structs, by Stack Overflow
		if entry, ok := dongerCategories[e.Request.URL.String()]; ok {
			entry.Dongers = dongers
			dongerCategories[e.Request.URL.String()] = entry
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	d.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting category", r.URL)
	})

	c.Visit(dongerListURL)
	fmt.Println(dongerCategories)
}
