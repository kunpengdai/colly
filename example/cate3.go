package example

import (
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/gocolly/colly"
)

// GetCate3 获取斗鱼所有三级分类
func GetCate3() {
	domain := "https://www.douyu.com"
	fName := "cate3.csv"
	file, err := os.Create(fName)
	beg := time.Now()
	workerNum := 1000
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"cid3", "cateName"})

	// Instantiate default collector
	c := colly.NewCollector()
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: workerNum})

	d := c.Clone()
	d.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: workerNum})

	c.OnHTML("#live-list-contentbox li", func(e *colly.HTMLElement) {
		// fmt.Println(e.ChildAttr("a.thumb", "href"))
		// fmt.Println(e.ChildText("p.title"))
		// #live-list-contentbox > li:nth-child(1) > a
		// print(e.ChildAttr(""))
		url := domain + e.ChildAttr("a.thumb", "href")
		d.Visit(url)
	})
	d.OnHTML("#js_item_tag li", func(e *colly.HTMLElement) {
		if e.ChildAttr("a", "data-sid") == "" {
			return
		}
		writer.Write([]string{
			e.ChildAttr("a", "data-sid"),
			e.ChildAttr("a", "data-live-list-type"),
		})
	})

	c.Visit("https://www.douyu.com/directory")
	c.Wait()
	d.Wait()

	log.Printf("Scraping finished, check file %q for results,duration:%v\n", fName, time.Since(beg))
}
