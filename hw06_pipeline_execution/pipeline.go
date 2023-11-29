package hw06pipelineexecution

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

	var out Bi
	for _, s := range stages {
		select {
		case <-done:
			return nil
		default:
			out = make(Bi)
			go func(in In, out Bi, s Stage) {
				defer close(out)
				for v := range s(in) {
					select {
					case <-done:
						return
					default:
						out <- v
					}
				}
			}(in, out, s)
			in = out
		}
	}

	return out
}
