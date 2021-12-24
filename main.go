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
	const dividerSpacing = 5
	for index, step := range solution {
		capped := ""
		if step.Capped {
			capped = "Capped!"
		}
		fmt.Fprintf(&b, "%4d: %12v × %v: %s -> %s %s\n", index+1, colorMap.StringFromColor(step.Color), step.Amount, step.SourceTubeName, step.DestinationTubeName, capped)
		if index%dividerSpacing == dividerSpacing-1 {
			fmt.Fprintln(&b, "     -----")
		}
	}
	return b.String()
}

func solutionListener(path string, channels solver.Channels, colorMap solver.ColorMap, printSolution bool) {
	remainingWorkers := 0
	workerCount := 0
	solutionCount := 0
	var shortestSolution []solver.Step
	shortestSolutionHeader := ""

	ticker := time.NewTicker(time.Second * 5)
	startTime := time.Now()

	printer := message.NewPrinter(language.English)

	for done := false; !done; {
		select {
		case now := <-ticker.C:
			if solutionCount > 0 {
				elapsed := int(now.Sub(startTime).Seconds())
				printer.Printf("%6d seconds elapsed.  %d Solutions found, %d workers outstanding.  Shortest found: %d\n", elapsed, solutionCount, remainingWorkers, len(shortestSolution))
			}
		case solution := <-channels.Solutions:
			solutionCount++

			if len(shortestSolution) == 0 || len(solution) < len(shortestSolution) {
				shortestSolution = solution
				shortestSolutionHeader = printer.Sprintf("%s Solution %d, %d steps\n", path, solutionCount, len(solution))
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
	ticker.Stop()
	printer.Print(shortestSolutionHeader)
	if printSolution {
		printer.Print(describeSolution(colorMap, shortestSolution))
	}
}

func processFile(inputPath string, printSolution bool) (elapsedMilliseconds int) {
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
	solutionListener(inputPath, channels, colorMap, printSolution)
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
	inputFiles := make([]string, 0)
	var totalTimeMillis int
	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Println(err)
		}
		inputFiles = append(inputFiles, matches...)
	}
	printSolution := len(inputFiles) == 1
	for _, inputfile := range inputFiles {
		totalTimeMillis += processFile(inputfile, printSolution)
	}
	if len(inputFiles) > 1 {
		printer := message.NewPrinter(language.English)
		printer.Printf("Processed %d files.  Total time: %d milliseconds.\n", len(inputFiles), totalTimeMillis)
	}
}
