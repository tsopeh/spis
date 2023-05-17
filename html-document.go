package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

func createHtmlDocumentCollector() *colly.Collector {
	byteOrderMarkReg := regexp.MustCompile("\uFEFF")
	nbspReg := regexp.MustCompile("[\u202F\u00A0]")
	multipleWhitespacesAndNewlines := regexp.MustCompile(`(\s*\n+(?:\n*|\s*)\n+\s*)`)
	c := colly.NewCollector(
		colly.MaxRequests(3),
	)

	c.OnHTML(`html`, func(h *colly.HTMLElement) {
		var contentEl = h.DOM.Find("#actContentPrimaryScroll")
		contentEl.Find("meta").Remove()
		contentEl.Find("link").Remove()
		contentEl.Find("style").Remove()
		contentEl.Find("script").Remove()
		var text = contentEl.Text()
		text = byteOrderMarkReg.ReplaceAllString(text, "")
		text = nbspReg.ReplaceAllString(text, " ")
		text = multipleWhitespacesAndNewlines.ReplaceAllString(text, "\n\n")
		text = strings.TrimSpace(text)

		var pageTitle = h.DOM.Find("title").Text()

		var hash = md5.Sum([]byte(text))
		var sanitazedName = colly.SanitizeFileName(pageTitle)
		var fileNameMaxLength = 20
		if len(sanitazedName) > fileNameMaxLength {
			sanitazedName = sanitazedName[:fileNameMaxLength]
		}
		var hashString = hex.EncodeToString(hash[:])
		sanitazedName = sanitazedName + "---" + hashString + ".txt"

		log.Println(sanitazedName)

		var outputDirPath = filepath.Join("./", "OUTPUT")
		check(os.MkdirAll(outputDirPath, os.ModePerm))
		outputFilePath := filepath.Join(outputDirPath, sanitazedName)
		f, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, 0600)
		check(err)
		defer func() { check(f.Close()) }()
		if _, err := f.WriteString(text); err != nil {
			panic(err)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error in `createHtmlDocumentCollector`: ", err)
	})

	return c
}
