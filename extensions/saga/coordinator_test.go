package saga_test

import (
	"context"
	"errors"
	"github.com/didopimentel/go-saga-poc/extensions/saga"
	"github.com/stretchr/testify/require"
	"testing"
)

type InitialPayload struct {
	field int64
}

type Step1Response struct {
	field int64
}

type Step2Response struct {
	field int64
}

type Step3Response struct {
	field int64
}

func TestCoordinator_Success(t *testing.T) {
	calledFirstStep := false
	calledSecondStep := false
	calledThirdStep := false
	finalField := int64(0)
	steps := []saga.Step{
		{
			Command: func(ctx context.Context) (interface{}, error) {
				f := ctx.Value(saga.ParamKey).(InitialPayload)
				calledFirstStep = true
				return Step1Response{field: f.field + 1}, nil
			},
			CompensationCommand: nil,
		},
		{
			Command: func(ctx context.Context) (interface{}, error) {
				f := ctx.Value(saga.ParamKey).(Step1Response)
				calledSecondStep = true
				return Step2Response{field: f.field + 1}, nil
			},
			CompensationCommand: nil,
		},
		{
			Command: func(ctx context.Context) (interface{}, error) {
				f := ctx.Value(saga.ParamKey).(Step2Response)
				calledThirdStep = true
				return Step3Response{field: f.field + 1}, nil
			},
			CompensationCommand: nil,
		},
		{
			Command: func(ctx context.Context) (interface{}, error) {
				f := ctx.Value(saga.ParamKey).(Step3Response)
				finalField = f.field
				return nil, nil
			},
			CompensationCommand: nil,
		},
	}

	aSaga := saga.NewSaga(steps)

	coordinator := saga.NewCoordinator(aSaga)

	ctx := context.WithValue(context.Background(), saga.ParamKey, InitialPayload{field: 0})
	coordinator.Execute(ctx)

	require.True(t, calledFirstStep)
	require.True(t, calledSecondStep)
	require.True(t, calledThirdStep)
	require.Equal(t, int64(3), finalField)
	require.Equal(t, 0, len(coordinator.GetErrors()))
	require.Equal(t, 0, len(coordinator.GetCompensationErrors()))
}

func TestCoordinator_Compensation_Success(t *testing.T) {
	calledFirstCompensation := false
	calledSecondCompensation := false
	calledThirdCompensation := false
	calledFourthCommand := false
	steps := []saga.Step{
		{
			Command: func(ctx context.Context) (interface{}, error) {
				return Step1Response{}, nil
			},
			CompensationCommand: func(ctx context.Context) (interface{}, error) {
				calledFirstCompensation = true
				return nil, nil
			},
		},
		{
			Command: func(ctx context.Context) (interface{}, error) {
				return Step2Response{}, nil
			},
			CompensationCommand: func(ctx context.Context) (interface{}, error) {
				calledSecondCompensation = true
				return nil, nil
			},
		},
		{
			Command: func(ctx context.Context) (interface{}, error) {
				return Step3Response{}, errors.New("error on step 3")
			},
			CompensationCommand: func(ctx context.Context) (interface{}, error) {
				calledThirdCompensation = true
				return nil, nil
			},
		},
		{
			Command: func(ctx context.Context) (interface{}, error) {
				calledFourthCommand = true
				return nil, nil
			},
			CompensationCommand: nil,
		},
	}

	aSaga := saga.NewSaga(steps)

	coordinator := saga.NewCoordinator(aSaga)

	ctx := context.WithValue(context.Background(), saga.ParamKey, InitialPayload{field: 0})
	coordinator.Execute(ctx)

	require.True(t, calledThirdCompensation)
	require.True(t, calledSecondCompensation)
	require.True(t, calledFirstCompensation)
	require.False(t, calledFourthCommand)
	require.Equal(t, 1, len(coordinator.GetErrors()))
	require.Equal(t, 0, len(coordinator.GetCompensationErrors()))
}
