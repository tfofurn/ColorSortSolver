package main

import (
	"colorsortsolver/colorsortsolver"
	"fmt"
	"io/ioutil"
	"math"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func solutionListener(solutionsChannel chan []string) {
	count := 0
	shortestSolution := math.MaxInt

	printer := message.NewPrinter(language.English)

	for solution := range solutionsChannel {
		count++

		if len(solution) < shortestSolution {
			printer.Printf("Solution %d, %d steps\n", count, len(solution))
			for index, step := range solution {
				fmt.Printf("%4d: %s\n", index+1, step)
			}
			shortestSolution = len(solution)
		}

		if count%1000000 == 0 {
			printer.Printf("Solution %d, %d steps\n", count, len(solution))
		}
	}
}

func main() {
	fileContents, err := ioutil.ReadFile("sample/inciting-incident.csv")
	if err != nil {
		fmt.Print(err)
	}
	fileContentsString := string(fileContents)
	baseRack := colorsortsolver.RackFromCSV(fileContentsString)

	solutionsChannel := make(chan []string, 100)

	baseRack.AttemptSolution(solutionsChannel)
	solutionListener(solutionsChannel)
}
