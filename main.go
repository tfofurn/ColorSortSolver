package main

import (
	"colorsortsolver/solver"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func describeSolution(colorMap solver.ColorMap, lastStep *solver.Step) string {
	var b strings.Builder
	const dividerSpacing = 5
	solution := make([]*solver.Step, lastStep.Index+1)
	for step := lastStep; step != nil; step = step.Previous {
		solution[step.Index] = step
	}

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
	var shortestSolution *solver.Step
	var solutionToPrint *solver.Step
	shortestSolutionHeader := ""
	maxTerminalDepth := uint(0)
	terminationCount := uint64(0)
	moveCount := uint64(0)

	ticker := time.NewTicker(time.Second * 5)
	startTime := time.Now()

	printer := message.NewPrinter(language.English)

	for done := false; !done; {
		select {
		case now := <-ticker.C:
			if solutionCount > 0 {
				elapsed := int(now.Sub(startTime).Seconds())
				printer.Printf("%6d seconds elapsed.  %d Solutions found, %d workers outstanding.  Max Depth: %d. Moves: %d\n",
					elapsed, solutionCount, remainingWorkers, maxTerminalDepth, moveCount)
			}
		case solution := <-channels.Solutions:
			solutionCount++

			if shortestSolution == nil || solution.Index < shortestSolution.Index {
				shortestSolution = solution
				shortestSolutionHeader = printer.Sprintf("%s Solution %d, %d steps\n", path, solutionCount, shortestSolution.Index+1)
				if printSolution && solutionCount == 1 {
					printer.Print(shortestSolutionHeader)
					printer.Print(describeSolution(colorMap, shortestSolution))
				} else {
					solutionToPrint = solution
				}

			}
		case increment := <-channels.WorkerCount:
			remainingWorkers += increment
			if increment > 0 {
				workerCount += 1
			}
			if remainingWorkers == 0 {
				done = true
			}
		case depth := <-channels.TerminalDepth:
			terminationCount += 1
			if depth > maxTerminalDepth {
				maxTerminalDepth = depth
			}

		case moveIncrement := <-channels.MovesTried:
			moveCount += uint64(moveIncrement)
		}

	}
	ticker.Stop()
	printer.Printf("%s: All solvers have exited.  %d workers found %d solutions.  Max Depth: %d. Moves: %d\n", path, workerCount, solutionCount, maxTerminalDepth, moveCount)
	printer.Print(shortestSolutionHeader)
	if printSolution && solutionToPrint != nil {
		printer.Print(describeSolution(colorMap, solutionToPrint))
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
	go baseRack.AttemptSolution(&channels)
	solutionListener(inputPath, channels, colorMap, printSolution)
	end := time.Now()
	elapsed := int(end.Sub(start) / 1000000)
	fmt.Printf("%s: elapsed time: %d milliseconds\n", inputPath, elapsed)
	return elapsed
}

func processFiles(inputFiles []string) {
	var totalTimeMillis int
	printSolution := len(inputFiles) == 1
	if !printSolution {
		fmt.Printf("%d files found.\n", len(inputFiles))
	}
	for _, inputfile := range inputFiles {
		totalTimeMillis += processFile(inputfile, printSolution)
	}
	if len(inputFiles) > 1 {
		printer := message.NewPrinter(language.English)
		printer.Printf("Processed %d files.  Total time: %d milliseconds.\n", len(inputFiles), totalTimeMillis)
	}
}

func getSamplePaths() []string {
	executablePath, _ := filepath.Abs(os.Args[0])
	sampleRoot := filepath.Join(filepath.Dir(executablePath), "sample/")
	result := []string{}
	filepath.WalkDir(sampleRoot, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && filepath.Ext(path) == ".csv" {
			result = append(result, path)
		}
		return nil
	})
	return result
}

func modTimeFromPath(path string) time.Time {
	stat, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Unable to stat %s\n", path)
		return time.Unix(0, 0)
	}
	return stat.ModTime()
}

func expandGlobs(globs []string) []string {
	result := make([]string, 0)
	for _, path := range globs {
		matches, err := filepath.Glob(path)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, matches...)
	}
	return result
}

func main() {
	paths := make([]string, 0)
	if len(os.Args) > 1 {
		if os.Args[1] == "latest" {
			latest := getLatestSample()
			paths = []string{latest}
		} else {
			paths = os.Args[1:]
		}
	}
	inputFiles := expandGlobs(paths)
	processFiles(inputFiles)
}

func getLatestSample() string {
	samples := expandGlobs(getSamplePaths())
	latest := samples[0]
	latestModTime := modTimeFromPath(latest)
	for _, candidate := range samples {
		candidateModTime := modTimeFromPath(candidate)
		if candidateModTime.After(latestModTime) {
			latest = candidate
			latestModTime = modTimeFromPath(candidate)
		}
	}
	return latest
}
