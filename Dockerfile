FROM golang:1.20.4-bullseye
RUN apt update && apt upgrade -y
RUN apt install tesseract-ocr tesseract-ocr-srp tesseract-ocr-srp-latn libtesseract-dev mupdf -y
WORKDIR spis
COPY go.mod .
RUN go mod download
COPY . .
RUN go build -o ./spis
CMD ["./spis"]
