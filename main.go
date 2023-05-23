package main

import (
	"github.com/cheggaaa/pb/v3"
	"os"
	"path"
	"path/filepath"
)

// TODO: Docker.

func main() {

	outputDirPath := filepath.Join("./", "OUTPUT")
	indexDirPath := filepath.Join(outputDirPath, "index")
	documentsDirPath := filepath.Join(outputDirPath, "documents")

	//check(os.RemoveAll(outputDirPath))
	//check(os.RemoveAll(documentsDirPath))

	check(os.MkdirAll(indexDirPath, os.ModePerm))
	check(os.MkdirAll(documentsDirPath, os.ModePerm))

	htmlDocumentUrls, pdfDocumentUrls, _ := fetchOrLoadUrlIndex(path.Join(indexDirPath, "document-urls.csv"))
	cache := createOrLoadScraperCache(path.Join(indexDirPath, "already-processed.csv"))

	barPool, err := pb.StartPool()
	barPool.Output = os.Stdout
	check(err)

	htmlCollector := *createHtmlLegislationCollector(htmlDocumentUrls, documentsDirPath, &cache, barPool)
	pdfCollector := *createPdfLegislationCollector(pdfDocumentUrls, documentsDirPath, &cache, barPool)

	htmlCollector.Wait()
	pdfCollector.Wait()
	cache.Close()

}
