package redis

import (
	"context"
)

type Params struct {
	cmd, cluster string
	db           int
}

// ExecQuery ...
func ExecQuery(ctx context.Context, req *Params) interface{} {
	resp, err := ExecReadTask(ctx, req)
	if err != nil {
		return err.Error()

	}
	return resp
}
