package colorsortsolver

type Channels struct {
	WorkerCount chan int
	Solutions   chan []string
}

func NewChannels() Channels {
	return Channels{make(chan int), make(chan []string)}
}
