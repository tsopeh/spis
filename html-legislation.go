package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/cheggaaa/pb/v3"
	"github.com/gocolly/colly/v2"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func createHtmlLegislationCollector(
	urls []DocumentUrl,
	documentsDirPath string,
	barPool *pb.Pool,
) *colly.Collector {

	bar := pb.New(len(urls))
	bar.Set("prefix", "HTML")
	bar.SetMaxWidth(80)
	barPool.Add(bar)

	byteOrderMarkReg := regexp.MustCompile("\uFEFF")
	nbspReg := regexp.MustCompile("[\u202F\u00A0]")
	multipleWhitespacesAndNewlines := regexp.MustCompile(`(\s*\n+(?:\n*|\s*)\n+\s*)`)

	c := colly.NewCollector(
		colly.Async(),
	)
	c.SetRequestTimeout(time.Duration(60) * time.Second)
	c.Limit(&colly.LimitRule{Delay: 50 * time.Millisecond, RandomDelay: 50 * time.Millisecond, Parallelism: 8, DomainGlob: "*"})

	c.OnHTML(`html`, func(h *colly.HTMLElement) {
		var contentEl = h.DOM.Find("#actContentPrimaryScroll")
		contentEl.Find("meta").Remove()
		contentEl.Find("link").Remove()
		contentEl.Find("style").Remove()
		contentEl.Find("script").Remove()
		contentEl.Find("table").Remove()
		contentEl.Find("thead").Remove()
		contentEl.Find("tbody").Remove()
		contentEl.Find("tr").Remove()
		contentEl.Find("td").Remove()
		var text = contentEl.Text()
		text = byteOrderMarkReg.ReplaceAllString(text, "")
		text = nbspReg.ReplaceAllString(text, " ")
		text = multipleWhitespacesAndNewlines.ReplaceAllString(text, "\n\n")
		text = strings.TrimSpace(text)

		var pageTitle = h.DOM.Find("title").Text()

		var hash = md5.Sum([]byte(text))
		var sanitizedName = colly.SanitizeFileName(pageTitle)
		var fileNameMaxLength = 20
		if len(sanitizedName) > fileNameMaxLength {
			sanitizedName = sanitizedName[:fileNameMaxLength]
		}
		var hashString = hex.EncodeToString(hash[:])
		sanitizedName = sanitizedName + "---" + "HTML" + "---" + hashString + ".txt"

		outputFilePath := filepath.Join(documentsDirPath, sanitizedName)
		f, err := os.OpenFile(outputFilePath, os.O_WRONLY|os.O_CREATE, 0600)
		check(err)
		defer func() { check(f.Close()) }()
		if _, err := f.WriteString(text); err != nil {
			panic(err)
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Error in `createHtmlLegislationCollector` for URL", r.Request.URL.String(), err)
	})

	c.OnScraped(func(response *colly.Response) {
		bar.Increment()
	})

	for _, url := range urls {
		c.Visit(url.url)
	}

	return c
}
