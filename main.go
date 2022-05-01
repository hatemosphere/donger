package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/mitchellh/go-homedir"
	"golang.design/x/clipboard"
	"golang.org/x/exp/maps"
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

	dongerFile, dongerFileOpenErr := os.OpenFile(dongersFilePath, os.O_RDWR, 0644)
	if os.IsNotExist(dongerFileOpenErr) {
		fmt.Println("Dongers file does not exist and will be generated")
		dongerCategories := scrapeDongers()
		mkDirErr := os.MkdirAll(dongersDirPath, os.ModePerm)
		if mkDirErr != nil {
			panic(mkDirErr)
		}
		file, dongerMarshalErr := json.Marshal(dongerCategories)
		if dongerMarshalErr != nil {
			panic(dongerMarshalErr)
		}
		_ = ioutil.WriteFile(dongersFilePath, file, 0644)
		fmt.Println("Donger file generated")
	} else {
		fmt.Println("Donger file already exists, skipping generation")
	}

	dongerFileBytes, dongerFileReadErr := ioutil.ReadAll(dongerFile)
	if dongerFileReadErr != nil {
		panic(dongerFileReadErr)
	}

	dongerCategories := map[string]dongerCategory{}
	json.Unmarshal([]byte(dongerFileBytes), &dongerCategories)

	dongerCategoriesSlice := maps.Values(dongerCategories)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s) // initialize local pseudorandom generator
	randDongerCategoriesSliceElement := r.Intn(len(dongerCategoriesSlice))

	ss := rand.NewSource(time.Now().Unix())
	rr := rand.New(ss) // initialize local pseudorandom generator
	dongersSlice := dongerCategoriesSlice[randDongerCategoriesSliceElement].Dongers
	randDongersSliceElement := rr.Intn(len(dongersSlice))
	randomDonger := (dongersSlice[randDongersSliceElement])

	clipboardInitErr := clipboard.Init()
	if clipboardInitErr != nil {
		panic(clipboardInitErr)
	}
	clipboard.Write(clipboard.FmtText, []byte(randomDonger))
}
