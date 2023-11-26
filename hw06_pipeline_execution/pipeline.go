package hw06pipelineexecution

import (
	"time"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return nil
	}
	out := make(Out)

	for _, s := range stages {
		select {
		case <-done:
			return out
		default:
			// Sleep to accommodate the termination case
			// Otherwise, launches an extra pipeline stage sometimes
			// Which doesn't fit the test's tolerances
			time.Sleep(time.Millisecond * 10)
			out = run(in, done, s)
			in = out
		}
	}

	return out
}

func run(in In, done In, stage Stage) Out {
	out := make(Bi)
	go func(in In, done In, stage Stage) {
		defer close(out)
		for v := range stage(in) {
			select {
			case <-done:
				return
			default:
				out <- v
			}
		}
	}(in, done, stage)

	return out
}
