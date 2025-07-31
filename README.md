# html2pdf

A Go package for converting HTML content to PDF using Chrome DevTools Protocol (CDP) via the `chromedp` library.

## Features

- Convert HTML strings to PDF
- Convert HTML files to PDF
- Customizable logging for debugging
- Context-aware operations with timeout support
- High-quality PDF output using Chrome's rendering engine
- Support for CSS styling and modern web features

## Installation

```bash
go get github.com/patipolchat/html2pdf
```

## Dependencies

This package requires:
- Go 1.24 or later
- Chrome/Chromium browser installed on the system
- `github.com/chromedp/chromedp` for Chrome DevTools Protocol communication

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "github.com/patipolchat/html2pdf"
    "os"
)

func main() {
    ctx := context.Background()
    
    // Convert HTML string to PDF
    htmlContent := `<html><body><h1>Hello World</h1></body></html>`
    pdfBytes, err := html2pdf.ConvertHtmlToPdf(ctx, htmlContent)
    if err != nil {
        panic(err)
    }
    
    // Save to file
    err = os.WriteFile("output.pdf", pdfBytes, 0644)
    if err != nil {
        panic(err)
    }
}
```

### Convert HTML File to PDF

```go
package main

import (
    "context"
    "github.com/patipolchat/html2pdf"
    "os"
)

func main() {
    ctx := context.Background()
    
    // Convert HTML file to PDF
    pdfBytes, err := html2pdf.ConvertHtmlFileToPdf(ctx, "input.html")
    if err != nil {
        panic(err)
    }
    
    // Save to file
    err = os.WriteFile("output.pdf", pdfBytes, 0644)
    if err != nil {
        panic(err)
    }
}
```

## API Reference

### Functions

#### `ConvertHtmlToPdf(ctx context.Context, htmlContent string, opts ...Option) ([]byte, error)`

Converts HTML content string to PDF.

**Parameters:**
- `ctx`: Context for cancellation and timeout control
- `htmlContent`: HTML content as string
- `opts`: Optional configuration options

**Returns:**
- `[]byte`: PDF content as bytes
- `error`: Error if conversion fails

#### `ConvertHtmlFileToPdf(ctx context.Context, fileName string, opts ...Option) ([]byte, error)`

Converts HTML file to PDF.

**Parameters:**
- `ctx`: Context for cancellation and timeout control
- `fileName`: Path to HTML file
- `opts`: Optional configuration options

**Returns:**
- `[]byte`: PDF content as bytes
- `error`: Error if conversion fails

### Options

#### `WithLogger(logger func(string, ...interface{})) Option`

Sets a custom logger function for debugging output.

```go
// Custom logger example
customLogger := func(format string, args ...interface{}) {
    fmt.Printf("[DEBUG] "+format, args...)
}

pdfBytes, err := html2pdf.ConvertHtmlToPdf(ctx, htmlContent, 
    html2pdf.WithLogger(customLogger))
```

### Error Types

- `ErrHTMLFileNotFound`: Returned when the specified HTML file does not exist

## Advanced Usage

### With Context Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

pdfBytes, err := html2pdf.ConvertHtmlToPdf(ctx, htmlContent)
if err != nil {
    log.Fatal(err)
}
```

### With Custom Logger

```go
// Silent mode (no logging)
silentLogger := func(format string, args ...interface{}) {
    // Do nothing
}

pdfBytes, err := html2pdf.ConvertHtmlToPdf(ctx, htmlContent,
    html2pdf.WithLogger(silentLogger))
```

### Complex HTML with CSS

```go
htmlContent := `
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        .highlight { background-color: yellow; }
    </style>
</head>
<body>
    <h1>Report Title</h1>
    <p class="highlight">This is highlighted text.</p>
    <table border="1">
        <tr><td>Data 1</td><td>Data 2</td></tr>
    </table>
</body>
</html>
`

pdfBytes, err := html2pdf.ConvertHtmlToPdf(ctx, htmlContent)
```

## Examples

See the `main.go` file in the root directory for a complete working example:

```go
package main

import (
    "context"
    "github.com/patipolchat/html2pdf"
    "os"
)

func main() {
    // Convert HTML file to PDF
    b, err := html2pdf.ConvertHtmlFileToPdf(context.Background(), "datatest/data1.html")
    if err != nil {
        panic(err)
    }

    // Save PDF to file
    err = os.WriteFile("tmp/report.pdf", b, 0644)
    if err != nil {
        panic(err)
    }
}
```

## Testing

Run the test suite:

```bash
go test ./html2pdf
```

Run benchmarks:

```bash
go test -bench=. ./html2pdf
```

## Requirements

- **Go**: 1.24 or later
- **Chrome/Chromium**: Must be installed on the system
- **Dependencies**: See `go.mod` for complete list

## Performance

The package uses Chrome DevTools Protocol for high-quality rendering. Performance depends on:
- HTML complexity
- CSS styling
- System resources
- Chrome/Chromium performance

For production use, consider:
- Using appropriate context timeouts
- Implementing retry logic for failed conversions
- Monitoring memory usage for large HTML documents

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

## Troubleshooting

### Common Issues

1. **Chrome not found**: Ensure Chrome/Chromium is installed and accessible
2. **Context timeout**: Increase timeout duration for complex HTML
3. **Memory issues**: Consider processing large documents in chunks
4. **Permission errors**: Ensure write permissions for output directory

### Debug Mode

Enable debug logging to troubleshoot issues:

```go
debugLogger := func(format string, args ...interface{}) {
    log.Printf("[DEBUG] "+format, args...)
}

pdfBytes, err := html2pdf.ConvertHtmlToPdf(ctx, htmlContent,
    html2pdf.WithLogger(debugLogger))
``` 