package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

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
	var skipPreview = flag.Bool("s", false, "Skip preview MD file in browser")
	flag.Parse()

	// no file, inform of usage
	if *fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*fileName, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

// Reads, then parses
func run(fileName string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	htmlData := parseContent(input)

	// create temp file
	temp, err := os.CreateTemp("", "mdp*.html") // temp == mdp<randomPattern>.html
	if err != nil {
		return err
	}

	if err := temp.Close(); err != nil {
		return err
	}

	// display created .html file
	outName := temp.Name()
	fmt.Println(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}
	return preview(outName)
	// return saveHTML(outName, htmlData)
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

func preview(fname string) error {
	var cName = ""
	var cParams = []string{}

	// define executable
	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}
	// append filename to params slice
	cParams = append(cParams, fname)
	// locate executable
	cPath, err := exec.LookPath(cName)

	if err != nil {
		return err
	}
	return exec.Command(cPath, cParams...).Run()
}
