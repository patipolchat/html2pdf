package html2pdf

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	// ErrHTMLFileNotFound is returned when the specified file does not exist.
	ErrHTMLFileNotFound = fmt.Errorf("html file not found")
)

// Option represents a configuration option for the PDF conversion functions.
type Option func(*options)

type options struct {
	logger func(string, ...interface{})
}

// WithLogger sets a custom logger function for debugging output.
func WithLogger(logger func(string, ...interface{})) Option {
	return func(o *options) {
		o.logger = logger
	}
}

// getDefaultOptions returns the default options.
func getDefaultOptions() *options {
	return &options{
		logger: log.Printf,
	}
}

// ConvertHtmlFileToPdf reads an HTML file and converts its content to PDF.
func ConvertHtmlFileToPdf(ctx context.Context, fileName string, opts ...Option) ([]byte, error) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrHTMLFileNotFound
		}
		return nil, fmt.Errorf("failed to read file %s: %w", fileName, err)
	}
	return ConvertHtmlToPdf(ctx, string(b), opts...)
}

// ConvertHtmlToPdf converts HTML content to PDF using chromedp.
func ConvertHtmlToPdf(ctx context.Context, htmlContent string, opts ...Option) ([]byte, error) {
	options := getDefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	ctx, cancel := chromedp.NewContext(ctx, chromedp.WithDebugf(options.logger))
	defer cancel()

	var buf []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var wg sync.WaitGroup
			wg.Add(1)
			chromedp.ListenTarget(ctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					wg.Done()
				}
			})
			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			if err := page.SetDocumentContent(frameTree.Frame.ID, htmlContent).Do(ctx); err != nil {
				return err
			}
			wg.Wait()
			return nil
		}),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, _, err = page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to convert HTML to PDF: %w", err)
	}
	return buf, nil
}
