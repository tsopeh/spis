package main

// TODO: Show progress for fetching and processing HTML and PDF documents.
// TODO: Capture stdout for OCR
// TODO: Experiment with goroutines to make thing more concurrent (execute faster).

func main() {

	htmlLegislationCollector := *createHtmlLegislationCollector()
	pdfLegislationCollector := *createPdfLegislationCollector()
	legislationUrlCollector := *createLegislationUrlCollector(&htmlLegislationCollector, &pdfLegislationCollector)
	menuCollector := *createMenuCollector(&legislationUrlCollector)

	//menuCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/api/reg/menu")

	// Test PDF & OCR
	legislationUrlCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/RegistarServlet?subareaid=545")

	menuCollector.Wait()
	legislationUrlCollector.Wait()
	pdfLegislationCollector.Wait()
	htmlLegislationCollector.Wait()

}
