package main

func main() {

	legislationDocumentCollector := *createHtmlDocumentCollector()
	legislationUrlCollector := *createLegislationUrlCollector(&legislationDocumentCollector)
	menuCollector := *createMenuCollector(&legislationUrlCollector)

	menuCollector.Visit("https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/api/reg/menu")

	menuCollector.Wait()
	legislationUrlCollector.Wait()
	legislationDocumentCollector.Wait()

	//ocrHelloWorld()

}

// curl 'https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/viewdoc?regactid=413504&doctype=reg&findpdfurl=true'   -H 'X-Referer: /SlGlasnikPortal/pdfjs/build/pdf.worker.js'   -H 'Referer: https://www.pravno-informacioni-sistem.rs/SlGlasnikPortal/pdfjs/build/pdf.worker.js'   -H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36'   --compressed --output - > test.pdf