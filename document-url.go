package main

import (
	"log"
	"regexp"

	"github.com/gocolly/colly/v2"
)

func createLegislationUrlCollector(legislationCollector *colly.Collector) *colly.Collector {

	uuidRegex := regexp.MustCompile(`'(.*)'`)

	c := colly.NewCollector()

	c.OnHTML(`a[ui-sref]`, func(aEl *colly.HTMLElement) {
		uuid := uuidRegex.FindStringSubmatch(aEl.Attr("ui-sref"))[1]
		legislationCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/reg/viewAct/" + uuid)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error in `createLegislationUrlCollector`: ", err)
	})

	return c
}
