package main

import (
	"github.com/gocolly/colly/v2"
	"log"
	"regexp"
)

func fetchDocumentUrls(
	menuUrls []string,
	htmlDocumentUrls *[]string,
	pdfDocumentUrls *[]string,
	unknownDocumentsFoundInMenuUrls *[]string,
) {

	uuidRegex := regexp.MustCompile(`'(.*)'`)
	findPdf := regexp.MustCompile(`findpdfurl=true`)

	c := colly.NewCollector()

	// PDF document
	c.OnHTML(`a[href]`, func(aEl *colly.HTMLElement) {

	})

	// Invalid
	c.OnHTML(`a`, func(aEl *colly.HTMLElement) {
		uuidMatch := uuidRegex.FindStringSubmatch(aEl.Attr("ui-sref"))
		hasHtmlDocumentUuid := uuidMatch != nil
		if hasHtmlDocumentUuid {
			htmlDocumentUuid := uuidMatch[1]
			htmlDocumentUrl := "https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/reg/viewAct/" + htmlDocumentUuid
			*htmlDocumentUrls = append(*htmlDocumentUrls, htmlDocumentUrl)
			return
		}

		href := aEl.Attr("href")
		isPdfHref := findPdf.MatchString(href)
		if isPdfHref {
			pdfDocumentUrl := "https://www.pravno-informacioni-sistem.rs" + href
			*pdfDocumentUrls = append(*pdfDocumentUrls, pdfDocumentUrl)
			return
		}

		*unknownDocumentsFoundInMenuUrls = append(*unknownDocumentsFoundInMenuUrls, aEl.Request.URL.String())

	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error in `fetchDocumentUrls` for URL", r.Request.URL.String(), err)
	})

	for _, url := range menuUrls {
		c.Visit(url)
	}
}
