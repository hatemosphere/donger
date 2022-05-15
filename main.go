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
// type dongerCategory struct {
// 	Name    string
// 	Dongers []string
// }

func scrapeDongers() map[string][]string {
	// aminoAcidsToCodons := map[rune]*[]string{}
	// for codon, aminoAcid := range utils.CodonsToAminoAcid {
	// 	mappedAminoAcid := aminoAcidsToCodons[aminoAcid]
	// 	*mappedAminoAcid = append(*mappedAminoAcid, codon)
	// }
	dongerCategories := make(map[string][]string)
	c := colly.NewCollector()
	d := c.Clone()

	c.OnHTML(`li[class=list-2-item]`, func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			if strings.HasPrefix(el.Attr("href"), fmt.Sprint(dongerListURL+categoryPath)) {
				// dongerCategories[el.Attr("href")] = dongerCategory{
				// 	Name: el.Text,
				// }
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

func main() {
	var dongerCategories = make(map[string][]string)
	// fmt.Println(dongerCategories)

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
	fmt.Println(dongerCategories)

	// dongerCategoriesSlice := maps.Values(dongerCategories)
	// s := rand.NewSource(time.Now().Unix())
	// r := rand.New(s) // initialize local pseudorandom generator
	// randDongerCategoriesSliceElement := r.Intn(len(dongerCategoriesSlice))

	// ss := rand.NewSource(time.Now().Unix())
	// rr := rand.New(ss) // initialize local pseudorandom generator
	// dongersSlice := dongerCategoriesSlice[randDongerCategoriesSliceElement].Dongers
	// randDongersSliceElement := rr.Intn(len(dongersSlice))
	// randomDonger := (dongersSlice[randDongersSliceElement])

	// clipboardInitErr := clipboard.Init()
	// if clipboardInitErr != nil {
	// 	panic(clipboardInitErr)
	// }
	// clipboard.Write(clipboard.FmtText, []byte(randomDonger))
}
