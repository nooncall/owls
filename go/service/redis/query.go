package redis

import (
	"context"
)

type Params struct {
	Cmd, Cluster string
	DB           int
}

// ExecQuery ...
func ExecQuery(ctx context.Context, req *Params) interface{} {
	resp, err := ExecReadTask(ctx, req)
	if err != nil {
		return err.Error()

	}
	return resp
}
