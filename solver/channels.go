package solver

type Channels struct {
	WorkerCount chan int
	Solutions   chan *Step
}

func NewChannels() Channels {
	return Channels{make(chan int), make(chan *Step, 16)}
}
