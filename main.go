package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type LinkStatus struct {
	URL    string `json:"url"`
	Status string `json:"status"`
}

func main() {

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
