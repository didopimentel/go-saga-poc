package saga

import "context"

type SagaContextKey string

// ParamKey must be used to pass down parameters through steps
// it will be overridden each step
const ParamKey = SagaContextKey("saga-context-param")

type Coordinator struct {
	saga               Saga
	errors             []error
	compensationErrors []error
	currentStep        int
	ctx                context.Context
	result             interface{}
}

func NewCoordinator(saga Saga) *Coordinator {
	return &Coordinator{
		saga: saga,
	}
}

func (c *Coordinator) Execute(ctx context.Context) (interface{}, bool) {
	c.ctx = ctx

	for i, step := range c.saga.Steps {
		c.currentStep = i
		ok := c.executeStep(step)
		if !ok {
			return nil, false
		}
	}

	return c.result, true
}

func (c *Coordinator) executeStep(step Step) bool {
	response, err := step.Command(c.ctx)
	if err != nil {
		c.errors = append(c.errors, err)
		c.compensateStep(c.currentStep)
		return false
	}
	c.ctx = context.WithValue(c.ctx, ParamKey, response)

	if c.currentStep == len(c.saga.Steps)-1 {
		c.result = response
	}

	return true
}

func (c *Coordinator) compensateStep(index int) {
	if index < 0 {
		return
	}

	_, err := c.saga.Steps[index].CompensationCommand(c.ctx)
	if err != nil {
		c.compensationErrors = append(c.compensationErrors, err)
	}

	index = index - 1
	c.compensateStep(index)
}

func (c *Coordinator) GetErrors() []error {
	return c.errors
}

func (c *Coordinator) GetCompensationErrors() []error {
	return c.compensationErrors
}
