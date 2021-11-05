package solver

import "testing"

func TestTopEmpty(t *testing.T) {
	tube := NewTube("T1")
	if !tube.IsEmpty() {
		t.Fatal("Empty tube should report IsEmpty")
	}
	if tube.IsSingleColor() || tube.IsCapped() {
		t.Fatalf("Empty tubes should not report IsSingleColor or IsCapped.")
	}

	color, count := tube.TopColor()
	if color != Empty {
		t.Fatalf("Empty tube's top color should be Empty, not %v", color)
	}
	if count != 0 {
		t.Fatalf("Empty tube's top count should be 0, not %v", count)
	}
	if tube.Slack() != TubeHeight {
		t.Fatalf(`Empty tube's slack should be %d, got %d`, TubeHeight, tube.Slack())
	}
}

func TestSingleColor(t *testing.T) {
	tube := NewTube("T1")
	color := Brown
	emptyTube := NewTube("E")

	slack := tube.Slack()
	for iPour := 0; iPour < TubeHeight; iPour++ {
		tube.PourIn(color, 1)
		if tube.IsEmpty() {
			t.Fatal("Tube we just poured into should not report IsEmpty")
		}
		if !tube.IsSingleColor() {
			t.Fatal("We've only added one color; tube should report IsSingleColor.")
		}

		if emptyTube.CanReceiveFrom(tube) {
			t.Fatal("Empty tubes should not receive from single-color tube")
		}

		topColor, topAmount := tube.TopColor()
		if topColor != color {
			t.Fatalf(`top color should be %v, got %v`, color, topColor)
		}
		if topAmount != iPour+1 {
			t.Fatalf(`top count should be %v, got %v`, iPour+1, topAmount)
		}

		if tube.Slack() >= slack {
			t.Fatalf(`tube slack %d should have decreased from %d`, tube.Slack(), slack)
		}
		slack = tube.Slack()
		if (slack == 0) != tube.IsCapped() {
			t.Fatalf("Tube with slack %v reported Capped of %v", slack, tube.IsCapped())
		}

		copy := tube.Copy()
		copyColor, copyAmount := copy.TopColor()
		if copyColor != topColor {
			t.Fatalf("Copied single-color tube color should be %v, got %v", topColor, copyColor)
		}
		if copyAmount != topAmount {
			t.Fatalf("Copied single-color tube amount should be %v, got %v", topAmount, copyAmount)
		}
	}
}

func TestInclusion(t *testing.T) {
	tube := NewTube("T17")
	firstColor, secondColor := MediumGreen, DarkGreen
	tube.PourIn(firstColor, 1)
	tube.PourIn(secondColor, 1)

	for expected := 1; expected <= 2; expected++ {
		tube.PourIn(firstColor, 1)

		topColor, topCount := tube.TopColor()
		if topColor != firstColor {
			t.Fatalf("Top color should be %v, got %v", firstColor, topColor)
		}
		if topCount != expected {
			t.Fatalf("Top count should be %v, got %v", expected, topCount)
		}
	}
}

func TestPourOutTop(t *testing.T) {
	tube := NewTube("T5")
	stillColor, movingColor := Cyan, Blue
	stillColorAmount, movingColorAmount := 1, 2
	tube.PourIn(stillColor, stillColorAmount)
	tube.PourIn(movingColor, movingColorAmount)

	color, count := tube.PourOutTop()
	if color != movingColor {
		t.Fatalf("poured out color should be %v, got %v", movingColor, color)
	}
	if count != movingColorAmount {
		t.Fatalf("poured out amount should be %v, got %v", movingColorAmount, count)
	}

	color, count = tube.PourOutTop()
	if color != stillColor {
		t.Fatalf("second poured out color should be %v, got %v", stillColor, color)
	}
	if count != stillColorAmount {
		t.Fatalf("second poured out amount should be %v, got %v", stillColorAmount, count)
	}

	if !tube.IsEmpty() {
		t.Fatal("Poured everything out, but didn't report IsEmpty!")
	}
}

func TestEmptyCanReceive(t *testing.T) {
	destination := NewTube("E")
	baseColor, topColor := LightGreen, Orange

	for iAmount := 1; iAmount < TubeHeight-1; iAmount++ {
		source := NewTube("S")
		source.PourIn(baseColor, 1)
		source.PourIn(topColor, iAmount)
		can := destination.CanReceiveFrom(source)
		if !can {
			t.Fatalf("Unable to pour into empty from tube with amount %v", iAmount)
		}
	}

	fullSource := NewTube("F")
	fullSource.PourIn(topColor, TubeHeight)
	can := destination.CanReceiveFrom(fullSource)
	if can {
		t.Fatal("Empty should not accept from a full source")
	}

	emptySource := NewTube("E2")
	can = destination.CanReceiveFrom(emptySource)
	if can {
		t.Fatal("Empty should not accept from empty")
	}
}

func TestMatchCanReceive(t *testing.T) {
	color := Pink
	for sourceAmount := 1; sourceAmount <= TubeHeight; sourceAmount++ {
		for destinationAmount := 1; destinationAmount <= TubeHeight; destinationAmount++ {
			sourceTube, destinationTube := NewTube("S"), NewTube("D")
			sourceTube.PourIn(color, sourceAmount)
			destinationTube.PourIn(color, destinationAmount)
			can := destinationTube.CanReceiveFrom(sourceTube)
			expected := (sourceAmount + destinationAmount) <= TubeHeight
			if can != expected {
				t.Fatalf("For matching color, source %v -> dest %v, CanReceiveFrom should be %v", sourceAmount, destinationAmount, expected)
			}
		}
	}
}

func TestMismatchCannotReceive(t *testing.T) {
	sourceColor, destinationColor := Red, Blue
	sourceTube, destinationTube := NewTube("S"), NewTube("D")
	sourceTube.PourIn(sourceColor, 1)
	destinationTube.PourIn(destinationColor, 1)
	if destinationTube.CanReceiveFrom(sourceTube) {
		t.Fatal("Tubes with mismatched colors should not allow transfer.")
	}
	emptySource := NewTube("E")
	if destinationTube.CanReceiveFrom(emptySource) {
		t.Fatal("Tubes should be be able to receive from Empty")
	}
}

func TestCopyEmpty(t *testing.T) {
	empty := NewTube("e")
	emptyCopy := empty.Copy()
	if emptyCopy.Slack() != TubeHeight {
		t.Fatalf("Copied empty tube should have slack %v, found %v", TubeHeight, emptyCopy.Slack())
	}
	copyColor, copyAmount := emptyCopy.TopColor()
	if copyColor != Empty {
		t.Fatalf("Copied empty should have top color %v, found %v", Empty, copyColor)
	}
	if copyAmount != 0 {
		t.Fatalf("Copied empty should have top amount %v, found %v", 0, copyAmount)
	}
}

func checkForColor(t *testing.T, tube Tube, expected Color) {
	color, _ := tube.PourOutTop()
	if color != expected {
		t.Fatalf("First from BottomFill should have been %v, got %v", expected, color)
	}
}

func TestBottomFill(t *testing.T) {
	tube := NewTube("t")
	tube.BottomFill(Brown)
	tube.BottomFill(Orange)
	tube.BottomFill(Pink)
	tube.BottomFill(Yellow)

	checkForColor(t, tube, Brown)
	checkForColor(t, tube, Orange)
	checkForColor(t, tube, Pink)
	checkForColor(t, tube, Yellow)
}