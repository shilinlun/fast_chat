package port

import "context"

type WebSocketHandler interface {
	Run(ctx context.Context) error
}
