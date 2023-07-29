package port

import (
	"context"
	"fast_chat/core/entity"
)

type Storage interface {
	Insert(ctx context.Context, msg *entity.Msg) error
}
