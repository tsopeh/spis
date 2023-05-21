package main

import (
	"os"
	"path/filepath"
)

// TODO: Show progress for fetching and processing HTML and PDF documents.
// TODO: Capture stdout for OCR
// TODO: Experiment with goroutines to make thing more concurrent (execute faster).

var outputDirPath = filepath.Join("./", "OUTPUT", "run01")

func main() {

	check(os.MkdirAll(outputDirPath, os.ModePerm))

	htmlLegislationCollector := *createHtmlLegislationCollector()
	pdfLegislationCollector := *createPdfLegislationCollector()
	legislationUrlCollector := *createLegislationUrlCollector(&htmlLegislationCollector, &pdfLegislationCollector)
	menuCollector := *createMenuCollector(&legislationUrlCollector)

	//menuCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/api/reg/menu")

	// Test PDF & OCR
	// Short PDF
	pdfLegislationCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/viewdoc?regactid=413516&doctype=reg&findpdfurl=true")
	// Long PDF
	//pdfLegislationCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/viewdoc?regactid=413518&doctype=reg&findpdfurl=true")

	menuCollector.Wait()
	legislationUrlCollector.Wait()
	pdfLegislationCollector.Wait()
	htmlLegislationCollector.Wait()

}
