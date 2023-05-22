package main

import (
	"github.com/gocolly/colly/v2"
	"log"
	"regexp"
)

func createLegislationUrlCollector(
	htmlLegislationCollector *colly.Collector,
	pdfLegislationCollector *colly.Collector,
) *colly.Collector {

	uuidRegex := regexp.MustCompile(`'(.*)'`)
	findPdf := regexp.MustCompile(`findpdfurl=true`)

	c := colly.NewCollector()

	// Html document
	c.OnHTML(`a[ui-sref]`, func(aEl *colly.HTMLElement) {
		uuid := uuidRegex.FindStringSubmatch(aEl.Attr("ui-sref"))[1]
		htmlLegislationCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/reg/viewAct/" + uuid)
	})

	// PDF document
	c.OnHTML(`a[href]`, func(aEl *colly.HTMLElement) {
		href := aEl.Attr("href")
		if !findPdf.MatchString(href) {
			return
		}
		pdfUrl := "https://www.pravno-informacioni-sistem.rs" + href
		pdfLegislationCollector.Visit(pdfUrl)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error in `createLegislationUrlCollector` for URL", r.Request.URL.String(), err)
	})

	return c
}
