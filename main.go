package main

import (
	"colorsortsolver/colorsortsolver"
	"fmt"
	"io/ioutil"
)

func main() {
	fileContents, err := ioutil.ReadFile("c:/Users/Scott/Downloads/SortPuz #255 - Sheet1.csv")
	if err != nil {
		fmt.Print(err)
	}
	fileContentsString := string(fileContents)
	// fmt.Print(fileContentsString)
	baseRack := colorsortsolver.RackFromCSV(fileContentsString)
	baseRack.AttemptSolution()
}
