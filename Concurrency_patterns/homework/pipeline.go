package executor

import (
	"context"
)

type (
	In  <-chan any
	Out = In
)

type Stage func(in In) (out Out)

func ExecutePipeline(ctx context.Context, in In, stages ...Stage) Out {
	channels := make([]chan any, len(stages))
	for i := range stages {
		channels[i] = make(chan any)
		if i == 0 {
			go executeStage(ctx, in, channels[i], stages[i])
		} else {
			go executeStage(ctx, channels[i-1], channels[i], stages[i])
		}
	}

	return channels[len(channels)-1]
}

func executeStage(ctx context.Context, in In, out chan any, stage Stage) {
	defer close(out)
	for v := range stage(in) {
		select {
		case out <- v:
		case <-ctx.Done():
			return
		}
	}
}
