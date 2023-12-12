package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

// create Header and Footer as 'blackfriday' does not include
const (
	header = `<!DOCTYPE html>
	<html>
	  <head>
		<meta http-equiv="content-type" content="text/html; charset=utf-8">
		<title>Markdown Preview Tool</title>
	  </head>
	 <body> `

	footer = `
	 </body>
	</html>
	`
)

func main() {
	// parse flags
	var fileName = flag.String("file", "", "Markdown file you wish to preview.")
	flag.Parse()

	// no file, inform of usage
	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

// Reads, then parses
func run(fileName string) error {
	input, err := os.ReadFile(fileName)

	if err != nil {
		return err
	}

	htmlData := parseContent(input)

	// display created .html file
	outName := fmt.Sprintf("%s.html", filepath.Base(fileName))
	fmt.Println(outName)

	return saveHTML(outName, htmlData)
}

func parseContent(input []byte) []byte {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	var buffer bytes.Buffer

	buffer.WriteString(header)
	buffer.Write(body)
	buffer.WriteString(footer)

	return buffer.Bytes()
}

func saveHTML(outName string, htmlData []byte) error {
	return os.WriteFile(outName, htmlData, 0644)
}
