# SPIS

SPIS — **S**craper for "**P**ravno **I**nformacioni **S**istem"

## The problem

We want to aggregate as much legislative text in serbian language as possible. The website [Pravno Informacioni Sistem](https://www.pravno-informacioni-sistem.rs/reg-overview) has a registry of many legislative documents in serbian. Of all documents, two thirds are in a form of HTML pages, which are easy enough to scrape using [Colly](https://github.com/gocolly/colly). The last third consists of PDF documents; where each of their pages is an image of text. We want to run an OCR tool on each page and combine the results afterward. 

### Project requirements

* Create a separate text file for each document (HTML and PDF).
* Ignore tabular data as much as possible in HTML documents.
* Create OCR pipeline for PDF documents.
* Stopping and resuming work (keep track of already completed work).
* Progress tracking.
* Parallel processing of documents.

## Usage (Docker)

1. Build image image.

```shell
docker build -t spis .
```

2. Instantiate a container.

```shell
docker run --rm -it --mount type=bind,src="$(pwd)/OUTPUT",target="/go/spis/OUTPUT" --name spis-con spis
```

## Local development

### Prerequisites

* Go [compiler](https://go.dev/dl/)
* Tesseract (for OCR) with language packs for _serbian-cyrillic_ and _serbian-latin_.

### Dev

```shell
go run *.go
```

### Build

Manually, without docker.

```shell
go build -o ./spis
```

Alternatively, see the [above section](#usage-docker).
