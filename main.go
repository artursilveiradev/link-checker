package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type LinkStatus struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

func main() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		fmt.Println("\nProcess interrupted by user.")
		os.Exit(1)
	}()

	inputFile := flag.String("file", "", "Path to the input HTML file")
	outputFile := flag.String("output", "report.json", "Path to the output JSON file")
	verbose := flag.Bool("verbose", false, "Print verbose output")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: Missing required flag -file")
		os.Exit(1)
	}

	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		fmt.Printf("Error: The file %s does not exist.\n", *inputFile)
		os.Exit(1)
	}

	start := time.Now()

	if *verbose {
		fmt.Printf("Extracting links from %s...\n", *inputFile)
	}
	links, err := extractLinks(*inputFile)
	if err != nil {
		fmt.Printf("Error: Failed to extract links from %s. Details: %v\n", *inputFile, err)
		os.Exit(1)
	}
	fmt.Printf("Found %d external links.\n", len(links))

	if *verbose {
		fmt.Println("Checking the status of each link...")
	}
	var results []LinkStatus
	for _, link := range links {
		if *verbose {
			fmt.Printf("Checking link: %s\n", link)
		}
		status := checkLink(link)
		results = append(results, LinkStatus{URL: link, Status: status})
	}

	if err := saveReport(results, *outputFile); err != nil {
		fmt.Printf("Error: Failed to save the report to %s. Details: %v\n", *outputFile, err)
		os.Exit(1)
	}
	fmt.Printf("Report successfully saved to: %s\n", *outputFile)

	elapsed := time.Since(start)
	fmt.Printf("Process completed in %.2f seconds.\n", elapsed.Seconds())
}

// extractLinks extracts links from an HTML file.
func extractLinks(filePath string) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	var links []string
	var crawler func(*html.Node)
	crawler = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					link := attr.Val
					if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
						links = append(links, link)
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			crawler(c)
		}
	}
	crawler(doc)
	return links, nil
}

// checkLink checks the status of a link.
func checkLink(url string) string {
	resp, err := http.Head(url)
	if err != nil {
		return "Erro"
	}
	defer resp.Body.Close()
	return resp.Status
}

// saveReport saves the report to a file.
func saveReport(report []LinkStatus, filePath string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, 0644)
}
