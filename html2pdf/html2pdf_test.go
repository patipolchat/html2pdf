package html2pdf

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestConvertHtmlToPdf(t *testing.T) {
	tests := []struct {
		name        string
		htmlContent string
		wantErr     bool
		opts        []Option
	}{
		{
			name:        "simple HTML",
			htmlContent: "<html><body><h1>Hello World</h1></body></html>",
			wantErr:     false,
		},
		{
			name:        "complex HTML with CSS",
			htmlContent: `<html><head><style>body{font-family:Arial;}</style></head><body><h1>Test</h1><p>This is a test.</p></body></html>`,
			wantErr:     false,
		},
		{
			name:        "empty HTML",
			htmlContent: "",
			wantErr:     false,
		},
		{
			name:        "with custom logger",
			htmlContent: "<html><body><h1>Test with Logger</h1></body></html>",
			wantErr:     false,
			opts: []Option{
				WithLogger(func(format string, args ...interface{}) {
					// Custom logger that does nothing (silent mode)
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			got, err := ConvertHtmlToPdf(ctx, tt.htmlContent, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertHtmlToPdf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("ConvertHtmlToPdf() returned empty PDF")
			}
			if !tt.wantErr && !strings.HasPrefix(string(got), "%PDF") {
				t.Errorf("ConvertHtmlToPdf() returned non-PDF content")
			}
		})
	}
}

func TestConvertHtmlFileToPdf(t *testing.T) {
	// Create a temporary HTML file for testing
	tempHTML := `<html><body><h1>Test File</h1><p>This is a test file.</p></body></html>`
	tempFile, err := os.CreateTemp("", "test-*.html")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(tempHTML); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tempFile.Close()

	tests := []struct {
		name     string
		fileName string
		wantErr  bool
		opts     []Option
	}{
		{
			name:     "valid HTML file",
			fileName: tempFile.Name(),
			wantErr:  false,
		},
		{
			name:     "non-existent file",
			fileName: "non-existent-file.html",
			wantErr:  true,
		},
		{
			name:     "valid HTML file with custom logger",
			fileName: tempFile.Name(),
			wantErr:  false,
			opts: []Option{
				WithLogger(func(format string, args ...interface{}) {
					// Custom logger that does nothing
				}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			got, err := ConvertHtmlFileToPdf(ctx, tt.fileName, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertHtmlFileToPdf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("ConvertHtmlFileToPdf() returned empty PDF")
			}
			if !tt.wantErr && !strings.HasPrefix(string(got), "%PDF") {
				t.Errorf("ConvertHtmlFileToPdf() returned non-PDF content")
			}
		})
	}
}

func TestWithLogger(t *testing.T) {
	var logCalled bool
	var logMessage string

	customLogger := func(format string, args ...interface{}) {
		logCalled = true
		logMessage = fmt.Sprintf(format, args...)
	}

	opts := &options{}
	WithLogger(customLogger)(opts)

	if opts.logger == nil {
		t.Error("WithLogger() did not set the logger function")
	}

	// Test that the logger function works
	opts.logger("test message %s", "value")
	if !logCalled {
		t.Error("Custom logger was not called")
	}
	if logMessage != "test message value" {
		t.Errorf("Expected log message 'test message value', got '%s'", logMessage)
	}
}

func TestGetDefaultOptions(t *testing.T) {
	opts := getDefaultOptions()

	if opts.logger == nil {
		t.Error("getDefaultOptions() returned nil logger")
	}

	// Test that the default logger works
	opts.logger("test default logger")
	// If we get here without panic, the default logger is working
}

func TestErrHTMLFileNotFound(t *testing.T) {
	if ErrHTMLFileNotFound == nil {
		t.Error("ErrHTMLFileNotFound should not be nil")
	}

	expectedMsg := "html file not found"
	if ErrHTMLFileNotFound.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, ErrHTMLFileNotFound.Error())
	}
}

func TestConvertHtmlToPdfWithContextTimeout(t *testing.T) {
	// Test with a very short timeout to ensure context cancellation works
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Use a complex HTML that might take longer to process
	htmlContent := `<html><body><h1>Test</h1><p>This is a test with timeout.</p></body></html>`

	_, err := ConvertHtmlToPdf(ctx, htmlContent)
	if err == nil {
		t.Error("Expected error due to context timeout, but got none")
	}
}

func TestConvertHtmlToPdfWithNilLogger(t *testing.T) {
	// Test with a nil logger option
	opts := &options{}
	WithLogger(nil)(opts)

	if opts.logger != nil {
		t.Error("WithLogger(nil) should set logger to nil")
	}
}

func BenchmarkConvertHtmlToPdf(b *testing.B) {
	htmlContent := `<html><body><h1>Benchmark Test</h1><p>This is a benchmark test.</p></body></html>`
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ConvertHtmlToPdf(ctx, htmlContent)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkConvertHtmlToPdfWithCustomLogger(b *testing.B) {
	htmlContent := `<html><body><h1>Benchmark Test</h1><p>This is a benchmark test with custom logger.</p></body></html>`
	ctx := context.Background()

	customLogger := func(format string, args ...interface{}) {
		// Silent logger for benchmarking
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ConvertHtmlToPdf(ctx, htmlContent, WithLogger(customLogger))
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
