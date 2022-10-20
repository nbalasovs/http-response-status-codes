package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	code          int
	path          string
	responseCodes ResponseCodes
)

type ResponseCodes []ResponseCode

type ResponseCode struct {
	Code        int    `json:"code"`
	Phrase      string `json:"phrase"`
	Description string `json:"description"`
}

type Info struct {
	Phrase      string
	Description string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "%s Invalid values are logged to stderr.\n", os.Args[0])
	}
	flag.IntVar(&code, "c", code, "`HTTP response status code`")
	flag.StringVar(&path, "f", path, "`Path to the HTTP status response codes file")

	flag.Parse()

	if code == 0 {
		log.Fatal("HTTP response code is required")
	} else if path == "" {
		log.Fatal("File path is required")
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open HTTP response status code file: %s", err)
	}

	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Unable to read JSON: %s", err)
	}

	uerr := json.Unmarshal(data, &responseCodes)
	if uerr != nil {
		log.Fatalf("Unable to unmarshal JSON: %s", uerr)
	}

	codeMap := make(map[int]Info, len(responseCodes))

	for _, responseCode := range responseCodes {
		codeMap[responseCode.Code] = Info{responseCode.Phrase, responseCode.Description}
	}

	if _, ok := codeMap[code]; !ok {
		fmt.Printf("[%d] Unable to find description for the given HTTP response code\n", code)
		os.Exit(0)
	}

	phrase := codeMap[code].Phrase
	description := codeMap[code].Description

	fmt.Printf("[%d] %s: %s\n", code, phrase, description)
}
