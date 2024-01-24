package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func stageWrapper(in In, done In, stage Stage) Out {
	inputCh := make(Bi)

	if stage == nil {
		return in
	}

	go func() {
		defer close(inputCh)
		for {
			select {
			case <-done:
				return
			case val, open := <-in:
				if !open {
					return
				}
				inputCh <- val
			}
		}
	}()

	return stage(inputCh)
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil || len(stages) < 1 {
		return nil
	}
	lastOut := in
	for _, stage := range stages {
		lastOut = stageWrapper(lastOut, done, stage)
	}

	return lastOut
}
