package main

import (
	"colorsortsolver/solver"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	var shortestSolution []solver.Step
	shortestSolutionHeader := ""

	printer := message.NewPrinter(language.English)
	printer.Printf("%s: Start!\n", path)

	for done := false; !done; {
		select {
		case solution := <-channels.Solutions:
			solutionCount++

			if len(shortestSolution) == 0 || len(solution) < len(shortestSolution) {
				shortestSolution = solution
				shortestSolutionHeader = printer.Sprintf("%s Solution %d, %d steps\n", path, solutionCount, len(solution))
			}

			if solutionCount%1000 == 0 {
				printer.Printf("Solution %d, %d workers outstanding\n", solutionCount, remainingWorkers)
			}
		case increment := <-channels.WorkerCount:
			remainingWorkers += increment
			if increment > 0 {
				workerCount += 1
			}
			if remainingWorkers == 0 {
				printer.Printf("%s: All solvers have exited.  %d workers found %d solutions.\n", path, workerCount, solutionCount)
				done = true
			}
		}
	}
	printer.Print(shortestSolutionHeader)
}

func processFile(inputPath string) (elapsedMilliseconds int) {
	fmt.Println(inputPath)
	fileContents, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Print(err)
	}
	fileContentsString := string(fileContents)
	colorMap := solver.NewColorMap()
	baseRack := solver.RackFromCSV(&colorMap, fileContentsString)

	channels := solver.NewChannels()

	start := time.Now()
	go baseRack.AttemptSolution(channels)
	solutionListener(inputPath, channels, colorMap)
	end := time.Now()
	elapsed := int(end.Sub(start) / 1000000)
	fmt.Printf("%s: elapsed time: %d milliseconds\n", inputPath, elapsed)
	return elapsed
}

func main() {
	executablePath, _ := filepath.Abs(os.Args[0])
	paths := []string{filepath.Join(filepath.Dir(executablePath), "sample", "*.csv")}
	if len(os.Args) > 1 {
		paths = os.Args[1:]
	}
	filesProcessed := 0
	var totalTimeMillis int
	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Println(err)
		}
		for _, file := range matches {
			totalTimeMillis += processFile(file)
			filesProcessed++
		}
	}
	if filesProcessed > 1 {
		printer := message.NewPrinter(language.English)
		printer.Printf("Processed %d files.  Total time: %d milliseconds.\n", filesProcessed, totalTimeMillis)
	}
}
