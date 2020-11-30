package deployer

import (
	"context"
)

type Executor interface {
	Exec(ctx context.Context,
		deployer Deployer,
		job Job) error
}
