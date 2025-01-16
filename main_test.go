package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestExtractLinks(t *testing.T) {
	htmlContent := `
	<html>
		<body>
			<a href="https://www.google.com">Google</a>
			<a href="https://www.example.com">Example</a>
			<a href="mailto:test@example.com">Email</a>
			<a href="#section1">Section 1</a>
		</body>
	</html>
	`

	tempFile, err := os.CreateTemp("", "test*.html")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString(htmlContent)
	if err != nil {
		t.Fatalf("failed to write content to temp file: %v", err)
	}
	tempFile.Close()

	got, err := extractLinks(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to extract links: %v", err)
	}

	want := []string{
		"https://www.google.com",
		"https://www.example.com",
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestCheckLink(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ok":
			w.WriteHeader(http.StatusOK)
		case "/notfound":
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	tests := []struct {
		url      string
		expected string
	}{
		{server.URL + "/ok", "200 OK"},
		{server.URL + "/notfound", "404 Not Found"},
		{server.URL + "/error", "500 Internal Server Error"},
		{"http://invalid-url", "Erro"},
	}

	for _, test := range tests {
		t.Run(test.url, func(t *testing.T) {
			result := checkLink(test.url)
			if result != test.expected {
				t.Errorf("checkLink(%q) = %q; want %q", test.url, result, test.expected)
			}
		})
	}
}

func TestSaveReport(t *testing.T) {
	report := []LinkStatus{
		{URL: "https://www.google.com", Status: "200 OK"},
		{URL: "https://www.example.com", Status: "404 Not Found"},
	}

	tempFile, err := os.CreateTemp("", "report*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	err = saveReport(report, tempFile.Name())
	if err != nil {
		t.Fatalf("failed to save report: %v", err)
	}

	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to read content from temp file: %v", err)
	}

	expected := `[
  {
    "url": "https://www.google.com",
    "status": "200 OK"
  },
  {
    "url": "https://www.example.com",
    "status": "404 Not Found"
  }
]`

	if string(content) != expected {
		t.Errorf("got: %s, want: %s", content, expected)
	}
}
