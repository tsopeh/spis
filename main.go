package main

import (
	"encoding/csv"
	"github.com/cheggaaa/pb/v3"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
)

func main() {

	outputDirPath := filepath.Join("./", "OUTPUT")
	indexDirPath := filepath.Join(outputDirPath, "index")
	documentsDirPath := filepath.Join(outputDirPath, "documents")

	//check(os.RemoveAll(outputDirPath))
	check(os.RemoveAll(documentsDirPath))

	check(os.MkdirAll(indexDirPath, os.ModePerm))
	check(os.MkdirAll(documentsDirPath, os.ModePerm))

	htmlDocumentUrls, pdfDocumentUrls, _ := fetchOrLoadUrlIndex(path.Join(indexDirPath, "document-urls.csv"))

	barPool, err := pb.StartPool()
	check(err)

	htmlCollector := *createHtmlLegislationCollector(htmlDocumentUrls[:10], documentsDirPath, barPool)
	pdfCollector := *createPdfLegislationCollector(pdfDocumentUrls[:2], documentsDirPath, barPool)

	htmlCollector.Wait()
	pdfCollector.Wait()

}

func fetchOrLoadUrlIndex(indexFilePath string) ([]DocumentUrl, []DocumentUrl, []DocumentUrl) {
	var htmlDocumentUrls []DocumentUrl
	var pdfDocumentUrls []DocumentUrl
	var unknownDocumentsFoundInMenuUrls []DocumentUrl
	if existingCsvFile, err := os.Open(indexFilePath); err != nil {
		log.Println("URL index not found. Fetching a new index.")
		var menuUrls []string
		fetchMenuUrls("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/api/reg/menu", &menuUrls)
		log.Println("Successfully fetched the list of " + strconv.Itoa(len(menuUrls)) + " menu URLs.")
		fetchDocumentUrls(menuUrls, &htmlDocumentUrls, &pdfDocumentUrls, &unknownDocumentsFoundInMenuUrls)
		log.Println("Found " + strconv.Itoa(len(htmlDocumentUrls)) + " HTML document URLs.")
		log.Println("Found " + strconv.Itoa(len(pdfDocumentUrls)) + " PDF document URLs.")
		log.Println("Found " + strconv.Itoa(len(unknownDocumentsFoundInMenuUrls)) + " menu items with UNKNOWN document URLs.")
		check(writeDocumentUrlsToCSV(indexFilePath, append(htmlDocumentUrls, append(pdfDocumentUrls, unknownDocumentsFoundInMenuUrls...)...)))
		return htmlDocumentUrls, pdfDocumentUrls, unknownDocumentsFoundInMenuUrls
	} else {
		log.Println("Discovered the existing URL index. Loading in progress.")
		reader := csv.NewReader(existingCsvFile)
		rows, err := reader.ReadAll()
		check(err)
		for _, row := range rows[1:] {
			documentUrl := DocumentUrl{
				url:           row[2],
				kind:          row[1],
				parentMenuUrl: row[0],
			}
			switch documentUrl.kind {
			case "HTML":
				{
					htmlDocumentUrls = append(htmlDocumentUrls, documentUrl)
				}
			case "PDF":
				{
					pdfDocumentUrls = append(pdfDocumentUrls, documentUrl)
				}
			default:
				{
					unknownDocumentsFoundInMenuUrls = append(unknownDocumentsFoundInMenuUrls, documentUrl)
				}
			}
		}
		log.Println("Found " + strconv.Itoa(len(htmlDocumentUrls)) + " HTML document URLs.")
		log.Println("Found " + strconv.Itoa(len(pdfDocumentUrls)) + " PDF document URLs.")
		log.Println("Found " + strconv.Itoa(len(unknownDocumentsFoundInMenuUrls)) + " menu items with UNKNOWN document URLs.")
		log.Println("Successfully loaded the existing index.")
		return htmlDocumentUrls, pdfDocumentUrls, unknownDocumentsFoundInMenuUrls
	}

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
