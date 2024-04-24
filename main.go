package main

import (
	// "bytes"
	"fmt"
	// "io"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("[USAGE]: go main . <file_name>") // file that contains the urls that is needed to be scraped
		return
	}

	filename := args[1]
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error opening file: %v", filename)
	}

	urls := strings.Split(string(fileContent), "\n")
	var wg sync.WaitGroup

	for _, url := range urls {
		if url == "" {
			continue
		}

		// will add one wait group for each valid url and wait for all goroutines to execute
		wg.Add(1)
		go func(url string) {
			defer wg.Done()

			res, err := http.Get(url)
			if err != nil {
				log.Fatalf("Error fetching url: %v", err)
			}

			defer res.Body.Close()
			ParseHtmlPage(res)
		}(url)
	}
	wg.Wait()
}

// this will handle the parsing of the html and extract needed text from it
func ParseHtmlPage(res *http.Response) {
	doc, err := html.Parse(res.Body)
	if err != nil {
		log.Fatalf("Error parsing html page: %v", err)
	}

	var processPage func(doc *html.Node)
	// get all data within the body tag
	processPage = func(doc *html.Node) {
		// check if body tag is found
		if doc.Type == html.ElementNode && doc.Data == "body" {
			// if body tag is found call the function
			PrintTextFromBody(doc, 0)
			return
		}

		for c := doc.FirstChild; c != nil; c = c.NextSibling {
			// Recursively search for the <body> tag in child nodes
			processPage(c)
		}
	}

	processPage(doc) // initial call to start the process
}

// TEMP: for now only print the text until further development
// this function is responsible to recursivly go through the body tag and print all text nodes inside
func PrintTextFromBody(n *html.Node, depth int) {
	if n == nil {
		return
	}

	indent := strings.Repeat(" ", depth*2)

	if n.Type == html.ElementNode {
		switch n.Data {
		case "h1", "h2", "h3":
			fmt.Printf("%sHeader: %s\n", indent, getTextContent(n))
		case "p", "span":
			fmt.Printf("%sParagraph: %s\n", indent, getTextContent(n))
		case "li":
			fmt.Printf("%sList Item: %s\n", indent, getTextContent(n))
		case "table":
			fmt.Printf("%sTable:\n", indent)
		case "tr":
			fmt.Printf("%sTableRow:\n", indent)
		case "th":
			fmt.Printf("%sTableHeader: %s\n", indent, getTextContent(n))
		case "td":
			fmt.Printf("%sTableCell: %s\n", indent, getTextContent(n))
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		PrintTextFromBody(c, depth+1)
	}
}

func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return strings.TrimSpace(n.Data)
	}
	content := ""
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		content += getTextContent(c) + " "
	}
	return strings.TrimSpace(content)
}
