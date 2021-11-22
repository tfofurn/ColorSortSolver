package solver

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

type Rack struct {
	steps []Step
	tubes []Tube
}

func RackFromCSV(colorMap *ColorMap, input string) Rack {
	csvReader := csv.NewReader(strings.NewReader(input))
	rawTubeNames, _ := csvReader.Read()
	tubes := make([]Tube, len(rawTubeNames))
	for iTube, rawName := range rawTubeNames {
		tubes[iTube] = NewTube(strings.TrimSpace(rawName))
	}

	colorCounts := make([]int, len(tubes))

	for {
		colors, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		for iColor, rawColorName := range colors {
			colorName := strings.TrimSpace(rawColorName)
			if colorName != "" {
				tube := &tubes[iColor]
				color := colorMap.ColorFromString(colorName)
				colorCounts[color] += 1
				tube.BottomFill(color)
			}
		}
	}

	for color, count := range colorCounts {
		if count == TubeHeight || count == 0 {
			continue
		}
		panic(fmt.Errorf("each color expected to appear %d times, found %d for %s", TubeHeight, count, colorMap.StringFromColor(Color(color))))
	}

	return Rack{make([]Step, 0), tubes}
}

func (r *Rack) Move(sourceIndex, destinationIndex int) Rack {
	color, amount := r.tubes[sourceIndex].TopColor()

	
	tubes := make([]Tube, len(r.tubes))
	copy(tubes, r.tubes)
	sourceTube := tubes[sourceIndex].Copy()
	destinationTube := tubes[destinationIndex].Copy()
	color, amount = sourceTube.PourOutTop()
	destinationTube.PourIn(color, amount)
	tubes[sourceIndex] = sourceTube
	tubes[destinationIndex] = destinationTube

	moveDescription := Step{color, amount, r.tubes[sourceIndex].name, r.tubes[destinationIndex].name, destinationTube.IsCapped()}

	steps := make([]Step, len(r.steps), len(r.steps)+1)
	copy(steps, r.steps)
	steps = append(steps, moveDescription)

	return Rack{steps: steps, tubes: tubes}
}

func (r *Rack) AttemptSolution(channels Channels) bool {
	if len(r.steps) == 1 {
		channels.WorkerCount <- 1
	}
	solved := false
	for srcIndex, srcTube := range r.tubes {
		sourceTube := r.tubes[srcIndex]
		if sourceTube.IsEmpty() || sourceTube.IsCapped() {
			continue
		}
		for destIndex, destTube := range r.tubes {
			if srcIndex == destIndex {
				continue
			}
			if destTube.CanReceiveFrom(srcTube) {
				postMoveRack := r.Move(srcIndex, destIndex)
				if postMoveRack.tubes[srcIndex].IsEmpty() && postMoveRack.tubes[destIndex].IsCapped() {
					solved = solved || postMoveRack.CheckSolved(channels)
				}
				if !solved {
					if len(r.steps) == 0 {
						go postMoveRack.AttemptSolution(channels)
					} else {
						solved = solved || postMoveRack.AttemptSolution(channels)
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
