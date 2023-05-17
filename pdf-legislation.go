package main

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/gocolly/colly/v2"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/gen2brain/go-fitz"
	"github.com/otiai10/gosseract/v2"
)

func createPdfLegislationCollector() *colly.Collector {
	c := colly.NewCollector()

	// Example of the request headers needed to fetch a PDF.
	// curl 'https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/viewdoc?regactid=413504&doctype=reg&findpdfurl=true'   -H 'X-Referer: /SlGlasnikPortal/pdfjs/build/pdf.worker.js'   -H 'Referer: https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/pdfjs/build/pdf.worker.js'   -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36'   --compressed --output - > test.pdf
	c.OnRequest(func(request *colly.Request) {
		request.Headers.Add("X-Referer", "/SlGlasnikPortal/pdfjs/build/pdf.worker.js")
		request.Headers.Add("Referer", "https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/pdfjs/build/pdf.worker.js")
		request.Headers.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36")
		request.Headers.Add("Accept-Encoding", "compressed")
		log.Println(request.URL.String())
		log.Println(request.Headers)
	})

	c.OnResponse(func(response *colly.Response) {
		processPdfWithOcr(response.Body)
	})

	return c
}

func processPdfWithOcr(pdfBuffer []byte) {

	doc, err := fitz.NewFromMemory(pdfBuffer)
	check(err)
	defer func() { check(doc.Close()) }()

	client := gosseract.NewClient()
	client.Languages = []string{"srp", "srp_latn", "eng"}
	defer func() { check(client.Close()) }()

	pdfTitle := doc.Metadata()["title"]
	var hash = md5.Sum(pdfBuffer)
	var sanitizedName = colly.SanitizeFileName(pdfTitle)
	var fileNameMaxLength = 20
	if len(sanitizedName) > fileNameMaxLength {
		sanitizedName = sanitizedName[:fileNameMaxLength]
	}
	var hashString = hex.EncodeToString(hash[:])
	sanitizedName = sanitizedName + "---" + "PDF" + "---" + hashString + ".txt"

	var outputDirPath = filepath.Join("./", "OUTPUT")
	check(os.MkdirAll(outputDirPath, os.ModePerm))
	outputFilePath := filepath.Join(outputDirPath, sanitizedName)

	f, err := os.Create(outputFilePath)
	check(err)
	defer func() { check(f.Close()) }()

	w := bufio.NewWriter(f)
	check(err)

	for n := 0; n < doc.NumPage(); n++ {
		// Page seg mode: 0=osd only, 1=auto+osd, 2=auto, 3=col, 4=block," " 5=line, 6=word, 7=char
		check(client.SetVariable("tessedit_pageseg_mode", "1"))
		img, err := doc.Image(n)
		check(err)

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			panic(err)
		}

		if err := client.SetImageFromBytes(buf.Bytes()); err != nil {
			panic(err)
		}
		text, err := client.Text()
		check(err)

		if _, err = w.WriteString(text); err != nil {
			panic(err)
		}

	}

	check(w.Flush())

}
