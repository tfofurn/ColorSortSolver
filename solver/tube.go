package solver

import "fmt"

type Tube struct {
	name   string
	colors []Color
}

const TubeHeight = 4

func NewTube(name string) Tube {
	result := Tube{name, make([]Color, 0, TubeHeight)}
	return result
}

func (t *Tube) TopSection() (color Color, amount int) {
	amount = 0
	for iColor := 0; iColor < len(t.colors); iColor++ {
		if t.colors[iColor] != t.colors[0] {
			break
		}
		amount++
	}
	return t.colors[0], amount
}

func (t *Tube) Slack() int {
	return TubeHeight - len(t.colors)
}

func (t *Tube) PourIn(color Color, amount int) {
	prefix := make([]Color, amount)
	for i := 0; i < amount; i++ {
		prefix[i] = color
	}
	t.colors = append(prefix, t.colors...)
}

func (t *Tube) BottomFill(color Color) {
	t.colors = append(t.colors, color)
	// fmt.Printf("%v: %v\n", t.name, t.colors)
}

func (t *Tube) PourOutTop() (Color, int) {
	color, amount := t.TopSection()
	t.colors = t.colors[amount:]
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
	return amount == len(t.colors)
}

func (t *Tube) IsEmpty() bool {
	return len(t.colors) == 0
}

func (t *Tube) Copy() Tube {
	result := NewTube(t.name)
	result.colors = make([]Color, len(t.colors))
	copy(result.colors, t.colors)
	return result
}

func (t *Tube) Describe() string {
	return fmt.Sprintf("{{%s: %v}}", t.name, t.colors)
}
