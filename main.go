package main

import (
	"context"
	"github.com/patipolchat/html2pdf/html2pdf"
	"os"
)

func main() {
	// log the CDP messages so that you can find the one to use.
	b, err := html2pdf.ConvertHtmlFileToPdf(context.Background(), "datatest/data1.html")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("tmp/report.pdf", b, 0644)
	if err != nil {
		panic(err)
	}
}
