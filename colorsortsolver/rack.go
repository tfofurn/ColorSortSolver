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
	fmt.Println(tubeNames)
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
			// fmt.Printf("BottomFill %v into %v\n", color, tube.Describe())
			if color != Empty && color != UnknownColor {
				tube.BottomFill(color)
			}
			tubes[iColor] = tube
		}
	}

	for _, tube := range tubes {
		fmt.Println(tube)
	}

	return Rack{make([]string, 0), tubes}
}

func (r *Rack) Move(sourceIndex, destinationIndex int) Rack {
	color, amount := r.tubes[sourceIndex].TopColor()
	moveDescription := fmt.Sprintf("%10vx%v:%s->%s", StringFromColor(color), amount, r.tubes[sourceIndex].name, r.tubes[destinationIndex].name)

	steps := make([]string, len(r.steps), len(r.steps)+1)
	copy(steps, r.steps)
	steps = append(steps, moveDescription)
	if len(steps) < 6 {
		fmt.Println(steps)
	}
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

func (r *Rack) AttemptSolution() bool {
	solved := false
	for srcIndex, srcTube := range r.tubes {
		for destIndex, destTube := range r.tubes {
			if srcIndex == destIndex {
				continue
			}
			if destTube.CanReceiveFrom(srcTube) {
				postMoveRack := r.Move(srcIndex, destIndex)
				solved = solved || postMoveRack.CheckSolved()
				if !solved {
					postMoveRack.AttemptSolution()
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
	return solved
}

func (r *Rack) CheckSolved() bool {
	capped, empty, mixed := 0, 0, 0
	for index, _ := range r.tubes {
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
		fmt.Printf("Solved! c %d e %d m %d: \n", capped, empty, mixed)
		for _, step := range r.steps {
			fmt.Printf("  %s\n", step)
		}
		return true
	}
	return false
}

func (r *Rack) TubeCount() int {
	return len(r.tubes)
}
