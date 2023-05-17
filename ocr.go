package main

import (
	"bytes"
	"image/png"
	"os"

	"github.com/gen2brain/go-fitz"
	"github.com/otiai10/gosseract/v2"
)

func ocrHelloWorld() {

	doc, err := fitz.New("test.pdf")
	if err != nil {
		panic(err)
	}
	defer func(doc *fitz.Document) {
		err := doc.Close()
		if err != nil {
			panic(err)
		}
	}(doc)

	client := gosseract.NewClient()
	client.Languages = []string{"srp", "srp_latn", "eng"}
	defer func() { check(client.Close()) }()

	for n := 0; n < doc.NumPage(); n++ {
		// Page seg mode: 0=osd only, 1=auto+osd, 2=auto, 3=col, 4=block," " 5=line, 6=word, 7=char
		check(client.SetVariable("tessedit_pageseg_mode", "1"))
		img, err := doc.Image(n)
		if err != nil {
			panic(err)
		}

		buf := new(bytes.Buffer)
		if err := png.Encode(buf, img); err != nil {
			panic(err)
		}

		if err := client.SetImageFromBytes(buf.Bytes()); err != nil {
			panic(err)
		}
		text, err := client.Text()
		if err != nil {
			panic(err)
		}

		f, err := os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			panic(err)
		}
		defer func() { check(f.Close()) }()

		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}

	}

}
