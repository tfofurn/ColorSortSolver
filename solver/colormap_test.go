package solver

import "testing"

func checkStringFromColor(t *testing.T, colorMap ColorMap, color Color, expectedName string) {
	actualName := colorMap.StringFromColor(color)
	if actualName != expectedName {
		t.Fatalf("ColorMap should have returned %s, returned %s", expectedName, actualName)
	}
}

func TestColorMap(t *testing.T) {
	colorMap := NewColorMap()
	firstName := "first"
	secondName := "second"
	firstColor := colorMap.ColorFromString(firstName)
	secondColor := colorMap.ColorFromString(secondName)

	if firstColor == secondColor {
		t.Fatalf("Two different names returned same color: %v", firstColor)
	}

	checkStringFromColor(t, colorMap, firstColor, firstName)
	checkStringFromColor(t, colorMap, secondColor, secondName)
}
