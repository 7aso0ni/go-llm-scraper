package utils

import (
	"log"
	"reflect"

	"github.com/gocolly/colly"
)

// This function will take link[s] and a struct given by the user which it will go and crawl the website and store it into the user defined struct
func PageCrawelr(urls []string, obj interface{}) {

	value := reflect.ValueOf(obj)

	// by using relect pacakge get the passed interface type
	// if interface type is not a pointer to a struct terminate program
	if value.Kind() != reflect.Pointer || value.Elem().Kind() != reflect.Struct {
		log.Fatalf("passed interface is not a pointer to a struct")
	}

	c := colly.NewCollector(
		// will allow to asynchronously crawl multiple pages at the same time
		colly.Async(true),
	)

	for _, url := range urls {
		c.Visit(url)

		// NOTE: TILL i FIGURE OUT A WAY TO AUTOMATE THIS IT WILL BE HARDCODED FOR NOW
		c.OnHTML("li.product", func(h *colly.HTMLElement) {

		})
	}
}
