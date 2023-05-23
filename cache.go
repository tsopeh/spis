package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"sync"
)

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

type ScraperCache struct {
	Add   func(FinishedDocument)
	Get   func(requestUrl string) (FinishedDocument, bool)
	Close func()
}

type FinishedDocument struct {
	requestUrl  string
	responseUrl string
	title       string
	fileName    string
	contentHash string
}

func createOrLoadScraperCache(
	cacheFilePath string,
) ScraperCache {

	memCache := make(map[string]FinishedDocument)

	fileCache, err := os.OpenFile(cacheFilePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	check(err)

	writer := csv.NewWriter(fileCache)

	stats, err := fileCache.Stat()
	check(err)
	isFileEmpty := stats.Size() == 0
	if isFileEmpty {
		check(writer.Write([]string{"request_url", "response_url", "title", "content_hash", "filename"}))
		writer.Flush()
	} else {
		// Read the cache
		reader := csv.NewReader(fileCache)
		rows, err := reader.ReadAll()
		check(err)
		for _, row := range rows[1:] {
			memCache[row[0]] = FinishedDocument{
				requestUrl:  row[0],
				responseUrl: row[1],
				title:       row[2],
				contentHash: row[3],
				fileName:    row[4],
			}
		}
	}

	sem := sync.RWMutex{}
	return ScraperCache{
		Add: func(document FinishedDocument) {
			sem.Lock()
			defer sem.Unlock()
			check(writer.Write([]string{document.requestUrl, document.responseUrl, document.title, document.contentHash, document.fileName}))
			writer.Flush()
			memCache[document.requestUrl] = document
		},
		Get: func(requestUrl string) (FinishedDocument, bool) {
			sem.RLock()
			defer sem.RUnlock()
			document, ok := memCache[requestUrl]
			return document, ok
		},
		Close: func() {
			check(fileCache.Close())
		},
	}

}
