package colorsortsolver

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

type Rack struct {
	steps []string
	tubes []Tube
}

func RackFromCSV(input string) Rack {

	csvReader := csv.NewReader(strings.NewReader(input))
	tubeNames, _ := csvReader.Read()
	tubes := make([]Tube, len(tubeNames))
	for iTube, name := range tubeNames {
		tubes[iTube] = NewTube(name)
	}

	for {
		colors, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		for iColor, colorName := range colors {
			tube := tubes[iColor]
			color := ColorFromString(colorName)
			if color != Empty && color != UnknownColor {
				tube.BottomFill(color)
			}
			tubes[iColor] = tube
		}
	}

	return Rack{make([]string, 0), tubes}
}

func (r *Rack) Move(sourceIndex, destinationIndex int) Rack {
	color, amount := r.tubes[sourceIndex].TopColor()
	moveDescription := fmt.Sprintf("%12v × %v: %s -> %s", StringFromColor(color), amount, r.tubes[sourceIndex].name, r.tubes[destinationIndex].name)

	steps := make([]string, len(r.steps), len(r.steps)+1)
	copy(steps, r.steps)
	steps = append(steps, moveDescription)
	tubes := make([]Tube, 0)
	for _, tube := range r.tubes {
		tubes = append(tubes, tube.Copy())
	}
	sourceTube := tubes[sourceIndex]
	destinationTube := tubes[destinationIndex]
	color, amount = sourceTube.PourOutTop()
	destinationTube.PourIn(color, amount)
	tubes[sourceIndex] = sourceTube
	tubes[destinationIndex] = destinationTube

	return Rack{steps: steps, tubes: tubes}
}

func (r *Rack) AttemptSolution(channels Channels) bool {
	if len(r.steps) == 1 {
		channels.WorkerCount <- 1
	}
	solved := false
	for srcIndex, srcTube := range r.tubes {
		for destIndex, destTube := range r.tubes {
			if srcIndex == destIndex {
				continue
			}
			if destTube.CanReceiveFrom(srcTube) {
				postMoveRack := r.Move(srcIndex, destIndex)
				solved = solved || postMoveRack.CheckSolved(channels)
				if !solved {
					if len(r.steps) == 0 {
						go postMoveRack.AttemptSolution(channels)
					} else {
						postMoveRack.AttemptSolution(channels)
					}
				}
			}
			if solved {
				continue
			}
		}
		if solved {
			continue
		}
	}

	if len(r.steps) == 1 {
		channels.WorkerCount <- -1
	}
	return solved
}

func (r *Rack) CheckSolved(channels Channels) bool {
	capped, empty, mixed := 0, 0, 0
	for index := range r.tubes {
		tube := r.tubes[index]
		switch {
		case tube.IsCapped():
			capped += 1
		case tube.IsEmpty():
			empty += 1
		default:
			mixed += 1
		}
	}
	if mixed == 0 {
		channels.Solutions <- r.steps
		return true
	}
	return false
}

func (r *Rack) TubeCount() int {
	return len(r.tubes)
}
