package main

import (
	"encoding/csv"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
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

	var htmlDocumentUrls []DocumentUrl
	var pdfDocumentUrls []DocumentUrl
	var unknownDocumentsFoundInMenuUrls []DocumentUrl
	fetchDocumentUrls(menuUrls, &htmlDocumentUrls, &pdfDocumentUrls, &unknownDocumentsFoundInMenuUrls)
	log.Println("Found " + strconv.Itoa(len(htmlDocumentUrls)) + " HTML document URLs.")
	log.Println("Found " + strconv.Itoa(len(pdfDocumentUrls)) + " PDF document URLs.")
	log.Println("Found " + strconv.Itoa(len(unknownDocumentsFoundInMenuUrls)) + " menu items with UNKNOWN document URLs.")
	check(writeDocumentUrlsToCSV(path.Join(indexDirPath, "document-urls.csv"), append(htmlDocumentUrls, append(pdfDocumentUrls, unknownDocumentsFoundInMenuUrls...)...)))

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

func writeDocumentUrlsToCSV(csvPath string, urls []DocumentUrl) error {
	if csvFile, err := os.OpenFile(csvPath, os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		return err
	} else {
		defer csvFile.Close()
		csvWriter := csv.NewWriter(csvFile)
		rows := [][]string{
			{"parent_menu_url", "kind", "document_url"},
		}
		for _, documentUrl := range urls {
			rows = append(rows, []string{documentUrl.parentMenuUrl, documentUrl.kind, documentUrl.url})
		}
		err = csvWriter.WriteAll(rows)
		return err
	}
}
