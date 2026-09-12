package backup

import "context"

// Downloader downloads a config backup blob.
type Downloader interface {
	Download(ctx context.Context, scope string) ([]byte, error)
}
