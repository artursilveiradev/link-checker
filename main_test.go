package main

import (
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
