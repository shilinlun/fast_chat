package port

import "context"

type HttpHandler interface {
	Run(ctx context.Context) error
}
