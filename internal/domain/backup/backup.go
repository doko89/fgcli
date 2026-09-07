package backup

import "context"

// Repository downloads a config backup blob.
type Repository interface {
	Download(ctx context.Context, scope string) ([]byte, error)
}
