package main

import (
	"log"

	"github.com/hhrutter/pdfcpu/pkg/api"
	"github.com/hhrutter/pdfcpu/pkg/pdfcpu"
	"github.com/namsral/flag"
)

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
	flag.Var(&inputFiles, "input", "PDFs to merge")
	flag.StringVar(&outputFile, "output", "", "Final PDF to write to")
	flag.Parse()

	if len(inputFiles) == 0 {
		// exit 1
	}

	if len(outputFile) == 0 {
		// exit 1
	}

	config := pdfcpu.NewDefaultConfiguration()
	err := api.Merge(inputFiles, outputFile, config)
	if err != nil {
		log.Panic(err)
	}
}
