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

func (t *Tube) TopColor() (Color, int) {
	if len(t.colors) == 0 {
		return Empty, 0
	}

	amount := 0
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
	color, amount := t.TopColor()
	t.colors = t.colors[amount:]
	return color, amount
}

func (destination *Tube) CanReceiveFrom(source Tube) bool {
	if source.IsCapped() || source.IsEmpty() {
		return false
	}
	sourceColor, sourceAmount := source.TopColor()
	destinationColor, _ := destination.TopColor()
	if sourceAmount > destination.Slack() {
		return false
	}
	if destination.IsEmpty() {
		return !source.IsSingleColor()
	}

	return sourceColor == destinationColor
}

func (t *Tube) IsCapped() bool {
	_, amount := t.TopColor()
	return amount == TubeHeight
}

func (t *Tube) IsSingleColor() bool {
	if t.IsEmpty() {
		return false
	}
	_, amount := t.TopColor()
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
