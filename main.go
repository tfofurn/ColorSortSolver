package main

import (
	"colorsortsolver/solver"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func describeSolution(colorMap solver.ColorMap, solution []solver.Step) string {
	var b strings.Builder
	for index, step := range solution {
		capped := ""
		if step.Capped {
			capped = "Capped!"
		}
		fmt.Fprintf(&b, "%4d: %12v × %v: %s -> %s %s\n", index+1, colorMap.StringFromColor(step.Color), step.Amount, step.SourceTubeName, step.DestinationTubeName, capped)
	}
	return b.String()
}

func solutionListener(path string, channels solver.Channels, colorMap solver.ColorMap) {
	remainingWorkers := 0
	workerCount := 0
	solutionCount := 0
	shortestSolutionLength := 1000000
	shortestSolutionDescription := ""

	printer := message.NewPrinter(language.English)
	printer.Printf("%s: Start!\n", path)

	for {
		select {
		case solution := <-channels.Solutions:
			solutionCount++

			if len(solution) < shortestSolutionLength {
				var b strings.Builder
				printer.Fprintf(&b, "%s Solution %d, %d steps\n", path, solutionCount, len(solution))
				b.WriteString(describeSolution(colorMap, solution))

				shortestSolutionDescription = b.String()
				shortestSolutionLength = len(solution)
			}

			if solutionCount%10000 == 0 {
				printer.Printf("Solution %d, %d workers outstanding\n", solutionCount, remainingWorkers)
			}
		case increment := <-channels.WorkerCount:
			remainingWorkers += increment
			if increment > 0 {
				workerCount += 1
			}
			if remainingWorkers == 0 {
				printer.Printf("%s: All solvers have exited.  %d workers found %d solutions.\n", path, workerCount, solutionCount)
				fmt.Print(shortestSolutionDescription)
				return
			}
		}
	}
}

func processFile(inputPath string) {
	fmt.Println(inputPath)
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
