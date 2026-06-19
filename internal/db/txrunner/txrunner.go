package txrunner

import "context"

type Runner interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}
