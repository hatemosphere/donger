package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/mitchellh/go-homedir"
)

const (
	dongerListURL  = "http://dongerlist.com/" // Fuck HTTP, but HTTPS does not work on this gods forsaken website for some reason
	categoryPath   = "category/"
	dongerFileName = "dongers.json"
)

// dongerCategory is the container for a dongers of single category
type dongerCategory struct {
	Name    string
	Dongers []string
}

func scrapeDongers() map[string]dongerCategory {
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
	return dongerCategories
}

func main() {
	// textPtr := flag.String("print", "", "Text to parse. (Required)")
	// metricPtr := flag.String("print", "random", "Print {random|[category]}")
	// uniquePtr := flag.Bool("unique", false, "Measure unique values of a metric.")

	homedir, homedirErr := homedir.Dir()
	if homedirErr != nil {
		panic(homedirErr)
	}

	dongersDirPath := path.Join(homedir, ".donger")
	dongersFilePath := fmt.Sprint(dongersDirPath + "/" + dongerFileName)

	dongerFile, err := os.OpenFile(dongersFilePath, os.O_RDWR, 0644)
	if os.IsNotExist(err) {
		fmt.Println("Dongers file does not exist and will be generated")
		dongerCategories := scrapeDongers()
		mkDirErr := os.MkdirAll(dongersDirPath, os.ModePerm)
		if mkDirErr != nil {
			panic(mkDirErr)
		}
		file, _ := json.Marshal(dongerCategories)
		_ = ioutil.WriteFile(dongersFilePath, file, 0644)
		dongerFile, openFileErr := os.OpenFile(dongersFilePath, os.O_RDWR, 0644)
		if openFileErr != nil {
			panic(mkDirErr)
		}
		fmt.Println("Donger file generated")
		fmt.Println(dongerFile.Name())
	} else {
		fmt.Println("Donger file already exists, skipping generation")
		fmt.Println(dongerFile.Name())
	}
}
