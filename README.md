# PDFManager

This tool serves as a baseline for operating with pdfs -> splitting and merging them without the need for online connections and 3rd party versions stealing your data.
The **Manager** package is the exposed API for how to use this tool, and simple CLI is implemented for easy use locally before integrations in other software.
The purpose for this tool is to add learned methods for automatically splitting PDFs based on content, and to be integrated in the **Harmoniq** suite.

## Prerequisites

You need Golang 1.26.3 or newer.

## How to use

1. Clone or download this repo
2. Run `go build -o pdfmanager .`
3. Run the binary with `./pdfmanager` and follow CLI intstructions

## How to use package in your own project

1. Run `go get github.com/brottum-brass/pdfmanager@latest`
