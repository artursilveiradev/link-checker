# Link Checker ⚓️

A lightweight Golang command-line application to extract and verify the HTTP status of external links in HTML files, with detailed reporting in JSON format

![CI](https://github.com/artursilveiradev/link-checker/actions/workflows/ci.yml/badge.svg)

## Build

Compile the application:
```bash
make build
```

## Run tests

Execute all unit tests:
```bash
make test
```

## Usage

### Simple check
```
link-checker -file example.html
```

### Custom report name
```
link-checker -file example.html -output custom-report-name.json
```

### Print verbose output
```
link-checker -file example.html -verbose
```

## Report example
```json
[
  {
    "url": "https://www.google.com",
    "status": "200 OK"
  },
  {
    "url": "https://www.example.com",
    "status": "200 OK"
  },
  {
    "url": "https://invalid-link.com",
    "status": "Erro"
  }
]
```