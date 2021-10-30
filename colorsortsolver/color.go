package colorsortsolver

import "fmt"

type Color int

const (
	UnknownColor Color = iota
	Empty
	Pink
	Red
	Brown
	Cyan
	DarkGreen
	MediumGreen
	LightGreen
	Yellow
	Orange
	Purple
	Lavendar
	Blue
)

func StringFromColor(input Color) string {
	switch input {
	case Pink:
		return "Pink"
	case Red:
		return "Red"
	case Brown:
		return "Brown"
	case Cyan:
		return "Cyan"
	case DarkGreen:
		return "DarkGreen"
	case MediumGreen:
		return "MediumGreen"
	case LightGreen:
		return "LightGreen"
	case Yellow:
		return "Yellow"
	case Orange:
		return "Orange"
	case Purple:
		return "Purple"
	case Lavendar:
		return "Lavendar"
	case Blue:
		return "Blue"
	}
	return "unknown"
}
func ColorFromString(input string) Color {
	switch input {
	case "":
		return Empty
	case "Empty":
		return Empty
	case "Pink":
		return Pink
	case "Red":
		return Red
	case "Brown":
		return Brown
	case "Cyan":
		return Cyan
	case "DarkGreen":
		return DarkGreen
	case "MedGreen":
		return MediumGreen
	case "LightGreen":
		return LightGreen
	case "Yellow":
		return Yellow
	case "Orange":
		return Orange
	case "Purple":
		return Purple
	case "Lavendar":
		return Lavendar
	case "Blue":
		return Blue
	default:
		fmt.Printf("Unknown color: %s", input)
		return UnknownColor
	}
}
