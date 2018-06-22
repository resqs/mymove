package main

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/namsral/flag"
)

const horizontalMargin = 0
const topMargin = 0
const totalWidth = 210.0
const bodyWidth = totalWidth - (horizontalMargin * 2)

type stringSlice []string

func (i *stringSlice) String() string {
	return ""
}

func (i *stringSlice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var inputFiles stringSlice
var outputFile string

func main() {
	flag.Var(&inputFiles, "input", "Image to add to PDF")
	flag.StringVar(&outputFile, "output", "", "File to write to")
	flag.Parse()

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(horizontalMargin, topMargin, horizontalMargin)

	if len(inputFiles) == 0 {
		// exit 1
	}

	if len(outputFile) == 0 {
		// exit 1
	}

	fmt.Printf("%v", inputFiles)
	var opt gofpdf.ImageOptions
	for _, path := range inputFiles {
		pdf.AddPage()
		pdf.ImageOptions(path, horizontalMargin, topMargin, bodyWidth, 0, false, opt, 0, "")
	}

	pdf.OutputFileAndClose(outputFile)
	// if err != nil {
	// 	fmt.Printf("error: %#v\n", err)
	// }
}
