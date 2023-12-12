package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

const (
	INPUT_FILE = "./testdata/test1.md"
	//RESULT_FILE   = "test1.md.html"
	EXPECTED_FILE = "./testdata/test1.md.html"
)

func TestParseContent(t *testing.T) {
	input, err := os.ReadFile(INPUT_FILE)

	if err != nil {
		t.Fatal(err)
	}

	result := parseContent(input)
	expected, err := os.ReadFile(EXPECTED_FILE)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("expected: \n%s\n", expected)
		t.Logf("result: \n%s\n", result)
		t.Error("Result content != expected contented")
	}
}

func TestRun(t *testing.T) {
	var mockStdOut bytes.Buffer // buffer acting as stdout

	if err := run(INPUT_FILE, &mockStdOut, true); err != nil {
		t.Fatal(err)
	}
	resFile := strings.TrimSpace(mockStdOut.String())

	result, err := os.ReadFile(resFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := os.ReadFile(EXPECTED_FILE)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(expected, result) {
		t.Logf("expected: \n%s\n", expected)
		t.Logf("result: \n%s\n", result)
		t.Error("Result content != expected contented")
	}

	os.Remove(resFile)
}
