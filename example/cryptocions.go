package example

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/gocolly/colly"
)

// GetBitCoin 比特币价钱挖掘
func GetBitCoin() {
	fName := "bitcoin.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	writer.Write([]string{"Name", "Symbol", "Price (USD)", "Volume (USD)", "Market capacity (USD)", "Change (1h)", "Change (24h)", "Change (7d)"})

	// Instantiate default collector
	c := colly.NewCollector()

	c.OnHTML("#currencies-all tbody tr", func(e *colly.HTMLElement) {
		changes := make([]string, 3)
		e.ForEach(".percent-change", func(i int, el *colly.HTMLElement) {
			if i < 3 {
				changes[i] = el.Text
			}
		})

		writer.Write([]string{
			e.ChildText(".currency-name-container"),
			e.ChildText(".col-symbol"),
			e.ChildAttr("a.price", "data-usd"),
			e.ChildAttr("a.volume", "data-usd"),
			e.ChildAttr(".market-cap", "data-usd"),
			changes[0],
			changes[1],
			changes[2],
		})
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished, check file %q for results\n", fName)
}
