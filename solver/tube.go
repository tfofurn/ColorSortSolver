package solver

import "fmt"

type Tube struct {
	name      string
	sections  []Color
	fillLevel int
}

const TubeHeight = 4

func NewTube(name string) Tube {
	result := Tube{name, make([]Color, TubeHeight), 0}
	return result
}

func (t *Tube) TopSection() (color Color, amount int) {
	amount = 0
	topIndex := t.fillLevel - 1
	for iSection := topIndex; iSection >= 0; iSection-- {
		if t.sections[iSection] != t.sections[topIndex] {
			break
		}
		amount++
	}
	return t.sections[topIndex], amount
}

func (t *Tube) Slack() int {
	return TubeHeight - t.fillLevel
}

func (t *Tube) PourIn(color Color, amount int) {
	for i := 0; i < amount; i++ {
		t.sections[t.fillLevel] = color
		t.fillLevel++
	}
}

func (t *Tube) BottomFill(color Color) {
	for iDestSection := t.fillLevel; iDestSection > 0; iDestSection-- {
		t.sections[iDestSection] = t.sections[iDestSection-1]
	}
	t.sections[0] = color
	t.fillLevel++
}

func (t *Tube) PourOutTop() (Color, int) {
	color, amount := t.TopSection()
	t.fillLevel -= amount
	return color, amount
}

func (destination *Tube) CanReceiveFrom(source Tube) bool {
	if source.IsEmpty() || source.IsCapped() {
		return false
	}
	if destination.IsEmpty() {
		return !source.IsSingleColor()
	}

	sourceColor, sourceAmount := source.TopSection()
	destinationColor, _ := destination.TopSection()
	if sourceAmount > destination.Slack() {
		return false
	}

	return sourceColor == destinationColor
}

func (t *Tube) IsCapped() bool {
	if t.IsEmpty() {
		return false
	}
	_, amount := t.TopSection()
	return amount == TubeHeight
}

func (t *Tube) IsSingleColor() bool {
	if t.IsEmpty() {
		return false
	}
	_, amount := t.TopSection()
	return amount == t.fillLevel
}

func (t *Tube) IsEmpty() bool {
	return t.fillLevel == 0
}

func (t *Tube) Copy() Tube {
	result := NewTube(t.name)
	result.sections = make([]Color, TubeHeight)
	copy(result.sections, t.sections)
	result.fillLevel = t.fillLevel
	return result
}

func (t *Tube) Describe() string {
	return fmt.Sprintf("{{%s: %v}}", t.name, t.sections)
}
