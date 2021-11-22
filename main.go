package main

import (
	"colorsortsolver/solver"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func solutionListener(path string, channels solver.Channels, colorMap solver.ColorMap) {
	remainingWorkers := 0
	workerCount := 0
	solutionCount := 0
	shortestSolution := 1000000

	printer := message.NewPrinter(language.English)

	for {
		select {
		case solution := <-channels.Solutions:
			solutionCount++

			if len(solution) < shortestSolution {
				printer.Printf("%s Solution %d, %d steps\n", path, solutionCount, len(solution))
				for index, step := range solution {
					capped := ""
					if step.Capped {
						capped = "Capped!"
					}
					fmt.Printf("%4d: %12v × %v: %s -> %s %s\n", index+1, colorMap.StringFromColor(step.Color), step.Amount, step.SourceTubeName, step.DestinationTubeName, capped)
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
				printer.Printf("%s: All solvers have exited.  %d/%d workers found a valid solution.", path, solutionCount, workerCount)
				return
			}
		}
	}
}

func processFile(inputPath string) {
	fileContents, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Print(err)
	}
	fileContentsString := string(fileContents)
	colorMap := solver.NewColorMap()
	baseRack := solver.RackFromCSV(&colorMap, fileContentsString)

	channels := solver.NewChannels()

	baseRack.AttemptSolution(channels)
	solutionListener(inputPath, channels, colorMap)
}

func main() {
	executablePath, _ := filepath.Abs(os.Args[0])
	paths := []string{filepath.Join(filepath.Dir(executablePath), "sample", "*.csv")}
	if len(os.Args) > 1 {
		paths = os.Args[1:]
	}
	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Println(err)
		}
		for _, file := range matches {
			processFile(file)
		}
	}
}
