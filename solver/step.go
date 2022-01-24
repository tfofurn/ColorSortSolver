package solver

type Step struct {
	Color                               Color
	Amount                              int
	SourceTubeName, DestinationTubeName string
	Capped                              bool
	Previous                            *Step
	Index                               int
}
