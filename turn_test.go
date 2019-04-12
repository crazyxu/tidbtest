package tidbtest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func fnCreactor(duration time.Duration, name string, collector *[]string, ret error) workerFn {
	return func(ctx context.Context) error {
		if ret != nil {
			return ret
		}
		select {
		case <-time.After(duration):
			*collector = append(*collector, name)
		case <-ctx.Done():
		}
		return nil
	}
}
func TestWorkInTurns(t *testing.T) {
	turns := []string{"A", "B", "C", "C", "B", "A"}
	collector := new([]string)
	fns := map[string]workerFn{
		"A": fnCreactor(0, "A", collector, nil),
		"B": fnCreactor(0, "B", collector, nil),
		"C": fnCreactor(0, "C", collector, nil),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err := workInTurns(ctx, turns, fns)

	assert.Equal(t, nil, err)
	assert.Equal(t, turns, *collector)
}

func TestWorkInTurns_WorkNotFound(t *testing.T) {
	turns := []string{"A"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := workInTurns(ctx, turns, nil)
	assert.NotEqual(t, nil, err)
}

func TestWorkInTurns_ManualCancel(t *testing.T) {
	turns := []string{"A", "B", "C"}
	collector := new([]string)
	fns := map[string]workerFn{
		"A": fnCreactor(time.Hour, "A", collector, nil),
		"B": fnCreactor(time.Hour, "B", collector, nil),
		"C": fnCreactor(time.Hour, "C", collector, nil),
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	startTime := time.Now()

	err := workInTurns(ctx, turns, fns)

	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(*collector))
	assert.Less(t, int(time.Since(startTime).Seconds()), 2)
}

func TestWorkInTurns_WorkError(t *testing.T) {
	turns := []string{"A", "B", "C", "A"}
	collector := new([]string)
	fns := map[string]workerFn{
		"A": fnCreactor(time.Hour, "A", collector, errors.New("something wrong")),
		"B": fnCreactor(time.Hour, "B", collector, nil),
		"C": fnCreactor(time.Hour, "C", collector, nil),
	}

	err := workInTurns(context.Background(), turns, fns)

	assert.NotEqual(t, nil, err)
	assert.Equal(t, 0, len(*collector))
}
