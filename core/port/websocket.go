package port

import "context"

type WebSocket interface {
	Run(ctx context.Context) error
}
