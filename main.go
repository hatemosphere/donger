package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/lukechampine/randmap"
	"github.com/mitchellh/go-homedir"
)

const (
	dongerListURL  = "http://dongerlist.com/" // Fuck HTTP, but HTTPS does not work on this gods forsaken website for some reason
	categoryPath   = "category/"
	dongerFileName = "dongers.json"
)

func scrapeDongers() map[string][]string {
	dongerCategories := make(map[string][]string)
	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))
	d := c.Clone()

	c.OnHTML(`li[class=list-2-item]`, func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			if strings.HasPrefix(el.Attr("href"), fmt.Sprint(dongerListURL+categoryPath)) {
				d.Visit(e.Request.AbsoluteURL(el.Attr("href")))
			}
		})
	})

	d.OnHTML(`ul[class=list-1]`, func(e *colly.HTMLElement) {

		// dongers := []string{}
		e.ForEach(`textarea[class=donger]`, func(_ int, el *colly.HTMLElement) {
			dongerCategoryName := strings.TrimPrefix(e.Request.URL.String(), fmt.Sprint(dongerListURL+categoryPath))
			dongerCategories[dongerCategoryName] = append(dongerCategories[dongerCategoryName], el.Text)

		})
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

func randomizeNumber(number int) int {
	randomSource := rand.NewSource(time.Now().Unix())
	randomizer := rand.New(randomSource) // initialize local pseudorandom generator
	randomNumber := randomizer.Intn(number)
	return randomNumber
}

func choseRandomDonger(chosenDongerCategory string, dongerCategories map[string][]string) string {
	if chosenDongerCategory == "random" {
		randomCategory := randmap.Val(dongerCategories).([]string)
		randomDongerIndex := randomizeNumber(len(randomCategory))
		return randomCategory[randomDongerIndex]
	} else {
		randomDongerIndex := randomizeNumber(len(chosenDongerCategory))
		dongerList := dongerCategories[chosenDongerCategory]
		return dongerList[randomDongerIndex]
	}
}

func main() {
	dongerChosenCategory := flag.String("category", "random", "donger category")
	flag.Parse()

	var dongerCategories = make(map[string][]string)

	homedir, homedirErr := homedir.Dir()
	if homedirErr != nil {
		panic(homedirErr)
	}

	dongersDirPath := path.Join(homedir, ".donger")
	dongersFilePath := fmt.Sprint(dongersDirPath + "/" + dongerFileName)

	dongerFile, dongerFileOpenErr := os.OpenFile(dongersFilePath, os.O_RDWR, 0644)
	if os.IsNotExist(dongerFileOpenErr) {
		fmt.Println("Dongers file does not exist and will be generated")
		dongerCategories = scrapeDongers()
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
		dongerFileBytes, dongerFileReadErr := ioutil.ReadAll(dongerFile)
		if dongerFileReadErr != nil {
			panic(dongerFileReadErr)
		}
		fmt.Println("Donger file already exists, skipping generation")
		json.Unmarshal([]byte(dongerFileBytes), &dongerCategories)
	}

	randomDonger := choseRandomDonger(*dongerChosenCategory, dongerCategories)

	// clipboardInitErr := clipboard.Init()
	// if clipboardInitErr != nil {
	// 	panic(clipboardInitErr)
	// }
	// clipboard.Write(clipboard.FmtText, []byte(randomDonger))
	fmt.Println("Donger: " + randomDonger + " was chosen and got copied to clipboard")
}
