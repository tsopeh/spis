package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/cheggaaa/pb/v3"
	"github.com/gen2brain/go-fitz"
	"github.com/gocolly/colly/v2"
	"github.com/otiai10/gosseract/v2"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func createPdfLegislationCollector(
	urls []DocumentUrl,
	documentsDirPath string,
	cache *ScraperCache,
	barPool *pb.Pool,
) *colly.Collector {

	documentBar := pb.New(len(urls))
	documentBar.Set("prefix", "PDF")
	documentBar.SetMaxWidth(80)

	pdfPageBar := pb.New(1)
	pdfPageBar.Set("prefix", "Pages (PDF)")
	pdfPageBar.SetMaxWidth(80)

	barPool.Add(documentBar, pdfPageBar)

	c := colly.NewCollector()

	// Example of the request headers needed to fetch a PDF.
	// curl 'https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/viewdoc?regactid=413504&doctype=reg&findpdfurl=true'   -H 'X-Referer: /SlGlasnikPortal/pdfjs/build/pdf.worker.js'   -H 'Referer: https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/pdfjs/build/pdf.worker.js'   -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36'   --compressed --output - > test.pdf
	c.OnRequest(func(request *colly.Request) {
		if _, ok := cache.Get(request.URL.String()); ok {
			documentBar.Increment()
			request.Abort()
		}
		request.Ctx.Put("requestUrl", request.URL.String())
		request.Headers.Add("X-Referer", "/SlGlasnikPortal/pdfjs/build/pdf.worker.js")
		request.Headers.Add("Referer", "https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/pdfjs/build/pdf.worker.js")
		request.Headers.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
		request.Headers.Add("Accept-Encoding", "compressed")
	})

	c.OnResponse(func(response *colly.Response) {
		processPdfWithOcr(response.Body, documentsDirPath, response, cache, pdfPageBar)
		documentBar.Increment()
	})

	for _, url := range urls {
		c.Visit(url.url)
	}

	return c
}

func processPdfWithOcr(
	pdfBuffer []byte,
	documentsDirPath string,
	response *colly.Response,
	cache *ScraperCache,
	bar *pb.ProgressBar,
) {
	doc, err := fitz.NewFromMemory(pdfBuffer)
	check(err)
	defer func() { check(doc.Close()) }()

	pdfTitle := strings.TrimSpace(strings.ReplaceAll(doc.Metadata()["title"], "\u0000", ""))
	var hash = md5.Sum(pdfBuffer)
	var sanitizedName = colly.SanitizeFileName(pdfTitle)
	var fileNameMaxLength = 20
	if len(sanitizedName) > fileNameMaxLength {
		sanitizedName = sanitizedName[:fileNameMaxLength]
	}
	var hashString = hex.EncodeToString(hash[:])
	sanitizedName = sanitizedName + "---" + "PDF" + "---" + hashString + ".txt"

	pageCount := doc.NumPage()
	result := make([]string, pageCount)
	poolSize := int(math.Max(1, math.Ceil(0.7*float64(runtime.NumCPU()))))
	wg := CreateSemaphoredWaitGroup(poolSize)

	bar.SetTotal(int64(pageCount))
	bar.SetCurrent(0)

	for i := 0; i < pageCount; i++ {
		// scoped var `pageIndex`
		pageIndex := i

		wg.AddJob(func() {
			client := gosseract.NewClient()
			defer func() { check(client.Close()) }()
			// Disable debug log messages.
			check(client.DisableOutput())
			client.Languages = []string{"srp", "srp_latn", "eng"}
			// Page seg mode: 0=osd only, 1=auto+osd, 2=auto, 3=col, 4=block," " 5=line, 6=word, 7=char
			check(client.SetVariable("tessedit_pageseg_mode", "1"))
			img, err := doc.Image(pageIndex)
			buf := new(bytes.Buffer)
			if err := png.Encode(buf, img); err != nil {
				panic(err)
			}
			if err := client.SetImageFromBytes(buf.Bytes()); err != nil {
				panic(err)
			}
			text, err := client.Text()
			check(err)
			result[pageIndex] = text
			bar.Increment()
		})
	}
	wg.WaitAll()
	outputFilePath := filepath.Join(documentsDirPath, sanitizedName)
	f, err := os.Create(outputFilePath)
	check(err)
	defer func() { check(f.Close()) }()
	if _, err := f.WriteString(strings.Join(result, "")); err != nil {
		panic(err)
	}
	cache.Add(FinishedDocument{
		requestUrl:  response.Ctx.Get("requestUrl"),
		responseUrl: response.Request.URL.String(),
		title:       pdfTitle,
		fileName:    sanitizedName,
		contentHash: hashString,
	})
}
