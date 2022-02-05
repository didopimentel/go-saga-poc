package saga

import "context"

type Saga struct {
	Steps []Step
}

type Step struct {
	Command             func(ctx context.Context) (interface{}, error)
	CompensationCommand func(ctx context.Context) (interface{}, error)
}

func NewSaga(steps []Step) Saga {
	return Saga{Steps: steps}
}
