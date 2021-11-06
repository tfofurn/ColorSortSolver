package main

import (
	"colorsortsolver/solver"
	"fmt"
	"io/ioutil"
	"math"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func solutionListener(channels solver.Channels, colorMap solver.ColorMap) {
	remainingWorkers := 0
	workerCount := 0
	solutionCount := 0
	shortestSolution := math.MaxInt

	printer := message.NewPrinter(language.English)

	for {
		select {
		case solution := <-channels.Solutions:
			solutionCount++

			if len(solution) < shortestSolution {
				printer.Printf("Solution %d, %d steps\n", solutionCount, len(solution))
				for index, step := range solution {
					fmt.Printf("%4d: %12v × %v: %s -> %s\n", index+1, colorMap.StringFromColor(step.Color), step.Amount, step.SourceTubeName, step.DestinationTubeName)
				}
				shortestSolution = len(solution)
			}

			if solutionCount%1000000 == 0 {
				printer.Printf("Solution %d, %d workers outstanding\n", solutionCount, remainingWorkers)
			}
		case increment := <-channels.WorkerCount:
			remainingWorkers += increment
			if increment > 0 {
				workerCount += 1
			}
			if remainingWorkers == 0 {
				printer.Printf("All solvers have exited.  %d/%d workers found a valid solution.", solutionCount, workerCount)
				return
			}
		}
	}
}

func main() {
	fileContents, err := ioutil.ReadFile("sample/inciting-incident.csv")
	if err != nil {
		fmt.Print(err)
	}
	fileContentsString := string(fileContents)
	colorMap := solver.NewColorMap()
	baseRack := solver.RackFromCSV(&colorMap, fileContentsString)

	channels := solver.NewChannels()

	baseRack.AttemptSolution(channels)
	solutionListener(channels, colorMap)
}
