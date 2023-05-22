package main

import (
	"github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly/v2"
	"log"
	"regexp"
)

type DocumentUrl struct {
	url           string
	kind          string
	parentMenuUrl string
}

func fetchDocumentUrls(
	menuUrls []string,
	htmlDocumentUrls *[]DocumentUrl,
	pdfDocumentUrls *[]DocumentUrl,
	unknownDocumentsFoundInMenuUrls *[]DocumentUrl,
) {

	log.Println("Processing menu items")
	bar := pb.StartNew(len(menuUrls))
	bar.SetMaxWidth(80)

	uuidRegex := regexp.MustCompile(`'(.*)'`)
	findPdf := regexp.MustCompile(`findpdfurl=true`)

	c := colly.NewCollector()

	// PDF document
	c.OnHTML(`a[href]`, func(aEl *colly.HTMLElement) {

	})

	// Invalid
	c.OnHTML(`a`, func(aEl *colly.HTMLElement) {
		parentMenuUrl := aEl.Request.URL.String()
		uuidMatch := uuidRegex.FindStringSubmatch(aEl.Attr("ui-sref"))
		hasHtmlDocumentUuid := uuidMatch != nil
		if hasHtmlDocumentUuid {
			htmlDocumentUuid := uuidMatch[1]
			htmlDocumentUrl := "https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/reg/viewAct/" + htmlDocumentUuid
			*htmlDocumentUrls = append(*htmlDocumentUrls, DocumentUrl{
				url:           htmlDocumentUrl,
				kind:          "HTML",
				parentMenuUrl: parentMenuUrl,
			})
			return
		}

		href := aEl.Attr("href")
		isPdfHref := findPdf.MatchString(href)
		if isPdfHref {
			pdfDocumentUrl := "https://www.pravno-informacioni-sistem.rs" + href
			*pdfDocumentUrls = append(*pdfDocumentUrls, DocumentUrl{
				url:           pdfDocumentUrl,
				kind:          "PDF",
				parentMenuUrl: parentMenuUrl,
			})
			return
		}

		*unknownDocumentsFoundInMenuUrls = append(*unknownDocumentsFoundInMenuUrls, DocumentUrl{
			url:           "",
			kind:          "UNKNOWN",
			parentMenuUrl: parentMenuUrl,
		})
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error in `fetchDocumentUrls` for URL", r.Request.URL.String(), err)
	})

	for _, url := range menuUrls {
		c.Visit(url)
		bar.Increment()
	}
	bar.Finish()
}
