package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// TODO: Show progress for fetching and processing HTML and PDF documents.
// TODO: Capture stdout for OCR
// TODO: Experiment with goroutines to make thing more concurrent (execute faster).

func main() {

	var outputDirPath = filepath.Join("./", "OUTPUT")
	var indexDirPath = filepath.Join(outputDirPath, "index")
	var documentsDirPath = filepath.Join(outputDirPath, "documents")

	//file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	//if err != nil {
	//	log.Fatalf("%s Fatal: Fatal Error Signal")
	//}
	//defer file.Close()
	//log.SetOutput(file)

	check(os.RemoveAll(outputDirPath))
	check(os.MkdirAll(indexDirPath, os.ModePerm))
	check(os.MkdirAll(documentsDirPath, os.ModePerm))

	var menuUrls []string
	fetchMenuUrls("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/api/reg/menu", &menuUrls)
	log.Println("Successfully retrieved the list of " + strconv.Itoa(len(menuUrls)) + " menu URLs.")
	check(writeIndexToFile(path.Join(indexDirPath, "menu-index.txt"), menuUrls))

	var htmlDocumentUrls []string
	var pdfDocumentUrls []string
	var unknownDocumentsFoundInMenuUrls []string
	fetchDocumentUrls(menuUrls, &htmlDocumentUrls, &pdfDocumentUrls, &unknownDocumentsFoundInMenuUrls)
	log.Println("Found " + strconv.Itoa(len(htmlDocumentUrls)) + " HTML document URLs.")
	log.Println("Found " + strconv.Itoa(len(pdfDocumentUrls)) + " PDF document URLs.")
	log.Println("Found " + strconv.Itoa(len(unknownDocumentsFoundInMenuUrls)) + " menu items with UNKNOWN document URLs.")

	check(writeIndexToFile(path.Join(indexDirPath, "html-index.txt"), htmlDocumentUrls))
	check(writeIndexToFile(path.Join(indexDirPath, "pdf-index.txt"), pdfDocumentUrls))
	check(writeIndexToFile(path.Join(indexDirPath, "unknown-menu-index.txt"), unknownDocumentsFoundInMenuUrls))

	//htmlLegislationCollector := *createHtmlLegislationCollector()
	//pdfLegislationCollector := *createPdfLegislationCollector()

	//menuCollector.Visit()
	//
	//// Test PDF & OCR
	//// Short PDF
	////pdfLegislationCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/viewdoc?regactid=413516&doctype=reg&findpdfurl=true")
	//// Long PDF
	////pdfLegislationCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/viewdoc?regactid=413518&doctype=reg&findpdfurl=true")
	//
	//menuCollector.Wait()
	//legislationUrlCollector.Wait()
	//pdfLegislationCollector.Wait()
	//htmlLegislationCollector.Wait()

}

func writeIndexToFile(filePath string, urls []string) error {
	if file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		return err
	} else {
		defer file.Close()
		_, err := file.WriteString(strings.Join(urls, "\n"))
		return err
	}
}
