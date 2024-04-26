package main

import (
	"llm_scraper/utils"
)

type Pokemon struct {
	URL, Name, Img string
	Price          float32
}

func main() {
	var pokemon Pokemon

	pages := []string{"https://scrapeme.live/shop/"}
	utils.PageCrawelr(pages, &pokemon)
}
