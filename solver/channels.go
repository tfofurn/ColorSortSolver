package solver

const expectedSolvers = 16

type Channels struct {
	WorkerCount   chan int
	Solutions     chan *Step
	TerminalDepth chan uint
	MovesTried    chan uint
}

func NewChannels() Channels {
	return Channels{
		make(chan int),
		make(chan *Step, expectedSolvers),
		make(chan uint, expectedSolvers),
		make(chan uint, expectedSolvers),
	}
}
