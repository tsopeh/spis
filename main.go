package main

import (
	"log"
	"os"
	"path/filepath"
)

// TODO: Show progress for fetching and processing HTML and PDF documents.
// TODO: Capture stdout for OCR
// TODO: Experiment with goroutines to make thing more concurrent (execute faster).

var outputDirPath = filepath.Join("/Volumes/USB_STORAGE/OUT", "OUTPUT", "run03")

func main() {

	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("%s Fatal: Fatal Error Signal")
	}
	defer file.Close()
	log.SetOutput(file)

	check(os.RemoveAll(outputDirPath))
	check(os.MkdirAll(outputDirPath, os.ModePerm))

	htmlLegislationCollector := *createHtmlLegislationCollector()
	pdfLegislationCollector := *createPdfLegislationCollector()
	legislationUrlCollector := *createLegislationUrlCollector(&htmlLegislationCollector, &pdfLegislationCollector)
	menuCollector := *createMenuCollector(&legislationUrlCollector)

	menuCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/api/reg/menu")

	// Test PDF & OCR
	// Short PDF
	//pdfLegislationCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/viewdoc?regactid=413516&doctype=reg&findpdfurl=true")
	// Long PDF
	//pdfLegislationCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/viewdoc?regactid=413518&doctype=reg&findpdfurl=true")

	menuCollector.Wait()
	legislationUrlCollector.Wait()
	pdfLegislationCollector.Wait()
	htmlLegislationCollector.Wait()

}
